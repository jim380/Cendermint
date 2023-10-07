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

	// Check if the response is valid JSON
	if !json.Valid(res) {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Response is not valid JSON"))
	}

	if err := json.Unmarshal(res, &upgradeInfo); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
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
