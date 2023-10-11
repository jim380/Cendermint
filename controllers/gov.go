package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetGovInfo(cfg config.Config, rd *types.RESTData) {
	rs.GovService.GetInfo(cfg, rd)
}
