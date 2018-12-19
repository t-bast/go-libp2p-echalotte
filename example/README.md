# Onion Relay Example App

This program demonstrates how you can use onion routing in a p2p network.

## Build

```bash
go build -o echalotte-node
```

## Usage

### Create your own network

To create your own network that supports onion routing, you'll need to start
many nodes.

Start a first bootstrap node:

```bash
./echalotte-node -listen /ip4/0.0.0.0/tcp/4001
```

You will see some warnings because the DHT is unable to find other nodes to
connect to. This is expected since you're currently the only node.

You should see in the logs all the addresses the bootstrap node listens on.
You can now choose one of these addresses and wrap it inside an ipfs one by
adding the node's peer ID.

You can now start many other nodes using the first one to bootstrap:

```bash
./echalotte-node -listen /ip4/0.0.0.0/tcp/4002 -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmQqd1ydWGBbwLhrR2Kjmvx46qLf3GPrjmj55EW6N5CpNU
./echalotte-node -listen /ip4/0.0.0.0/tcp/4003 -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmQqd1ydWGBbwLhrR2Kjmvx46qLf3GPrjmj55EW6N5CpNU
./echalotte-node -listen /ip4/0.0.0.0/tcp/4004 -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmQqd1ydWGBbwLhrR2Kjmvx46qLf3GPrjmj55EW6N5CpNU
./echalotte-node -listen /ip4/0.0.0.0/tcp/4005 -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmQqd1ydWGBbwLhrR2Kjmvx46qLf3GPrjmj55EW6N5CpNU
./echalotte-node -listen /ip4/0.0.0.0/tcp/4006 -peer /ip4/127.0.0.1/tcp/4001/ipfs/QmQqd1ydWGBbwLhrR2Kjmvx46qLf3GPrjmj55EW6N5CpNU
```

## Tips

You can change the logging level from 4 (INFO) to 5 (DEBUG) if you want to see
all the debug information.
