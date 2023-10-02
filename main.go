/*
Copyright © 2022 Jay Jie

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
	"github.com/jim380/Cendermint/dashboard"
	"github.com/jim380/Cendermint/exporter"
	"github.com/jim380/Cendermint/logging"
	"github.com/jim380/Cendermint/rest"
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
		RpcAddr:          os.Getenv("RPC_ADDR"),
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
	rpcAddr = cfg.RpcAddr
	listeningPort = cfg.ListeningPort
	logOutput = cfg.LogOutput
	logLevel = config.GetLogLevel(cfg.LogLevel)
	logger = logging.InitLogger(logOutput, logLevel)
	zap.ReplaceGlobals(logger)

	cfg.SetSDKConfig()
	rest.RESTAddr = restAddr
	rest.RPCAddr = rpcAddr
	rest.OperAddr = operAddr

	// run dashboard in a separate thread in enabled
	if strings.ToLower(cfg.DashboardEnabled) == "true" {
		dashboard.StartDashboard()
	}

	exporter.Start(&cfg, listeningPort, logger)
}
