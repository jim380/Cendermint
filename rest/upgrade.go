package rest

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/jim380/Cendermint/config"
	"go.uber.org/zap"
)

type upgradeInfo struct {
	Planned bool
	Plan    struct {
		Name   string `json:"name"`
		Time   string `json:"time"`
		Height string `json:"height"`
		Info   string `json:"info"`
	}
}

func (rd *RESTData) getUpgradeInfo(cfg config.Config) {
	var upgradeInfo upgradeInfo

	route := getUpgradeCurrentPlanRoute(cfg)
	res, err := HttpQuery(RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &upgradeInfo)
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
