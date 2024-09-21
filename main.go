/*
Copyright © 2023 Jay Jie

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"strconv"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/controllers"
	"github.com/jim380/Cendermint/dashboard"
	"github.com/jim380/Cendermint/exporter"
	"github.com/jim380/Cendermint/fetcher"
	"github.com/jim380/Cendermint/logging"
	"github.com/jim380/Cendermint/models"
	"github.com/jim380/Cendermint/types"
	"go.uber.org/zap"
)

var (
	appConfig types.AppConfig
)

func main() {
	cfg := config.LoadConfig()
	appConfig = cfg.ValidateConfig()

	// initialize logger
	logger := logging.InitLogger(cfg.LogOutput, appConfig.LogLevel)
	zap.ReplaceGlobals(logger)

	// initialize sdk config
	cfg.SetSDKConfig()

	// initialize constants
	constants.RESTAddr = appConfig.RestAddr
	constants.RPCAddr = appConfig.RpcAddr
	constants.OperAddr = appConfig.OperAddr
	constants.PollIntervalChain, _ = strconv.Atoi(appConfig.PollIntervalChain)
	constants.PollIntervalAsync, _ = strconv.Atoi(appConfig.PollIntervalAsync)
	constants.PollIntervalBackfill, _ = strconv.Atoi(appConfig.PollIntervalBackfill)
	constants.LastUpdatedMoreThan, _ = strconv.Atoi(appConfig.LastUpdatedMoreThan)

	// setup a db connection
	db := models.SetupDatabase()
	defer db.Close()

	// db migration
	models.MigrateDatabase(db)

	// initialize services
	rpcServicesController := controllers.InitializeRpcServices(db)
	restServicesController := controllers.InitializeRestServices(db)

	// start dashboard in a separate thread in enabled
	if strings.ToLower(cfg.DashboardEnabled) == "true" {
		go dashboard.StartDashboard()
	}

	denomList, err := config.GetDenomList(cfg.Chain.Name, cfg.ChainList)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	// start the data fetcher in a separate thread
	fetcher.Start(&cfg, restServicesController, rpcServicesController, denomList, logger)

	// start the exporter in a separate thread
	exporter.Start(&cfg, cfg.ListeningPort, logger, restServicesController, rpcServicesController)

	// start the metrics server
	exporter.StartMetricsHttpServer(cfg.ListeningPort)
}
