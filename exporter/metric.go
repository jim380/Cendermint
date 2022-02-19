package exporter

import (
	"time"

	"go.uber.org/zap"

	"github.com/jim380/Cendermint/rest"
	utils "github.com/jim380/Cendermint/utils"
)

func SetMetric(currentBlock int64, restData *rest.RESTData, log *zap.Logger) {
	operAddr := rest.OperAddr
	consPubKey := restData.Validators.ConsPubKey
	consAddr := restData.Validatorsets[consPubKey.Key][0]

	// chain
	metricData.Network.BlockHeight = currentBlock
	metricData.Network.BlockInterval = restData.BlockInterval

	// minting
	metricData.Network.Minting.Inflation = restData.Inflation
	metricData.Network.Minting.ActualInflation = metricData.Network.Minting.Inflation / metricData.Network.Staking.BondedRatio

	// staking
	metricData.Network.Staking.NotBondedTokens = utils.StringToFloat64(restData.StakingPool.Pool.Not_bonded_tokens)
	metricData.Network.Staking.BondedTokens = utils.StringToFloat64(restData.StakingPool.Pool.Bonded_tokens)
	metricData.Network.Staking.TotalSupply = restData.StakingPool.Pool.Total_supply
	metricData.Network.Staking.BondedRatio = metricData.Network.Staking.BondedTokens / metricData.Network.Staking.TotalSupply

	// slashing
	metricData.Network.Slashing.SignedBlocksWindow = utils.StringToFloat64(restData.Slashing.Params.SignedBlocksWindow)
	metricData.Network.Slashing.MinSignedPerWindow = utils.StringToFloat64(restData.Slashing.Params.MinSignedPerWindow)
	jailedDurationTime, _ := time.ParseDuration(restData.Slashing.Params.DowntimeJailDuration)
	metricData.Network.Slashing.DowntimeJailDuration = float64(jailedDurationTime.Seconds())
	metricData.Network.Slashing.SlashFractionDoubleSign = utils.StringToFloat64(restData.Slashing.Params.SlashFractionDoubleSign)
	metricData.Network.Slashing.SlashFractionDowntime = utils.StringToFloat64(restData.Slashing.Params.SlashFractionDowntime)
	metricData.Network.Slashing.StartHeight = utils.StringToFloat64(restData.Slashing.ValSigning.StartHeight)
	metricData.Network.Slashing.IndexOffset = utils.StringToFloat64(restData.Slashing.ValSigning.IndexOffset)
	jailedUntilTime, _ := time.Parse("2006-01-02T15:04:05Z", restData.Slashing.ValSigning.JailedUntil)
	metricData.Network.Slashing.JailedUntil = float64(jailedUntilTime.Unix())
	metricData.Network.Slashing.Tombstoned = utils.BoolToFloat64(restData.Slashing.ValSigning.Tombstoned)
	metricData.Network.Slashing.MissedBlocksCounter = utils.StringToFloat64(restData.Slashing.ValSigning.MissedBlocksCounter)

	// gov
	metricData.Network.Gov.TotalProposalCount = restData.Gov.TotalProposalCount
	metricData.Network.Gov.VotingProposalCount = restData.Gov.VotingProposalCount

	// validator info
	metricData.Validator.VotingPower = utils.StringToFloat64(restData.Validatorsets[consPubKey.Key][1])
	metricData.Validator.JailStatus = utils.BoolToFloat64(restData.Validators.Jailed)
	metricData.Validator.MinSelfDelegation = utils.StringToFloat64(restData.Validators.MinSelfDelegation)
	// validator delegation
	metricData.Validator.Delegation.Shares = utils.StringToFloat64(restData.Validators.DelegatorShares)
	metricData.Validator.Delegation.Ratio = metricData.Validator.Delegation.Shares / metricData.Network.Staking.BondedTokens
	// metricData.Validator.Delegation.DelegatorCount = restData.Delegations.DelegationCount
	// metricData.Validator.Delegation.Self = restData.Delegations.SelfDelegation
	// validator commission
	metricData.Validator.Commission.Rate = utils.StringToFloat64(restData.Validators.Commission.Commission.Rate)
	metricData.Validator.Commission.MaxRate = utils.StringToFloat64(restData.Validators.Commission.Commission.Max_rate)
	metricData.Validator.Commission.MaxChangeRate = utils.StringToFloat64(restData.Validators.Commission.Commission.Max_change_rate)
	// validator signing
	metricData.Validator.Commit.PrecommitStatus = restData.Commit.ValidatorPrecommitStatus
	metricData.Validator.Proposer.Status = restData.Commit.ValidatorProposingStatus
	metricData.Validator.Commit.LastSigned = float64(restData.Commit.LastSigned)
	metricData.Validator.Commit.MissConsecutive = restData.Commit.MissConsecutive
	metricData.Validator.Commit.MissThreshold = restData.Commit.MissThreshold
	metricData.Validator.Commit.MissedCount = float64(restData.Commit.MissedCount)

	// upgrade
	metricData.Upgrade.Planned = utils.BoolToFloat64(restData.UpgradeInfo.Planned)

	// ibc
	metricData.IBC.IBCChannels.Total = float64(len(restData.IBC.IBCChannels))
	metricData.IBC.IBCChannels.Open = float64(restData.IBC.IBCInfo.OpenChannels)
	metricData.IBC.IBCConnections.Total = float64(len(restData.IBC.IBCConnections))
	metricData.IBC.IBCConnections.Open = float64(restData.IBC.IBCInfo.OpenConnections)

	// tx
	metricData.Tx.TPS = restData.TxInfo.TPS
	metricData.Tx.GasWantedTotal = restData.TxInfo.Result.GasWantedTotal
	metricData.Tx.GasUsedTotal = restData.TxInfo.Result.GasUsedTotal
	// tx events default
	metricData.Tx.EventsTotal = restData.TxInfo.Result.Default.EventsTotal
	metricData.Tx.DelegateTotal = restData.TxInfo.Result.Default.DelegateTotal
	metricData.Tx.MessageTotal = restData.TxInfo.Result.Default.MessageTotal
	metricData.Tx.TransferTotal = restData.TxInfo.Result.Default.TransferTotal
	metricData.Tx.UnbondTotal = restData.TxInfo.Result.Default.UnbondTotal
	metricData.Tx.WithdrawRewardsTotal = restData.TxInfo.Result.Default.WithdrawRewardsTotal
	metricData.Tx.CreateValidatorTotal = restData.TxInfo.Result.Default.CreateValidatorTotal
	metricData.Tx.RedelegateTotal = restData.TxInfo.Result.Default.RedelegateTotal
	metricData.Tx.ProposalVote = restData.TxInfo.Result.Default.ProposalVote
	// tx events ibc
	metricData.Tx.FungibleTokenPacketTotal = restData.TxInfo.Result.IBC.FungibleTokenPacketTotal
	metricData.Tx.IbcTransferTotal = restData.TxInfo.Result.IBC.IbcTransferTotal
	metricData.Tx.UpdateClientTotal = restData.TxInfo.Result.IBC.UpdateClientTotal
	metricData.Tx.AckPacketTotal = restData.TxInfo.Result.IBC.AckPacketTotal
	metricData.Tx.SendPacketTotal = restData.TxInfo.Result.IBC.SendPacketTotal
	metricData.Tx.RecvPacketTotal = restData.TxInfo.Result.IBC.RecvPacketTotal
	metricData.Tx.TimeoutTotal = restData.TxInfo.Result.IBC.TimeoutTotal
	metricData.Tx.TimeoutPacketTotal = restData.TxInfo.Result.IBC.TimeoutPacketTotal
	metricData.Tx.DenomTraceTotal = restData.TxInfo.Result.IBC.DenomTraceTotal
	// tx events swap
	metricData.Tx.SwapWithinBatchTotal = restData.TxInfo.Result.Swap.SwapWithinBatchTotal
	metricData.Tx.WithdrawWithinBatchTotal = restData.TxInfo.Result.Swap.WithdrawWithinBatchTotal
	metricData.Tx.DepositWithinBatchTotal = restData.TxInfo.Result.Swap.DepositWithinBatchTotal
	// tx events others
	metricData.Tx.OthersTotal = restData.TxInfo.Result.OthersTotal

	// gravity
	metricData.Gravity.gravityParams.SignedValsetsWindow = utils.StringToFloat64(restData.GravityInfo.SignedValsetsWindow)
	metricData.Gravity.gravityParams.SignedBatchesWindow = utils.StringToFloat64(restData.GravityInfo.SignedBatchesWindow)
	metricData.Gravity.gravityParams.TargetBatchTimeout = utils.StringToFloat64(restData.GravityInfo.TargetBatchTimeout)
	metricData.Gravity.gravityParams.SlashFractionValset = utils.StringToFloat64(restData.GravityInfo.SlashFractionValset)
	metricData.Gravity.gravityParams.SlashFractionBatch = utils.StringToFloat64(restData.GravityInfo.SlashFractionBatch)
	metricData.Gravity.gravityParams.SlashFractionBadEthSig = utils.StringToFloat64(restData.GravityInfo.SlashFractionBadEthSig)
	metricData.Gravity.gravityParams.ValsetReward.Amount = utils.StringToFloat64(restData.GravityInfo.ValsetReward.Amount)
	metricData.Gravity.GravityActive = restData.GravityInfo.GravityActive
	metricData.Gravity.ValSetCount = float64(restData.GravityInfo.ValSetCount)
	metricData.Gravity.ValSetActive = restData.GravityInfo.ValActive
	metricData.Gravity.EventNonce = utils.StringToFloat64(restData.GravityInfo.EventNonce)
	metricData.Gravity.Erc20Price = restData.GravityInfo.UMEEPrice
	metricData.Gravity.BatchFees = restData.GravityInfo.BatchFees
	metricData.Gravity.BatchesFees = restData.GravityInfo.BatchesFees
	metricData.Gravity.BridgeFees = restData.GravityInfo.BridgeFees

	// labels node
	metricData.Network.ChainID = restData.Commit.ChainId
	metricData.Validator.Moniker = restData.Validators.Description.Moniker
	metricData.Network.NodeInfo.Moniker = restData.NodeInfo.Default.Moniker
	metricData.Network.NodeInfo.NodeID = restData.NodeInfo.Default.NodeID
	metricData.Network.NodeInfo.TMVersion = restData.NodeInfo.Default.TMVersion
	metricData.Network.NodeInfo.AppName = restData.NodeInfo.Application.AppName
	metricData.Network.NodeInfo.Name = restData.NodeInfo.Application.Name
	metricData.Network.NodeInfo.Version = restData.NodeInfo.Application.Version
	metricData.Network.NodeInfo.GitCommit = restData.NodeInfo.Application.GitCommit
	metricData.Network.NodeInfo.GoVersion = restData.NodeInfo.Application.GoVersion
	metricData.Network.NodeInfo.SDKVersion = restData.NodeInfo.Application.SDKVersion
	// labels addr
	metricData.Validator.Address.Operator = operAddr
	metricData.Validator.Address.Account = utils.GetAccAddrFromOperAddr(operAddr)
	metricData.Validator.Address.ConsensusHex = utils.Bech32AddrToHexAddr(consAddr)
	// labels upgrade
	metricData.Upgrade.Name = restData.UpgradeInfo.Plan.Name
	metricData.Upgrade.Time = restData.UpgradeInfo.Plan.Time
	metricData.Upgrade.Height = restData.UpgradeInfo.Plan.Height
	metricData.Upgrade.Info = restData.UpgradeInfo.Plan.Info

	// denom gauges
	metricData.Validator.Account.Balances = restData.Balances
	metricData.Validator.Account.Commission = restData.Commission
	metricData.Validator.Account.Rewards = restData.Rewards
}

func GetMetric() *metric {
	return &metricData
}
