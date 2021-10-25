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

func getBalances(accAddr string, log *zap.Logger) []Coin {
	var b balances

	res, _ := runRESTCommand("/cosmos/bank/v1beta1/balances/" + accAddr)
	json.Unmarshal(res, &b)
	if strings.Contains(string(res), "not found") {
		// handle error
		log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		log.Info("\t", zap.Bool("Success", true), zap.String("Balances", fmt.Sprint(b.Balances)))
	}

	return b.Balances
}
