package rest

import (
	"github.com/jim380/Cendermint/config"
)

/***********************
 * SDK Routes
************************/
func GetBlockInfoRoute(cfg config.Config) string {
	if cfg.IsLegacySDKVersion() {
		return "/blocks/latest"
	} else {
		return "/cosmos/base/tendermint/v1beta1/blocks/latest"
	}
}

func GetBlockByHeightRoute(cfg config.Config) string {
	if cfg.IsLegacySDKVersion() {
		return "/blocks/"
	} else {
		return "/cosmos/base/tendermint/v1beta1/blocks/"
	}
}

func GetValidatorSetByHeightRoute(cfg config.Config) string {
	if cfg.IsLegacySDKVersion() {
		return "/validatorsets/"
	} else {
		return "/cosmos/base/tendermint/v1beta1/validatorsets/"
	}
}

func GetValidatorDistributionByAddressRoute(cfg config.Config) string {
	if cfg.IsLegacySDKVersion() {
		return "/distribution/validators/"
	} else {
		return "/cosmos/distribution/v1beta1/validators/"
	}
}

func GetInflationRoute(cfg config.Config) string {
	if cfg.Chain.Name == "irisnet" {
		return "/irishub/mint/params"
	} else if cfg.IsLegacySDKVersion() {
		return "/minting/inflation"

	} else {
		return "/cosmos/mint/v1beta1/inflation"
	}
}

func GetBalancesByAddressRoute(cfg config.Config) string {
	return "/cosmos/bank/v1beta1/balances/"
}

func GetStakingPoolRoute(cfg config.Config) string {
	return "/cosmos/staking/v1beta1/pool"
}

func GetSupplyRoute(cfg config.Config) string {
	return "/cosmos/bank/v1beta1/supply/"
}

func GetValidatorByAddressRoute(cfg config.Config) string {
	return "/cosmos/staking/v1beta1/validators/"
}

func GetValidatorsRoute() string {
	return "/cosmos/staking/v1beta1/validators"
}

func GetTxByHeightRoute(cfg config.Config) string {
	return "/cosmos/tx/v1beta1/txs?events=tx.height="
}

func GetSlashingParamsRoute(cfg config.Config) string {
	return "/cosmos/slashing/v1beta1/params"
}

func GetSigningInfoByAddressRoute(cfg config.Config) string {
	return "/cosmos/slashing/v1beta1/signing_infos/"
}

func GetProposalsRoute(cfg config.Config) string {
	return "/cosmos/gov/v1beta1/proposals"
}

func GetNodeInfoRoute() string {
	return "/cosmos/base/tendermint/v1beta1/node_info"
}

func GetUpgradeCurrentPlanRoute(cfg config.Config) string {
	return "/cosmos/upgrade/v1beta1/current_plan"
}

/***********************
 * IBC Routes
************************/
func GetIBCChannelsRoute(cfg config.Config) string {
	return "/ibc/core/channel/v1/channels"
}

func GetIBCConnectionsRoute(cfg config.Config) string {
	return "/ibc/core/connection/v1/connections"
}

/***********************
 * Gravity Bridge Routes
************************/
func GetBatchFeesRoute() string {
	return "/gravity/v1beta/batchfees"
}

func GetBatchesFeesRoute() string {
	return "/gravity/v1beta1/batch/outgoingtx"
}

func GetBridgeFeesRoute() string {
	return "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"
}

func GetBridgeParamsRoute() string {
	return "/gravity/v1beta/params"
}

func GetOracleEventNonceByAddressRoute() string {
	return "/gravity/v1beta/oracle/eventnonce/"
}

func GetCurrentValidatorSetRoute() string {
	return "/gravity/v1beta/valset/current"
}

/***********************
 * Akash Routes
************************/
func GetDeploymentsRoute() string {
	return "/akash/deployment/v1beta3/deployments/list"
}

func GetProvidersRoute() string {
	return "/akash/provider/v1beta3/providers"
}

func GetAuditorForProviderOwnerRoute(owner string) string {
	return "/akash/audit/v1beta3/audit/attributes/" + owner + "/list"
}

/***********************
 * Oracle Routes
************************/
func GetMissedCounterRoute() string {
	return "/refractedlabs/oracle/oracle/miss_counter"
}

func GetPrevoteRoute() string {
	return "/refractedlabs/oracle/oracle/oracle_pre_vote"
}

func GetVoteRoute() string {
	return "/refractedlabs/oracle/oracle/oracle_vote"
}
