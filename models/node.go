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

type NodeService struct {
	Node *types.NodeInfo
	DB   *sql.DB
}

func (ns *NodeService) GetInfo(cfg *config.Config, rd *types.RESTData) {
	var nodeInfo types.NodeInfo

	route := rest.GetNodeInfoRoute()
	res, err := utils.HTTPQuery(constants.RESTAddr + route)
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
