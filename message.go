package echalotte

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/json"

	"github.com/pkg/errors"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"

	"gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	"gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
)

// Errors used by the messaging layer.
const (
	ErrCouldNotDecrypt    = "could not decrypt ciphertext"
	ErrCouldNotSign       = "could not sign message"
	ErrDecapsulateKey     = "cannot decapsulate: private key doesn't match recipient"
	ErrDecapsulateLastHop = "cannot decapsulate: last hop reached"
	ErrInvalidRecipient   = "invalid message recipient: ID doesn't match ours"
	ErrInvalidSender      = "invalid message sender"
	ErrInvalidSignature   = "invalid message signature"
	ErrMarshal            = "could not marshal/unmarshal message"
)

// OnionMessage contains the next recipient and some bytes supposedly encrypted
// for that next recipient.
type OnionMessage struct {
	To            []byte
	From          []byte
	FromPublicKey []byte
	Content       []byte
	Signature     []byte
}

// NewMessage creates a new signed message.
func NewMessage(from peer.ID, sk crypto.PrivKey, content []byte) (*OnionMessage, error) {
	fromKey, err := sk.GetPublic().Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	m := &OnionMessage{
		From:          []byte(from),
		FromPublicKey: fromKey,
		Content:       content,
	}

	// Signing before encryption provides forward secrecy.
	sig, err := sk.Sign(content)
	if err != nil {
		return nil, errors.Wrap(err, ErrCouldNotSign)
	}

	m.Signature = sig
	return m, nil
}

// Validate the onion message.
// Verifies that we're the correct recipient for the current layer.
func (l *OnionMessage) Validate(id peer.ID) error {
	if l.IsLastHop() {
		return l.validateLastHop()
	}

	return l.validateIntermediateHop(id)
}

func (l *OnionMessage) validateLastHop() error {
	from, err := peer.IDFromBytes(l.From)
	if err != nil {
		return errors.Wrap(err, ErrInvalidSender)
	}

	fromKey, err := crypto.UnmarshalPublicKey(l.FromPublicKey)
	if err != nil {
		return errors.Wrap(err, ErrInvalidSender)
	}

	if !from.MatchesPublicKey(fromKey) {
		return errors.New(ErrInvalidSender)
	}

	ok, err := fromKey.Verify(l.Content, l.Signature)
	if err != nil {
		return errors.Wrap(err, ErrInvalidSignature)
	}

	if !ok {
		return errors.New(ErrInvalidSignature)
	}

	return nil
}

func (l *OnionMessage) validateIntermediateHop(id peer.ID) error {
	to, err := peer.IDFromBytes(l.To)
	if err != nil {
		return errors.WithStack(err)
	}

	if to != id {
		return errors.New(ErrInvalidRecipient)
	}

	return nil
}

// IsLastHop returns true when we reached the last hop.
// The message should contain the plaintext content, the sender and a sender
// signature.
func (l *OnionMessage) IsLastHop() bool {
	return len(l.To) == 0
}

// Encapsulate adds another layer of onion encryption.
// The recipient key should be an NaCl public key (curve25519 point).
func (l *OnionMessage) Encapsulate(to peer.ID, publicKey *[32]byte) (*OnionMessage, error) {
	// TODO: should use a better encoding than JSON.
	b, err := json.Marshal(l)
	if err != nil {
		return nil, errors.Wrap(err, ErrMarshal)
	}

	ciphertext, err := seal(publicKey, b)
	if err != nil {
		return nil, err
	}

	return &OnionMessage{
		To:      []byte(to),
		Content: ciphertext,
	}, nil
}

// Decapsulate decrypts the content as another onion layer.
// The first argument is the peer's signing private key.
// The second argument is the peer's encryption private key (curve25519 point).
func (l *OnionMessage) Decapsulate(signingKey crypto.PrivKey, encryptionPrivKey *[32]byte) (*OnionMessage, error) {
	if l.IsLastHop() {
		return nil, errors.New(ErrDecapsulateLastHop)
	}

	to, err := peer.IDFromBytes(l.To)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !to.MatchesPrivateKey(signingKey) {
		return nil, errors.New(ErrDecapsulateKey)
	}

	content, err := open(encryptionPrivKey, l.Content)
	if err != nil {
		return nil, err
	}

	// TODO: should use a better encoding than JSON.
	var onionContent OnionMessage
	err = json.Unmarshal(content, &onionContent)
	if err != nil {
		return nil, errors.Wrap(err, ErrMarshal)
	}

	return &onionContent, nil
}

// Seal a message to a given peer with no sender authentication.
// We use throw-away ephemeral keys to provide such a feature.
func seal(to *[32]byte, message []byte) ([]byte, error) {
	// Ephemeral throw-away key pair.
	epk, esk, err := box.GenerateKey(crand.Reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var nonce [24]byte
	keysHash := sha256.Sum256(append(epk[:], to[:]...))
	copy(nonce[:], keysHash[:])

	var ciphertext []byte
	ciphertext = box.Seal(epk[:], message, &nonce, to, esk)

	return ciphertext, nil
}

// Open a sealed unauthenticated message.
// An NaCl private key (curve25519 point) is required.
func open(key *[32]byte, ciphertext []byte) ([]byte, error) {
	var pubKey [32]byte
	curve25519.ScalarBaseMult(&pubKey, key)

	// Extract ephemeral key and nonce from ciphertext.
	var epk [32]byte
	copy(epk[:], ciphertext[:32])

	expectedNonce := sha256.Sum256(append(epk[:], pubKey[:]...))
	var nonce [24]byte
	copy(nonce[:], expectedNonce[:])

	decrypted, ok := box.Open(nil, ciphertext[32:], &nonce, &epk, key)
	if !ok {
		return nil, errors.New(ErrCouldNotDecrypt)
	}

	return decrypted, nil
}
