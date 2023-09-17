package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"

	"github.com/jim380/Cendermint/config"
	utils "github.com/jim380/Cendermint/utils"
)

type stakingPool struct {
	Pool struct {
		Not_bonded_tokens string `json:"not_bonded_tokens"`
		Bonded_tokens     string `json:"bonded_tokens"`
		Total_supply      float64
	}
}

type totalSupply struct {
	Amount Coin
}

func (rd *RESTData) getStakingPool(cfg config.Config, denom string) {
	var sp stakingPool

	route := getStakingPoolRoute(cfg)
	res, err := HttpQuery(RESTAddr + route)
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

	route := getSupplyRoute(cfg)
	res, err := HttpQuery(RESTAddr + route + denom)
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
