package models

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type UpgradeService struct {
	DB *sql.DB
}

func (us *UpgradeService) GetInfo(cfg config.Config, rd *types.RESTData) {
	var upgradeInfo types.UpgradeInfo

	route := rest.GetUpgradeCurrentPlanRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &upgradeInfo); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	if upgradeInfo.Plan.Name != "" {
		upgradeInfo.Planned = true
	}
	zap.L().Info("", zap.Bool("Success", true), zap.String("Upgrade planned", strconv.FormatBool(upgradeInfo.Planned)))
	rd.UpgradeInfo = upgradeInfo
}
