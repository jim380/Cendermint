package models

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type GravityService struct {
	DB *sql.DB
}

func GetUmeePrice(rd *types.RESTData) {
	var p types.UmeePrice

	res, err := utils.HttpQuery("https://api.coingecko.com/api/v3/simple/price?ids=umee&vs_currencies=usd")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &p); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}

	rd.GravityInfo.UMEEPrice = p.UMEEPrice
}

func (gs *GravityService) GetBatchFees(cfg config.Config, rd *types.RESTData) {
	if !cfg.IsGravityBridgeEnabled() {
		return
	}
	var b types.BatchFees

	route := rest.GetBatchFeesRoute()
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &b); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}

	for _, bf := range b.BatchFees {
		b.Fees += utils.StringToFloat64(bf.TotalFees)
	}

	GetUmeePrice(rd)
	feesTotal := rd.GravityInfo.UmeePrice.UMEEPrice * (b.Fees / 1000000)
	rd.GravityInfo.BatchFees = feesTotal
}

func (gs *GravityService) GetBatchesFees(cfg config.Config, rd *types.RESTData) {
	if !cfg.IsGravityBridgeEnabled() {
		return
	}
	var b types.Batches

	route := rest.GetBatchesFeesRoute()
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &b); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}

	for _, batch := range b.Batches {
		for _, tx := range batch.Transactions {
			b.Fees += utils.StringToFloat64(tx.ERC20Fee.Amount)
		}
	}

	GetUmeePrice(rd)
	feesTotal := rd.GravityInfo.UmeePrice.UMEEPrice * (b.Fees / 1000000)
	rd.GravityInfo.BatchesFees = feesTotal
}

func (gs *GravityService) GetBridgeFees(cfg config.Config, rd *types.RESTData) {
	if !cfg.IsGravityBridgeEnabled() {
		return
	}
	var p types.EthPrice
	var bf float64

	route := rest.GetBridgeFeesRoute()
	res, err := utils.HttpQuery(route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &p); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}

	rd.GravityInfo.ETHPrice = p.ETHPrice
	GetUmeePrice(rd)
	bf = (0.00225 * rd.GravityInfo.ETHPrice) / (100 * rd.GravityInfo.UmeePrice.UMEEPrice)
	rd.GravityInfo.BridgeFees = bf
}

func (gs *GravityService) GetBridgeParams(cfg config.Config, rd *types.RESTData) {
	if !cfg.IsGravityBridgeEnabled() {
		return
	}
	var params types.GravityInfo

	rd.GravityInfo.GravityActive = 0.0
	route := rest.GetBridgeParamsRoute()
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &params); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}

	rd.GravityInfo.BridgeParams = params.BridgeParams

	if params.BridgeActive {
		rd.GravityInfo.GravityActive = 1.0
	} else {
		rd.GravityInfo.GravityActive = 0.0
	}
}

func (gs *GravityService) GetOracleEventNonce(cfg config.Config, rd *types.RESTData) {
	if !cfg.IsGravityBridgeEnabled() {
		return
	}
	var evt types.OracleEventNonce

	orchAddr := os.Getenv("UMEE_ORCH_ADDR")
	route := rest.GetOracleEventNonceByAddressRoute()
	res, err := utils.HttpQuery(constants.RESTAddr + route + orchAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &evt); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}

	rd.GravityInfo.EventNonce = evt.EventNonce
}

func (gs *GravityService) GetValSet(cfg config.Config, rd *types.RESTData) {
	if !cfg.IsGravityBridgeEnabled() {
		return
	}
	var vs types.ValSetInfo

	var vsResult map[string]string = make(map[string]string)

	route := rest.GetCurrentValidatorSetRoute()
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &vs); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
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
