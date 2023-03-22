package ingest_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/xeptore/wireuse/ingest"
	"github.com/xeptore/wireuse/ingest/mocks"
)

func TestEngineSingleStaticPeer(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	store := mocks.NewMockStore(ctrl)
	store.EXPECT().LoadBeforeRestartUsage(ctx).Times(0)
	gomock.InOrder(
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 10, Download: 30, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 20, Download: 60, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 30, Download: 90, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 40, Download: 120, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 50, Download: 150, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 60, Download: 180, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 70, Download: 210, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 80, Download: 240, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 90, Download: 270, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 100, Download: 300, PublicKey: "xyz"}}).Return(nil).Times(1),
	)

	readRestartMarkFile := mocks.NewMockRestartMarkFileReadRemover(ctrl)
	readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(10)
	readRestartMarkFile.EXPECT().Remove("TODO").Return(nil).Times(0)

	readWGPeersUsage := mocks.NewMockWgPeers(ctrl)
	gomock.InOrder(
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 10, Download: 30, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 20, Download: 60, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 30, Download: 90, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 40, Download: 120, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 50, Download: 150, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 60, Download: 180, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 70, Download: 210, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 80, Download: 240, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 90, Download: 270, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 100, Download: 300, PublicKey: "xyz"}}, nil).Times(1),
	)

	e := ingest.NewEngine(readRestartMarkFile, readWGPeersUsage, store)

	ticker := make(chan struct{})

	var runErr error
	wait := make(chan struct{})
	go func() {
		defer func() {
			wait <- struct{}{}
		}()
		runErr = e.Run(ctx, ticker, "TODO")
	}()

	for i := 0; i < 10; i++ {
		select {
		case ticker <- struct{}{}:
		case <-wait:
			t.Fatal("unexpected engine run termination")
		}
	}

	close(ticker)
	<-wait
	require.Nil(t, runErr)
}

func TestEngineSingleStaticPeerWithWgReadPeersUsageFailure(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	store := mocks.NewMockStore(ctrl)
	store.EXPECT().LoadBeforeRestartUsage(ctx).Times(0)
	gomock.InOrder(
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 10, Download: 30, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 20, Download: 60, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 40, Download: 120, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 50, Download: 150, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 60, Download: 180, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 70, Download: 210, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 90, Download: 270, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 100, Download: 300, PublicKey: "xyz"}}).Return(nil).Times(1),
	)

	readRestartMarkFile := mocks.NewMockRestartMarkFileReadRemover(ctrl)
	readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(8)
	readRestartMarkFile.EXPECT().Remove("TODO").Return(nil).Times(0)

	readWGPeersUsage := mocks.NewMockWgPeers(ctrl)
	gomock.InOrder(
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 10, Download: 30, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 20, Download: 60, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return(nil, errors.New("unknown error")).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 40, Download: 120, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 50, Download: 150, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 60, Download: 180, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 70, Download: 210, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return(nil, errors.New("network error")).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 90, Download: 270, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 100, Download: 300, PublicKey: "xyz"}}, nil).Times(1),
	)

	e := ingest.NewEngine(readRestartMarkFile, readWGPeersUsage, store)

	ticker := make(chan struct{})

	var runErr error
	wait := make(chan struct{})
	go func() {
		defer func() {
			wait <- struct{}{}
		}()
		runErr = e.Run(ctx, ticker, "TODO")
	}()

	for i := 0; i < 10; i++ {
		select {
		case ticker <- struct{}{}:
		case <-wait:
			t.Fatal("unexpected engine run termination")
		}
	}

	close(ticker)
	<-wait
	require.Nil(t, runErr)
}

