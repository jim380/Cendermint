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

func getRewards(log *zap.Logger) []Coin {
	var r rewards

	res, _ := runRESTCommand("/cosmos/distribution/v1beta1/delegators/" + AccAddr + "/rewards/" + OperAddr)
	json.Unmarshal(res, &r)
	if strings.Contains(string(res), "not found") {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		log.Info("\t", zap.Bool("Success", true), zap.String("Rewards", fmt.Sprint(r.Rewards)))
	}

	return r.Rewards
}
