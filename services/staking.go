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

type StakingService struct {
	DB *sql.DB
}

func (stks *StakingService) Init(db *sql.DB) {
	stks.DB = db
}

type totalSupply struct {
	Amount types.Coin
}

func (ss *StakingService) GetInfo(cfg config.Config, denom string, rd *types.RESTData) {
	var sp types.StakingPool

	route := rest.GetStakingPoolRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
	}
	json.Unmarshal(res, &sp)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Bonded tokens", sp.Pool.Bonded_tokens))
	}

	sp.Pool.Total_supply = getTotalSupply(cfg, denom, zap.L())
	rd.StakingPool = sp
}

func getTotalSupply(cfg config.Config, denom string, log *zap.Logger) float64 {
	var ts totalSupply

	route := rest.GetSupplyRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route + denom)
	if err != nil {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &ts)
	if strings.Contains(string(res), "not found") {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		log.Info("", zap.Bool("Success", true), zap.String("Total supply", ts.Amount.Amount))
	}

	return utils.StringToFloat64(ts.Amount.Amount)
}
