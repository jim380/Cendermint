package rest

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type commission struct {
	Commission []Coin
}

func getCommission(log *zap.Logger) []Coin {
	var c commission

	res, _ := runRESTCommand("/cosmos/distribution/v1beta1/validators/" + OperAddr + "/commission")
	json.Unmarshal(res, &c)
	if strings.Contains(string(res), "not found") {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		log.Info("\t", zap.Bool("Success", true), zap.String("Commission", fmt.Sprint(c.Commission)))
	}

	return c.Commission
}
