package rest

import (
	"github.com/jim380/Cendermint/config"
)

/***********************
 * SDK Routes
************************/
func getBlockInfoRoute(cfg config.Config) string {
	if cfg.IsLegacySDKVersion() {
		return "/blocks/latest"
	} else {
		return "/cosmos/base/tendermint/v1beta1/blocks/latest"
	}
}

func getBlockByHeightRoute(cfg config.Config) string {
	if cfg.IsLegacySDKVersion() {
		return "/blocks/"
	} else {
		return "/cosmos/base/tendermint/v1beta1/blocks/"
	}
}

func getValidatorSetByHeightRoute(cfg config.Config) string {
	if cfg.IsLegacySDKVersion() {
		return "/validatorsets/"
	} else {
		return "/cosmos/base/tendermint/v1beta1/validatorsets/"
	}
}

func getValidatorDistributionByAddressRoute(cfg config.Config) string {
	if cfg.IsLegacySDKVersion() {
		return "/distribution/validators/"
	} else {
		return "/cosmos/distribution/v1beta1/validators/"
	}
}

func getInflationRoute(cfg config.Config) string {
	if cfg.Chain == "iris" {
		return "/irishub/mint/params"
	} else if cfg.IsLegacySDKVersion() {
		return "/minting/inflation"

	} else {
		return "/cosmos/mint/v1beta1/inflation"
	}
}

func getBalancesByAddressRoute(cfg config.Config) string {
	return "/cosmos/bank/v1beta1/balances/"
}

func getStakingPoolRoute(cfg config.Config) string {
	return "/cosmos/staking/v1beta1/pool"
}

func getSupplyRoute(cfg config.Config) string {
	return "/cosmos/bank/v1beta1/supply/"
}

func getValidatorByAddressRoute(cfg config.Config) string {
	return "/cosmos/staking/v1beta1/validators/"
}

func getValidatorsRoute(cfg config.Config) string {
	return "/cosmos/staking/v1beta1/validators"
}

func getTxByHeightRoute(cfg config.Config) string {
	return "/cosmos/tx/v1beta1/txs?events=tx.height="
}

func getSlashingParamsRoute(cfg config.Config) string {
	return "/cosmos/slashing/v1beta1/params"
}

func getSigningInfoByAddressRoute(cfg config.Config) string {
	return "/cosmos/slashing/v1beta1/signing_infos/"
}

func getProposalsRoute(cfg config.Config) string {
	return "/cosmos/gov/v1beta1/proposals"
}

func getNodeInfoRoute() string {
	return "/cosmos/base/tendermint/v1beta1/node_info"
}

func getUpgradeCurrentPlanRoute(cfg config.Config) string {
	return "/cosmos/upgrade/v1beta1/current_plan"
}

/***********************
 * IBC Routes
************************/
func getIBCChannelsRoute(cfg config.Config) string {
	return "/ibc/core/channel/v1/channels"
}

func getIBCConnectionsRoute(cfg config.Config) string {
	return "/ibc/core/connection/v1/connections"
}

/***********************
 * Gravity Bridge Routes
************************/
func getBatchFeesRoute() string {
	return "/gravity/v1beta/batchfees"
}

func getBatchesFeesRoute() string {
	return "/gravity/v1beta1/batch/outgoingtx"
}

func getBridgeFeesRoute() string {
	return "https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"
}

func getBridgeParamsRoute() string {
	return "/gravity/v1beta/params"
}

func getOracleEventNonceByAddressRoute() string {
	return "/gravity/v1beta/oracle/eventnonce/"
}

func getCurrentValidatorSetRoute() string {
	return "/gravity/v1beta/valset/current"
}

/***********************
 * Akash Routes
************************/
func getDeploymentsRoute() string {
	return "/akash/deployment/v1beta2/deployments/list"
}
