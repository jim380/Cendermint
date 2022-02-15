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
		float64(metricData.Network.BlockInterval),

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

		// validator info
		metricData.Validator.VotingPower,
		metricData.Validator.MinSelfDelegation,
		metricData.Validator.JailStatus,
		// vadalidator delegation
		metricData.Validator.Delegation.Shares,
		metricData.Validator.Delegation.Ratio,
		// vadalidator commission
		metricData.Validator.Commission.Rate,
		metricData.Validator.Commission.MaxRate,
		metricData.Validator.Commission.MaxChangeRate,
		// vadalidator signing
		metricData.Validator.Commit.PrecommitStatus,
		metricData.Validator.Proposer.Status,
		metricData.Validator.Commit.LastSigned,
		metricData.Validator.Commit.MissConsecutive,
		metricData.Validator.Commit.MissThreshold,
		metricData.Validator.Commit.MissedCount,

		// upgrade
		metricData.Upgrade.Planned,

		// ibc
		metricData.IBC.IBCChannels.Total,
		metricData.IBC.IBCChannels.Open,
		metricData.IBC.IBCConnections.Total,
		metricData.IBC.IBCConnections.Open,

		// tx
		metricData.Tx.TPS,
		metricData.Tx.GasWantedTotal,
		metricData.Tx.GasUsedTotal,
		// tx events default
		metricData.Tx.EventsTotal,
		metricData.Tx.DelegateTotal,
		metricData.Tx.MessageTotal,
		metricData.Tx.TransferTotal,
		metricData.Tx.UnbondTotal,
		metricData.Tx.WithdrawRewardsTotal,
		metricData.Tx.CreateValidatorTotal,
		metricData.Tx.RedelegateTotal,
		metricData.Tx.ProposalVote,
		// tx events ibc
		metricData.Tx.FungibleTokenPacketTotal,
		metricData.Tx.IbcTransferTotal,
		metricData.Tx.UpdateClientTotal,
		metricData.Tx.AckPacketTotal,
		metricData.Tx.SendPacketTotal,
		metricData.Tx.RecvPacketTotal,
		metricData.Tx.TimeoutTotal,
		metricData.Tx.TimeoutPacketTotal,
		metricData.Tx.DenomTraceTotal,
		// tx events swap
		metricData.Tx.SwapWithinBatchTotal,
		metricData.Tx.WithdrawWithinBatchTotal,
		metricData.Tx.DepositWithinBatchTotal,
		// tx events others
		metricData.Tx.OthersTotal,

		// graivty
		metricData.Gravity.gravityParams.SignedValsetsWindow,
		metricData.Gravity.gravityParams.SignedBatchesWindow,
		metricData.Gravity.gravityParams.TargetBatchTimeout,
		metricData.Gravity.gravityParams.SlashFractionValset,
		metricData.Gravity.gravityParams.SlashFractionBatch,
		metricData.Gravity.gravityParams.SlashFractionBadEthSig,
		metricData.Gravity.gravityParams.ValsetReward.Amount,
		metricData.Gravity.GravityActive,
		metricData.Gravity.ValSetCount,
		metricData.Gravity.ValSetActive,
		metricData.Gravity.EventNonce,
		// metricData.GravityInfo.LastClaimHeight,
		metricData.Gravity.Erc20Price,
		metricData.Gravity.BatchFees,
		metricData.Gravity.BatchesFees,
		metricData.Gravity.BridgeFees,
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

func setUpgradeLabels(metricData *metric, labels prometheus.CounterVec) {
	labels.WithLabelValues(
		metricData.Upgrade.Name,
		metricData.Upgrade.Time,
		metricData.Upgrade.Height,
		metricData.Upgrade.Info,
	).Add(0)
}
