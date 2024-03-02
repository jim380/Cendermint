package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetAkashData(cfg config.Config, data *types.AsyncData) {
	rs.AkashService.GetAkashDeployments(cfg, data)
}
