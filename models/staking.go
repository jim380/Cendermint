package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type StakingService struct {
	DB *sql.DB
}

type totalSupply struct {
	Amount types.Coin
}

func (ss *StakingService) GetInfo(cfg config.Config, denom string, rd *types.RESTData) {
	var sp types.StakingPool

	route := rest.GetStakingPoolRoute(cfg)
	res, err := utils.HTTPQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &sp); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Bonded tokens", sp.Pool.BondedTokens))
	}

	totalSupply, err := getTotalSupply(cfg, denom, zap.L())
	if err != nil {
		// Handle the error here, such as logging or returning an error response.
		zap.L().Error("Failed to get total supply", zap.Error(err))
	}

	sp.Pool.TotalSupply = totalSupply
	rd.StakingPool = sp
}

func getTotalSupply(cfg config.Config, denom string, log *zap.Logger) (float64, error) {
	var ts totalSupply

	route := rest.GetSupplyRoute(cfg)
	res, err := utils.HTTPQuery(constants.RESTAddr + route + denom)
	if err != nil {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		return 0, err
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return 0, fmt.Errorf("response is not valid JSON")
	}
	if err := json.Unmarshal(res, &ts); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return 0, err
	}
	switch {
	case strings.Contains(string(res), "not found"):
		log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
		return 0, fmt.Errorf("resource not found")
	case strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":"):
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
		return 0, fmt.Errorf("error in response")
	default:
		log.Info("", zap.Bool("Success", true), zap.String("total supply", ts.Amount.Amount))
	}

	return utils.StringToFloat64(ts.Amount.Amount), nil
}
