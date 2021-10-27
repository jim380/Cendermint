package rest

import (
	"encoding/json"
	"fmt"
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

	res, err := runRESTCommand("/cosmos/bank/v1beta1/balances/" + AccAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
	}
	json.Unmarshal(res, &b)
	if strings.Contains(string(res), "not found") {
		// handle error
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Balances", fmt.Sprint(b.Balances)))
	}

	rd.Balances = b.Balances
}