func TestEngineSingleStaticPeerWithIngestFailure(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	store := mocks.NewMockStore(ctrl)
	store.EXPECT().LoadBeforeRestartUsage(ctx).Times(0)
	gomock.InOrder(
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 10, Download: 30, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 20, Download: 60, PublicKey: "xyz"}}).Return(errors.New("unknown error")).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 30, Download: 90, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 40, Download: 120, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 50, Download: 150, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 60, Download: 180, PublicKey: "xyz"}}).Return(errors.New("network error")).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 70, Download: 210, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 80, Download: 240, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 90, Download: 270, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 100, Download: 300, PublicKey: "xyz"}}).Return(nil).Times(1),
	)

	readRestartMarkFile := mocks.NewMockRestartMarkFileReadRemover(ctrl)
	readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(10)
	readRestartMarkFile.EXPECT().Remove("TODO").Return(nil).Times(0)

	readWGPeersUsage := mocks.NewMockWgPeers(ctrl)
	gomock.InOrder(
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 10, Download: 30, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 20, Download: 60, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 30, Download: 90, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 40, Download: 120, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 50, Download: 150, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 60, Download: 180, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 70, Download: 210, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 80, Download: 240, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 90, Download: 270, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 100, Download: 300, PublicKey: "xyz"}}, nil).Times(1),
	)

	e := ingest.NewEngine(readRestartMarkFile, readWGPeersUsage, store)

	ticker := make(chan struct{})

	var runErr error
	wait := make(chan struct{})
	go func() {
		defer func() {
			wait <- struct{}{}
		}()
		runErr = e.Run(ctx, ticker, "TODO")
	}()

	for i := 0; i < 10; i++ {
		select {
		case ticker <- struct{}{}:
		case <-wait:
			t.Fatal("unexpected engine run termination")
		}
	}

	close(ticker)
	<-wait
	require.Nil(t, runErr)
}

func TestEngineSingleStaticPeerWithRestartMarkFileReadFailure(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	store := mocks.NewMockStore(ctrl)
	store.EXPECT().LoadBeforeRestartUsage(ctx).Times(0)
	gomock.InOrder(
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 10, Download: 30, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 20, Download: 60, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 30, Download: 90, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 40, Download: 120, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 50, Download: 150, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 60, Download: 180, PublicKey: "xyz"}}).Return(nil).Times(1),
	)

	readRestartMarkFile := mocks.NewMockRestartMarkFileReadRemover(ctrl)
	gomock.InOrder(
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(3),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{2}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(2),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrPermission).Times(1),
	)
	readRestartMarkFile.EXPECT().Remove("TODO").Return(nil).Times(0)

	readWGPeersUsage := mocks.NewMockWgPeers(ctrl)
	gomock.InOrder(
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 10, Download: 30, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 20, Download: 60, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 30, Download: 90, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 40, Download: 120, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 50, Download: 150, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 60, Download: 180, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 70, Download: 210, PublicKey: "xyz"}}, nil).Times(1),
	)

	e := ingest.NewEngine(readRestartMarkFile, readWGPeersUsage, store)

	ticker := make(chan struct{})

	var runErr error
	wait := make(chan struct{})
	go func() {
		defer func() {
			wait <- struct{}{}
		}()
		runErr = e.Run(ctx, ticker, "TODO")
	}()

	for i := 0; i < 7; i++ {
		select {
		case ticker <- struct{}{}:
		case <-wait:
			t.Fatal("unexpected engine run termination")
		}
	}

	close(ticker)
	<-wait
	require.NotNil(t, runErr)
	require.ErrorIs(t, runErr, os.ErrPermission)
}

