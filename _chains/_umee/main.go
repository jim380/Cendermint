package umee

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/jim380/Cosmos-IE/chains/umee/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	rpc "github.com/jim380/Cosmos-IE/chains/umee/getData/rpc"
)

func Main(port string) {
	log, _ := zap.NewDevelopment()
	defer log.Sync()

	rpc.OpenSocket(log)

	http.Handle("/metrics", promhttp.Handler())
	go exporter.Start(log)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("HTTP Handle", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	} else {
		log.Info("HTTP Handle", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("Listen&Serve", "Prometheus Handler(Port: "+port+")"))
	}
}
