package controllers

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"

	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

func (rs RestServices) GetChainData(cfg *config.Config, rpcService RpcServices, blockHeight int64, blockData types.Blocks, denom string) *types.RESTData {
	// rpc
	var rpcData types.RPCData
	rpc := rpcData.New()

	// REST
	var restData types.RESTData
	constants.AccAddr = utils.GetAccAddrFromOperAddr(constants.OperAddr)

	rd := restData.New(blockHeight)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		rpcService.GetRpcInfo(*cfg, rpc)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetNodeInfo(cfg, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetStakingInfo(*cfg, denom, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetSlashingInfo(*cfg, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetInflationInfo(*cfg, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetGovInfo(*cfg, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		valInfo := rpcService.GetValidatorInfo(*cfg, blockHeight, rd)
		rs.GetSigningInfo(*cfg, valInfo[0], rd)
		consHexAddr := utils.Bech32AddrToHexAddr(valInfo[0])
		rs.GetCommitInfo(*cfg, rd, blockData, consHexAddr)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetBalanceInfo(*cfg, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetRewardsCommissionInfo(*cfg, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetIbcChannelInfo(*cfg, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetIbcConnectionInfo(*cfg, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetTxnInfo(*cfg, blockHeight, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		computerTPS(blockData, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetUpgradeInfo(*cfg, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetGravityBridgeInfo(*cfg, rd)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		rs.GetOracleInfo(*cfg, rd)
		wg.Done()
	}()

	wg.Wait()

	zap.L().Info("", zap.Bool("Success", true), zap.String("Moniker", rd.Validator.Description.Moniker))
	zap.L().Info("", zap.Bool("Success", true), zap.String("Precommit", fmt.Sprintf("%f", rd.Commit.ValidatorPrecommitStatus)))
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Balances", fmt.Sprint(rd.Balances)))
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Rewards", fmt.Sprint(rd.Rewards)))
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Commission", fmt.Sprint(rd.Commission)))

	return rd
}

func (rs RestServices) GetAsyncData(cfg *config.Config) *types.AsyncData {
	var data types.AsyncData
	dt := data.New()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		rs.GetAkashInfo(*cfg, dt)
		wg.Done()
	}()

	wg.Wait()

	return dt
}

func computerTPS(blockData types.Blocks, rd *types.RESTData) {
	lastTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", blockData.Block.Header.LastTimestamp)
	currentTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", blockData.Block.Header.Timestamp)
	interval := (currentTimestamp.UnixMilli() - lastTimestamp.UnixMilli()) / 1000 // ms -> s
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Block interval", strconv.Itoa(int(interval))))
	rd.BlockInterval = interval
	totalTxs, _ := strconv.Atoi(rd.TxInfo.Total)
	tps := float64(totalTxs) / float64(interval)
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("TPS", fmt.Sprintf("%.2f", tps)))
	rd.TxInfo.TPS = tps
}