func TestEngineSingleStaticPeerWithRestart(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)

	store := mocks.NewMockStore(ctrl)
	gomock.InOrder(
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 154, Download: 215, PublicKey: "xyz"}}, nil).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 5852, Download: 43146, PublicKey: "xyz"}}, nil).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 6406, Download: 43888, PublicKey: "xyz"}}, nil).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 8555, Download: 67015, PublicKey: "xyz"}}, nil).Times(1),
	)
	gomock.InOrder(
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 120, Download: 169, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 122, Download: 170, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 141, Download: 176, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 142, Download: 186, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 150, Download: 194, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 151, Download: 198, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 154, Download: 215, PublicKey: "xyz"}}).Return(nil).Times(1),

		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 4130, Download: 29607, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 4230, Download: 31776, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 4562, Download: 32532, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 4631, Download: 35744, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 4889, Download: 38013, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 5659, Download: 38580, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 5852, Download: 43146, PublicKey: "xyz"}}).Return(nil).Times(1),

		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6365, Download: 43847, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6383, Download: 43848, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6392, Download: 43854, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6394, Download: 43867, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6396, Download: 43876, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6403, Download: 43884, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6406, Download: 43888, PublicKey: "xyz"}}).Return(nil).Times(1),

		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 7776, Download: 57184, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 7896, Download: 57947, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 7916, Download: 58876, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 8319, Download: 59953, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 8366, Download: 63397, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 8434, Download: 64510, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 8555, Download: 67015, PublicKey: "xyz"}}).Return(nil).Times(1),

		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15057, Download: 143589, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15061, Download: 149515, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15637, Download: 150101, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15915, Download: 157102, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15972, Download: 158860, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 16083, Download: 163393, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 16545, Download: 166809, PublicKey: "xyz"}}).Return(nil).Times(1),
	)

	readRestartMarkFile := mocks.NewMockRestartMarkFileReadRemover(ctrl)
	readRestartMarkFile.EXPECT().Remove("TODO").Return(nil).Times(4)
	gomock.InOrder(
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(7),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(6),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(6),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(6),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(6),
	)

	readWGPeersUsage := mocks.NewMockWgPeers(ctrl)
	gomock.InOrder(
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 120, Download: 169, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 122, Download: 170, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 141, Download: 176, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 142, Download: 186, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 150, Download: 194, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 151, Download: 198, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 154, Download: 215, PublicKey: "xyz"}}, nil).Times(1),

		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 3976, Download: 29392, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 4076, Download: 31561, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 4408, Download: 32317, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 4477, Download: 35529, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 4735, Download: 37798, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 5505, Download: 38365, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 5698, Download: 42931, PublicKey: "xyz"}}, nil).Times(1),

		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 513, Download: 701, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 531, Download: 702, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 540, Download: 708, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 542, Download: 721, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 544, Download: 730, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 551, Download: 738, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 554, Download: 742, PublicKey: "xyz"}}, nil).Times(1),

		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 1370, Download: 13296, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 1490, Download: 14059, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 1510, Download: 14988, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 1913, Download: 16065, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 1960, Download: 19509, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 2028, Download: 20622, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 2149, Download: 23127, PublicKey: "xyz"}}, nil).Times(1),

		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 6502, Download: 76574, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 6506, Download: 82500, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 7082, Download: 83086, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 7360, Download: 90087, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 7417, Download: 91845, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 7528, Download: 96378, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 7990, Download: 99794, PublicKey: "xyz"}}, nil).Times(1),
	)

	e := ingest.NewEngine(readRestartMarkFile, readWGPeersUsage, store)

	ticker := make(chan struct{})

	var runErr error
	wait := make(chan struct{})
	go func() {
		defer func() {
			wait <- struct{}{}
		}()
		runErr = e.Run(ctx, ticker, "TODO")
	}()

	for i := 0; i < 35; i++ {
		select {
		case ticker <- struct{}{}:
		case <-wait:
			t.Fatal("unexpected engine run termination")
		}
	}

	close(ticker)
	<-wait
	require.Nil(t, runErr)
}

