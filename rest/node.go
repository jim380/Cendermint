package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"
)

type nodeInfo struct {
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

func (rd *RESTData) getNodeInfo() {
	var nodeInfo nodeInfo

	res, err := RESTQuery("/cosmos/base/tendermint/v1beta1/node_info")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &nodeInfo)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		// zap.L().Info("", zap.Bool("Success", true), zap.String("SDK Version:", nodeInfo.Application.SDKVersion))
	}

	rd.NodeInfo = nodeInfo
}
