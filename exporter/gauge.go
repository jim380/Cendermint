package exporter

import (
	utils "github.com/jim380/Cendermint/utils"
	"github.com/prometheus/client_golang/prometheus"
)

func setDenomGauges(metricData *metric, denomList []string) {
	count := 0
	for i := 0; i < len(denomList)*3; i += 3 {
		for _, value := range metricData.Validator.Account.Balances {
			if value.Denom == denomList[count] {
				gaugesDenom[i].Set(utils.StringToFloat64(value.Amount))
			}
		}
		for _, value := range metricData.Validator.Account.Commission {
			if value.Denom == denomList[count] {
				gaugesDenom[i+1].Set(utils.StringToFloat64(value.Amount))
			}
		}
		for _, value := range metricData.Validator.Account.Rewards {
			if value.Denom == denomList[count] {
				gaugesDenom[i+2].Set(utils.StringToFloat64(value.Amount))
			}
		}
		count++
	}
}

func setNormalGauges(metricData *metric, defaultGauges []prometheus.Gauge) {
	// set values for normal guages
	gaugesValue := [...]float64{
		// IMPORTANT!!! order needs to match with gaugesNamespaceList

		// chain
		float64(metricData.Network.BlockHeight),

		// minting
		metricData.Network.Minting.Inflation,
		metricData.Network.Minting.ActualInflation,

		// staking
		metricData.Network.Staking.NotBondedTokens,
		metricData.Network.Staking.BondedTokens,
		metricData.Network.Staking.TotalSupply,
		metricData.Network.Staking.BondedRatio,

		// slashing
		metricData.Network.Slashing.SignedBlocksWindow,
		metricData.Network.Slashing.MinSignedPerWindow,
		metricData.Network.Slashing.DowntimeJailDuration,
		metricData.Network.Slashing.SlashFractionDoubleSign,
		metricData.Network.Slashing.SlashFractionDowntime,
		metricData.Network.Slashing.StartHeight,
		metricData.Network.Slashing.IndexOffset,
		metricData.Network.Slashing.JailedUntil,
		metricData.Network.Slashing.Tombstoned,
		metricData.Network.Slashing.MissedBlocksCounter,

		// gov
		metricData.Network.Gov.TotalProposalCount,
		metricData.Network.Gov.VotingProposalCount,

		// validator
		metricData.Validator.VotingPower,
		metricData.Validator.MinSelfDelegation,
		metricData.Validator.JailStatus,
		// vadalidator_delegation
		metricData.Validator.Delegation.Shares,
		metricData.Validator.Delegation.Ratio,
		// vadalidator_commission
		metricData.Validator.Commission.Rate,
		metricData.Validator.Commission.MaxRate,
		metricData.Validator.Commission.MaxChangeRate,
		// vadalidator_signing
		metricData.Validator.Commit.PrecommitStatus,
		metricData.Validator.Proposer.Status,
		metricData.Validator.Commit.LastSigned,
		metricData.Validator.Commit.MissConsecutive,
		metricData.Validator.Commit.MissThreshold,
		metricData.Validator.Commit.MissedCount,

		// ibc
		metricData.IBC.IBCChannels.Total,
		metricData.IBC.IBCChannels.Open,
		metricData.IBC.IBCConnections.Total,
		metricData.IBC.IBCConnections.Open,
	}
	for i := 0; i < len(gaugesNamespaceList); i++ {
		defaultGauges[i].Set(gaugesValue[i])
	}
}

// {"chain_id", "node_moniker", "node_id", "tm_version", "app_name", "binary_name", "app_version", "git_commit", "go_version", "sdk_version"}
func setNodeLabels(metricData *metric, labels prometheus.CounterVec) {
	labels.WithLabelValues(
		metricData.Network.ChainID,
		metricData.Network.NodeInfo.Moniker,
		metricData.Network.NodeInfo.NodeID,
		metricData.Network.NodeInfo.TMVersion,
		metricData.Network.NodeInfo.AppName,
		metricData.Network.NodeInfo.Name,
		metricData.Network.NodeInfo.Version,
		metricData.Network.NodeInfo.GitCommit,
		metricData.Network.NodeInfo.GoVersion,
		metricData.Network.NodeInfo.SDKVersion,
	).Add(0)
}

func setAddrLabels(metricData *metric, labels prometheus.CounterVec) {
	labels.WithLabelValues(
		metricData.Validator.Address.Operator,
		metricData.Validator.Address.Account,
		metricData.Validator.Address.ConsensusHex,
	).Add(0)
}
