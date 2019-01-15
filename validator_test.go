package echalotte_test

import (
	rand "crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/t-bast/go-libp2p-echalotte"
	pb "github.com/t-bast/go-libp2p-echalotte/pb"

	"gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	"gx/ipfs/QmW7VUmSvhvSGbYbdsh7uRjhGmsYkc9fL8aJ5CorxxrU5N/go-crypto/nacl/box"
	"gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
	"gx/ipfs/QmdxUuburamoF6zF9qjeQC4WYcWGbWuRmdLacMEsW8ioD8/gogo-protobuf/proto"
	ptypes "gx/ipfs/QmdxUuburamoF6zF9qjeQC4WYcWGbWuRmdLacMEsW8ioD8/gogo-protobuf/types"
)

func TestPublicKeyValidator(t *testing.T) {
	aliceSigPrivKey, aliceSigPubKey, err := crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)

	alice, err := peer.IDFromPublicKey(aliceSigPubKey)
	require.NoError(t, err)

	aliceEncPubKey, _, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	pkv := &echalotte.PublicKeyValidator{}
	aliceRecord, err := pkv.CreateRecord(aliceSigPrivKey, aliceEncPubKey)
	require.NoError(t, err)

	t.Run("Validate()", func(t *testing.T) {
		t.Run("Invalid key format", func(t *testing.T) {
			err := pkv.Validate("/charles-baudelaire-rocks", aliceRecord)
			assert.EqualError(t, err, echalotte.ErrInvalidKeyFormat)
		})

		t.Run("Invalid key namespace", func(t *testing.T) {
			err := pkv.Validate(fmt.Sprintf("/charles/%s", alice.Pretty()), aliceRecord)
			assert.EqualError(t, err, echalotte.ErrInvalidNamespace)
		})

		t.Run("Invalid key peer ID", func(t *testing.T) {
			err := pkv.Validate("/enc/MyIDIsB4tm4n", aliceRecord)
			assert.Error(t, err)
		})

		t.Run("Invalid message format", func(t *testing.T) {
			err := pkv.Validate(pkv.CreateKey(alice), []byte{42})
			assert.Error(t, err)
		})

		t.Run("Invalid signature key", func(t *testing.T) {
			invalidSig := &pb.PublicKey{
				Type: pb.KeyType_Curve25519,
				Data: aliceEncPubKey[:],
			}

			toSign, err := proto.Marshal(invalidSig)
			require.NoError(t, err)

			invalidSig.Signature, err = aliceSigPrivKey.Sign(toSign)
			require.NoError(t, err)

			invalidSig.SignatureKey = []byte{42}

			invalidSigRecord, err := proto.Marshal(invalidSig)
			require.NoError(t, err)

			err = pkv.Validate(pkv.CreateKey(alice), invalidSigRecord)
			assert.Error(t, err)
		})

		t.Run("Signature key mismatch", func(t *testing.T) {
			validRecord, err := pkv.CreateRecord(aliceSigPrivKey, aliceEncPubKey)
			require.NoError(t, err)

			_, pk, err := crypto.GenerateEd25519Key(rand.Reader)
			require.NoError(t, err)

			otherPeerID, err := peer.IDFromPublicKey(pk)
			require.NoError(t, err)

			err = pkv.Validate(pkv.CreateKey(otherPeerID), validRecord)
			assert.Error(t, err)
		})

		t.Run("Invalid signature", func(t *testing.T) {
			invalidSig := &pb.PublicKey{
				Type:      pb.KeyType_Curve25519,
				CreatedAt: ptypes.TimestampNow(),
				Data:      aliceEncPubKey[:],
				Signature: []byte{42},
			}

			invalidSig.SignatureKey, err = aliceSigPubKey.Bytes()
			require.NoError(t, err)

			invalidRecord, err := proto.Marshal(invalidSig)
			require.NoError(t, err)

			err = pkv.Validate(pkv.CreateKey(alice), invalidRecord)
			assert.EqualError(t, err, echalotte.ErrInvalidSenderSignature)
		})

		t.Run("Valid record", func(t *testing.T) {
			err := pkv.Validate(pkv.CreateKey(alice), aliceRecord)
			assert.NoError(t, err)
		})
	})

	t.Run("Select()", func(t *testing.T) {
		key1, key2, err := box.GenerateKey(rand.Reader)
		require.NoError(t, err)

		record1, err := pkv.CreateRecord(aliceSigPrivKey, key1)
		require.NoError(t, err)

		<-time.After(2 * time.Second)

		record2, err := pkv.CreateRecord(aliceSigPrivKey, key2)
		require.NoError(t, err)

		t.Run("Skips invalid values", func(t *testing.T) {
			i, err := pkv.Select(pkv.CreateKey(alice), [][]byte{record1, []byte{42}, record2})
			require.NoError(t, err)
			assert.Equal(t, 2, i)
		})

		t.Run("Selects most recent key", func(t *testing.T) {
			i, err := pkv.Select(pkv.CreateKey(alice), [][]byte{record1, record2})
			require.NoError(t, err)
			assert.Equal(t, 1, i)
		})
	})
}
