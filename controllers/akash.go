package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
	"go.uber.org/zap"
)

func (rs RestServices) GetAkashInfo(cfg config.Config, data *types.AsyncData) {
	deployments := rs.AkashService.GetAkashDeployments(cfg, data)
	// index
	rs.IndexAkashDeployments(cfg, deployments)
}

func (rs RestServices) IndexAkashDeployments(cfg config.Config, deployments types.Deployments) {
	err := rs.AkashService.IndexDeployments(cfg, deployments)
	if err != nil {
		zap.L().Error("Error indexing akash deployments", zap.String("Error", err.Error()))
		return
	} else {
		zap.L().Debug("Akash deployments successfully indexed", zap.String("", ""))
	}
}
