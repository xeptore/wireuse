package ingest_test

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/xeptore/wireuse/ingest"
	"github.com/xeptore/wireuse/ingest/mocks"
)

func TestEngineSingleStaticPeer(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	store := mocks.NewMockStore(ctrl)
	store.EXPECT().LoadBeforeRestartUsage(ctx).Times(0)

	store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(90), Download: uint(270), PublicKey: "xyz"}}).Return(nil).Times(1).
		After(store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(80), Download: uint(240), PublicKey: "xyz"}}).Return(nil).Times(1)).
		After(store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(70), Download: uint(210), PublicKey: "xyz"}}).Return(nil).Times(1)).
		After(store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(60), Download: uint(180), PublicKey: "xyz"}}).Return(nil).Times(1)).
		After(store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(50), Download: uint(150), PublicKey: "xyz"}}).Return(nil).Times(1)).
		After(store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(40), Download: uint(120), PublicKey: "xyz"}}).Return(nil).Times(1)).
		After(store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(30), Download: uint(90), PublicKey: "xyz"}}).Return(nil).Times(1)).
		After(store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(20), Download: uint(60), PublicKey: "xyz"}}).Return(nil).Times(1)).
		After(store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(10), Download: uint(30), PublicKey: "xyz"}}).Return(nil).Times(1)).
		After(store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: uint(0), Download: uint(0), PublicKey: "xyz"}}).Return(nil).Times(1))

	peersUsage := make(chan []ingest.PeerUsage)
	readWGPeersUsage := func(ctx context.Context) ([]ingest.PeerUsage, error) {
		return <-peersUsage, nil
	}
	readRestartMarkFile := func(filename string) ([1]byte, error) {
		return [1]byte{0}, os.ErrNotExist
	}
	removeRestartMarkFile := func(filename string) error {
		return nil
	}
	e := ingest.NewEngine(readWGPeersUsage, readRestartMarkFile, removeRestartMarkFile, store)
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
			peersUsage <- []ingest.PeerUsage{{Upload: uint(i * 10), Download: uint(i * 30), PublicKey: "xyz"}}
		case <-wait:
			t.Fatal("unexpected engine run termination")
			wait <- struct{}{}
		}
	}
	close(ticker)
	<-wait
	require.Nil(t, runErr)
}

type EngineSingleStaticPeerWithRestartSuite struct {
	suite.Suite
}

func (suite *EngineSingleStaticPeerWithRestartSuite) TestUsualRestart() {
	t := suite.T()
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	store := mocks.NewMockStore(ctrl)

	store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 10, Download: 20, PublicKey: "xyz"}}, nil).Times(1).
		After(store.EXPECT().LoadBeforeRestartUsage(ctx).Times(1))

	store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 15, Download: 30, PublicKey: "xyz"}}).Return(nil).Times(1).
		After(store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 10, Download: 20, PublicKey: "xyz"}}).Return(nil).Times(1))

	peersUsage := make(chan []ingest.PeerUsage)
	readWGPeersUsage := func(ctx context.Context) ([]ingest.PeerUsage, error) {
		return <-peersUsage, nil
	}
	var mux sync.Mutex
	readRestartMarkFileFn := func() ([1]byte, error) {
		return [1]byte{0}, os.ErrNotExist
	}
	readRestartMarkFile := func(filename string) ([1]byte, error) {
		mux.Lock()
		defer mux.Unlock()
		return readRestartMarkFileFn()
	}
	removeRestartMarkFile := func(filename string) error {
		return nil
	}
	e := ingest.NewEngine(readWGPeersUsage, readRestartMarkFile, removeRestartMarkFile, store)
	ticker := make(chan struct{})

	var runErr error
	wait := make(chan struct{})
	go func() {
		defer func() {
			wait <- struct{}{}
		}()
		runErr = e.Run(ctx, ticker, "TODO")
	}()
	for i := 0; i < 1; i++ {
		select {
		case ticker <- struct{}{}:
			peersUsage <- []ingest.PeerUsage{{Upload: uint(10), Download: uint(20), PublicKey: "xyz"}}
		case <-wait:
			t.Fatal("unexpected engine run termination")
			wait <- struct{}{}
		}
	}

	mux.Lock()
	readRestartMarkFileFn = func() ([1]byte, error) {
		return [1]byte{1}, nil
	}
	mux.Unlock()

	for i := 0; i < 1; i++ {
		select {
		case ticker <- struct{}{}:
			peersUsage <- []ingest.PeerUsage{{Upload: uint(5), Download: uint(10), PublicKey: "xyz"}}
		case <-wait:
			t.Fatal("unexpected engine run termination")
			wait <- struct{}{}
		}
	}
	close(ticker)
	<-wait
	require.Nil(t, runErr)
}

func (suite *EngineSingleStaticPeerWithRestartSuite) TestHugeImmediateDownloadRestart() {
	t := suite.T()
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	store := mocks.NewMockStore(ctrl)

	store.EXPECT().LoadBeforeRestartUsage(ctx).Return(map[string]ingest.PeerUsage{"xyz": {Upload: 10, Download: 20, PublicKey: "xyz"}}, nil).Times(1).
		After(store.EXPECT().LoadBeforeRestartUsage(ctx).Times(1))

	store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 510, Download: 1020, PublicKey: "xyz"}}).Return(nil).Times(1).
		After(store.EXPECT().IngestUsage(ctx, []ingest.PeerUsage{{Upload: 10, Download: 20, PublicKey: "xyz"}}).Return(nil).Times(1))

	peersUsage := make(chan []ingest.PeerUsage)
	readWGPeersUsage := func(ctx context.Context) ([]ingest.PeerUsage, error) {
		return <-peersUsage, nil
	}
	var mux sync.Mutex
	readRestartMarkFileFn := func() ([1]byte, error) {
		return [1]byte{0}, os.ErrNotExist
	}
	readRestartMarkFile := func(filename string) ([1]byte, error) {
		mux.Lock()
		defer mux.Unlock()
		return readRestartMarkFileFn()
	}
	removeRestartMarkFile := func(filename string) error {
		return nil
	}
	e := ingest.NewEngine(readWGPeersUsage, readRestartMarkFile, removeRestartMarkFile, store)
	ticker := make(chan struct{})

	var runErr error
	wait := make(chan struct{})
	go func() {
		defer func() {
			wait <- struct{}{}
		}()
		runErr = e.Run(ctx, ticker, "TODO")
	}()
	for i := 0; i < 1; i++ {
		select {
		case ticker <- struct{}{}:
			peersUsage <- []ingest.PeerUsage{{Upload: uint(10), Download: uint(20), PublicKey: "xyz"}}
		case <-wait:
			t.Fatal("unexpected engine run termination")
			wait <- struct{}{}
		}
	}

	mux.Lock()
	readRestartMarkFileFn = func() ([1]byte, error) {
		return [1]byte{1}, nil
	}
	mux.Unlock()

	for i := 0; i < 1; i++ {
		select {
		case ticker <- struct{}{}:
			peersUsage <- []ingest.PeerUsage{{Upload: uint(500), Download: uint(1000), PublicKey: "xyz"}}
		case <-wait:
			t.Fatal("unexpected engine run termination")
			wait <- struct{}{}
		}
	}
	close(ticker)
	<-wait
	require.Nil(t, runErr)
}

func TestEngineSingleStaticPeerWithRestart(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(EngineSingleStaticPeerWithRestartSuite))
}
