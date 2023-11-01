package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RpcServices) GetRpcInfo(cfg config.Config, rpc *types.RPCData) {
	conspubMonikerMap := rs.ConsensusService.GetConsensusDump(cfg, rpc)
	for k, v := range conspubMonikerMap {
		rs.IndexValidator(k, v)
	}
}
