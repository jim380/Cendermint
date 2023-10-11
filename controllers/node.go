package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetNodeInfo(cfg *config.Config, rd *types.RESTData) {
	rs.NodeService.GetInfo(cfg, rd)
}
