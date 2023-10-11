package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RPCServices) GetRPCInfo(cfg config.Config, rpc *types.RPCData) {
	rs.ConsensusService.GetConsensusDump(cfg, rpc)
}
