package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetIbcChannelInfo(cfg config.Config, rd *types.RESTData) {
	rs.IbcServices.GetChannelInfo(cfg, rd)
}

func (rs RestServices) GetIbcConnectionInfo(cfg config.Config, rd *types.RESTData) {
	rs.IbcServices.GetConnectionInfo(cfg, rd)
}
