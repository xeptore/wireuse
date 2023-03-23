package ingest

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type PeerUsage struct {
	Upload    uint
	Download  uint
	PublicKey string
}

type Store interface {
	LoadBeforeRestartUsage(ctx context.Context) (map[string]PeerUsage, error)
	IngestUsage(ctx context.Context, peersUsage []PeerUsage, gatheredAt time.Time) error
}

type WgPeers interface {
	Usage(ctx context.Context) (peersUsage []PeerUsage, gatheredAt time.Time, err error)
}

type RestartMarkFileReadRemover interface {
	Read(filename string) ([1]byte, error)
	Remove(filename string) error
}

type Engine struct {
	restartMarkFile RestartMarkFileReadRemover
	wgPeers         WgPeers
	store           Store
	logger          zerolog.Logger
}

func NewEngine(
	restartMarkFile RestartMarkFileReadRemover,
	wgPeers WgPeers,
	store Store,
	logger zerolog.Logger,
) Engine {
	return Engine{
		restartMarkFile: restartMarkFile,
		wgPeers:         wgPeers,
		store:           store,
		logger:          logger,
	}
}

func (e *Engine) Run(ctx context.Context, tick <-chan struct{}, restartMarkFileName string) error {
	var previousPeersUsage map[string]PeerUsage
	for range tick {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			peersUsage, gatheredAt, err := e.wgPeers.Usage(ctx)
			if nil != err {
				if !errors.Is(err, os.ErrNotExist) {
					e.logger.Error().Err(err).Msg("failed to get wireguard peers usage data")
				}
				continue
			}

			mustDeleteRestartMarkFile := false
			content, err := e.restartMarkFile.Read(restartMarkFileName)
			if nil != err && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to read restart-mark file: %w", err)
			} else if content == [1]byte{1} {
				previousPeersUsage, err = e.store.LoadBeforeRestartUsage(ctx)
				if nil != err {
					e.logger.Error().Err(err).Msg("failed to load before restart peers usage data")
					continue
				}
				mustDeleteRestartMarkFile = true
			}

			if nil != previousPeersUsage {
				for i := 0; i < len(peersUsage); i++ {
					if prevUsage, exists := previousPeersUsage[peersUsage[i].PublicKey]; exists {
						peersUsage[i].Download += prevUsage.Download
						peersUsage[i].Upload += prevUsage.Upload
					}
				}
			}

			if len(peersUsage) > 0 {
				if err := e.store.IngestUsage(ctx, peersUsage, gatheredAt); nil != err {
					e.logger.Error().Err(err).Msg("failed to ingest peers usage data")
					continue
				}
			}

			if mustDeleteRestartMarkFile {
				if err := e.restartMarkFile.Remove(restartMarkFileName); nil != err && !errors.Is(err, os.ErrNotExist) {
					e.logger.Error().Err(err).Msg("failed to remove restart-mark file")
				}
			}
		}
	}

	return nil
}
