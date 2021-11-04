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

## Available metrics
| Name | Type | Tags | Description |
|------|------|------|-------------|
| `cendermint_chain_blockHeight` | Gauge | - | Current blockchain height|
| `cendermint_staking_bonded_ratio` | Gauge | - | Bonded stake ratio in the network |
| `cendermint_staking_bonded_tokens` | Gauge | - | Bonded stake amount in the network |
| `cendermint_staking_not_bonded_tokens` | Gauge | - | Unbonded stake amount in the network |
| `cendermint_staking_total_supply` | Gauge | - | Total token supply in the network |
| `cendermint_slashing_downtime_jail_duration` | Gauge | - | Downtime duration before getting jailed|
| `cendermint_slashing_min_signed_per_window` | Gauge | - | Minimum number of blocks that need to be signed per signing window before getting jailed |
| `cendermint_slashing_signed_blocks_window` | Gauge | - | Number of blocks in a signing window |
| `cendermint_slashing_slash_fraction_double_sign` | Gauge | - | % of stake to be slashed in the event of a double sign |
| `cendermint_slashing_slash_fraction_downtime` | Gauge | - | % of stake to be slashed in the event of downtime  |
| `cendermint_slashing_start_Height` | Gauge | - | The first block the validator signed on the current chain |
| `cendermint_slashing_index_offset` | Gauge | - | The index used to check if the validator has crossed below the liveness threshold over a sliding window |
| `cendermint_slashing_jailed_until` | Gauge | - | Most recent `jailed_until` date/time of the validator recorded on chain |
| `cendermint_slashing_tombstoned` | Gauge | - | Whether the validator is tombstoned (i.e. double sign) [0] False - [1] True |
| `cendermint_slashing_missed_blocks_counter` | Gauge | - | Total number of blocks the validator missed since last unjail |
| `cendermint_minting_actual_inflation` | Gauge | - | Actual inflation in the network |
| `cendermint_minting_inflation` | Gauge | - | Default inflation in the network |
| `cendermint_gov_total_proposal_count` | Gauge | - | Total number of proposals ever submitted in the network |
| `cendermint_gov_voting_proposal_count` | Gauge | - | Number of proposals currently in voting |
| `cendermint_validator_voting_power` | Gauge | - | Voting power of the validator |
| `cendermint_validator_min_self_delegation` | Gauge | - | Minimum self delegation amount of the validator |
| `cendermint_validator_jail_status` | Gauge | - | Jail status of the validator<br>[0] Active - [1] Jailed |
| `cendermint_validator_delegation_shares` | Gauge | - | Total number of delegated tokens of the validator |
| `cendermint_validator_delegation_ratio` | Gauge | - | Ratio of the validator's bonded stake to the network's total bonded stake |
| `cendermint_validator_commission_rate` | Gauge | - | Commission rate of the validator |
| `cendermint_validator_commission_max_rate` | Gauge | - | Maximum commission rate of the validator |
| `cendermint_validator_commission_max_change_rate` | Gauge | - | Maximum change rate of the validator's commission |
| `cendermint_validator_balances_uatom` | Gauge | - | Available balance of the validator |
| `cendermint_validator_commission_uatom` | Gauge | - | Available commission of the validator |
| `cendermint_validator_rewards_uatom` | Gauge | - | Available self-delegation rewards of the validator |
| `cendermint_validator_precommit_status` | Gauge | - | Precommit status of the validator<br>[0] Missed - [1] Signed |
| `cendermint_validator_proposer_status` | Gauge | - | Proposer status of the validator<br>[0] Not the proposer - [1] Proposer |
| `cendermint_validator_last_signed_height` | Gauge | - | The last height the validator signed |
| `cendermint_validator_miss_count` | Gauge | - | Number of blocks missed since the validator last signed |
| `cendermint_validator_miss_consecutive` | Gauge | - | The validator has missed two blocks in a row |
| `cendermint_validator_miss_threshold` | Gauge | - | The validator has missed `>= threshold` block since s/he last signed |
| `cendermint_ibc_channels_total` | Gauge | - | Total number of ibc channels in the network |
| `cendermint_ibc_channels_open` | Gauge | - | Total number of open ibc channels in the network |