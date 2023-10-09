package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetGravityBridgeInfo(cfg config.Config, rd *types.RESTData) {
	rs.GravityService.GetBatchFees(cfg, rd)
	rs.GravityService.GetBatchesFees(cfg, rd)
	rs.GravityService.GetBridgeFees(cfg, rd)
	rs.GravityService.GetBridgeParams(cfg, rd)
	rs.GravityService.GetOracleEventNonce(cfg, rd)
	rs.GravityService.GetValSet(cfg, rd)
}
