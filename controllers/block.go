package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
	"go.uber.org/zap"
)

func (rs RestServices) IndexBlock(height int, hash string, timestamp time.Time, txnCount int) {
	block, err := rs.BlockService.Index(height, hash, timestamp, txnCount)
	if err != nil {
		zap.L().Error("Error indexing block", zap.Error(err))
		return
	} else {
		zap.L().Debug("Block successfully indexed", zap.String("Height", strconv.Itoa(block.Height)))
	}
}

func (rs RestServices) GetBlockInfo(cfg config.Config) types.Blocks {
	block := rs.BlockService.GetInfo(cfg)

	// index
	height, err := strconv.Atoi(block.Block.Header.Height)
	if err != nil {
		log.Fatal("Failed to convert height to int: ", err)
	}
	timestamp, err := time.Parse(time.RFC3339, block.Block.Header.Timestamp)
	if err != nil {
		log.Fatal("Failed to parse timestamp: ", err)
	}
	txnCount := len(block.Block.Data.Txs)

	rs.IndexBlock(height, block.BlockId.Hash, timestamp, txnCount)

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
