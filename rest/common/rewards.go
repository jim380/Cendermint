package rest

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type rewards struct {
	Rewards []Coin
}

func (rd *RESTData) getRewards() {
	var r rewards

	res, err := runRESTCommand("/cosmos/distribution/v1beta1/delegators/" + AccAddr + "/rewards/" + OperAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
	}
	json.Unmarshal(res, &r)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Rewards", fmt.Sprint(r.Rewards)))
	}

	rd.Rewards = r.Rewards
}
