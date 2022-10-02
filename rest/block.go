package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"
)

type Blocks struct {
	Block struct {
		Header     header     `json:"header"`
		LastCommit lastCommit `json:"last_commit"`
	}
}

type header struct {
	ChainID          string `json:"chain_id"`
	Height           string `json:"height"`
	Proposer_address string `json:"proposer_address"`
	Timestamp        string `json:"time"`
}

type lastCommit struct {
	Signatures []struct {
		Validator_address string `json:"validator_address"`
	}
}

func (b *Blocks) GetInfo() Blocks {
	res, err := HttpQuery(RESTAddr + "/cosmos/base/tendermint/v1beta1/blocks/latest")
	if err != nil {
		zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &b)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	return *b
}
