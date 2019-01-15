package main

import (
	"context"
	crand "crypto/rand"
	"time"

	"github.com/t-bast/go-libp2p-echalotte"

	"golang.org/x/crypto/nacl/box"

	"gx/ipfs/QmNTCey11oxhb1AxDnQBRHtdhap6Ctud872NjAYPYYXPuc/go-multiaddr"
	"gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	"gx/ipfs/QmPiemjiKBC9VA7vZF82m4x1oygtg2c2YVqag8PX7dN1BD/go-libp2p-peerstore"
	"gx/ipfs/QmSQE3LqUVq8YvnmCCZHwkSDrcyQecfEWTjcpsUzH8iHtW/go-libp2p-kad-dht"
	"gx/ipfs/QmSQE3LqUVq8YvnmCCZHwkSDrcyQecfEWTjcpsUzH8iHtW/go-libp2p-kad-dht/opts"
	"gx/ipfs/QmTiRqrF5zkdZyrdsL5qndG1UbeWi8k8N2pYxCtXWrahR2/go-libp2p-routing"
	"gx/ipfs/QmaoXrM4Z41PD48JY36YqQGKQpLGjyLA2cKcLsES7YddAq/go-libp2p-host"
	logging "gx/ipfs/QmcuXC5cxs79ro2cUuHs4HQ2bkDLJUYokwL8aivcX6HW3C/go-log"
	"gx/ipfs/QmdJdFQc5U3RAKgJQGmWR7SSM7TLuER5FWz5Wq6Tzs2CnS/go-libp2p"
	"gx/ipfs/QmemYsfqwAbyvqwFiApk1GfLKhDkMm8ZQK6fCvzDbaRNyX/go-libp2p-discovery"
)

var log = logging.Logger("echalotte")

func main() {
	// The context governs the lifetime of the libp2p node.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set loggers to be verbose (INFO level).
	logging.SetAllLoggers(4)
	log.Info("Starting node...")

	config, err := ParseFlags()
	if err != nil {
		log.Error(err)
		return
	}

	peerKey, _, err := crypto.GenerateEd25519Key(crand.Reader)
	if err != nil {
		log.Error(err)
		return
	}

	// The DHT will initialize when we bootstrap the host.
	// It is used for peer discovery internally by our host.
	dhtValidator := &echalotte.PublicKeyValidator{}
	var kadDHT *dht.IpfsDHT

	options := []libp2p.Option{
		libp2p.ListenAddrs(config.ListenAddresses...),
		libp2p.Identity(peerKey),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			var err error

			// Start a DHT, for use in peer discovery.
			// We can't just make a new DHT client because we want each peer to
			// maintain its own local copy of the DHT, so that the bootstrapping node
			// of the DHT can go down without inhibitting future peer discovery.
			kadDHT, err = dht.New(ctx, h, dhtopts.Validator(dhtValidator))
			if err != nil {
				return nil, err
			}

			// Bootstrap the DHT.
			// In the default configuration, this spawns a background thread that will
			// refresh the peer table every five minutes.
			log.Info("Bootstrapping the DHT...")
			if err = kadDHT.Bootstrap(ctx); err != nil {
				return nil, err
			}

			return kadDHT, nil
		}),
	}

	host, err := libp2p.New(ctx, options...)
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("Host created. We are: %s", host.ID().Pretty())

	err = bootstrapConnections(ctx, host, config.BootstrapPeers)
	if err != nil {
		log.Error(err)
		return
	}

	// TODO: move that to a node init/warmup method in echalotte.
	// The encryption private key needs to be stored permanently (PeerStore?).
	// Most of this file's logic should be handled by a function exported by
	// the echalotte package.
	encryptionPublicKey, _, _ := box.GenerateKey(crand.Reader)
	dhtRecord, err := dhtValidator.CreateRecord(peerKey, encryptionPublicKey)
	if err != nil {
		log.Error(err)
		return
	}

	err = kadDHT.PutValue(ctx, dhtValidator.CreateKey(host.ID()), dhtRecord)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Initializing circuit builder...")
	circuitBuilder, err := echalotte.NewCircuitBuilder(
		ctx,
		discovery.NewRoutingDiscovery(kadDHT),
	)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Circuit builder ready.")

	log.Info("Generating circuit...")
	for {
		circuit, err := circuitBuilder.Build(ctx, echalotte.CircuitSize(2))
		if err != nil {
			log.Warningf("Could not generate circuit: %s", err.Error())
			log.Info("Waiting before retrying (peers are likely not ready)...")
			<-time.After(1 * time.Minute)
			log.Info("Retrying...")
			continue
		}

		log.Infof("Circuit generated: %s", circuit)
		break
	}

	log.Info("Press Ctrl+C to stop the node.")
	select {}
}

func bootstrapConnections(ctx context.Context, host host.Host, bootstrapPeers []multiaddr.Multiaddr) error {
	if len(bootstrapPeers) == 0 {
		log.Info("No bootstrap nodes configured.")
		log.Info("Use one of the following bootstrap addresses to connect other nodes to the network:")
		for _, addr := range host.Addrs() {
			log.Info(addr.String())
		}
	} else {
		for _, peerAddr := range bootstrapPeers {
			peerInfo, err := peerstore.InfoFromP2pAddr(peerAddr)
			if err != nil {
				return err
			}

			err = host.Connect(ctx, *peerInfo)
			if err != nil {
				return err
			}

			log.Infof("Connection established with bootstrap node: %s", peerInfo.ID.Pretty())
		}
	}

	log.Info("Waiting for enough connected peers...")
	for {
		log.Infof("%d peer(s) connected...", len(host.Network().Conns()))

		if len(host.Network().Conns()) >= 6 {
			log.Info("Network bootstrapped.")
			return nil
		}

		<-time.After(30 * time.Second)
	}
}
