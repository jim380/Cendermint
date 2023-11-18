package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/jim380/Cendermint/logging"
	"github.com/jim380/Cendermint/types"

	"github.com/jim380/Cendermint/utils"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Chain            Chain
	ChainList        map[string][]string
	SDKVersion       string
	OperatorAddr     string
	RestAddr         string
	RpcAddr          string
	ListeningPort    string
	MissThreshold    string
	MissConsecutive  string
	LogOutput        string
	PollInterval     string
	LogLevel         string
	DashboardEnabled string
}

type Chain struct {
	Name   string `json:"chain"`
	Assets []struct {
		Denom string `json:"denom"`
	} `json:"assets"`
}

func (cfg Config) SetSDKConfig() {
	// Bech32MainPrefix is the common prefix of all prefixes
	Bech32MainPrefix := utils.GetPrefix(cfg.Chain.Name)
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
	// TO-DO add more robust checks
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
	jsonFile, err := os.Open("chains.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var chains []Chain
	json.Unmarshal(byteValue, &chains)

	chainList := make(map[string][]string)
	for _, chain := range chains {
		for _, asset := range chain.Assets {
			chainList[chain.Name] = append(chainList[chain.Name], asset.Denom)
		}
	}

	return chainList
}

func (config Config) IsLegacySDKVersion() bool {
	var legacy bool = false

	if strings.Contains(config.SDKVersion, "0.45") {
		legacy = true
	}

	return legacy
}

func (config Config) IsGravityBridgeEnabled() bool {
	var enabled bool = false

	if config.Chain.Name == "gravity" || config.Chain.Name == "umee" {
		enabled = true
	}

	return enabled
}

func LoadConfig() Config {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("CHAIN") == "" {
		log.Fatal("Chain was not provided.")
	}

	cfg := Config{
		Chain:            Chain{Name: os.Getenv("CHAIN")},
		OperatorAddr:     os.Getenv("OPERATOR_ADDR"),
		RestAddr:         os.Getenv("REST_ADDR"),
		RpcAddr:          os.Getenv("RPC_ADDR"),
		ListeningPort:    os.Getenv("LISTENING_PORT"),
		MissThreshold:    os.Getenv("MISS_THRESHOLD"),
		MissConsecutive:  os.Getenv("MISS_CONSECUTIVE"),
		LogOutput:        os.Getenv("LOG_OUTPUT"),
		PollInterval:     os.Getenv("POLL_INTERVAL"),
		LogLevel:         os.Getenv("LOG_LEVEL"),
		DashboardEnabled: os.Getenv("DASHBOARD_ENABLED"),
	}

	return cfg
}

func (cfg *Config) ValidateConfig() types.AppConfig {
	chainList := GetChainList()
	cfg.ChainList = chainList
	supportedChains := make([]string, 0, len(chainList))
	for key := range chainList {
		supportedChains = append(supportedChains, key)
	}
	var found bool
	if _, found = chainList[cfg.Chain.Name]; found {
		cfg.Chain = Chain{Name: cfg.Chain.Name}
	}
	if !found {
		log.Fatal(fmt.Sprintf("%s is not supported", cfg.Chain.Name) + fmt.Sprint("\nList of supported chains: ", supportedChains))
	}

	cfg.CheckInputs(chainList)

	appConfig := types.AppConfig{
		Chain:         cfg.Chain.Name,
		OperAddr:      cfg.OperatorAddr,
		RestAddr:      cfg.RestAddr,
		RpcAddr:       cfg.RpcAddr,
		ListeningPort: cfg.ListeningPort,
		LogOutput:     cfg.LogOutput,
		LogLevel:      GetLogLevel(cfg.LogLevel),
		Logger:        logging.InitLogger(cfg.LogOutput, GetLogLevel(cfg.LogLevel)),
	}

	return appConfig
}
