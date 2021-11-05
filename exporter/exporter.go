package exporter

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/jim380/Cendermint/rest"
	utils "github.com/jim380/Cendermint/utils"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	defaultGauges []prometheus.Gauge
	gaugesDenom   []prometheus.Gauge
)

func Run(chain string, log *zap.Logger) {
	denomList := getDenomList(chain)

	defaultGauges = make([]prometheus.Gauge, len(gaugesNamespaceList))
	gaugesDenom = make([]prometheus.Gauge, len(denomList)*3)

	// register nomal guages
	for i := 0; i < len(gaugesNamespaceList); i++ {
		defaultGauges[i] = utils.NewGauge("cendermint", gaugesNamespaceList[i], "")
		prometheus.MustRegister(defaultGauges[i])
	}

	// register denom guages
	count := 0
	for i := 0; i < len(denomList)*3; i += 3 {
		gaugesDenom[i] = utils.NewGauge("cendermint_validator_balances", denomList[count], "")
		gaugesDenom[i+1] = utils.NewGauge("cendermint_validator_commission", denomList[count], "")
		gaugesDenom[i+2] = utils.NewGauge("cendermint_validator_rewards", denomList[count], "")
		prometheus.MustRegister(gaugesDenom[i], gaugesDenom[i+1], gaugesDenom[i+2])
		count++
	}

	// register labels
	labelNode := []string{"chain_id", "node_moniker", "node_id", "tm_version", "app_name", "binary_name", "app_version", "git_commit", "go_version", "sdk_version"}
	counterVecNode := utils.NewCounterVec("cendermint", "labels_node_info", "", labelNode)
	labelAddr := []string{"operator_address", "account_address", "cons_address_hex"}
	counterVecAddr := utils.NewCounterVec("cendermint", "labels_addr", "", labelAddr)

	prometheus.MustRegister(counterVecNode, counterVecAddr)

	pollInterval, _ := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	ticker := time.NewTicker(1 * time.Second).C
	// ticker2 := time.NewTicker(40 * time.Second).C

	go func() {
		for {
			var block rest.Blocks
			block.GetInfo()

			currentBlockHeight, _ := strconv.ParseInt(block.Block.Header.Height, 10, 64)
			if previousBlockHeight != currentBlockHeight {
				fmt.Println("")
				zap.L().Info("\t", zap.Bool("Success", true), zap.String("Block Height", fmt.Sprint(currentBlockHeight)))
				select {
				case <-ticker:
					// fetch info from REST
					restData := rest.GetData(chain, currentBlockHeight, block, denomList[0])
					SetMetric(currentBlockHeight, restData, log)
					// case <-ticker2:
					// takes ~5-6 blocks to return results per request
					// tends to halt the node too. Caution !!!
					// restData := rest.GetDelegationsData(chain, currentBlockHeight, block, denomList[0])
					// SetMetric(currentBlockHeight, restData, log)
				}

				metricData := GetMetric()

				setDenomGauges(metricData, denomList)

				setNormalGauges(metricData, defaultGauges)

				setNodeLabels(metricData, counterVecNode)
				setAddrLabels(metricData, counterVecAddr)

				previousBlockHeight = currentBlockHeight
			}
		}
	}()
	time.Sleep(time.Duration(pollInterval) * time.Second)
}
