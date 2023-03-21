package ingest_test

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/xeptore/wireuse/ingest"
	"github.com/xeptore/wireuse/ingest/mocks"
)

func chain(calls ...*gomock.Call) {
	if len(calls) == 0 {
		panic("at least one call is required")
	}
	if len(calls) == 1 {
		return
	}

	for i := len(calls) - 1; i > 0; i-- {
		calls[i].After(calls[i-1])
	}
}

func TestEngineSingleStaticPeer(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	store := mocks.NewMockStore(ctrl)
	store.EXPECT().LoadBeforeRestartUsage(ctx).Times(0)
	chain(
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(0), Download: uint(0), PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(10), Download: uint(30), PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(20), Download: uint(60), PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(30), Download: uint(90), PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(40), Download: uint(120), PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(50), Download: uint(150), PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(60), Download: uint(180), PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(70), Download: uint(210), PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(80), Download: uint(240), PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(90), Download: uint(270), PublicKey: "xyz"}}).Return(nil).Times(1),
	)

	readRestartMarkFile := mocks.NewMockRestartMarkFileReaderRemover(ctrl)
	readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(10)
	readRestartMarkFile.EXPECT().Remove("TODO").Return(nil).Times(10)

	readWGPeersUsage := mocks.NewMockWgPeers(ctrl)
	chain(
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 0, Download: 0, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 10, Download: 30, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 20, Download: 60, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 30, Download: 90, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 40, Download: 120, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 50, Download: 150, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 60, Download: 180, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 70, Download: 210, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 80, Download: 240, PublicKey: "xyz"}}, nil).Times(1),
		readWGPeersUsage.EXPECT().Usage(ctx).Return([]ingest.PeerUsage{{Upload: 90, Download: 270, PublicKey: "xyz"}}, nil).Times(1),
	)

	e := ingest.NewEngine(readRestartMarkFile, readWGPeersUsage, store)

	ticker := make(chan struct{}, 1)

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

func TestEngineSingleStaticPeerWithRestart(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)

	store := mocks.NewMockStore(ctrl)
	chain(
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 154, Download: 215, PublicKey: "xyz"}}, nil).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 5852, Download: 43146, PublicKey: "xyz"}}, nil).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 13842, Download: 142940, PublicKey: "xyz"}}, nil).Times(1),
		store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 15991, Download: 166067, PublicKey: "xyz"}}, nil).Times(1),
	)
	chain(
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
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 12354, Download: 119720, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 12358, Download: 125646, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 12934, Download: 126232, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 13212, Download: 133233, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 13269, Download: 134991, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 13380, Download: 139524, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 13842, Download: 142940, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15212, Download: 156236, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15332, Download: 156999, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15352, Download: 157928, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15755, Download: 159005, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15802, Download: 162449, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15870, Download: 163562, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15991, Download: 166067, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 16504, Download: 166768, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 16522, Download: 166769, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 16531, Download: 166775, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 16533, Download: 166788, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 16535, Download: 166797, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 16542, Download: 166805, PublicKey: "xyz"}}).Return(nil).Times(1),
		store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 16545, Download: 166809, PublicKey: "xyz"}}).Return(nil).Times(1),
	)

	readRestartMarkFile := mocks.NewMockRestartMarkFileReaderRemover(ctrl)
	readRestartMarkFile.EXPECT().Remove("TODO").Return(nil).AnyTimes()
	chain(
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{1}, nil).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
		readRestartMarkFile.EXPECT().Read("TODO").Return([1]byte{0}, os.ErrNotExist).Times(1),
	)

	readWGPeersUsage := mocks.NewMockWgPeers(ctrl)
	chain(
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
