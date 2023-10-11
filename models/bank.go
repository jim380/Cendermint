package models

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type rewardsAndCommisson struct {
	Height string `json:"height"`
	Result struct {
		Operator_Address string       `json:"operator_address"`
		Selfbond_Rewards []types.Coin `json:"self_bond_rewards"`
		Commission       `json:"val_commission"`
	}
}

type Commission struct {
	Commission []types.Coin `json:"commission"`
}

type BankService struct {
	DB *sql.DB
}

func (bs *BankService) GetBalanceInfo(cfg config.Config, rd *types.RESTData) {
	var b types.Balances

	route := rest.GetBalancesByAddressRoute(cfg)

	res, err := utils.HttpQuery(constants.RESTAddr + route + constants.AccAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &b); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.Balances = b.Balances
}

func (bs *BankService) GetRewardsCommissionInfo(cfg config.Config, rd *types.RESTData) {
	var rc rewardsAndCommisson

	route := rest.GetValidatorDistributionByAddressRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route + constants.OperAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &rc); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.Rewards = rc.Result.Selfbond_Rewards
	rd.Commission = rc.Result.Commission.Commission
}
