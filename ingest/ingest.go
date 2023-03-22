package ingest

import (
	"context"
	"errors"
	"fmt"
	"os"
)

type PeerUsage struct {
	Upload    uint
	Download  uint
	PublicKey string
}

type Store interface {
	LoadBeforeRestartUsage(ctx context.Context) (map[string]PeerUsage, error)
	IngestUsage(ctx context.Context, peersUsage []PeerUsage) error
}

type WgPeers interface {
	Usage(ctx context.Context) ([]PeerUsage, error)
}

type RestartMarkFileReadRemover interface {
	Read(filename string) ([1]byte, error)
	Remove(filename string) error
}

type Engine struct {
	restartMarkFile RestartMarkFileReadRemover
	wgPeers         WgPeers
	store           Store
}

func NewEngine(
	restartMarkFile RestartMarkFileReadRemover,
	wgPeers WgPeers,
	store Store,
) Engine {
	return Engine{
		restartMarkFile: restartMarkFile,
		wgPeers:         wgPeers,
		store:           store,
	}
}

func (e *Engine) Run(ctx context.Context, tick <-chan struct{}, restartMarkFileName string) error {
	var previousPeersUsage map[string]PeerUsage
	for range tick {
		peersUsage, err := e.wgPeers.Usage(ctx)
		if nil != err {
			// TODO: log error
			continue
		}

		mustDeleteRestartMarkFile := false
		content, err := e.restartMarkFile.Read(restartMarkFileName)
		if nil != err && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to read restart-mark file: %w", err)
		} else if content == [1]byte{1} {
			previousPeersUsage, err = e.store.LoadBeforeRestartUsage(ctx)
			if nil != err {
				return fmt.Errorf("failed to load last before restart usage records: %v", err)
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
			if err := e.store.IngestUsage(ctx, peersUsage); nil != err {
				// TODO: log error
				continue // not gonna remove the restart-mark file as it might succeed in the next tick.
			}
		}

		// ignore the non-existing restart-mark file error in the removal operation
		// as it's either the case when it doesn't exist at all, or it's been removed in the previous tick.
		if mustDeleteRestartMarkFile {
			if err := e.restartMarkFile.Remove(restartMarkFileName); nil != err && !errors.Is(err, os.ErrNotExist) {
				// log.Error().Err(err).Msg("failed to remove restart-mark file")
			}
		}
	}

	return nil
}
