package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetInflationInfo(cfg config.Config, rd *types.RESTData) {
	rs.InflationService.GetInfo(cfg, rd)
}
