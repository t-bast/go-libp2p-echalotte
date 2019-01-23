package echalotte_test

import (
	"context"
	crand "crypto/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/t-bast/go-libp2p-echalotte"
	"github.com/t-bast/go-libp2p-echalotte/echalottetesting"
	pb "github.com/t-bast/go-libp2p-echalotte/pb"

	"gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	"gx/ipfs/QmPiemjiKBC9VA7vZF82m4x1oygtg2c2YVqag8PX7dN1BD/go-libp2p-peerstore"
	"gx/ipfs/QmW7VUmSvhvSGbYbdsh7uRjhGmsYkc9fL8aJ5CorxxrU5N/go-crypto/nacl/box"
	"gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
	"gx/ipfs/QmdxUuburamoF6zF9qjeQC4WYcWGbWuRmdLacMEsW8ioD8/gogo-protobuf/proto"
)

func TestHost(t *testing.T) {
	alicePrivateKey, _, err := crypto.GenerateEd25519Key(crand.Reader)
	require.NoError(t, err)

	alice, err := peer.IDFromPrivateKey(alicePrivateKey)
	require.NoError(t, err)

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

			connectChan := make(chan struct{})
			go func() {
				echalotte.Connect(
					ctx,
					h,
					echalottetesting.NewInMemoryDHT(),
					echalottetesting.NewDummyCircuitBuilder(t).StartFailing(),
				)
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

	t.Run("SendMessage()", func(t *testing.T) {
		t.Run("circuit error", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			cb := echalottetesting.NewDummyCircuitBuilder(t)
			h, err := echalotte.Connect(
				ctx,
				echalottetesting.RandomHost(ctx, t),
				echalottetesting.NewInMemoryDHT(),
				cb,
			)
			require.NoError(t, err)

			cb.StartFailing()

			err = h.SendMessage(ctx, alice, []byte("Rappelez-vous l'objet que nous vîmes, mon âme,"))
			assert.EqualError(t, err, echalottetesting.ErrBuildCircuit)
		})

		t.Run("encryption key not found", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			h, err := echalotte.Connect(
				ctx,
				echalottetesting.RandomHost(ctx, t),
				echalottetesting.NewInMemoryDHT(),
				echalottetesting.NewDummyCircuitBuilder(t),
			)
			require.NoError(t, err)

			err = h.SendMessage(ctx, alice, []byte("Ce beau matin d'été si doux :"))
			assert.Error(t, err)
			assert.True(t, strings.HasPrefix(err.Error(), "could not get encryption key"))
		})

		t.Run("peer not responding", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			dht := echalottetesting.NewInMemoryDHT()

			var relays []peer.ID
			for i := 0; i < 5; i++ {
				sk, _, _ := crypto.GenerateEd25519Key(crand.Reader)
				pk, _, _ := box.GenerateKey(crand.Reader)
				relayID, _ := peer.IDFromPrivateKey(sk)
				relays = append(relays, relayID)

				v := &echalotte.PublicKeyValidator{}
				record, _ := v.CreateRecord(sk, pk)
				dht.PutValue(ctx, v.CreateKey(relayID), record)
			}

			cb := echalottetesting.NewDummyCircuitBuilderFromNetwork(t, relays)

			h, err := echalotte.Connect(
				ctx,
				echalottetesting.RandomHost(ctx, t),
				dht,
				cb,
			)
			require.NoError(t, err)

			err = h.SendMessage(ctx, alice, []byte("Au détour d'un sentier une charogne infâme"))
			assert.Error(t, err)
			assert.True(t, strings.HasPrefix(err.Error(), "dial attempt failed"))
		})

		t.Run("success", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			dht := echalottetesting.NewInMemoryDHT()

			var relays []peer.ID
			var relaysKey []crypto.PrivKey
			for i := 0; i < 5; i++ {
				sk, _, _ := crypto.GenerateEd25519Key(crand.Reader)
				pk, _, _ := box.GenerateKey(crand.Reader)
				relayID, _ := peer.IDFromPrivateKey(sk)
				relays = append(relays, relayID)
				relaysKey = append(relaysKey, sk)

				v := &echalotte.PublicKeyValidator{}
				record, _ := v.CreateRecord(sk, pk)
				dht.PutValue(ctx, v.CreateKey(relayID), record)
			}

			cb := echalottetesting.NewDummyCircuitBuilderFromNetwork(t, relays)

			h1, err := echalotte.Connect(
				ctx,
				echalottetesting.RandomHost(ctx, t),
				dht,
				cb,
			)
			require.NoError(t, err)

			h2, err := echalotte.Connect(
				ctx,
				echalottetesting.HostWithIdentity(ctx, t, relaysKey[4]),
				dht,
				cb,
			)
			require.NoError(t, err)

			h1.Peerstore().AddAddrs(h2.ID(), h2.Addrs(), peerstore.AddressTTL)

			err = h1.SendMessage(ctx, h2.ID(), []byte("Sur un lit semé de cailloux,"))
			assert.NoError(t, err)
		})
	})
}
