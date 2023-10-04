package rest

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/jim380/Cendermint/config"
	"go.uber.org/zap"
)

type LastBlock struct {
	Block struct {
		Header lastBlockHeader `json:"header"`
	}
}

type lastBlockHeader struct {
	Timestamp string `json:"time"`
}

type Blocks struct {
	BlockId struct {
		Hash string `json:"hash"`
	} `json:"block_id"`
	Block struct {
		Header     header     `json:"header"`
		LastCommit lastCommit `json:"last_commit"`
	} `json:"block"`
	MissingValidators []struct {
		Moniker     string
		ConsHexAddr string
	} // not part of the response so no json tag
}

type header struct {
	ChainID          string `json:"chain_id"`
	Height           string `json:"height"`
	Proposer_address string `json:"proposer_address"`
	Timestamp        string `json:"time"`
	LastTimestamp    string // not part of the response so no json tag
}

type lastCommit struct {
	Signatures []struct {
		Validator_address string `json:"validator_address"`
		Signature         string `json:"signature"`
	} `json:"signatures"`
}

func (b *Blocks) GetInfo(cfg config.Config) Blocks {
	route := GetBlockInfoRoute(cfg)
	res, err := HttpQuery(RESTAddr + route)
	if err != nil {
		zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	// Unmarshal the JSON response and check for errors
	if err := json.Unmarshal(res, &b); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	return *b
}

func (b *Blocks) GetLastBlockTimestamp(cfg config.Config, currentHeight int64) Blocks {
	var lastBlock LastBlock
	route := getBlockByHeightRoute(cfg)
	res, err := HttpQuery(RESTAddr + route + strconv.Itoa(int(currentHeight-1)))
	if err != nil {
		zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	// Unmarshal the JSON response and check for errors
	if err := json.Unmarshal(res, &lastBlock); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	b.Block.Header.LastTimestamp = lastBlock.Block.Header.Timestamp

	return *b
}
