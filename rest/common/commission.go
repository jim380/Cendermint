package rest

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type commission struct {
	commissionInner
}

type commissionInner struct {
	Commission []Coin
}

func (rd *RESTData) getCommission() {
	var c commission

	res, err := runRESTCommand("/cosmos/distribution/v1beta1/validators/" + OperAddr + "/commission")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
	}
	json.Unmarshal(res, &c)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Commission", fmt.Sprint(c.Commission)))
	}

	rd.Commission = c.Commission
}
