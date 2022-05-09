package config

import (
	"fmt"
	"log"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Chain           string
	OperatorAddr    string
	RestAddr        string
	RpcAddr         string
	ListeningPort   string
	MissThreshold   string
	MissConsecutive string
	LogOutput       string
	PollInterval    string
	LogLevel        string
}

func (cfg Config) SetSDKConfig() {
	// Bech32MainPrefix is the common prefix of all prefixes
	Bech32MainPrefix := utils.GetPrefix(cfg.Chain)
	// Bech32PrefixAccAddr is the prefix of account addresses
	Bech32PrefixAccAddr := Bech32MainPrefix
	// Bech32PrefixAccPub is the prefix of account public keys
	Bech32PrefixAccPub := Bech32MainPrefix + sdktypes.PrefixPublic
	// Bech32PrefixValAddr is the prefix of validator operator addresses
	Bech32PrefixValAddr := Bech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixOperator
	// Bech32PrefixValPub is the prefix of validator operator public keys
	Bech32PrefixValPub := Bech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixOperator + sdktypes.PrefixPublic
	// Bech32PrefixConsAddr is the prefix of consensus node addresses
	Bech32PrefixConsAddr := Bech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixConsensus
	// Bech32PrefixConsPub is the prefix of consensus node public keys
	Bech32PrefixConsPub := Bech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixConsensus + sdktypes.PrefixPublic
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.Seal()
}

func (config Config) CheckInputs(chainList map[string][]string) {
	if config.Chain == "" {
		log.Fatal("Chain was not provided.")
	} else {
		valid := false
		if _, found := chainList[config.Chain]; found {
			valid = true
		}
		if !valid {
			log.Fatal(fmt.Sprintf("%s is not supported", config.Chain) + fmt.Sprint("\nList of supported chains: ", chainList))
		}
	}

	// TODO add more robust checks
	if config.OperatorAddr == "" {
		log.Fatal("Operator address was not provided")
	}

	if config.RestAddr == "" {
		log.Fatal("REST address was not provided")
	}

	if config.RpcAddr == "" {
		log.Fatal("RPC address was not provided")
	}

	if config.ListeningPort == "" {
		log.Fatal("Listening port was not provided")
	}

	if config.MissThreshold == "" {
		log.Fatal("Threshold to trigger missing block alerts was not provided")
	}

	if config.MissConsecutive == "" {
		log.Fatal("Threshold to trigger consecutively-missing block alerts was not provided")
	}

	if config.LogOutput == "" {
		log.Fatal("Log output was not provided")
	}

	if config.PollInterval == "" {
		log.Fatal("Poll interval was not provided")
	}

	if config.LogLevel == "" {
		log.Fatal("Log level was not provided")
	}
}

func GetLogLevel(lvl string) zapcore.Level {
	switch lvl {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "log level not supported"))
		return -2
	}
}

func GetDenomList(chain string, chainList map[string][]string) []string {
	var found bool

	for k, v := range chainList {
		if k == chain {
			found = true
			return v
		}
	}
	if !found {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "chain("+chain+") denom not supported"))
	}
	return []string{}
}

func GetChainList() map[string][]string {
	var chainList = map[string][]string{}

	chainList["cosmos"] = []string{"uatom"}
	chainList["iris"] = []string{"uiris"}
	chainList["umee"] = []string{"uumee"}
	chainList["osmosis"] = []string{"uosmo"}
	chainList["juno"] = []string{"ujuno"}
	chainList["akash"] = []string{"uakt"}
	chainList["regen"] = []string{"uregen"}
	chainList["microtick"] = []string{"utick"}
	chainList["nyx"] = []string{"unyx"}
	chainList["evmos"] = []string{"aevmos"}
	chainList["assetMantle"] = []string{"umntl"}
	chainList["rizon"] = []string{"uatolo"}
	chainList["stargaze"] = []string{"ustars"}
	chainList["chihuahua"] = []string{"uhuahua"}
	chainList["gravity"] = []string{"ugraviton"}
	chainList["lum"] = []string{"ulum"}
	chainList["provenance"] = []string{"nhash"}
	chainList["crescent"] = []string{"ucre"}
	chainList["sifchain"] = []string{"urowan"}

	return chainList
}
