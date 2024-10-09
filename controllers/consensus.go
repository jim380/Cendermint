package controllers

import (
	"fmt"
	"time"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

func (rs RpcServices) GetRpcInfo(cfg config.Config, rpc *types.RPCData) {
	validatorsetMap := rs.ConsensusService.GetConsensusDump(cfg, rpc)
	lastActive := time.Now().UTC()
	for consAddrHex, v := range validatorsetMap {
		consPubKey := v[0]
		consAddr, err := utils.HexToBase64(consAddrHex)
		if err != nil {
			zap.L().Fatal("GetRpcInfo", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
		}
		moniker := v[5]

		rs.IndexValidator(consPubKey, consAddr, consAddrHex, moniker, lastActive)
	}
}
