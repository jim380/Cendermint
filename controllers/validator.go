package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
	"go.uber.org/zap"
)

func (rpc RpcServices) IndexValidator(consHexAddr, moniker string) {
	validator, err := rpc.ValidatorService.Index(consHexAddr, moniker)

	if err != nil {
		zap.L().Error("Error indexing validator", zap.Error(err))
		return
	} else {
		zap.L().Debug("Validator successfully indexed", zap.String("ConsHexAddress", validator.ConsHexAddress))
	}
}

func (rs RestServices) IndexAbsentValidator(height int, consAddrBase64 string) {
	absentValidator, err := rs.AbsentValidatorService.Index(height, consAddrBase64)
	if err != nil {
		zap.L().Error("Error indexing abscent validator", zap.Error(err))
		return
	} else {
		zap.L().Debug("Absent validator successfully indexed", zap.String("ConsHexAddress", absentValidator.ConsAddrBase64))
	}
}

func (rpc RpcServices) GetValidatorInfo(cfg config.Config, currentBlockHeight int64, rd *types.RESTData) []string {
	return rpc.ValidatorService.GetValidatorInfo(cfg, currentBlockHeight, rd)
}
