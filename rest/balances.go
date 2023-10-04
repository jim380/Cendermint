package rest

import (
	"encoding/json"
	"strings"

	"github.com/jim380/Cendermint/config"
	"go.uber.org/zap"
)

type balances struct {
	Balances []Coin
}

type Coin struct {
	Denom  string
	Amount string
}

func (rd *RESTData) getBalances(cfg config.Config) {
	var b balances

	route := getBalancesByAddressRoute(cfg)

	res, err := HttpQuery(RESTAddr + route + AccAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	// Unmarshal the JSON response and check for errors
	if err := json.Unmarshal(res, &b); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.Balances = b.Balances
}
