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
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/controllers"
	"github.com/jim380/Cendermint/dashboard"
	"github.com/jim380/Cendermint/exporter"
	"github.com/jim380/Cendermint/logging"
	"github.com/jim380/Cendermint/models"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	chain, restAddr, rpcAddr, listeningPort, operAddr, logOutput string
	logLevel                                                     zapcore.Level
	logger                                                       *zap.Logger
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("CHAIN") == "" {
		log.Fatal("Chain was not provided.")
	}

	providedChain := os.Getenv("CHAIN")

	cfg := config.Config{
		OperatorAddr:     os.Getenv("OPERATOR_ADDR"),
		RestAddr:         os.Getenv("REST_ADDR"),
		RPCAddr:          os.Getenv("RPC_ADDR"),
		ListeningPort:    os.Getenv("LISTENING_PORT"),
		MissThreshold:    os.Getenv("MISS_THRESHOLD"),
		MissConsecutive:  os.Getenv("MISS_CONSECUTIVE"),
		LogOutput:        os.Getenv("LOG_OUTPUT"),
		PollInterval:     os.Getenv("POLL_INTERVAL"),
		LogLevel:         os.Getenv("LOG_LEVEL"),
		DashboardEnabled: os.Getenv("DASHBOARD_ENABLED"),
	}

	chainList := config.GetChainList()
	cfg.ChainList = chainList
	supportedChains := make([]string, 0, len(chainList))
	for key := range chainList {
		supportedChains = append(supportedChains, key)
	}
	var found bool
	if _, found = chainList[providedChain]; found {
		cfg.Chain = config.Chain{Chain: providedChain}
	}
	if !found {
		log.Fatal(fmt.Sprintf("%s is not supported", providedChain) + fmt.Sprint("\nList of supported chains: ", supportedChains))
	}

	cfg.CheckInputs(chainList)

	chain = cfg.Chain.Chain
	operAddr = cfg.OperatorAddr
	restAddr = cfg.RestAddr
	rpcAddr = cfg.RPCAddr
	listeningPort = cfg.ListeningPort
	logOutput = cfg.LogOutput
	logLevel = config.GetLogLevel(cfg.LogLevel)
	logger = logging.InitLogger(logOutput, logLevel)
	zap.ReplaceGlobals(logger)

	cfg.SetSDKConfig()
	constants.RESTAddr = restAddr
	constants.RPCAddr = rpcAddr
	constants.OperAddr = operAddr

	// Setup a db connection
	dbConfig := models.DefaultPostgresConfig()
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Using db config", dbConfig.String()))
	db, err := models.Open(dbConfig)
	if err != nil {
		zap.L().Fatal("\t", zap.Bool("Success", false), zap.String("Database connection", "failed with error:"+err.Error()))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Database connection", "ok"))
	}
	defer db.Close()

	// initialize rpc services
	consensusService := models.ConsensusService{
		DB: db,
	}

	rpcServicesController := controllers.RPCServices{
		ConsensusService: &consensusService,
	}

	// initialize rest services
	blockService := models.BlockService{
		DB: db,
	}
	validatorService := models.ValidatorService{
		DB: db,
	}
	absentValidatorService := models.AbsentValidatorService{
		DB: db,
	}
	nodeService := models.NodeService{
		DB: db,
	}
	stakingService := models.StakingService{
		DB: db,
	}
	slashingService := models.SlashingService{
		DB: db,
	}
	inflationService := models.InflationService{
		DB: db,
	}
	govService := models.GovService{
		DB: db,
	}
	bankService := models.BankService{
		DB: db,
	}
	delegationService := models.DelegationService{
		DB: db,
	}
	upgradeService := models.UpgradeService{
		DB: db,
	}
	ibcService := models.IbcService{
		DB: db,
	}
	gravityService := models.GravityService{
		DB: db,
	}
	akashService := models.AkashService{
		DB: db,
	}

	restServicesController := controllers.RestServices{
		BlockService:           &blockService,
		ValidatorService:       &validatorService,
		AbsentValidatorService: &absentValidatorService,
		NodeService:            &nodeService,
		StakingService:         &stakingService,
		SlashingService:        &slashingService,
		InflationService:       &inflationService,
		GovService:             &govService,
		BankService:            &bankService,
		DelegationService:      &delegationService,
		UpgradeService:         &upgradeService,
		IbcServices:            &ibcService,
		GravityService:         &gravityService,
		AkashService:           &akashService,
	}

	// run dashboard in a separate thread in enabled
	if strings.ToLower(cfg.DashboardEnabled) == "true" {
		dashboard.StartDashboard()
	}

	exporter.Start(&cfg, listeningPort, logger, restServicesController, rpcServicesController)
}
