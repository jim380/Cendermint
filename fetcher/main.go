package fetcher

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/controllers"
	"github.com/jim380/Cendermint/exporter"
	"go.uber.org/zap"
)

// TO-DO store in db
var previousBlockHeight int64

func Start(cfg *config.Config, restService controllers.RestServices, rpcService controllers.RpcServices, mutex *sync.Mutex, ticker <-chan time.Time, denomList []string, log *zap.Logger) {
	go FetchRESTData(cfg, restService, rpcService, mutex, ticker, denomList, log)
}

func FetchRESTData(cfg *config.Config, restService controllers.RestServices, rpcService controllers.RpcServices, mutex *sync.Mutex, ticker <-chan time.Time, denomList []string, log *zap.Logger) {
	for {
		block := restService.GetBlockInfo(*cfg)
		currentBlockHeight, _ := strconv.ParseInt(block.Block.Header.Height, 10, 64)

		fmt.Println("--------------------------- Started fetching REST data ---------------------------")
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Last block timestamp", block.Block.Header.LastTimestamp))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Current block timestamp", block.Block.Header.Timestamp))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Current block height", fmt.Sprint(currentBlockHeight)))

		mutex.Lock()
		if previousBlockHeight != currentBlockHeight {
			block = restService.GetLastBlockTimestamp(*cfg, currentBlockHeight)
			select {
			case <-ticker:
				// fetch data with block info via REST
				restData := restService.GetData(cfg, rpcService, currentBlockHeight, block, denomList[0])

				exporter.SetMetric(currentBlockHeight, restData, log)
				// case <-ticker2:
				// takes ~5-6 blocks to return results per request
				// tends to halt the node too. Caution !!!
				// restService.DelegationService.GetInfo(*cfg, restData)
				// SetMetric(currentBlockHeight, restData, log)
			}
			previousBlockHeight = currentBlockHeight
			fmt.Println("--------------------------- Finished fetching REST data ---------------------------")
			fmt.Println("")
			fmt.Println("")
			fmt.Println("")
		}
		mutex.Unlock()

		time.Sleep(time.Duration(constants.PollInterval) * time.Second)
	}
}
