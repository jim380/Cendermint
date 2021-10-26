package rest

import (
	"encoding/json"
	"fmt"
	"strings"

	utils "github.com/jim380/Cosmos-IE/utils"
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
	accAddr := utils.GetAccAddrFromOperAddr(OperAddr)

	res, err := runRESTCommand("/cosmos/bank/v1beta1/balances/" + accAddr)
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
