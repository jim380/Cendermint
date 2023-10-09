package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/rest/types"
	"go.uber.org/zap"
)

func (rs RestServices) IndexBlock(height int, hash string, timestamp time.Time) {
	block, err := rs.BlockService.Index(height, hash, timestamp)
	if err != nil {
		zap.L().Error("Error indexing block", zap.Error(err))
		return
	} else {
		zap.L().Info("Block successfully indexed", zap.String("Height", strconv.Itoa(block.Height)))
	}
}

func (rs RestServices) GetBlockInfo(cfg config.Config) types.Blocks {
	block := rs.BlockService.GetInfo(cfg)

	return block
}

func (rs RestServices) GetLastBlockTimestamp(cfg config.Config, currentHeight int64) types.Blocks {
	block := rs.BlockService.GetLastBlockTimestamp(cfg, currentHeight)

	fmt.Println("--------------------------- Start ---------------------------")
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Last block timestamp", block.Block.Header.LastTimestamp))
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Current block timestamp", block.Block.Header.Timestamp))
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Current block height", fmt.Sprint(currentHeight)))

	return block
}
