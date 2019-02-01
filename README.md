# go-libp2p-echalotte

[![IPFS](https://img.shields.io/badge/project-IPFS-blue.svg?style=flat-square)](http://libp2p.io/)
[![GoDoc](https://godoc.org/github.com/t-bast/go-libp2p-echalotte?status.svg)](https://godoc.org/github.com/t-bast/go-libp2p-echalotte)
[![Go Report Card](https://goreportcard.com/badge/github.com/t-bast/go-libp2p-echalotte)](https://goreportcard.com/report/github.com/t-bast/go-libp2p-echalotte)
[![codecov](https://codecov.io/gh/t-bast/go-libp2p-echalotte/branch/master/graph/badge.svg)](https://codecov.io/gh/t-bast/go-libp2p-echalotte)
[![Travis CI](https://travis-ci.org/t-bast/go-libp2p-echalotte.svg?branch=master)](https://travis-ci.org/t-bast/go-libp2p-echalotte)

An onion routing implementation for libp2p.

The goal of this package is to provide onion routing through a libp2p network.
It might not be useful at all (Tor is only secure as a whole), but we'll see
once the project matures.

## Status

The current code has reached a phase where it's a working prototype.
Here are the areas that require work to productize it:

- [ ] Cleaner separation of all components (separate packages/repos)
- [ ] Improve encryption key management (deprecation, update, etc)
- [ ] Add performance benchmarks
- [ ] Investigate circuit caching
- [ ] Introduce mocks to get better test coverage (in particular stream testing)
- [ ] Investigate known attacks on plain onion routing (outside of Tor)
- [ ] Investigate MorphMix, Tarzan and other p2p onion-routing experiments (as opposed to Tor's client/server model)
- [ ] Long-lived heterogeneous test network
- [ ] Try attacking via network analysis
- [ ] Investigate implementation at the transport level
