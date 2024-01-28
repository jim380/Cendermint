# Cendermint

![CreatePlan](https://img.shields.io/badge/release-v0.1.0-red) ![CreatePlan](https://img.shields.io/badge/go-1.15%2B-blue) ![CreatePlan](https://img.shields.io/badge/license-Apache--2.0-green)

A sophisticated Prometheus Exporter designed specifically for Cosmos SDK chains. It comes with an optional lightweight dashboard, providing a user-friendly interface for real-time data visualization and monitoring. This powerful tool is an essential asset for any developer working with Cosmos SDK chains, offering a seamless integration with Prometheus and enhancing the overall development experience.

## Disclaimer

Cendermint originated as a derivative of the [Cosmos-IE](https://github.com/node-a-team/Cosmos-IE) project by [Node A-Team](https://github.com/node-a-team). We extend our profound gratitude and appreciation to the Node A-Team for their pioneering efforts and their generosity in making their remarkable work open-source.

Since its inception, Cendermint has undergone numerous substantial revisions, refactoring, and design transformations. While it has achieved a degree of stability, it is important to note that Cendermint remains in an active development phase. We encourage users to exercise discretion while utilizing this evolving software.

## Architecture

<details>

<summary>Design</summary>

![architecture](assets/design.png)

</details>

<details>

<summary>Infrastructure</summary>

![architecture](assets/arch.png)

</details>

## Supported chains

See [`chains.json`](/chains.json) .

</details>

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

#### Local Dev

##### Start Up

```bash
$ docker-compose up -d db adminer prometheus grafana
$ modd
```

##### Tear Down

```bash
$ docker compose down && docker system prune --volumes -f
```

#### Deploy

```bash
$ docker run --name cendermint -dt --restart unless-stopped -v <your_dir>:/root --net="host" --env-file ./config.env ghcr.io/jim380/cendermint:<tag> Cendermint run && docker logs cendermint -f --since 1m
```

Again, remember to create a `config.env` under `<your_dir>` and have it filled out.
