package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"

	utils "github.com/jim380/Cosmos-IE/utils"
)

type inflation struct {
	Height string `json:"height"`
	Result string `json:"result"`
}

type inflation_iris struct {
	Params struct {
		Mint_Denom string
		Inflation  string
	}
}

func getInflation(chain string, denom string, log *zap.Logger) float64 {
	var result string

	switch chain {
	case "iris":
		var i inflation_iris

		res, _ := runRESTCommand("/irishub/mint/params")
		json.Unmarshal(res, &i)
		if strings.Contains(string(res), "not found") {
			log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
		} else {
			result = i.Params.Inflation
			log.Info("\t", zap.Bool("Success", true), zap.String("Inflation", result))
		}
	default:
		var i inflation
		// Does not work
		// res, _ := runRESTCommand("/cosmos/mint/v1beta1/inflation")
		res, _ := runRESTCommand("/minting/inflation")
		json.Unmarshal(res, &i)
		if strings.Contains(string(res), "not found") {
			log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
		} else {
			result = i.Result
			log.Info("\t", zap.Bool("Success", true), zap.String("Inflation", result))
		}
	}

	return utils.StringToFloat64(result)
}
