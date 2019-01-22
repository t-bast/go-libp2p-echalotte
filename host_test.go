package echalotte_test

import (
	"context"
	crand "crypto/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/t-bast/go-libp2p-echalotte"
	"github.com/t-bast/go-libp2p-echalotte/echalottetesting"
	pb "github.com/t-bast/go-libp2p-echalotte/pb"

	"gx/ipfs/QmW7VUmSvhvSGbYbdsh7uRjhGmsYkc9fL8aJ5CorxxrU5N/go-crypto/nacl/box"
	"gx/ipfs/QmdxUuburamoF6zF9qjeQC4WYcWGbWuRmdLacMEsW8ioD8/gogo-protobuf/proto"
)

func TestHost(t *testing.T) {
	t.Run("Encryption keys", func(t *testing.T) {
		circuitBuilder := echalottetesting.NewDummyCircuitBuilder(t, echalotte.CircuitSize(3))

		t.Run("already generated", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			h := echalottetesting.RandomHost(ctx, t)
			encryptionPublicKey, encryptionPrivateKey, _ := box.GenerateKey(crand.Reader)

			err := h.Peerstore().Put(h.ID(), "/encryption/publickey", encryptionPublicKey)
			require.NoError(t, err)

			err = h.Peerstore().Put(h.ID(), "/encryption/privatekey", encryptionPrivateKey)
			require.NoError(t, err)

			eh, err := echalotte.Connect(ctx, h, nil, circuitBuilder)
			require.NoError(t, err)

			decryptionKey, err := eh.DecryptionKey()
			require.NoError(t, err)
			assert.Equal(t, encryptionPrivateKey, decryptionKey)

			encryptionKey, err := eh.EncryptionKey()
			require.NoError(t, err)
			assert.Equal(t, encryptionPublicKey, encryptionKey)
		})

		t.Run("corrupted store", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			h := echalottetesting.RandomHost(ctx, t)

			err := h.Peerstore().Put(h.ID(), "/encryption/publickey", 42)
			require.NoError(t, err)

			err = h.Peerstore().Put(h.ID(), "/encryption/privatekey", 4242)
			require.NoError(t, err)

			eh, err := echalotte.Connect(ctx, h, echalottetesting.NewInMemoryDHT(), circuitBuilder)
			require.NoError(t, err)

			decryptionKey, err := eh.DecryptionKey()
			require.NoError(t, err)
			require.NotNil(t, decryptionKey)

			encryptionKey, err := eh.EncryptionKey()
			require.NoError(t, err)
			require.NotNil(t, encryptionKey)
		})

		t.Run("generate when connecting", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			h := echalottetesting.RandomHost(ctx, t)
			dht := echalottetesting.NewInMemoryDHT()

			eh, err := echalotte.Connect(ctx, h, dht, circuitBuilder)
			require.NoError(t, err)

			decryptionKey, err := eh.DecryptionKey()
			require.NoError(t, err)
			require.NotNil(t, decryptionKey)

			encryptionKey, err := eh.EncryptionKey()
			require.NoError(t, err)
			require.NotNil(t, encryptionKey)

			dhtKey := echalotte.PublicKeyValidator{}.CreateKey(h.ID())
			dhtValue, err := dht.GetValue(ctx, dhtKey)
			require.NoError(t, err)

			var dhtEncryptionKey pb.PublicKey
			err = proto.Unmarshal(dhtValue, &dhtEncryptionKey)
			require.NoError(t, err)

			assert.Equal(t, encryptionKey[:], dhtEncryptionKey.Data)
		})
	})

	t.Run("Connect()", func(t *testing.T) {
		t.Run("succeeds after generating sample circuit", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			h := echalottetesting.RandomHost(ctx, t)
			dht := echalottetesting.NewInMemoryDHT()

			_, err := echalotte.Connect(ctx, h, dht, echalottetesting.NewDummyCircuitBuilder(t))
			require.NoError(t, err)
		})

		t.Run("blocks until circuit can be generated", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			h := echalottetesting.RandomHost(ctx, t)
			dht := echalottetesting.NewInMemoryDHT()

			connectChan := make(chan struct{})
			go func() {
				echalotte.Connect(ctx, h, dht, echalottetesting.NewFailingCircuitBuilder())
				connectChan <- struct{}{}
			}()

			select {
			case <-connectChan:
				assert.Fail(t, "connect should not succeed")
			case <-time.After(200 * time.Millisecond):
				return
			}
		})
	})
}
