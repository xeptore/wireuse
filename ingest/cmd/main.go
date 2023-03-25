package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/xeptore/wireuse/ingest"
	"github.com/xeptore/wireuse/pkg/env"
	"github.com/xeptore/wireuse/pkg/funcutils"
)

var (
	watchDir     string
	wgDeviceName string
)

func main() {
	ctx := context.Background()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if err := godotenv.Load(); nil != err {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal().Err(err).Msg("unexpected error while loading .env file")
		}
		log.Warn().Msg(".env file not found")
	}

	tz := env.MustGet("TZ")
	if tz != "UTC" {
		log.Fatal().Msg("TZ environment variable must be set to UTC")
	}

	flag.StringVar(&watchDir, "w", "", "directory path to watch for wireguard device up and dump files")
	flag.StringVar(&wgDeviceName, "i", "", "wireguard interface")

	flag.Parse()
	if nonFlagArgs := flag.Args(); len(nonFlagArgs) > 0 {
		log.Fatal().Msgf("expected no additional flags, got: %s", strings.Join(nonFlagArgs, ","))
	}
	if watchDir == "" {
		log.Fatal().Msg("directory path to watch for wireguard device up and dump files option is required and cannot be empty")
	}
	if wgDeviceName == "" {
		log.Fatal().Msg("wireguard device name option is required and cannot be empty")
	}

	uri := env.MustGet("MONGODB_URI")
	uriOption := options.Client().ApplyURI(uri)
	if err := uriOption.Validate(); nil != err {
		log.Fatal().Err(err).Msg("invalid value is set for 'MONGODB_URI' environment variable")
	}
	client, err := mongo.Connect(ctx, uriOption.SetMaxConnIdleTime(time.Minute).SetMaxConnecting(4).SetServerSelectionTimeout(5*time.Second).SetSocketTimeout(3*time.Second).SetRetryReads(true).SetRetryWrites(true))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	if err := client.Ping(ctx, readpref.Primary()); nil != err {
		log.Fatal().Err(err).Msg("failed to verify database connectivity")
	}
	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		if err := client.Disconnect(disconnectCtx); err != nil {
			log.Err(err).Msg("failed to disconnect from database")
			return
		}
		log.Info().Msg("successfully disconnected from database")
	}()
	cs, _ := connstring.Parse(uri)
	collection := client.Database(cs.Database).Collection(wgDeviceName)

	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "publicKey", Value: "hashed"}},
		},
		{
			Keys:    bson.D{{Key: "publicKey", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	names, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database indexes")
	}
	log.Info().Strs("index_names", names).Msg("successfully inserted database indexes")

	wg, err := wgctrl.New()
	if nil != err {
		log.Fatal().Err(err).Msg("failed to initialize wg control client")
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	runCtx, cancelEngineRun := context.WithCancelCause(ctx)
	stopSignalErr := errors.New("stop signal received")
	go func() {
		<-signals
		cancelEngineRun(stopSignalErr)
	}()

	wgUpEvents, wgDownEvents := make(chan ingest.WgUpEvent), make(chan ingest.WgDownEvent)
	go func() {
		log = log.With().Str("app", "watcher").Logger()
		w, err := fsnotify.NewWatcher()
		if nil != err {
			log.Fatal().Err(err).Msg("failed to initialize")
		}
		if err := w.Add(watchDir); nil != err {
			log.Fatal().Err(err).Msg("failed to add directory")
		}
		log.Debug().Str("dir", watchDir).Msg("started watching directory")
		for {
			select {
			case <-runCtx.Done():
				log.Info().Msg("existing loop due to context done")
				return
			case err, open := <-w.Errors:
				if !open {
					log.Info().Msg("exiting loop due to closure")
					return
				}
				if nil != err {
					log.Error().Err(err).Msg("received error")
					return
				}
			case event, open := <-w.Events:
				if !open {
					log.Info().Msg("exiting loop due to closure")
					return
				}

				if !event.Has(fsnotify.Create) && !event.Has(fsnotify.Write) {
					log.Debug().Msg("ignoring event with irrelevant op")
				}

				if filename := filepath.Base(event.Name); strings.HasSuffix(filename, fmt.Sprintf("%s.up", wgDeviceName)) {
					wgUpEvents <- ingest.WgUpEvent{}
				} else if strings.HasSuffix(filename, fmt.Sprintf("%s.down", wgDeviceName)) {
					wgDownEvents <- ingest.WgDownEvent{ChangedAt: time.Now(), FileName: event.Name}
				}
			}
		}
	}()

	ticker := make(chan ingest.None)
	go func() {
		timeTicker := time.NewTicker(5 * time.Second)
		for range timeTicker.C {
			ticker <- ingest.None{}
		}
	}()

	wp := wgPeers{wg}
	store := storeMongo{collection}
	engine := ingest.NewEngine(&wp, &store, log.With().Str("app", "engine").Logger())
	if err := engine.Run(runCtx, ticker, wgUpEvents, wgDownEvents); nil != err {
		if err := runCtx.Err(); nil != err {
			if errors.Is(err, context.Canceled) {
				if errors.Is(context.Cause(runCtx), stopSignalErr) {
					log.Info().Msg("root context was canceled due to receiving an interrupt signal")
					return
				}

				log.Info().Err(err).Msg("root context was canceled due to unexpected cause")
				return
			}

			log.Error().Err(err).Msg("root context was canceled with unexpected error")
			return
		}

		log.Error().Msg("root context was canceled unexpectedly with no errors")
		return
	}
}

type storeMongo struct {
	collection *mongo.Collection
}

func (m *storeMongo) LoadUsage(ctx context.Context) (map[string]ingest.PeerUsage, error) {
	cursor, err := m.collection.Aggregate(ctx, bson.A{
		bson.M{"$project": bson.M{"lastUsage": bson.M{"$last": "$usage"}, "_id": 0, "publicKey": 1}},
	})
	if nil != err {
		return nil, fmt.Errorf("failed to query last usage data: %v", err)
	}

	var results []struct {
		publicKey string `bson:"publicKey"`
		lastUsage struct {
			upload   uint `bson:"upload"`
			download uint `bson:"download"`
		} `bson:"lastUsage"`
	}
	if err := cursor.All(ctx, &results); nil != err {
		return nil, fmt.Errorf("failed to read all documents: %v", err)
	}

	out := make(map[string]ingest.PeerUsage, len(results))
	for _, v := range results {
		out[v.publicKey] = ingest.PeerUsage{
			Upload:    v.lastUsage.upload,
			Download:  v.lastUsage.download,
			PublicKey: v.publicKey,
		}
	}

	return out, nil
}

func (m *storeMongo) IngestUsage(ctx context.Context, peersUsage []ingest.PeerUsage, gatheredAt time.Time) error {
	models := funcutils.Map(peersUsage, func(p ingest.PeerUsage) mongo.WriteModel {
		return mongo.NewUpdateOneModel().
			SetFilter(bson.M{"publicKey": p.PublicKey}).
			SetUpdate(bson.M{"$push": bson.M{"usage": bson.M{"upload": p.Upload, "download": p.Download, "at": gatheredAt.UnixMilli()}}}).
			SetUpsert(true)
	})
	opts := options.BulkWrite().SetOrdered(false).SetBypassDocumentValidation(true)
	if _, err := m.collection.BulkWrite(ctx, models, opts); nil != err {
		return fmt.Errorf("failed to upsert peer models: %v", err)
	}

	return nil
}

type wgPeers struct {
	ctrl *wgctrl.Client
}

func (wg *wgPeers) Usage(ctx context.Context) ([]ingest.PeerUsage, time.Time, error) {
	dev, err := wg.ctrl.Device(wgDeviceName)
	gatheredAt := time.Now()
	if nil != err {
		return nil, gatheredAt, err
	}

	out := funcutils.Map(dev.Peers, func(p wgtypes.Peer) ingest.PeerUsage {
		return ingest.PeerUsage{
			Upload:    uint(p.ReceiveBytes),
			Download:  uint(p.TransmitBytes),
			PublicKey: p.PublicKey.String(),
		}
	})

	return out, gatheredAt, nil
}
