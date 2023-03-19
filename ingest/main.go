package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/xeptore/wireuse/pkg/env"
	"github.com/xeptore/wireuse/pkg/funcutils"
)

var (
	restartMarkFileName string
	wgDeviceName        string
)

type ingestFunc = func(ctx context.Context, peers []wgtypes.Peer) error

func main() {
	ctx := context.Background()

	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

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

	wg, err := wgctrl.New()
	if nil != err {
		log.Fatal().Err(err).Msg("failed to initialize wg control client")
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	ctx, done := context.WithCancelCause(ctx)
	go func() {
		<-signals
		done(stopSignalErr)
	}()

	ticker := time.NewTicker(5 * time.Second)

	var ingest ingestFunc

	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); nil != err {
				if errors.Is(err, context.Canceled) {
					if errors.Is(context.Cause(ctx), stopSignalErr) {
						log.Info().Msg("root context was canceled due to receiving a interrupt signal")
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
		default:
			for range ticker.C {
				dev, err := wg.Device(wgDeviceName)
				if nil != err {
					log.Error().Err(err).Msg("failed to get wg device info")
					return
				}

				content, err := os.ReadFile(restartMarkFileName)
				if nil != err {
					// it's ok that the restart-mark file doesn't exist as it means the wg server hasn't been restarted since the previous tick.
					if !errors.Is(err, os.ErrNotExist) {
						log.Error().Err(err).Msg("failed to read restart-mark file")
						return
					}
					ingest = ingestData(collection)
				} else if bytes.Equal(content, []byte{1}) {
					ingest = ingestRestartedServerData(collection)
				}

				if nil := ingest(ctx, dev.Peers); nil != err {
					log.Error().Err(err).Msg("failed to ingest data")
					continue // not gonna remove the restart-mark file as it might succeed in the next tick.
				}

				// ignore the non-existing restart-mark file error in the removal operation
				// as it's either the case when it doesn't exist at all, or it's been removed in the previous tick.
				if err := os.Remove(restartMarkFileName); !errors.Is(err, os.ErrNotExist) {
					log.Error().Err(err).Msg("failed to remove restart-mark file")
				}
			}
		}
	}
}

var (
	stopSignalErr = errors.New("stop signal received")
)

func ingestRestartedServerData(collection *mongo.Collection) ingestFunc {
	return func(ctx context.Context, peers []wgtypes.Peer) error {
		// get last usage of every peer
		// if number of retrieved peer list items > requested ones, return error
		// in a loop:
		//     if peer exists:
		//         current usage += previous usage
		// append usage record to respective peers
		return nil
	}
}

func ingestData(collection *mongo.Collection) ingestFunc {
	return func(ctx context.Context, peers []wgtypes.Peer) error {
		models := funcutils.MapSlice(peers, func(p wgtypes.Peer) mongo.WriteModel {
			return mongo.NewUpdateOneModel().
				SetFilter(bson.M{"publicKey": p.PublicKey}).
				SetUpdate(bson.M{"$push": bson.M{"uploadBytes": p.ReceiveBytes, "downloadBytes": p.TransmitBytes}}). // TODO: make sure upload and download byte fields are mapped correctly
				SetUpsert(true)
		})
		opts := options.BulkWrite().SetOrdered(false).SetBypassDocumentValidation(true)
		if _, err := collection.BulkWrite(ctx, models, opts); nil != err {
			return fmt.Errorf("failed to upsert peer models: %v", err)
		}

		return nil
	}
}
