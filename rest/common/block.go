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
			Height           string
			Proposer_address string
		}

		Last_commit struct {
			Signatures []struct {
				Block_id_flag     string
				Validator_address string
			}
		}
	}
}

func (b *Blocks) GetInfo() Blocks {
	res, err := runRESTCommand("/blocks/latest")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
	}
	json.Unmarshal(res, &b)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("Info", zap.Bool("Success", true), zap.String("Block Info ", "successfully fetched."))
	}

	return *b
}
