package fetcher

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jim380/Cendermint/cache"
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/controllers"
	"github.com/jim380/Cendermint/exporter"
	"go.uber.org/zap"
)

func Start(cfg *config.Config, restService controllers.RestServices, rpcService controllers.RpcServices, denomList []string, log *zap.Logger) {
	go FetchChainData(cfg, restService, rpcService, denomList, log)
	go FetchAsyncData(cfg, restService, log)
}

func FetchChainData(cfg *config.Config, restService controllers.RestServices, rpcService controllers.RpcServices, denomList []string, log *zap.Logger) {
	for {
		block := restService.GetBlockInfo(*cfg)
		currentBlockHeight, _ := strconv.ParseInt(block.Block.Header.Height, 10, 64)

		lastCachedHeight, err := cache.Get("lastCachedHeight")
		if err != nil {
			lastCachedHeight = 0
		}

		shouldPoll := lastCachedHeight != currentBlockHeight

		if shouldPoll {
			fmt.Println("--------------------------- Started fetching chain data ---------------------------")
			zap.L().Info("\t", zap.Bool("Success", true), zap.String("Last block timestamp", block.Block.Header.LastTimestamp))
			zap.L().Info("\t", zap.Bool("Success", true), zap.String("Current block timestamp", block.Block.Header.Timestamp))
			zap.L().Info("\t", zap.Bool("Success", true), zap.String("Current block height", fmt.Sprint(currentBlockHeight)))

			block = restService.GetLastBlockTimestamp(*cfg, currentBlockHeight)

			// fetch data with block info via REST
			restData := restService.GetChainData(cfg, rpcService, currentBlockHeight, block, denomList[0])
			exporter.SetMetricChain(currentBlockHeight, restData, log)
			// case <-ticker2:
			// takes ~5-6 blocks to return results per request
			// tends to halt the node too. Caution !!!
			// restService.DelegationService.GetInfo(*cfg, restData)
			// SetMetricSetMetricChain(currentBlockHeight, restData, log)

			err := cache.Set("lastCachedHeight", currentBlockHeight)
			if err != nil {
				panic(err)
			}

			fmt.Println("--------------------------- Finished fetching chain data ---------------------------")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
		}

		time.Sleep(time.Duration(constants.PollIntervalChain) * time.Second)
	}
}

func FetchAsyncData(cfg *config.Config, restService controllers.RestServices, log *zap.Logger) {
	for {
		fmt.Println("--------------------------- Started fetching async data ---------------------------")
		data := restService.GetAsyncData(cfg)
		exporter.SetMetricAsync(data, log)
		fmt.Println("--------------------------- Finished fetching async data ---------------------------")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")

		time.Sleep(time.Duration(constants.PollIntervalAsync) * time.Second)
	}
}

func BackfillData(cfg *config.Config, restService controllers.RestServices, log *zap.Logger) {
	for {
		fmt.Println("--------------------------- Started backfilling data ---------------------------")
		restService.BackfillData(cfg)
		fmt.Println("--------------------------- Finished backfilling data ---------------------------")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")

		time.Sleep(time.Duration(constants.PollIntervalAsync) * time.Second)
	}
}
