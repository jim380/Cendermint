package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetUpgradeInfo(cfg config.Config, rd *types.RESTData) {
	rs.UpgradeService.GetInfo(cfg, rd)
}
