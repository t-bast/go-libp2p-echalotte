package echalotte_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/t-bast/go-libp2p-echalotte"
	"github.com/t-bast/go-libp2p-echalotte/mocks"
)

func TestCircuitBuilder(t *testing.T) {
	t.Run("New()", func(t *testing.T) {
		t.Run("forwards advertiser error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			discover := mocks.NewMockDiscovery(ctrl)
			discover.EXPECT().Advertise(gomock.Any(), echalotte.OnionRelay).Return(
				42*time.Millisecond,
				errors.New("fatal"),
			)

			cb, err := echalotte.NewCircuitBuilder(context.Background(), discover)
			assert.EqualError(t, err, "fatal")
			assert.Nil(t, cb)
		})

		t.Run("advertises onion relay", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			discover := mocks.NewMockDiscovery(ctrl)
			discover.EXPECT().Advertise(gomock.Any(), echalotte.OnionRelay)

			cb, err := echalotte.NewCircuitBuilder(context.Background(), discover)
			require.NoError(t, err)
			assert.NotNil(t, cb)
		})
	})
}
