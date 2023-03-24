package ingest

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/xeptore/wireuse/pkg/dump"
	"github.com/xeptore/wireuse/pkg/funcutils"
)

type WGFileEventKind string

type WGFileEvent struct {
	Kind WGFileEventKind
	Name string
}

const (
	WGFileEventKindUp   WGFileEventKind = "up"
	WGFileEventKindDown WGFileEventKind = "down"
)

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

func (e *Engine) Run(ctx context.Context, tick <-chan struct{}, wgFileEvents <-chan WGFileEvent) error {
	var previousPeersUsage map[string]PeerUsage
	var mustLoadUsage bool
loop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-wgFileEvents:
			switch event.Kind {
			case WGFileEventKindUp:
				var err error
				previousPeersUsage, err = e.store.LoadUsage(ctx)
				if nil != err {
					e.logger.Error().Err(err).Msg("failed to load before restart peers usage data")
					mustLoadUsage = true
					continue loop
				}
				mustLoadUsage = false
			case WGFileEventKindDown:
				file, err := os.Open(event.Name)
				if nil != err {
					e.logger.Error().Err(err).Msg("failed to open wireguard down file for read")
					continue loop
				}

				peers, err := dump.Parse(file)
				if nil != err {
					e.logger.Error().Err(err).Msg("failed to parse wireguard down dump file content")
					continue loop
				}

				gatheredAt := func() time.Time {
					stat, err := os.Lstat(event.Name)
					if nil != err {
						e.logger.Error().Err(err).Msg("failed to get wireguard down dump file stat; using current time")
						return time.Now()
					}

					return stat.ModTime()
				}()

				peersUsage := funcutils.Map(peers, func(p dump.Peer) PeerUsage {
					return PeerUsage{
						Upload:    uint(p.ReceiveBytes),
						Download:  uint(p.TransmitBytes),
						PublicKey: p.PublicKey.String(),
					}
				})

				if len(peers) > 0 {
					if err := e.store.IngestUsage(ctx, peersUsage, gatheredAt); nil != err {
						e.logger.Error().Err(err).Msg("failed to ingest peers usage data")
						continue loop
					}
				}
			default:
				return fmt.Errorf("unknown wg file event kind received: %s", event.Kind)
			}
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
