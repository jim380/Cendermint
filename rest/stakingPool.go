package rest

import (
	"encoding/json"
	"strconv"
	"strings"

	"go.uber.org/zap"

	utils "github.com/jim380/Cendermint/utils"
)

type stakingPool struct {
	Pool struct {
		Not_bonded_tokens string `json:"not_bonded_tokens"`
		Bonded_tokens     string `json:"bonded_tokens"`
		Total_supply      float64
	}
}

type supplyResult struct {
	Supply     `json:"supply"`
	Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	} `json:"pagination"`
}

type Supply []struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

func (rd *RESTData) getStakingPool(denom string) {
	var sp stakingPool

	res, err := HttpQuery(RESTAddr + "/cosmos/staking/v1beta1/pool")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
	}
	json.Unmarshal(res, &sp)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Bonded tokens", sp.Pool.Bonded_tokens))
	}

	sp.Pool.Total_supply = getTotalSupply(denom, zap.L())
	rd.StakingPool = sp
}

func getTotalSupply(denom string, log *zap.Logger) float64 {
	var ts supplyResult

	res, err := HttpQuery(RESTAddr + "/cosmos/bank/v1beta1/supply")
	if err != nil {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &ts)
	length, err := strconv.Atoi(ts.Pagination.Total)
	if err != nil {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if strings.Contains(string(res), "not found") {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		log.Info("", zap.Bool("Success", true), zap.String("Total supply", ts.Supply[length-1].Amount))
	}

	return utils.StringToFloat64(ts.Supply[length-1].Amount)
}
