package models

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
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

type InflationService struct {
	DB *sql.DB
}

func (is *InflationService) GetInfo(cfg config.Config, rd *types.RESTData) {
	var result string

	route := rest.GetInflationRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route)

	switch cfg.Chain.Chain {
	case "irisnet":
		var i inflation_iris
		if err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
		}
		if !json.Valid(res) {
			zap.L().Error("Response is not valid JSON")
			return
		}
		if err := json.Unmarshal(res, &i); err != nil {
			zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
			return
		}
		if strings.Contains(string(res), "not found") {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
		} else {
			result = i.Params.Inflation
			zap.L().Info("\t", zap.Bool("Success", true), zap.String("Inflation", result))
		}
	default:
		var i inflation

		res, err := utils.HttpQuery(constants.RESTAddr + route) // route does not existing in osmosis
		if err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}
		if !json.Valid(res) {
			zap.L().Error("Response is not valid JSON")
			return
		}
		if err := json.Unmarshal(res, &i); err != nil {
			zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
			return
		}
		switch {
		case strings.Contains(string(res), "not found"):
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
		case strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":"):
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
		default:
			result = i.Result
			zap.L().Info("\t", zap.Bool("Success", true), zap.String("Inflation", result))
		}
	}

	rd.Inflation = utils.StringToFloat64(result)
}
