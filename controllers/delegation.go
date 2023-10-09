package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetDelegationsInfo(cfg config.Config, rd *types.RESTData) {
	rs.DelegationService.GetInfo(cfg, rd)
}
