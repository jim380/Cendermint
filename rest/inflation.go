package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"

	"github.com/jim380/Cendermint/config"
	utils "github.com/jim380/Cendermint/utils"
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

func (rd *RESTData) getInflation(cfg config.Config, denom string) {
	var result string

	route := getInflationRoute(cfg)
	res, err := HttpQuery(RESTAddr + route)

	switch cfg.Chain {
	case "iris":
		var i inflation_iris
		if err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
		}
		json.Unmarshal(res, &i)
		if strings.Contains(string(res), "not found") {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
		} else {
			result = i.Params.Inflation
			zap.L().Info("\t", zap.Bool("Success", true), zap.String("Inflation", result))
		}
	default:
		var i inflation

		res, err := HttpQuery(RESTAddr + route) // route does not existing in osmosis
		if err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}
		json.Unmarshal(res, &i)
		if strings.Contains(string(res), "not found") {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
		} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
		} else {
			result = i.Result
			zap.L().Info("\t", zap.Bool("Success", true), zap.String("Inflation", result))
		}
	}

	rd.Inflation = utils.StringToFloat64(result)
}
