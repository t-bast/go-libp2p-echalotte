package echalotte_test

import (
	crand "crypto/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/t-bast/go-libp2p-echalotte"

	"golang.org/x/crypto/nacl/box"

	"gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	"gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
)

func TestOnionMessage(t *testing.T) {
	aliceSignKey, _, err := crypto.GenerateEd25519Key(crand.Reader)
	require.NoError(t, err)

	alice, err := peer.IDFromPrivateKey(aliceSignKey)
	require.NoError(t, err)

	bobPubKey, bobPrivKey, err := box.GenerateKey(crand.Reader)
	require.NoError(t, err)

	bobSignKey, _, err := crypto.GenerateEd25519Key(crand.Reader)
	require.NoError(t, err)

	bob, err := peer.IDFromPrivateKey(bobSignKey)
	require.NoError(t, err)

	carolPubKey, carolPrivKey, err := box.GenerateKey(crand.Reader)
	require.NoError(t, err)

	carolSignKey, _, err := crypto.GenerateEd25519Key(crand.Reader)
	require.NoError(t, err)

	carol, err := peer.IDFromPrivateKey(carolSignKey)
	require.NoError(t, err)

	t.Run("NewMessage()", func(t *testing.T) {
		content := []byte("Ce ne seront jamais ces beautés de vignettes,")
		m, err := echalotte.NewMessage(alice, aliceSignKey, content)
		require.NoError(t, err)
		assert.NotNil(t, m)

		assert.Nil(t, m.To)
		assert.Equal(t, []byte(alice), m.From)

		senderKey, err := crypto.UnmarshalPublicKey(m.FromPublicKey)
		require.NoError(t, err)
		assert.True(t, alice.MatchesPublicKey(senderKey))

		assert.Equal(t, content, m.Content)
		assert.NotNil(t, m.Signature)
	})

	t.Run("Validate()", func(t *testing.T) {
		t.Run("Last Hop", func(t *testing.T) {
			t.Run("Valid message", func(t *testing.T) {
				m, err := echalotte.NewMessage(alice, aliceSignKey, []byte("Produits avariés, nés d’un siècle vaurien,"))
				require.NoError(t, err)
				require.NoError(t, m.Validate(alice))
			})

			t.Run("Missing from", func(t *testing.T) {
				m, err := echalotte.NewMessage(alice, aliceSignKey, []byte("Ces pieds à brodequins, ces doigts à castagnettes,"))
				require.NoError(t, err)

				m.From = nil
				err = m.Validate(alice)
				assert.Error(t, err)
				assert.True(t, strings.HasPrefix(err.Error(), echalotte.ErrInvalidSender))
			})

			t.Run("Invalid signature", func(t *testing.T) {
				m, err := echalotte.NewMessage(alice, aliceSignKey, []byte("Je laisse à Gavarni, poëte des chloroses,"))
				require.NoError(t, err)

				m.Signature[13]++
				err = m.Validate(alice)
				assert.Error(t, err)
				assert.True(t, strings.HasPrefix(err.Error(), echalotte.ErrInvalidSignature))
			})
		})

		t.Run("Intermediate Hop", func(t *testing.T) {
			t.Run("Valid Message", func(t *testing.T) {
				m := echalotte.OnionMessage{
					To:      []byte(alice),
					Content: []byte("Encrypted garbage."),
				}

				assert.NoError(t, m.Validate(alice))
			})

			t.Run("Missing recipient", func(t *testing.T) {
				m := echalotte.OnionMessage{
					Content: []byte("Encrypted garbage."),
				}

				assert.Error(t, m.Validate(alice))
			})

			t.Run("Recipient mismatch", func(t *testing.T) {
				m := echalotte.OnionMessage{
					To:      []byte(bob),
					Content: []byte("Encrypted garbage."),
				}

				err := m.Validate(alice)
				assert.EqualError(t, err, echalotte.ErrInvalidRecipient)
			})
		})
	})

	t.Run("IsLastHop()", func(t *testing.T) {
		t.Run("True", func(t *testing.T) {
			m := &echalotte.OnionMessage{
				Content: []byte("Qui sauront satisfaire un cœur comme le mien."),
			}

			assert.True(t, m.IsLastHop())
		})

		t.Run("False", func(t *testing.T) {
			m := &echalotte.OnionMessage{
				To:      []byte(alice),
				Content: []byte("Encrypted garbage."),
			}

			assert.False(t, m.IsLastHop())
		})
	})

	t.Run("Encapsulate()/Decapsulate()", func(t *testing.T) {
		t.Run("Encapsulate and decapsulate", func(t *testing.T) {
			content := []byte("C’est vous, Lady Macbeth, âme puissante au crime,")
			m, err := echalotte.NewMessage(alice, aliceSignKey, content)
			require.NoError(t, err)

			m, err = m.Encapsulate(carol, carolPubKey)
			require.NoError(t, err)
			require.Equal(t, []byte(carol), m.To)
			require.NoError(t, m.Validate(carol))

			m, err = m.Encapsulate(bob, bobPubKey)
			require.NoError(t, err)
			require.Equal(t, []byte(bob), m.To)
			require.NoError(t, m.Validate(bob))

			m, err = m.Decapsulate(bobSignKey, bobPrivKey)
			require.NoError(t, err)
			require.Equal(t, []byte(carol), m.To)
			require.NoError(t, m.Validate(carol))

			m, err = m.Decapsulate(carolSignKey, carolPrivKey)
			require.NoError(t, err)
			require.NoError(t, m.Validate(carol))
			require.Equal(t, content, m.Content)
		})

		t.Run("Decapsulate wrong recipient", func(t *testing.T) {
			content := []byte("Rêve d’Eschyle éclos au climat des autans;")
			m, err := echalotte.NewMessage(alice, aliceSignKey, content)
			require.NoError(t, err)

			m, err = m.Encapsulate(bob, bobPubKey)
			require.NoError(t, err)

			_, err = m.Decapsulate(carolSignKey, carolPrivKey)
			assert.EqualError(t, err, echalotte.ErrDecapsulateKey)
		})

		t.Run("Decapsulate key mismatch", func(t *testing.T) {
			content := []byte("Ou bien toi, grande Nuit, fille de Michel-Ange,")
			m, err := echalotte.NewMessage(alice, aliceSignKey, content)
			require.NoError(t, err)

			m, err = m.Encapsulate(bob, bobPubKey)
			require.NoError(t, err)

			_, err = m.Decapsulate(bobSignKey, bobPubKey)
			assert.EqualError(t, err, echalotte.ErrCouldNotDecrypt)
		})

		t.Run("Decapsulate last hop", func(t *testing.T) {
			content := []byte("Ou bien toi, grande Nuit, fille de Michel-Ange,")
			m, err := echalotte.NewMessage(alice, aliceSignKey, content)
			require.NoError(t, err)

			_, err = m.Decapsulate(aliceSignKey, nil)
			assert.EqualError(t, err, echalotte.ErrDecapsulateLastHop)
		})

		t.Run("Decapsulate altered message", func(t *testing.T) {
			content := []byte("Ou bien toi, grande Nuit, fille de Michel-Ange,")
			m, err := echalotte.NewMessage(alice, aliceSignKey, content)
			require.NoError(t, err)

			m, err = m.Encapsulate(bob, bobPubKey)
			require.NoError(t, err)

			m.Content[13]++
			m.Content[42]++

			_, err = m.Decapsulate(bobSignKey, bobPrivKey)
			assert.EqualError(t, err, echalotte.ErrCouldNotDecrypt)
		})
	})
}
