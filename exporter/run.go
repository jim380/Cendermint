package exporter

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	sdk "github.com/cosmos/cosmos-sdk/types"
	iris "github.com/irisnet/irishub/address"
	"github.com/jim380/Cosmos-IE/common"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Go(chain string, port string) {
	log, _ := zap.NewDevelopment()
	defer log.Sync()

	setConfig(chain)

	http.Handle("/metrics", promhttp.Handler())
	go Start(chain, log)

	err := http.ListenAndServe(":"+port, nil)
	// log
	if err != nil {
		// handle error
		log.Fatal("HTTP Handle", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	} else {
		log.Info("HTTP Handle", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("Listen&Serve", "Prometheus Handler(Port: "+port+")"))
	}

}

// set custom configs
func setConfig(chain string) {
	config := sdk.GetConfig()

	switch chain {
	case "iris":
		iris.ConfigureBech32Prefix()
	case "umee":
		config := sdk.GetConfig()
		config.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)
		config.SetBech32PrefixForValidator(common.Bech32PrefixValAddr, common.Bech32PrefixValPub)
		config.SetBech32PrefixForConsensusNode(common.Bech32PrefixConsAddr, common.Bech32PrefixConsPub)
		config.Seal()
	}

	config.Seal()
}
