package echalotte

import (
	"context"
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/pkg/errors"

	"gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
	"gx/ipfs/QmYJtCabf3prS3HKQUGgqDLVxvbT9iDx5mfeVfhtCcJxxE/go-libp2p-discovery"
	"gx/ipfs/QmZ9zH2FnLcxv1xyzFeUpDUeo55xEhZQHgveZijcxr7TLj/go-libp2p-peerstore"
)

const (
	// OnionRelay is the name of the namespace advertized by onion relays.
	OnionRelay = "/libp2p/onion"

	// DefaultCircuitSize is the default size of the onion circuit.
	// This is configurable.
	DefaultCircuitSize = 5

	// DefaultCircuitTimeout is the default timeout for circuit building.
	DefaultCircuitTimeout = 30 * time.Second
)

// Errors used by the CircuitBuilder.
const (
	ErrAdvertise          = "failed to advertise onion service"
	ErrFindRelays         = "failed to find enough onion relays in the network"
	ErrInvalidCircuitSize = "circuit size should be strictly positive"
)

// CircuitOption is a single circuit option.
type CircuitOption func(opts *CircuitOptions) error

// CircuitOptions is a set of circuit options.
type CircuitOptions struct {
	Size    int
	Timeout time.Duration
}

// Apply the given options to this CircuitOptions.
func (opts *CircuitOptions) Apply(options ...CircuitOption) error {
	for _, o := range options {
		if err := o(opts); err != nil {
			return err
		}
	}

	return nil
}

// CircuitSize is an option to choose the size of the circuit.
func CircuitSize(size int) CircuitOption {
	return func(opts *CircuitOptions) error {
		if size <= 0 {
			return errors.New(ErrInvalidCircuitSize)
		}

		opts.Size = size
		return nil
	}
}

// CircuitTimeout is an option to choose the timeout for circuit building.
func CircuitTimeout(timeout time.Duration) CircuitOption {
	return func(opts *CircuitOptions) error {
		opts.Timeout = timeout
		return nil
	}
}

// Circuit can be used to create an onion-routed message.
type Circuit []peer.ID

// CircuitBuilder lets you build random circuits for onion routing.
type CircuitBuilder struct {
	discover discovery.Discoverer
}

// NewCircuitBuilder creates a new circuit builder that leverages the given
// discovery component to find other peers that provide onion relays.
func NewCircuitBuilder(ctx context.Context, discover discovery.Discovery) (*CircuitBuilder, error) {
	// TODO: set appropriate TTL option and auto-refresh in go routine.
	// This will prevent peers that aren't responsive from being chosen in an
	// onion circuit.
	// Test how it behaves in a real network.
	_, err := discover.Advertise(ctx, OnionRelay)
	if err != nil {
		return nil, errors.Wrap(err, ErrAdvertise)
	}

	return &CircuitBuilder{
		discover: discover,
	}, nil
}

// Build a random circuit between network relay peers.
func (cb *CircuitBuilder) Build(ctx context.Context, opts ...CircuitOption) (Circuit, error) {
	options := &CircuitOptions{
		Size:    DefaultCircuitSize,
		Timeout: DefaultCircuitTimeout,
	}
	err := options.Apply(opts...)
	if err != nil {
		return nil, err
	}

	// Randomize the limit to prevent attackers from discovering the circuit
	// size by analyzing DHT requests.
	rlimit, _ := crand.Int(crand.Reader, big.NewInt(int64(2*options.Size)))
	limit := 4*options.Size + int(rlimit.Int64())

	// TODO: test how this behaves in a real network with an underlying kademlia DHT.
	peerChan, err := cb.discover.FindPeers(ctx, OnionRelay, discovery.Limit(limit))
	if err != nil {
		return nil, errors.Wrap(err, ErrFindRelays)
	}

	// Collect more peers than the circuit size.
	// Randomize the number of peers chosen to obfuscate circuit size.
	rcount, _ := crand.Int(crand.Reader, big.NewInt(int64(options.Size)))
	minRelaysCount := 2*options.Size + int(rcount.Int64())
	relays, err := cb.findRelays(peerChan, minRelaysCount, options.Timeout)
	if err != nil {
		return nil, err
	}

	circuitRelays := cb.selectRelays(relays, options.Size)
	return circuitRelays, nil
}

// findRelays synchronously finds the request number of relay peers.
// If it can't find enough peers before the timeout expires, it will return an
// error.
func (cb *CircuitBuilder) findRelays(
	peerChan <-chan peerstore.PeerInfo,
	count int,
	timeout time.Duration,
) ([]peerstore.PeerInfo, error) {
	relays := make([]peerstore.PeerInfo, count)

	errChan := make(chan error, count)
	wg := sync.WaitGroup{}

	for i := 0; i < count; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			select {
			case peerInfo, ok := <-peerChan:
				if !ok {
					errChan <- errors.New("peers channel closed")
					return
				}

				relays[i] = peerInfo
			case <-time.After(timeout):
				errChan <- errors.New("peers channel timed out")
			}
		}(i)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		return nil, errors.Wrap(err, ErrFindRelays)
	default:
		break
	}

	return relays, nil
}

// selectRelays randomly selects a subset of the available relays.
func (cb *CircuitBuilder) selectRelays(relays []peerstore.PeerInfo, count int) Circuit {
	seed, _ := crand.Int(crand.Reader, big.NewInt(1<<62))
	rand.Seed(seed.Int64())
	rand.Shuffle(len(relays), func(i, j int) { relays[i], relays[j] = relays[j], relays[i] })

	c := make([]peer.ID, count)
	for i := 0; i < count; i++ {
		c[i] = relays[i].ID
	}

	return c
}
