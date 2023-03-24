package ingest

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/xeptore/wireuse/pkg/dump"
	"github.com/xeptore/wireuse/pkg/funcutils"
)

type WgDownEvent struct {
	ChangedAt time.Time
	FileName  string
}

type PeerUsage struct {
	Upload    uint
	Download  uint
	PublicKey string
}

type Store interface {
	LoadUsage(ctx context.Context) (map[string]PeerUsage, error)
	IngestUsage(ctx context.Context, peersUsage []PeerUsage, gatheredAt time.Time) error
}

type WgPeers interface {
	Usage(ctx context.Context) (peersUsage []PeerUsage, gatheredAt time.Time, err error)
}

type Engine struct {
	wgPeers WgPeers
	store   Store
	logger  zerolog.Logger
}

func NewEngine(wgPeers WgPeers, store Store, logger zerolog.Logger) Engine {
	return Engine{
		wgPeers: wgPeers,
		store:   store,
		logger:  logger,
	}
}

func (e *Engine) Run(ctx context.Context, tick <-chan struct{}, wgUpEvents <-chan struct{}, wgDownEvents <-chan WgDownEvent) error {
	var previousPeersUsage map[string]PeerUsage
	var mustLoadUsage bool
loop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-wgDownEvents:
			file, err := os.Open(event.FileName)
			if nil != err {
				e.logger.Error().Err(err).Msg("failed to open wireguard down file for read")
				continue loop
			}

			peers, err := dump.Parse(file)
			if nil != err {
				e.logger.Error().Err(err).Msg("failed to parse wireguard down dump file content")
				continue loop
			}

			peersUsage := funcutils.Map(peers, func(p dump.Peer) PeerUsage {
				return PeerUsage{
					Upload:    uint(p.ReceiveBytes),
					Download:  uint(p.TransmitBytes),
					PublicKey: p.PublicKey.String(),
				}
			})

			if len(peers) > 0 {
				if err := e.store.IngestUsage(ctx, peersUsage, event.ChangedAt); nil != err {
					e.logger.Error().Err(err).Msg("failed to ingest peers usage data")
					continue loop
				}
			}
		case <-wgUpEvents:
			var err error
			previousPeersUsage, err = e.store.LoadUsage(ctx)
			if nil != err {
				e.logger.Error().Err(err).Msg("failed to load before restart peers usage data")
				mustLoadUsage = true
				continue loop
			}
			mustLoadUsage = false
		case <-tick:
			peersUsage, gatheredAt, err := e.wgPeers.Usage(ctx)
			if nil != err {
				if errors.Is(err, os.ErrNotExist) {
					e.logger.Warn().Msg("could not find specified wireguard device")
				} else {
					e.logger.Error().Err(err).Msg("failed to get specified wireguard usage info")
				}
				continue loop
			}

			if mustLoadUsage {
				previousPeersUsage, err = e.store.LoadUsage(ctx)
				if nil != err {
					e.logger.Error().Err(err).Msg("failed to load before restart peers usage data")
					mustLoadUsage = true
					continue loop
				}
				mustLoadUsage = false
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
					continue loop
				}
			}
		}
	}
}
