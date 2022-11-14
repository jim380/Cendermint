package exporter

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/rest"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(chain string, port string, logger *zap.Logger) {
	http.Handle("/metrics", promhttp.Handler())
	go Run(chain, logger)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		zap.L().Fatal("\t", zap.Bool("Success", false), zap.String("HTTP error", fmt.Sprint(err)))
	}
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Serving at port", port))

}

func Run(chain string, log *zap.Logger) {
	cl := config.GetChainList()
	denomList := config.GetDenomList(chain, cl)

	registerGauges(denomList)
	counterVecs := registerLabels()

	pollInterval, _ := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	ticker := time.NewTicker(1 * time.Second).C
	// ticker2 := time.NewTicker(40 * time.Second).C

	go func() {
		for {
			var block, blockSputnik, blockApollo rest.Blocks
			block.GetInfo(rest.RESTAddr)
			blockSputnik.GetInfo(rest.RESTAddrSputnik)
			blockApollo.GetInfo(rest.RESTAddrApollo)

			currentHeightProvider, _ := strconv.ParseInt(block.Block.Header.Height, 10, 64)
			currentHeightSputnik, _ := strconv.ParseInt(blockSputnik.Block.Header.Height, 10, 64)
			currentHeightApollo, _ := strconv.ParseInt(blockApollo.Block.Header.Height, 10, 64)

			if previousBlockHeight != currentHeightProvider {
				fmt.Println("--------------------------- Start ---------------------------")
				block.GetLastBlockTimestamp(currentHeightProvider)
				zap.L().Info("\t", zap.Bool("Success", true), zap.String("Last block timestamp", block.Block.Header.LastTimestamp))
				zap.L().Info("\t", zap.Bool("Success", true), zap.String("Current block timestamp", block.Block.Header.Timestamp))
				zap.L().Info("\t", zap.Bool("Success", true), zap.String("Provider block height", fmt.Sprint(currentHeightProvider)))
				zap.L().Info("\t", zap.Bool("Success", true), zap.String("Sputnik block height", fmt.Sprint(currentHeightSputnik)))
				zap.L().Info("\t", zap.Bool("Success", true), zap.String("Apollo block height", fmt.Sprint(currentHeightApollo)))
				select {
				case <-ticker:
					// fetch info from REST
					restData := rest.GetData(chain, currentHeightProvider, currentHeightSputnik, currentHeightApollo, block, denomList[0])
					SetMetric(currentHeightProvider, restData, log)
					// case <-ticker2:
					// takes ~5-6 blocks to return results per request
					// tends to halt the node too. Caution !!!
					// restData := rest.GetDelegationsData(chain, currentBlockHeight, block, denomList[0])
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

				previousBlockHeight = currentHeightProvider
				fmt.Println("--------------------------- End ---------------------------")
				fmt.Println("")
				fmt.Println("")
				fmt.Println("")
			}
		}
	}()
	time.Sleep(time.Duration(pollInterval) * time.Second)
}
