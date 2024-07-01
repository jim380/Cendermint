package services

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
	Value string `json:"inflation"`
}

type InflationService struct {
	DB *sql.DB
}

func (is *InflationService) Init(db *sql.DB) {
	is.DB = db
}

func (is *InflationService) GetInfo(cfg config.Config, rd *types.RESTData) {
	var result string
	var i inflation

	route := rest.GetInflationRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		rd.Inflation = 0
		return
	}

	json.Unmarshal(res, &i)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		result = i.Value
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Inflation", result))
	}

	rd.Inflation = utils.StringToFloat64(result)
}
