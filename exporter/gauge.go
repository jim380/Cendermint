package exporter

import (
	utils "github.com/jim380/Cendermint/utils"
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

func setNormalGauges(metricData *metric) {
	// set values for normal guages
	gaugesValue := [...]float64{
		float64(metricData.Network.BlockHeight),

		metricData.Network.Staking.NotBondedTokens,
		metricData.Network.Staking.BondedTokens,
		metricData.Network.Staking.TotalSupply,
		metricData.Network.Staking.BondedRatio,

		metricData.Network.Gov.TotalProposalCount,
		metricData.Network.Gov.VotingProposalCount,

		metricData.Validator.VotingPower,
		metricData.Validator.MinSelfDelegation,
		metricData.Validator.JailStatus,

		metricData.Validator.Proposer.Ranking,
		metricData.Validator.Proposer.Status,

		metricData.Validator.Delegation.Shares,
		metricData.Validator.Delegation.Ratio,
		metricData.Validator.Delegation.DelegatorCount,
		metricData.Validator.Delegation.Self,

		metricData.Validator.Commission.Rate,
		metricData.Validator.Commission.MaxRate,
		metricData.Validator.Commission.MaxChangeRate,
		// metricData.Validator.Commit.VoteType,
		metricData.Validator.Commit.PrecommitStatus,

		metricData.Network.Minting.Inflation,
		metricData.Network.Minting.ActualInflation,
	}
	for i := 0; i < len(gaugesNamespaceList); i++ {
		defaultGauges[i].Set(gaugesValue[i])
	}
}
