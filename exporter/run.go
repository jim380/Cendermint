package exporter

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(chain string, port string, logger *zap.Logger) {
	http.Handle("/metrics", promhttp.Handler())
	go Run(chain, logger)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Fatal("HTTP Handle", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	} else {
		logger.Info("HTTP Handle", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("Listen&Serve", "Prometheus Handler(Port: "+port+")"))
	}

}
