package echalotte_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/t-bast/go-libp2p-echalotte"
	"github.com/t-bast/go-libp2p-echalotte/mocks"

	"gx/ipfs/QmPiemjiKBC9VA7vZF82m4x1oygtg2c2YVqag8PX7dN1BD/go-libp2p-peerstore"
	"gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
)

// Create a test circuit builder using the given discovery mock.
func newTestCircuitBuilder(t *testing.T, discover *mocks.MockDiscovery) *echalotte.CircuitBuilder {
	discover.EXPECT().Advertise(gomock.Any(), echalotte.OnionRelay)

	cb, err := echalotte.NewCircuitBuilder(context.Background(), discover)
	require.NoError(t, err)

	return cb
}

func TestCircuitBuilder(t *testing.T) {
	t.Run("New()", func(t *testing.T) {
		t.Run("wraps advertiser error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			discover := mocks.NewMockDiscovery(ctrl)
			discover.EXPECT().Advertise(gomock.Any(), echalotte.OnionRelay).Return(
				42*time.Millisecond,
				errors.New("fatal"),
			)

			cb, err := echalotte.NewCircuitBuilder(context.Background(), discover)
			assert.EqualError(t, errors.Cause(err), "fatal")
			assert.Nil(t, cb)
		})

		t.Run("advertises onion relay", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			discover := mocks.NewMockDiscovery(ctrl)
			cb := newTestCircuitBuilder(t, discover)
			assert.NotNil(t, cb)
		})
	})

	t.Run("Build()", func(t *testing.T) {
		t.Run("rejects invalid circuit size", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			discover := mocks.NewMockDiscovery(ctrl)
			cb := newTestCircuitBuilder(t, discover)

			c, err := cb.Build(context.Background(), echalotte.CircuitSize(0))
			assert.EqualError(t, err, echalotte.ErrInvalidCircuitSize)
			assert.Nil(t, c)
		})

		t.Run("wraps find peers error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			discover := mocks.NewMockDiscovery(ctrl)
			cb := newTestCircuitBuilder(t, discover)

			discover.EXPECT().FindPeers(gomock.Any(), echalotte.OnionRelay, gomock.Any()).Return(nil, errors.New("fatal"))

			c, err := cb.Build(context.Background())
			assert.EqualError(t, errors.Cause(err), "fatal")
			assert.Nil(t, c)
		})

		t.Run("returns error if not enough relays discovered", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			discover := mocks.NewMockDiscovery(ctrl)
			cb := newTestCircuitBuilder(t, discover)

			// Simulate only 4 relays found whereas default circuit size is 5.
			relaysChan := make(chan peerstore.PeerInfo)
			go func() {
				for i := 0; i < 4; i++ {
					relaysChan <- peerstore.PeerInfo{ID: peer.ID(i)}
				}

				close(relaysChan)
			}()

			discover.EXPECT().FindPeers(gomock.Any(), echalotte.OnionRelay, gomock.Any()).Return(relaysChan, nil)

			c, err := cb.Build(context.Background(), echalotte.CircuitTimeout(10*time.Millisecond))
			assert.Error(t, err)
			assert.True(t, strings.HasPrefix(err.Error(), echalotte.ErrFindRelays))
			assert.True(t, strings.Contains(err.Error(), "channel"))
			assert.Nil(t, c)
		})

		t.Run("times out when finding relay peers", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			discover := mocks.NewMockDiscovery(ctrl)
			cb := newTestCircuitBuilder(t, discover)

			// Simulate slow rate of finding relay peers.
			// Timeout should kick in.
			relaysChan := make(chan peerstore.PeerInfo)
			go func() {
				for i := 0; i < 100; i++ {
					relaysChan <- peerstore.PeerInfo{ID: peer.ID(i)}
					<-time.After(3 * time.Millisecond)
				}

				close(relaysChan)
			}()

			discover.EXPECT().FindPeers(gomock.Any(), echalotte.OnionRelay, gomock.Any()).Return(relaysChan, nil)

			c, err := cb.Build(context.Background(), echalotte.CircuitTimeout(10*time.Millisecond))
			assert.Error(t, err)
			assert.True(t, strings.HasPrefix(err.Error(), echalotte.ErrFindRelays))
			assert.True(t, strings.Contains(err.Error(), "timed out"))
			assert.Nil(t, c)
		})

		t.Run("builds random circuit", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			discover := mocks.NewMockDiscovery(ctrl)
			cb := newTestCircuitBuilder(t, discover)

			relaysChan := make(chan peerstore.PeerInfo)
			go func() {
				for i := 0; i < 100; i++ {
					relaysChan <- peerstore.PeerInfo{ID: peer.ID(i)}
				}

				close(relaysChan)
			}()

			discover.EXPECT().FindPeers(gomock.Any(), echalotte.OnionRelay, gomock.Any()).Return(relaysChan, nil)

			c, err := cb.Build(context.Background(),
				echalotte.CircuitSize(5),
				echalotte.CircuitTimeout(10*time.Millisecond),
			)
			require.NoError(t, err)
			require.Len(t, c, 5)
			assert.NotSubset(t, c, []peer.ID{peer.ID(0), peer.ID(1), peer.ID(2), peer.ID(3), peer.ID(4)})
		})
	})
}
