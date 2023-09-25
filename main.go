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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/exporter"
	"github.com/jim380/Cendermint/logging"
	"github.com/jim380/Cendermint/rest"
	"github.com/joho/godotenv"
	"github.com/kyoto-framework/kyoto/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	chain, restAddr, rpcAddr, listeningPort, operAddr, logOutput string
	logLevel                                                     zapcore.Level
	logger                                                       *zap.Logger
)

type PIndexState struct {
	Block *kyoto.ComponentF[rest.Blocks]
}

/*
Page
  - A page is a top-level component, which attaches components and
    defines rendering
*/
func PIndex(ctx *kyoto.Context) (state PIndexState) {
	// Define rendering
	kyoto.Template(ctx, "page.index.html")

	// Attach components
	state.Block = kyoto.Use(ctx, GetBlockInfo)

	return
}

/*
Component
  - Each component is a context receiver, which returns its state
  - Each component becomes a part of the page or top-level component,
    which executes component asynchronously and gets a state future object
  - Context holds common objects like http.ResponseWriter, *http.Request, etc
*/
func GetBlockInfo(ctx *kyoto.Context) (state rest.Blocks) {
	route := "/cosmos/base/tendermint/v1beta1/blocks/latest" //TO-DO refactor this
	fetchBlockInfo := func() rest.Blocks {
		var state rest.Blocks
		resp, err := rest.HttpQuery(rest.RESTAddr + route)
		if err != nil {
			zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return rest.Blocks{}
		}

		err = json.Unmarshal(resp, &state)
		if err != nil {
			zap.L().Fatal("Failed to unmarshal response", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return rest.Blocks{}
		}

		return state
	}

	/*
		Handle Actions
			- To call an action of parent component, use $ prefix in action name
			- To call an action of component by id, use <id:action> as an action name
		    - To push multiple component UI updates during a single action call,
		        call kyoto.ActionFlush(ctx, state) to initiate an update
	*/
	handled := kyoto.Action(ctx, "Reload Block", func(args ...any) {
		// add logic here
		state = fetchBlockInfo()
		log.Println("New block info fetched on block", state.Block.Header.Height)
	})
	// Prevent further execution if action handled
	if handled {
		return
	}
	// Default loading behavior if not handled
	state = fetchBlockInfo()

	return
}

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
		go func() {
			port := os.Getenv("DASHBOARD_PORT")
			// Register page
			kyoto.HandlePage("/", PIndex)
			// Client
			kyoto.HandleAction(GetBlockInfo)
			// Serve
			kyoto.Serve(":" + port)
		}()
	}

	exporter.Start(&cfg, listeningPort, logger)
}
