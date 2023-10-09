package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type Block struct {
	Height    int
	BlockHash string
	Timestamp time.Time
}

type BlockService struct {
	Block *types.Blocks
	DB    *sql.DB
}

func (bs *BlockService) GetInfo(cfg config.Config) types.Blocks {
	route := rest.GetBlockInfoRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &bs.Block)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	return *bs.Block
}

func (bs *BlockService) GetLastBlockTimestamp(cfg config.Config, currentHeight int64) types.Blocks {
	var lastBlock types.LastBlock
	route := rest.GetBlockByHeightRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route + strconv.Itoa(int(currentHeight-1)))
	if err != nil {
		zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &lastBlock)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	bs.Block.Block.Header.LastTimestamp = lastBlock.Block.Header.Timestamp

	return *bs.Block
}

func (bs *BlockService) Index(height int, hash string, timestamp time.Time) (*Block, error) {
	block := Block{
		Height:    height,
		BlockHash: hash,
		Timestamp: timestamp,
	}
	row := bs.DB.QueryRow(`
		INSERT INTO blocks (height, block_hash, timestamp)
		VALUES ($1, $2, $3) RETURNING block_hash`, height, hash, timestamp)
	err := row.Scan(&block.BlockHash)
	if err != nil {
		return nil, fmt.Errorf("error indexing block: %w", err)
	}
	return &block, nil
}
