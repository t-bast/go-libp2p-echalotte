package main

import (
	"context"

	"gx/ipfs/QmNTCey11oxhb1AxDnQBRHtdhap6Ctud872NjAYPYYXPuc/go-multiaddr"
	"gx/ipfs/QmPiemjiKBC9VA7vZF82m4x1oygtg2c2YVqag8PX7dN1BD/go-libp2p-peerstore"
	"gx/ipfs/QmSQE3LqUVq8YvnmCCZHwkSDrcyQecfEWTjcpsUzH8iHtW/go-libp2p-kad-dht"
	"gx/ipfs/QmTiRqrF5zkdZyrdsL5qndG1UbeWi8k8N2pYxCtXWrahR2/go-libp2p-routing"
	"gx/ipfs/QmaoXrM4Z41PD48JY36YqQGKQpLGjyLA2cKcLsES7YddAq/go-libp2p-host"
	logging "gx/ipfs/QmcuXC5cxs79ro2cUuHs4HQ2bkDLJUYokwL8aivcX6HW3C/go-log"
	"gx/ipfs/QmdJdFQc5U3RAKgJQGmWR7SSM7TLuER5FWz5Wq6Tzs2CnS/go-libp2p"
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

	options := []libp2p.Option{
		libp2p.ListenAddrs(config.ListenAddresses...),
	}

	host, err := libp2p.New(ctx, options...)
	if err != nil {
		panic(err)
	}

	log.Infof("Host created. We are: %s", host.ID().Pretty())

	_, err = startDHT(ctx, host)
	if err != nil {
		log.Error(err)
		return
	}

	err = bootstrapConnections(ctx, host, config.BootstrapPeers)
	if err != nil {
		log.Error(err)
		return
	}

	select {}
}

func startDHT(ctx context.Context, host host.Host) (routing.ContentRouting, error) {
	// Start a DHT, for use in peer discovery.
	// We can't just make a new DHT client because we want each peer to
	// maintain its own local copy of the DHT, so that the bootstrapping node
	// of the DHT can go down without inhibitting future peer discovery.
	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		return nil, err
	}

	// Bootstrap the DHT.
	// In the default configuration, this spawns a background thread that will
	// refresh the peer table every five minutes.
	log.Info("Bootstrapping the DHT...")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}

	log.Info("DHT bootstrap complete.")
	return kademliaDHT, nil
}

func bootstrapConnections(ctx context.Context, host host.Host, bootstrapPeers []multiaddr.Multiaddr) error {
	if len(bootstrapPeers) == 0 {
		log.Info("No bootstrap nodes configured.")
		log.Info("Use one of the following bootstrap addresses to connect other nodes to the network:")
		for _, addr := range host.Addrs() {
			log.Info(addr.String())
		}

		return nil
	}

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

	return nil
}
