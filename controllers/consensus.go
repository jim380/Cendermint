package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
)

func (rs RpcServices) GetRpcInfo(cfg config.Config, rpc *types.RPCData) {
	validatorsetMap := rs.ConsensusService.GetConsensusDump(cfg, rpc)
	for k, v := range validatorsetMap {
		consPubKey := v[0]
		consAddrHex := k
		consAddr := utils.HexToBase64(consAddrHex)
		moniker := v[5]

		rs.IndexValidator(consPubKey, consAddr, consAddrHex, moniker)
	}
}
