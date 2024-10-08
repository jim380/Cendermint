package exporter

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/controllers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(config *config.Config, port string, logger *zap.Logger, restService controllers.RestServices, rpcService controllers.RpcServices) {
	go CollectMetrics(config, logger, restService, rpcService)
}

func CollectMetrics(cfg *config.Config, log *zap.Logger, restService controllers.RestServices, rpcService controllers.RpcServices) {
	denomList, err := config.GetDenomList(cfg.Chain.Name, cfg.ChainList)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	registerGauges(denomList)
	counterVecs := registerLabels()

	for {
		metricData := GetMetric()
		if metricData == nil {
			continue
		}

		// set gauges
		metricData.setDenomGauges(denomList)
		metricData.setNormalGauges(defaultGauges)

		// set labels
		metricData.setNodeLabels(counterVecs[0])
		metricData.setAddrLabels(counterVecs[1])
		metricData.setUpgradeLabels(counterVecs[2])

		time.Sleep(time.Duration(constants.PollIntervalChain) * time.Second)
	}
}

func StartMetricsHttpServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		zap.L().Fatal("\t", zap.Bool("Success", false), zap.String("HTTP error", fmt.Sprint(err)))
	}
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Serving at port", port))
}
