package exporter

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/controllers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(config *config.Config, port string, logger *zap.Logger, restService controllers.RestServices, rpcService controllers.RPCServices) {
	http.Handle("/metrics", promhttp.Handler())
	go Run(config, logger, restService, rpcService)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		zap.L().Fatal("\t", zap.Bool("Success", false), zap.String("HTTP error", fmt.Sprint(err)))
	}
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Serving at port", port))
}

func Run(cfg *config.Config, log *zap.Logger, restService controllers.RestServices, rpcService controllers.RPCServices) {
	denomList := config.GetDenomList(cfg.Chain.Chain, cfg.ChainList)

	registerGauges(denomList)
	counterVecs := registerLabels()

	pollInterval, _ := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	ticker := time.NewTicker(1 * time.Second).C

	go func() {
		for {
			block := restService.GetBlockInfo(*cfg)
			currentBlockHeight, _ := strconv.ParseInt(block.Block.Header.Height, 10, 64)
			if previousBlockHeight != currentBlockHeight {
				block = restService.GetLastBlockTimestamp(*cfg, currentBlockHeight)
				select {
				case <-ticker:
					// fetch data with block info via REST
					restData := restService.GetData(cfg, rpcService, currentBlockHeight, block, denomList[0])
					SetMetric(currentBlockHeight, restData, log)
					// case <-ticker2:
					// takes ~5-6 blocks to return results per request
					// tends to halt the node too. Caution !!!
					// restService.DelegationService.GetInfo(*cfg, restData)
					// SetMetric(currentBlockHeight, restData, log)
				}

				metricData := GetMetric()

				// set gauges
				metricData.setDenomGauges(denomList)
				metricData.setNormalGauges(defaultGauges)

				// set labels
				metricData.setNodeLabels(counterVecs[0])
				metricData.setAddrLabels(counterVecs[1])
				metricData.setUpgradeLabels(counterVecs[2])

				previousBlockHeight = currentBlockHeight
				fmt.Println("--------------------------- End ---------------------------")
				fmt.Println("")
				fmt.Println("")
				fmt.Println("")
			}
		}
	}()
	time.Sleep(time.Duration(pollInterval) * time.Second)
}
