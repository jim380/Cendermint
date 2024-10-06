package controllers

import (
	"time"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
	"go.uber.org/zap"
)

func (rpc RpcServices) IndexValidator(consPubKey, consAddr, consAddrHex, moniker string, lastActive time.Time) {
	validator, err := rpc.ValidatorService.Index(consPubKey, consAddr, consAddrHex, moniker, lastActive)

	if err != nil {
		zap.L().Error("Error indexing validator", zap.Error(err))
		return
	} else {
		zap.L().Debug("Validator successfully indexed", zap.String("ConsPubAddress", validator.ConsPubKey))
	}
}

func (rs RestServices) IndexAbsentValidator(height int, consAddrBase64 string) {
	absentValidator, err := rs.AbsentValidatorService.Index(height, consAddrBase64)
	if err != nil {
		zap.L().Error("Error indexing absent validator", zap.String("ConsAddrBase64", consAddrBase64), zap.Error(err))
		return
	} else {
		zap.L().Debug("Absent validator successfully indexed", zap.String("ConsHexAddress", absentValidator.ConsAddrBase64))
	}
}

func (rpc RpcServices) GetValidatorInfo(cfg config.Config, currentBlockHeight int64, rd *types.RESTData) []string {
	return rpc.ValidatorService.GetValidatorInfo(cfg, currentBlockHeight, rd)
}
