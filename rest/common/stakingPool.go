package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"

	utils "github.com/jim380/Cosmos-IE/utils"
)

var ()

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

func getStakingPool(denom string, log *zap.Logger) stakingPool {
	var sp stakingPool

	res, err := runRESTCommand("/cosmos/staking/v1beta1/pool")
	if err != nil {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
	}
	json.Unmarshal(res, &sp)
	if strings.Contains(string(res), "not found") {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		log.Info("", zap.Bool("Success", true), zap.String("Bonded tokens", sp.Pool.Bonded_tokens))
	}

	sp.Pool.Total_supply = getTotalSupply(denom, log)

	return sp
}

func getTotalSupply(denom string, log *zap.Logger) float64 {
	var ts totalSupply

	res, err := runRESTCommand("/cosmos/bank/v1beta1/supply/" + denom)
	if err != nil {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
	}
	json.Unmarshal(res, &ts)
	if strings.Contains(string(res), "not found") {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		log.Info("", zap.Bool("Success", true), zap.String("Total Supply", ts.Amount.Amount))
	}

	return utils.StringToFloat64(ts.Amount.Amount)
}
