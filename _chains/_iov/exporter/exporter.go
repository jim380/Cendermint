package exporter

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	metric "github.com/jim380/Cosmos-IE/chains/iov/exporter/metric"
	rpc "github.com/jim380/Cosmos-IE/chains/iov/getData/rpc"
	utils "github.com/jim380/Cosmos-IE/utils"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	previousBlockHeight int64
)

func Start(log *zap.Logger) {

	gaugesNamespaceList := metric.GaugesNamespaceList

	var gauges []prometheus.Gauge = make([]prometheus.Gauge, len(gaugesNamespaceList))

	// nomal guages
	for i := 0; i < len(gaugesNamespaceList); i++ {
		gauges[i] = utils.NewGauge("exporter", gaugesNamespaceList[i], "")
		prometheus.MustRegister(gauges[i])
	}

	// labels
	labels := []string{"chainId"}
	gaugesForLabel := utils.NewCounterVec("exporter", "labels", "", labels)

	prometheus.MustRegister(gaugesForLabel)

	for {
		func() {
			defer func() {

				if r := recover(); r != nil {
					//Error Log
				}

				time.Sleep(500 * time.Millisecond)

			}()

			currentBlockHeight := rpc.BlockHeight()

			if previousBlockHeight != currentBlockHeight {

				fmt.Println("")
				log.Info("RPC-Server", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("Get Data", "Block Height: "+fmt.Sprint(currentBlockHeight)))

				rpcData := rpc.GetData(currentBlockHeight, rpc.ConsHexAddr, log)

				metric.SetMetric(currentBlockHeight, rpcData, log)

				metricData := metric.GetMetric()

				gaugesValue := [...]float64{
					float64(metricData.Network.BlockHeight),

					metricData.Validator.Commit.VoteType,
					metricData.Validator.Commit.PrecommitStatus,
				}

				for i := 0; i < len(gaugesNamespaceList); i++ {
					gauges[i].Set(gaugesValue[i])
				}

				gaugesForLabel.WithLabelValues(metricData.Network.ChainID).Add(0)
			}

			previousBlockHeight = currentBlockHeight
		}()
	}
}
