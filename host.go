package echalotte

import (
	"context"
	crand "crypto/rand"
	"encoding/json"
	"time"

	pb "github.com/t-bast/go-libp2p-echalotte/pb"

	inet "gx/ipfs/QmNgLg1NTw37iWbYPKcyK85YJ9Whs1MkPtJwhfqbNYAyKg/go-libp2p-net"
	ropts "gx/ipfs/QmTiRqrF5zkdZyrdsL5qndG1UbeWi8k8N2pYxCtXWrahR2/go-libp2p-routing/options"
	"gx/ipfs/QmVmDhyTTUcQXFD1rRQ64fGLMSAoaQvNH3hwuaCFAPq2hy/errors"
	"gx/ipfs/QmW7VUmSvhvSGbYbdsh7uRjhGmsYkc9fL8aJ5CorxxrU5N/go-crypto/nacl/box"
	"gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmaoXrM4Z41PD48JY36YqQGKQpLGjyLA2cKcLsES7YddAq/go-libp2p-host"
	"gx/ipfs/QmdxUuburamoF6zF9qjeQC4WYcWGbWuRmdLacMEsW8ioD8/gogo-protobuf/proto"
)

const (
	// ProtocolID is the ID for the echalotte protocol.
	ProtocolID = protocol.ID("/echalotte/v1.0.0")

	publicKeyStoreKey  = "/encryption/publickey"
	privateKeyStoreKey = "/encryption/privatekey"
)

// Errors used by the host.
const (
	ErrInvalidEncryptionKey = "invalid key: not a curve25519 key"
)

// DHT interface needed to advertise encryption keys in the network.
type DHT interface {
	PutValue(context.Context, string, []byte, ...ropts.Option) error
	GetValue(context.Context, string, ...ropts.Option) ([]byte, error)
}

// Host wraps a standard host with onion routing capabilities.
type Host struct {
	host.Host

	dht            DHT
	circuitBuilder CircuitBuilder
	validator      *PublicKeyValidator
}

// Connect to the echalotte network.
// This will block until enough peers have been discovered.
// It then returns a super-powered host instance that can use onion routing.
func Connect(ctx context.Context, host host.Host, dht DHT, cb CircuitBuilder) (*Host, error) {
	h := &Host{
		Host:           host,
		dht:            dht,
		circuitBuilder: cb,
		validator:      &PublicKeyValidator{},
	}

	_, err := h.DecryptionKey()
	if err != nil {
		err = h.registerEncryptionKey(ctx)
		if err != nil {
			return nil, err
		}
	}

	h.SetStreamHandler(ProtocolID, func(stream inet.Stream) {
		ctx := context.Background()
		err := h.HandleMessage(ctx, stream)
		if err != nil {
			log.Errorf("Message error: %s", err.Error())
		}
	})

	// Test the network readiness by generating a sample circuit.
	for {
		_, err = cb.Build(ctx)
		if err == nil {
			break
		}

		select {
		case <-ctx.Done():
			return h, errors.WithStack(ctx.Err())
		case <-time.After(30 * time.Second):
			continue
		}
	}

	return h, nil
}

// EncryptionKey that other peers can use to encrypt messages for the current
// host.
func (h *Host) EncryptionKey() (*[32]byte, error) {
	k, err := h.Peerstore().Get(h.ID(), publicKeyStoreKey)
	if err != nil {
		return nil, err
	}

	publicKey, ok := k.(*[32]byte)
	if !ok {
		return nil, errors.New(ErrInvalidEncryptionKey)
	}

	return publicKey, nil
}

// DecryptionKey that the current host can use to decrypt messages.
func (h *Host) DecryptionKey() (*[32]byte, error) {
	k, err := h.Peerstore().Get(h.ID(), privateKeyStoreKey)
	if err != nil {
		return nil, err
	}

	privateKey, ok := k.(*[32]byte)
	if !ok {
		return nil, errors.New(ErrInvalidEncryptionKey)
	}

	return privateKey, nil
}

