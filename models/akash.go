package models

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type AkashService struct {
	DB *sql.DB
}

func (as *AkashService) GetAkashDeployments(cfg config.Config, rd *types.RESTData) {
	if cfg.Chain.Chain != "akash" {
		return
	}
	var deployments, activeDeployments types.AkashDeployments

	route := rest.GetDeploymentsRoute()
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &deployments)

	// get total deployments count
	totalDeploymentsCount, err := strconv.Atoi(deployments.Pagination.Total)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	rd.AkashInfo.TotalDeployments = totalDeploymentsCount

	// get active deployments count
	resActive, err := utils.HttpQuery(constants.RESTAddr + route + "?filters.state=active")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(resActive, &activeDeployments)

	activeDeploymentsCount, err := strconv.Atoi(activeDeployments.Pagination.Total)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	rd.AkashInfo.ActiveDeployments = activeDeploymentsCount

	// get closed deployments count
	rd.AkashInfo.ClosedDeployments = totalDeploymentsCount - activeDeploymentsCount
}
