package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetStakingInfo(cfg config.Config, denom string, rd *types.RESTData) {
	rs.StakingService.GetInfo(cfg, denom, rd)
}
