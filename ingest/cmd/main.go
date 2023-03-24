package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
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
	restartMarkFileName   string
	wgDeviceName          string
	wgPreDownDumpFileName string
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

	flag.StringVar(&restartMarkFileName, "r", "", "restart-mark file name")
	flag.StringVar(&wgDeviceName, "i", "", "wireguard interface")
	flag.StringVar(&wgPreDownDumpFileName, "d", "", "wireguard pre-down dumped output file")

	flag.Parse()
	if nonFlagArgs := flag.Args(); len(nonFlagArgs) > 0 {
		log.Fatal().Msgf("expected no additional flags, got: %s", strings.Join(nonFlagArgs, ","))
	}
	if restartMarkFileName == "" {
		log.Fatal().Msg("restart-mark file name option is required and cannot be empty")
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
		if err := client.Disconnect(ctx); err != nil {
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
	ctx, cancel := context.WithCancelCause(ctx)
	stopSignalErr := errors.New("stop signal received")
	go func() {
		<-signals
		cancel(stopSignalErr)
	}()

	go func() {
		if wgPreDownDumpFileName == "" {
			return
		}

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize wireguard pre-down dump file watcher")
		}
		defer watcher.Close()
		watcher.Add("")
	}()

	timeTicker := time.NewTicker(5 * time.Second)
	engineTicker := make(chan struct{})
	go func() {
		for range timeTicker.C {
			engineTicker <- struct{}{}
		}
	}()

	rmf := restartMarkFileReadRemover{}
	wp := wgPeers{wg}
	store := storeMongo{collection}
	engine := ingest.NewEngine(&rmf, &wp, &store, log)
	if err := engine.Run(ctx, engineTicker, restartMarkFileName); nil != err {
		if err := ctx.Err(); nil != err {
			if errors.Is(err, context.Canceled) {
				if errors.Is(context.Cause(ctx), stopSignalErr) {
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

func (m *storeMongo) LoadBeforeRestartUsage(ctx context.Context) (map[string]ingest.PeerUsage, error) {
	cursor, err := m.collection.Aggregate(ctx, bson.A{
		bson.M{"$project": bson.M{"lastUsage": bson.M{"$last": "$usage"}, "_id": 0, "publicKey": 1}},
	})
	if nil != err {
		return nil, fmt.Errorf("failed to query before restart last usage data: %v", err)
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
			Upload:    uint(p.TransmitBytes),
			Download:  uint(p.ReceiveBytes),
			PublicKey: p.PublicKey.String(),
		}
	})

	return out, gatheredAt, nil
}

type restartMarkFileReadRemover struct{}

func (*restartMarkFileReadRemover) Read(filename string) ([1]byte, error) {
	file, err := os.Open(restartMarkFileName)
	if nil != err {
		return [1]byte{0}, fmt.Errorf("failed to open restart-mark file: %w", err)
	}

	buf := make([]byte, 1)
	n, err := file.Read(buf)
	if nil != err {
		return [1]byte{0}, fmt.Errorf("failed to read first byte of restart-mark file: %w", err)
	}
	if n > 1 {
		return [1]byte{0}, fmt.Errorf("expected to read at most 1 byte from file read: %d", n)
	}
	if n == 0 {
		return [1]byte{0}, nil
	}

	return [1]byte{buf[0]}, nil
}

func (*restartMarkFileReadRemover) Remove(filename string) error {
	return os.Remove(filename)
}