func (h *Host) registerEncryptionKey(ctx context.Context) error {
	encryptionPublicKey, encryptionPrivateKey, err := box.GenerateKey(crand.Reader)
	if err != nil {
		return errors.WithStack(err)
	}

	err = h.Peerstore().Put(h.ID(), privateKeyStoreKey, encryptionPrivateKey)
	if err != nil {
		return errors.WithStack(err)
	}

	err = h.Peerstore().Put(h.ID(), publicKeyStoreKey, encryptionPublicKey)
	if err != nil {
		return errors.WithStack(err)
	}

	dhtRecord, err := h.validator.CreateRecord(h.Peerstore().PrivKey(h.ID()), encryptionPublicKey)
	if err != nil {
		return errors.WithStack(err)
	}

	err = h.dht.PutValue(ctx, h.validator.CreateKey(h.ID()), dhtRecord)
	if err != nil {
		return errors.WithStack(err)
	}

	log.Info("Encryption key registered")

	return nil
}

// SendMessage sends a private message to the given peer.
// It leverages onion routing through the echalotte network.
func (h *Host) SendMessage(ctx context.Context, to peer.ID, message []byte) error {
	circuit, err := h.circuitBuilder.Build(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	m, err := NewMessage(h.ID(), h.Peerstore().PrivKey(h.ID()), message)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, relay := range circuit {
		key, err := h.peerEncryptionKey(ctx, relay)
		if err != nil {
			return errors.Wrapf(err, "could not get encryption key for %s", relay.Pretty())
		}

		m, err = m.Encapsulate(relay, key)
		if err != nil {
			return errors.Wrapf(err, "could not encapsulate to peer %s", relay.Pretty())
		}
	}

	firstHop := circuit[len(circuit)-1]
	stream, err := h.NewStream(ctx, firstHop, ProtocolID)
	if err != nil {
		return errors.WithStack(err)
	}
	defer stream.Close()

	enc := json.NewEncoder(stream)
	return enc.Encode(m)
}

func (h *Host) peerEncryptionKey(ctx context.Context, peerID peer.ID) (*[32]byte, error) {
	record, err := h.dht.GetValue(ctx, h.validator.CreateKey(peerID))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var peerKey pb.PublicKey
	err = proto.Unmarshal(record, &peerKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var pubKey [32]byte
	copy(pubKey[:], peerKey.Data)

	return &pubKey, nil
}

// HandleMessage receives an onion message and forwards it.
// If we are the message recipient we print it.
func (h *Host) HandleMessage(ctx context.Context, stream inet.Stream) error {
	defer stream.Close()

	var message *OnionMessage

	dec := json.NewDecoder(stream)
	err := dec.Decode(&message)
	if err != nil {
		return errors.WithStack(err)
	}

	err = message.Validate(h.ID())
	if err != nil {
		return errors.WithStack(err)
	}

	decryptionKey, err := h.DecryptionKey()
	if err != nil {
		return errors.WithStack(err)
	}

	message, err = message.Decapsulate(h.Peerstore().PrivKey(h.ID()), decryptionKey)
	if err != nil {
		return errors.WithStack(err)
	}

	if message.IsLastHop() {
		from, _ := peer.IDFromBytes(message.From)
		log.Infof("Private message received from %s: %s", from.Pretty(), message.Content)
		return nil
	}

	go func() {
		err := h.forwardMessage(context.Background(), message)
		if err != nil {
			log.Errorf("Could not forward message: %s", err.Error())
		}
	}()

	return nil
}

// Forward a message to the next recipient.
func (h *Host) forwardMessage(ctx context.Context, message *OnionMessage) error {
	to, _ := peer.IDFromBytes(message.To)
	stream, err := h.NewStream(ctx, to, ProtocolID)
	if err != nil {
		return errors.WithStack(err)
	}
	defer stream.Close()

	enc := json.NewEncoder(stream)
	return enc.Encode(message)
}
