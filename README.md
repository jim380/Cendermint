# Cendermint
![CreatePlan](https://img.shields.io/badge/release-v0.1.0-red) ![CreatePlan](https://img.shields.io/badge/go-1.15%2B-blue) ![CreatePlan](https://img.shields.io/badge/license-Apache--2.0-green)  
Prometheus Exporter for Tendermint based blockchains.

## Disclaimer
This project started out as a fork of [Cosmos-IE](https://github.com/node-a-team/Cosmos-IE) by [Node A-Team](https://github.com/node-a-team). I'd like to express my greatest gratitude and appreciation to them for initiating and open-sourcing their awesome work. Since the fork it has undergone a few significant rewrites, refactors and design changes. Though semi-stable, Cendermint is still very much a work-in-progress so please proceed with caution.

## Architecture
![](assets/cendermint.png)

## Supported chains
- Cosmos(`cosmoshub-4`)
- NYM (`testnet-milhon`)
- Umme (`umeevengers-1c`)

## Get Up and Running

### Build from Source
```bash
$ cd $GOPATH/src/githb.com
$ git clone https://github.com/jim380/Cendermint.git
$ cd $HOME/Cendermint
$ go build
# Important!!! Remember to fill out config.env
$ ./Cendermint run
```

### Docker
```bash
$ docker run --name cendermint -dt --restart on-failure -v <your_dir>:/root --net="host" --env-file ./config.env ghcr.io/jim380/cendermint:master /bin/sh -c "Cendermint run"
```
Again, remember to create a `config.env` under `<your_dir>` and have it filled out.