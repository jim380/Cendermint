package controllers

import (
	"strconv"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
	"go.uber.org/zap"
)

func (rs RestServices) GetSlashingInfo(cfg config.Config, rd *types.RESTData) {
	rs.SlashingService.GetSlashingParams(cfg, rd)
}

func (rs RestServices) GetSigningInfo(cfg config.Config, consAddr string, rd *types.RESTData) {
	rs.SlashingService.GetSigningInfo(cfg, consAddr, rd)
}

func (rs RestServices) GetCommitInfo(cfg config.Config, rd *types.RESTData, blockData types.Blocks, consAddrHex string) {
	missingValidators := rs.SlashingService.GetCommitInfo(cfg, rd, blockData, consAddrHex)
	height := blockData.Block.Header.Height
	// convert height to int
	heightInt, err := strconv.Atoi(height)
	if err != nil {
		zap.L().Error("Failed to convert height to int: ", zap.Error(err))
	}
	// index missing validators
	for _, v := range missingValidators {
		rs.IndexAbsentValidator(heightInt, v.ConsPubAddr, v.Moniker)
	}
}
