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
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}

	// Unmarshal the JSON response and check for errors
	if err := json.Unmarshal(res, &deployments); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}

	// get total deployments count
	if deployments.Pagination.Total == "" {
		zap.L().Error("Total deployments count is empty")
		return
	}
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
	if !json.Valid(resActive) {
		zap.L().Error("Response is not valid JSON")
		return
	}

	// Unmarshal the JSON response and check for errors
	if err := json.Unmarshal(resActive, &activeDeployments); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}
	if activeDeployments.Pagination.Total == "" {
		zap.L().Error("Active deployments count is empty")
		return
	}
	activeDeploymentsCount, err := strconv.Atoi(activeDeployments.Pagination.Total)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	rd.AkashInfo.ActiveDeployments = activeDeploymentsCount

	// get closed deployments count
	rd.AkashInfo.ClosedDeployments = totalDeploymentsCount - activeDeploymentsCount
}
