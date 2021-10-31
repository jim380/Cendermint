package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"
)

type Blocks struct {
	Block struct {
		Header struct {
			ChainID          string `json:"chain_id"`
			Height           string `json:"height"`
			Proposer_address string `json:"proposer_address"`
		}

		Last_commit struct {
			Signatures []struct {
				Validator_address string `json:"validator_address"`
			}
		}
	}
}

func (b *Blocks) GetInfo() Blocks {
	res, err := RESTQuery("/blocks/latest")
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
