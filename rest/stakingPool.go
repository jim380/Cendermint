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
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &sp); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}
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
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return utils.StringToFloat64(ts.Amount.Amount)
	}
	if err := json.Unmarshal(res, &ts); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return utils.StringToFloat64(ts.Amount.Amount)
	}

	resStr := string(res)
	switch {
	case strings.Contains(resStr, "not found"):
		log.Fatal("", zap.Bool("Success", false), zap.String("err", resStr))
	case strings.Contains(resStr, "error:") || strings.Contains(resStr, "error\\\":"):
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", resStr))
	default:
		log.Info("", zap.Bool("Success", true), zap.String("Total supply", ts.Amount.Amount))
	}

	return utils.StringToFloat64(ts.Amount.Amount)
}
