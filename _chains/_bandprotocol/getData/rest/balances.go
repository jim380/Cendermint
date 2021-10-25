package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"
)

type balances struct {
	Height string `json:"height"`
	Result []Coin
}

type Coin struct {
	Denom  string
	Amount string
}

func getBalances(accAddr string, log *zap.Logger) []Coin {

	var b balances

	res, err := runRESTCommand("/bank/balances/" + accAddr)
	if err != nil {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
	}
	json.Unmarshal(res, &b)
	// log
	if strings.Contains(string(res), "not found") {
		// handle error
		log.Fatal("REST-Server", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		log.Info("REST-Server", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("Get Data", "Staking Pool"))
	}

	return b.Result
}