func TestEngineSingleStaticPeerWithRestartStartingWithRestartMarkFileExistence(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)

	store := mocks.NewMockStore(ctrl)
	gomock.InOrder(
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 1, Download: 2, PublicKey: "xyz"}}, nil).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 155, Download: 217, PublicKey: "xyz"}}, nil).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 5853, Download: 43148, PublicKey: "xyz"}}, nil).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 6407, Download: 43890, PublicKey: "xyz"}}, nil).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 8556, Download: 67017, PublicKey: "xyz"}}, nil).Times(1),
	)
	gomock.InOrder(
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 121, Download: 171, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 123, Download: 172, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 142, Download: 178, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 143, Download: 188, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 151, Download: 196, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 152, Download: 200, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 155, Download: 217, PublicKey: "xyz"}}).Return(nil).Times(1),

		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 4131, Download: 29609, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 4231, Download: 31778, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 4563, Download: 32534, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 4632, Download: 35746, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 4890, Download: 38015, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 5660, Download: 38582, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 5853, Download: 43148, PublicKey: "xyz"}}).Return(nil).Times(1),

		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6366, Download: 43849, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6384, Download: 43850, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6393, Download: 43856, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6395, Download: 43869, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6397, Download: 43878, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6404, Download: 43886, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 6407, Download: 43890, PublicKey: "xyz"}}).Return(nil).Times(1),

		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 7777, Download: 57186, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 7897, Download: 57949, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 7917, Download: 58878, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 8320, Download: 59955, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 8367, Download: 63399, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 8435, Download: 64512, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 8556, Download: 67017, PublicKey: "xyz"}}).Return(nil).Times(1),

		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15058, Download: 143591, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15062, Download: 149517, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15638, Download: 150103, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15916, Download: 157104, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15973, Download: 158862, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 16084, Download: 163395, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 16546, Download: 166811, PublicKey: "xyz"}}).Return(nil).Times(1),
	)

	readRestartMarkFile := mocks.NewMockRestartMarkFileReadRemover(ctrl)
	readRestartMarkFile.EXPECT().Remove("TODO").Return(nil).Times(5)
	gomock.InOrder(
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(6),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(6),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(6),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(6),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(6),
	)

	readWGPeersUsage := mocks.NewMockWgPeers(ctrl)
	gomock.InOrder(
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 120, Download: 169, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 122, Download: 170, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 141, Download: 176, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 142, Download: 186, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 150, Download: 194, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 151, Download: 198, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 154, Download: 215, PublicKey: "xyz"}}, nil).Times(1),

		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 3976, Download: 29392, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 4076, Download: 31561, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 4408, Download: 32317, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 4477, Download: 35529, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 4735, Download: 37798, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 5505, Download: 38365, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 5698, Download: 42931, PublicKey: "xyz"}}, nil).Times(1),

		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 513, Download: 701, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 531, Download: 702, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 540, Download: 708, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 542, Download: 721, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 544, Download: 730, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 551, Download: 738, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 554, Download: 742, PublicKey: "xyz"}}, nil).Times(1),

		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 1370, Download: 13296, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 1490, Download: 14059, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 1510, Download: 14988, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 1913, Download: 16065, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 1960, Download: 19509, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 2028, Download: 20622, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 2149, Download: 23127, PublicKey: "xyz"}}, nil).Times(1),

		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 6502, Download: 76574, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 6506, Download: 82500, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 7082, Download: 83086, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 7360, Download: 90087, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 7417, Download: 91845, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 7528, Download: 96378, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 7990, Download: 99794, PublicKey: "xyz"}}, nil).Times(1),
	)

	e := ingest.NewEngine(readRestartMarkFile, readWGPeersUsage, store)

	ticker := make(chan struct{})

	var runErr error
	wait := make(chan struct{})
	go func() {
		defer func() {
			wait <- struct{}{}
		}()
		runErr = e.Run(ctx, ticker, "TODO")
	}()

	for i := 0; i < 35; i++ {
		select {
		case ticker <- struct{}{}:
		case <-wait:
			t.Fatal("unexpected engine run termination")
		}
	}

	close(ticker)
	<-wait
	require.Nil(t, runErr)
}

