package rest

import (
	"encoding/json"
	"os"

	"github.com/jim380/Cendermint/config"
	utils "github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type gravityInfo struct {
	gravityParams
	oracleEventNonce
	ValSetCount   int
	ValActive     float64 // [0]: false, [1]: true
	GravityActive float64 // [0]: false, [1]: true
	umeePrice
	ethPrice
	BatchFees   float64
	BatchesFees float64
	BridgeFees  float64
}

type gravityParams struct {
	BridgeParams `json:"params"`
}

type BridgeParams struct {
	SignedValsetsWindow    string `json:"signed_valsets_window"`
	SignedBatchesWindow    string `json:"signed_batches_window"`
	TargetBatchTimeout     string `json:"target_batch_timeout"`
	SlashFractionValset    string `json:"slash_fraction_valset"`
	SlashFractionBatch     string `json:"slash_fraction_batch"`
	SlashFractionBadEthSig string `json:"slash_fraction_bad_eth_signature"`
	ValsetReward           `json:"valset_reward"`
	BridgeActive           bool `json:"bridge_active"`
}

type ValsetReward struct {
	Amount string `json:"amount"`
}

type ethPrice struct {
	ETHUSD `json:"ethereum"`
}

type ETHUSD struct {
	ETHPrice float64 `json:"usd"`
}

type umeePrice struct {
	UMEEUSD `json:"umee"`
}

type UMEEUSD struct {
	UMEEPrice float64 `json:"usd"`
}

type batches struct {
	Batches []batch `json:"batches"`
	Fees    float64
}

type batch struct {
	BatchNonce   string        `json:"batch_nonce"`
	BatchTimeout string        `json:"batch_timeout"`
	Transactions []transaction `json:"transactions"`
}

type transaction struct {
	ID         string     `json:"id"`
	Sender     string     `json:"sender"`
	DestAddr   string     `json:"dest_address"`
	ERC20Token erc20Token `json:"erc20_token"`
	ERC20Fee   erc20Fee   `json:"erc20_fee"`
}

type erc20Token struct {
	Contract string `json:"contract"`
	Amount   string `json:"amount"`
}

type erc20Fee struct {
	Contract string `json:"contract"`
	Amount   string `json:"amount"`
}

type batchFees struct {
	BatchFees []batchFee `json:"batchFees"`
	Fees      float64
}

type batchFee struct {
	Token     string `json:"token"`
	TotalFees string `json:"total_fees"`
}

type oracleEventNonce struct {
	EventNonce string `json:"event_nonce"`
}

type lastClaimEvent struct {
	EventNonce  string `json:"ethereum_event_nonce"`
	EventHeight string `json:"ethereum_event_height"`
}

type valSetInfo struct {
	ValSet   valSet `json:"valset"`
	ValCount int
}

type valSet struct {
	Nonce   string   `json:"nonce"`
	Members []member `json:"members"`
}

type member struct {
	Power   string `json:"power"`
	ETHAddr string `json:"ethereum_address"`
}

func (rd *RESTData) getUmeePrice() {
	var p umeePrice

	res, err := HttpQuery("https://api.coingecko.com/api/v3/simple/price?ids=umee&vs_currencies=usd")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if err := json.Unmarshal(res, &p); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Bool("Success", false), zap.String("err", err.Error()))
		return
	}

	rd.GravityInfo.UMEEPrice = p.UMEEPrice
}

func (rd *RESTData) getBatchFees(cfg config.Config) {
	if !cfg.IsGravityBridgeEnabled() {
		return
	}
	var b batchFees

	route := getBatchFeesRoute()
	res, err := HttpQuery(RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if err := json.Unmarshal(res, &b); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Bool("Success", false), zap.String("err", err.Error()))
		return
	}

	for _, bf := range b.BatchFees {
		b.Fees += utils.StringToFloat64(bf.TotalFees)
	}

	rd.getUmeePrice()
	feesTotal := rd.GravityInfo.umeePrice.UMEEPrice * (b.Fees / 1000000)
	rd.GravityInfo.BatchFees = feesTotal
}

func (rd *RESTData) getBatchesFees(cfg config.Config) {
	if !cfg.IsGravityBridgeEnabled() {
		return
	}
	var b batches

	route := getBatchesFeesRoute()
	res, err := HttpQuery(RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if err := json.Unmarshal(res, &b); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Bool("Success", false), zap.String("err", err.Error()))
		return
	}

	for _, batch := range b.Batches {
		for _, tx := range batch.Transactions {
			b.Fees += utils.StringToFloat64(tx.ERC20Fee.Amount)
		}
	}

	rd.getUmeePrice()
	feesTotal := rd.GravityInfo.umeePrice.UMEEPrice * (b.Fees / 1000000)
	rd.GravityInfo.BatchesFees = feesTotal
}

func (rd *RESTData) getBridgeFees(cfg config.Config) {
	if !cfg.IsGravityBridgeEnabled() {
		return
	}
	var p ethPrice
	var bf float64

	route := getBridgeFeesRoute()
	res, err := HttpQuery(route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if err := json.Unmarshal(res, &p); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Bool("Success", false), zap.String("err", err.Error()))
		return
	}

	rd.GravityInfo.ETHPrice = p.ETHPrice
	rd.getUmeePrice()
	bf = (0.00225 * rd.GravityInfo.ETHPrice) / (100 * rd.GravityInfo.umeePrice.UMEEPrice)
	rd.GravityInfo.BridgeFees = bf
}

func (rd *RESTData) getBridgeParams(cfg config.Config) {
	if !cfg.IsGravityBridgeEnabled() {
		return
	}
	var params gravityParams

	rd.GravityInfo.GravityActive = 0.0
	route := getBridgeParamsRoute()
	res, err := HttpQuery(RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if err := json.Unmarshal(res, &params); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Bool("Success", false), zap.String("err", err.Error()))
		return
	}

	rd.GravityInfo.BridgeParams = params.BridgeParams

	if params.BridgeActive {
		rd.GravityInfo.GravityActive = 1.0
	} else {
		rd.GravityInfo.GravityActive = 0.0
	}
}

func (rd *RESTData) getOracleEventNonce(cfg config.Config) {
	if !cfg.IsGravityBridgeEnabled() {
		return
	}
	var evt oracleEventNonce

	orchAddr := os.Getenv("UMEE_ORCH_ADDR")
	route := getOracleEventNonceByAddressRoute()
	res, err := HttpQuery(RESTAddr + route + orchAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if err := json.Unmarshal(res, &evt); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Bool("Success", false), zap.String("err", err.Error()))
		return
	}

	rd.GravityInfo.EventNonce = evt.EventNonce
}

func (rd *RESTData) getValSet(cfg config.Config) {
	if !cfg.IsGravityBridgeEnabled() {
		return
	}
	var vs valSetInfo

	var vsResult map[string]string = make(map[string]string)

	route := getCurrentValidatorSetRoute()
	res, err := HttpQuery(RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if err := json.Unmarshal(res, &vs); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Bool("Success", false), zap.String("err", err.Error()))
		return
	}

	for _, member := range vs.ValSet.Members {
		vsResult[member.ETHAddr] = member.Power
	}

	rd.GravityInfo.ValSetCount = len(vs.ValSet.Members)

	_, found := vsResult[os.Getenv("ETH_ORCH_ADDR")]
	if found {
		rd.GravityInfo.ValActive = 1.0
	}
}
