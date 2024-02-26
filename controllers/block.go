package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
	"go.uber.org/zap"
)

func (rs RestServices) IndexBlock(height int, hash string, timestamp time.Time, proposer string, txnCount int) {
	block, err := rs.BlockService.Index(height, hash, timestamp, proposer, txnCount)
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
	propser := block.Block.Header.ProposerAddress

	rs.IndexBlock(height, block.BlockId.Hash, timestamp, propser, txnCount)
	// rs.IndexTxnsInBlock(cfg, int64(height))

	return block
}

func (rs RestServices) GetLastBlockTimestamp(cfg config.Config, currentHeight int64) types.Blocks {
	block := rs.BlockService.GetLastBlockTimestamp(cfg, currentHeight)

	return block
}