func TestEngineMultipleDynamicPeers(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	store := mocks.NewMockStore(ctrl)

	gomock.InOrder(
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(
			map[string]ingest.PeerUsage{
				"abc": {Upload: 25, Download: 65, PublicKey: "abc"},
				"xyz": {Upload: 20, Download: 60, PublicKey: "xyz"},
				"123": {Upload: 17, Download: 53, PublicKey: "123"},
			},
			nil,
		).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(
			map[string]ingest.PeerUsage{
				"xyz": {Upload: 50, Download: 150, PublicKey: "xyz"},
				"852": {Upload: 186580512, Download: 995098551, PublicKey: "852"},
				"abc": {Upload: 128001692, Download: 186202004, PublicKey: "abc"},
				"123": {Upload: 57, Download: 153, PublicKey: "123"},
			},
			nil,
		).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(
			map[string]ingest.PeerUsage{
				"xyz": {Upload: 50, Download: 150, PublicKey: "xyz"},
				"852": {Upload: 186580512, Download: 995098551, PublicKey: "852"},
				"abc": {Upload: 128001692, Download: 186202004, PublicKey: "abc"},
				"123": {Upload: 107, Download: 263, PublicKey: "123"},
			},
			nil,
		).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(
			map[string]ingest.PeerUsage{
				"852": {Upload: 186580512, Download: 995098551, PublicKey: "852"},
				"qwe": {Upload: 40, Download: 85, PublicKey: "qwe"},
				"xyz": {Upload: 37 + 50, Download: 98 + 150, PublicKey: "xyz"},
				"123": {Upload: 45 + 107, Download: 120 + 263, PublicKey: "123"},
				"456": {Upload: 124728866, Download: 155917550, PublicKey: "456"},
				"abc": {Upload: 49 + 128001692, Download: 137 + 186202004, PublicKey: "abc"},
			},
			nil,
		).Times(1),
	)

	gomock.InOrder(
		store.EXPECT().IngestUsage(
			ctx,
			[]ingest.PeerUsage{
				{Upload: 10, Download: 30, PublicKey: "xyz"},
				{Upload: 15, Download: 35, PublicKey: "abc"},
			},
		).Return(nil).Times(1),
		store.EXPECT().IngestUsage(
			ctx,
			[]ingest.PeerUsage{
				{Upload: 20, Download: 60, PublicKey: "xyz"},
				{Upload: 17, Download: 53, PublicKey: "123"},
				{Upload: 25, Download: 65, PublicKey: "abc"},
			},
		).Return(nil).Times(1),
		store.EXPECT().IngestUsage(
			ctx,
			[]ingest.PeerUsage{
				{Upload: 6531511 + 25, Download: 55474931 + 65, PublicKey: "abc"},
				{Upload: 12 + 20, Download: 33 + 60, PublicKey: "xyz"},
			},
		).Return(nil).Times(1),
		store.EXPECT().IngestUsage(
			ctx,
			[]ingest.PeerUsage{
				{Upload: 50, Download: 150, PublicKey: "xyz"},
				{Upload: 186580512, Download: 995098551, PublicKey: "852"},
				{Upload: 128001692, Download: 186202004, PublicKey: "abc"},
				{Upload: 57, Download: 153, PublicKey: "123"},
			},
		).Return(nil).Times(1),
		store.EXPECT().IngestUsage(
			ctx,
			[]ingest.PeerUsage{
				{Upload: 50 + 57, Download: 110 + 153, PublicKey: "123"},
			},
		).Return(nil).Times(1),
		store.EXPECT().IngestUsage(
			ctx,
			[]ingest.PeerUsage{
				{Upload: 20, Download: 40, PublicKey: "qwe"},
				{Upload: 37 + 50, Download: 98 + 150, PublicKey: "xyz"},
				{Upload: 45 + 57 + 50, Download: 120 + 110 + 153, PublicKey: "123"},
				{Upload: 49 + 128001692, Download: 137 + 186202004, PublicKey: "abc"},
			},
		).Return(nil).Times(1),
		store.EXPECT().IngestUsage(
			ctx,
			[]ingest.PeerUsage{
				{Upload: 30, Download: 75, PublicKey: "qwe"},
				{Upload: 53 + 107, Download: 132 + 263, PublicKey: "123"},
			},
		).Return(nil).Times(1),
		store.EXPECT().IngestUsage(
			ctx,
			[]ingest.PeerUsage{
				{Upload: 40, Download: 85, PublicKey: "qwe"},
				{Upload: 37 + 50, Download: 98 + 150, PublicKey: "xyz"},
				{Upload: 45 + 107, Download: 120 + 263, PublicKey: "123"},
				{Upload: 124728866, Download: 155917550, PublicKey: "456"},
				{Upload: 49 + 128001692, Download: 137 + 186202004, PublicKey: "abc"},
			},
		).Return(nil).Times(1),
		store.EXPECT().IngestUsage(
			ctx,
			[]ingest.PeerUsage{
				{Upload: 203658612 + 124728866, Download: 220351578 + 155917550, PublicKey: "456"},
				{Upload: 50 + 40, Download: 95 + 85, PublicKey: "qwe"},
				{Upload: 76 + 37 + 50, Download: 168 + 98 + 150, PublicKey: "xyz"},
			},
		).Return(nil).Times(1),
	)

	readRestartMarkFile := mocks.NewMockRestartMarkFileReadRemover(ctrl)
	readRestartMarkFile.EXPECT().Remove("TODO").Return(nil).Times(4)
	gomock.InOrder(
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(2),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(2),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(3),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
	)

	readWGPeersUsage := mocks.NewMockWgPeers(ctrl)
	gomock.InOrder(
		readWGPeersUsage.EXPECT().Usage(ctx).Return(
			[]ingest.PeerUsage{
				{Upload: 10, Download: 30, PublicKey: "xyz"},
				{Upload: 15, Download: 35, PublicKey: "abc"},
			},
			nil,
		).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return(
			[]ingest.PeerUsage{
				{Upload: 20, Download: 60, PublicKey: "xyz"},
				{Upload: 17, Download: 53, PublicKey: "123"},
				{Upload: 25, Download: 65, PublicKey: "abc"},
			},
			nil,
		).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return(
			[]ingest.PeerUsage{
				{Upload: 6531511, Download: 55474931, PublicKey: "abc"},
				{Upload: 12, Download: 33, PublicKey: "xyz"},
			},
			nil,
		).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return(
			[]ingest.PeerUsage{
				{Upload: 30, Download: 90, PublicKey: "xyz"},
				{Upload: 186580512, Download: 995098551, PublicKey: "852"},
				{Upload: 128001667, Download: 186201939, PublicKey: "abc"},
				{Upload: 40, Download: 100, PublicKey: "123"},
			},
			nil,
		).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return(
			[]ingest.PeerUsage{
				{Upload: 50, Download: 110, PublicKey: "123"},
			},
			nil,
		).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return(
			[]ingest.PeerUsage{},
			nil,
		).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return(
			[]ingest.PeerUsage{
				{Upload: 20, Download: 40, PublicKey: "qwe"},
				{Upload: 37, Download: 98, PublicKey: "xyz"},
				{Upload: 45, Download: 120, PublicKey: "123"},
				{Upload: 49, Download: 137, PublicKey: "abc"},
			},
			nil,
		).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return(
			[]ingest.PeerUsage{
				{Upload: 30, Download: 75, PublicKey: "qwe"},
				{Upload: 53, Download: 132, PublicKey: "123"},
			},
			nil,
		).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return(
			[]ingest.PeerUsage{
				{Upload: 40, Download: 85, PublicKey: "qwe"},
				{Upload: 37, Download: 98, PublicKey: "xyz"},
				{Upload: 45, Download: 120, PublicKey: "123"},
				{Upload: 124728866, Download: 155917550, PublicKey: "456"},
				{Upload: 49, Download: 137, PublicKey: "abc"},
			},
			nil,
		).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return(
			[]ingest.PeerUsage{
				{Upload: 203658612, Download: 220351578, PublicKey: "456"},
				{Upload: 50, Download: 95, PublicKey: "qwe"},
				{Upload: 76, Download: 168, PublicKey: "xyz"},
			},
			nil,
		).Times(1),
	)

	e := ingest.NewEngine(readRestartMarkFile, readWGPeersUsage, store)

	ticker := make(chan struct{})

	var runErr error
	wait := make(chan struct{})
	go func() {
		defer func() {
			wait <- struct{}{}
		}()
		runErr = e.Run(ctx, ticker, "TODO")
	}()

	for i := 0; i < 10; i++ {
		select {
		case ticker <- struct{}{}:
		case <-wait:
			t.Fatal("unexpected engine run termination")
		}
	}

	close(ticker)
	<-wait
	require.Nil(t, runErr)
}
