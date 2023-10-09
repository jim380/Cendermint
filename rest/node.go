package rest

import (
	"encoding/json"
	"strings"

	"github.com/jim380/Cendermint/config"
	"go.uber.org/zap"
)

type NodeInfo struct {
	Default     DefaultInfo `json:"default_node_info"`
	Application appVersion  `json:"application_version"`
}

type DefaultInfo struct {
	NodeID    string `json:"default_node_id"`
	TMVersion string `json:"version"`
	Moniker   string `json:"moniker"`
}

type appVersion struct {
	AppName    string `json:"name"`
	Name       string `json:"app_name"`
	Version    string `json:"version"`
	GitCommit  string `json:"git_commit"`
	GoVersion  string `json:"go_version"`
	SDKVersion string `json:"cosmos_sdk_version"`
}

func (rd *RESTData) getNodeInfo(cfg *config.Config) {
	var nodeInfo NodeInfo

	route := GetNodeInfoRoute()
	res, err := HttpQuery(RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &nodeInfo); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	zap.L().Info("", zap.Bool("Success", true), zap.String("App version", nodeInfo.Application.Version))
	zap.L().Info("", zap.Bool("Success", true), zap.String("SDK version", nodeInfo.Application.SDKVersion))

	rd.NodeInfo = nodeInfo
	cfg.SDKVersion = nodeInfo.Application.SDKVersion
}
