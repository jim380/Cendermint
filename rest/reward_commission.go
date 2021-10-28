package rest

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type rewardsAndCommisson struct {
	Height string `json:"height"`
	Result struct {
		Operator_Address string `json:"operator_address"`
		Selfbond_Rewards []Coin `json:"self_bond_rewards"`
		Commission       `json:"val_commission"`
	}
}

type Commission struct {
	Commission []Coin `json:"commission"`
}

func (rd *RESTData) getRewardsCommission() {
	var rc rewardsAndCommisson

	res, err := runRESTCommand("/distribution/validators/" + OperAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
	}
	json.Unmarshal(res, &rc)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Rewards:", fmt.Sprint(rc.Result.Selfbond_Rewards)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Commission:", fmt.Sprint(rc.Result.Commission.Commission)))
	}

	rd.Rewards = rc.Result.Selfbond_Rewards
	rd.Commission = rc.Result.Commission.Commission
}
