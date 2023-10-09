package controllers

import (
	"strconv"
	"time"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/models"
	"github.com/jim380/Cendermint/rest"
	"go.uber.org/zap"
)

type RestServices struct {
	BlockService           *models.BlockService
	ValidatorService       *models.ValidatorService
	AbsentValidatorService *models.AbsentValidatorService
}

func (rs RestServices) IndexBlock(height int, hash string, timestamp time.Time) {
	block, err := rs.BlockService.Index(height, hash, timestamp)
	if err != nil {
		zap.L().Error("Error indexing block", zap.Error(err))
		return
	} else {
		zap.L().Info("Block successfully indexed", zap.String("Height", strconv.Itoa(block.Height)))
	}
}

func (rs RestServices) GetBlockInfo(cfg config.Config) rest.Blocks {
	block := rs.BlockService.GetBlockInfo(cfg)

	return block
}

func (rs RestServices) IndexValidator(consHexAddr, moniker string) {
	validator, err := rs.ValidatorService.Index(consHexAddr, moniker)
	if err != nil {
		zap.L().Error("Error indexing validator", zap.Error(err))
		return
	} else {
		zap.L().Info("Validator successfully indexed", zap.String("ConsHexAddress", validator.ConsHexAddress))
	}
}

func (rs RestServices) IndexAbsentValidator(height int, consHexAddr string) {
	absentValidator, err := rs.AbsentValidatorService.Index(height, consHexAddr)
	if err != nil {
		zap.L().Error("Error indexing abscent validator", zap.Error(err))
		return
	} else {
		zap.L().Info("Absent validator successfully indexed", zap.String("ConsHexAddress", absentValidator.ConsHexAddress))
	}
}
