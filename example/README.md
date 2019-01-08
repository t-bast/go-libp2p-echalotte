# Onion Relay Example App

This program demonstrates how you can use onion routing in a p2p network.

## Build

```bash
go build -o echalotte-node
```

## Usage

### Create your own network

To create your own network that supports onion routing, you'll need to start
many nodes (at least 20 for this example).

Start a first bootstrap node:

```bash
./echalotte-node -listen /ip4/0.0.0.0/tcp/4001
```

You should see in the logs all the addresses the bootstrap node listens on.
You can now choose one of these addresses and wrap it inside an ipfs one by
adding the node's peer ID.

You can now start many other nodes using the first one to bootstrap:

```bash
./echalotte-node -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmdJQmo3gYqTXbGL3R1j8Zgw4svP6zLkkFqP6fBXR9Gnt1 -listen /ip4/0.0.0.0/tcp/4002
./echalotte-node -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmdJQmo3gYqTXbGL3R1j8Zgw4svP6zLkkFqP6fBXR9Gnt1 -listen /ip4/0.0.0.0/tcp/4003
./echalotte-node -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmdJQmo3gYqTXbGL3R1j8Zgw4svP6zLkkFqP6fBXR9Gnt1 -listen /ip4/0.0.0.0/tcp/4004
./echalotte-node -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmdJQmo3gYqTXbGL3R1j8Zgw4svP6zLkkFqP6fBXR9Gnt1 -listen /ip4/0.0.0.0/tcp/4005
```

Once enough nodes have been added and are connected to each other, the onion
routing layer will initialize and generate a sample circuit that will be
printed to the console.

## Tips

You can change the logging level from 4 (INFO) to 5 (DEBUG) if you want to see
all the debug information.
