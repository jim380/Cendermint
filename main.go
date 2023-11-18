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
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/controllers"
	"github.com/jim380/Cendermint/dashboard"
	"github.com/jim380/Cendermint/exporter"
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

	// Setup a db connection
	db := models.SetupDatabase()
	defer db.Close()

	// DB migration
	models.MigrateDatabase(db)

	// initialize services
	rpcServicesController := controllers.InitializeRpcServices(db)
	restServicesController := controllers.InitializeRestServices(db)

	// run dashboard in a separate thread in enabled
	if strings.ToLower(cfg.DashboardEnabled) == "true" {
		dashboard.StartDashboard()
	}

	exporter.Start(&cfg, cfg.ListeningPort, logger, restServicesController, rpcServicesController)
}
