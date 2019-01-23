# Onion Relay Example App

This program demonstrates how you can use onion routing in a p2p network.

## Build

```bash
go build -o echalotte-node
```

## Usage

### Create your own network

To create your own network that supports onion routing, you'll need to start
many nodes (at least 7 for this example).

Start a first bootstrap node:

```bash
./echalotte-node -listen /ip4/0.0.0.0/tcp/4001
```

You should see in the logs all the addresses the bootstrap node listens on.
You can now choose one of these addresses and wrap it inside an ipfs one by
adding the node's peer ID.

You can now start many other nodes using the first one to bootstrap (for the
demo to work you should initialize at least 7 nodes):

```bash
./echalotte-node -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmW5fMEsusmL8H598reQgSCPmvv4UZ1Q5JeVArVvXSBqUh -listen /ip4/0.0.0.0/tcp/4002
./echalotte-node -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmW5fMEsusmL8H598reQgSCPmvv4UZ1Q5JeVArVvXSBqUh -listen /ip4/0.0.0.0/tcp/4003
./echalotte-node -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmW5fMEsusmL8H598reQgSCPmvv4UZ1Q5JeVArVvXSBqUh -listen /ip4/0.0.0.0/tcp/4004
./echalotte-node -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmW5fMEsusmL8H598reQgSCPmvv4UZ1Q5JeVArVvXSBqUh -listen /ip4/0.0.0.0/tcp/4005
./echalotte-node -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmW5fMEsusmL8H598reQgSCPmvv4UZ1Q5JeVArVvXSBqUh -listen /ip4/0.0.0.0/tcp/4006
./echalotte-node -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmW5fMEsusmL8H598reQgSCPmvv4UZ1Q5JeVArVvXSBqUh -listen /ip4/0.0.0.0/tcp/4007
```

Once enough nodes have been added and are connected to each other, the onion
routing layer will initialize and generate a sample circuit that will be
printed to the console.

The bootstrapping process can take a few minutes (to setup enough nodes in the
network and bootstrap the DHT content), so be patient!

Once the bootstrapping process completes, you will get a prompt that lets you
send messages between nodes.

## Tips

You can change the logging level from 4 (INFO) to 5 (DEBUG) if you want to see
all the debug information.
