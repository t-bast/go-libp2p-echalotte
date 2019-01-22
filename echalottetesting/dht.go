package echalottetesting

import (
	"context"
	"errors"
	"sync"

	ropts "gx/ipfs/QmTiRqrF5zkdZyrdsL5qndG1UbeWi8k8N2pYxCtXWrahR2/go-libp2p-routing/options"
)

// InMemoryDHT provides an in-memory implementation of a DHT.
// It's a simply key-value pair without network capabilities.
type InMemoryDHT struct {
	lock   sync.RWMutex
	values map[string][]byte
}

// NewInMemoryDHT creates a new DHT.
func NewInMemoryDHT() *InMemoryDHT {
	return &InMemoryDHT{
		values: make(map[string][]byte),
	}
}

// PutValue adds a value to the DHT.
func (dht *InMemoryDHT) PutValue(_ context.Context, key string, value []byte, _ ...ropts.Option) error {
	dht.lock.Lock()
	defer dht.lock.Unlock()

	dht.values[key] = value
	return nil
}

// GetValue reads a value from the DHT.
func (dht *InMemoryDHT) GetValue(_ context.Context, key string, _ ...ropts.Option) ([]byte, error) {
	dht.lock.RLock()
	defer dht.lock.RUnlock()

	v, ok := dht.values[key]
	if !ok {
		return v, errors.New("not found")
	}

	return v, nil
}
