package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetAkashInfo(cfg config.Config, rd *types.RESTData) {
	rs.AkashService.GetAkashDeployments(cfg, rd)
}
