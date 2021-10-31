package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"
)

type balances struct {
	Balances []Coin
}

type Coin struct {
	Denom  string
	Amount string
}

func (rd *RESTData) getBalances() {
	var b balances

	res, err := RESTQuery("/cosmos/bank/v1beta1/balances/" + AccAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &b)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.Balances = b.Balances
}
