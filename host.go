package echalotte

import (
	"context"
	crand "crypto/rand"
	"time"

	inet "gx/ipfs/QmNgLg1NTw37iWbYPKcyK85YJ9Whs1MkPtJwhfqbNYAyKg/go-libp2p-net"
	ropts "gx/ipfs/QmTiRqrF5zkdZyrdsL5qndG1UbeWi8k8N2pYxCtXWrahR2/go-libp2p-routing/options"
	"gx/ipfs/QmVmDhyTTUcQXFD1rRQ64fGLMSAoaQvNH3hwuaCFAPq2hy/errors"
	"gx/ipfs/QmW7VUmSvhvSGbYbdsh7uRjhGmsYkc9fL8aJ5CorxxrU5N/go-crypto/nacl/box"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmaoXrM4Z41PD48JY36YqQGKQpLGjyLA2cKcLsES7YddAq/go-libp2p-host"
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

	h.SetStreamHandler(ProtocolID, h.handleStream)

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

func (h *Host) handleStream(stream inet.Stream) {
	// TODO: handle streams ;)
}
