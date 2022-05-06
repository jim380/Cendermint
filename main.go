/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"os"

	"github.com/jim380/Cendermint/cmd"
	"github.com/jim380/Cendermint/exporter"
	"github.com/jim380/Cendermint/logging"
	"github.com/jim380/Cendermint/rest"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	chainList                                                 = []string{"cosmos", "umee", "nym", "osmosis", "juno", "akash", "regen", "microtick", "evmos"}
	chain, restAddr, rpcAddr, listenPort, operAddr, logOutput string
	logger                                                    *zap.Logger
)

func main() {
	_ = godotenv.Load("config.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	inputs := []string{os.Getenv("CHAIN"), os.Getenv("OPERATOR_ADDR"), os.Getenv("REST_ADDR"), os.Getenv("RPC_ADDR"), os.Getenv("LISTENING_PORT"), os.Getenv("LOG_OUTPUT"), os.Getenv("POLL_INTERVAL"), os.Getenv("MISS_THRESHOLD"), os.Getenv("MISS_CONSECUTIVE")}
	cmd.CheckInputs(inputs, chainList)

	chain = inputs[0]
	operAddr = inputs[1]
	restAddr = inputs[2]
	rpcAddr = inputs[3]
	listenPort = inputs[4]
	logOutput = inputs[5]

	logger = logging.InitLogger(logOutput)
	zap.ReplaceGlobals(logger)

	cmd.SetSDKConfig(chain)
	rest.RESTAddr = restAddr
	rest.RPCAddr = rpcAddr
	rest.OperAddr = operAddr
	startExporter()
}

func startExporter() {
	exporter.Start(chain, listenPort, logger)
}
