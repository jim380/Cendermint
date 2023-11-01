package controllers

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/models"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type RestServices struct {
	BlockService           *models.BlockService
	TxnService             *models.TxnService
	AbsentValidatorService *models.AbsentValidatorService
	NodeService            *models.NodeService
	StakingService         *models.StakingService
	SlashingService        *models.SlashingService
	InflationService       *models.InflationService
	GovService             *models.GovService
	BankService            *models.BankService
	DelegationService      *models.DelegationService
	UpgradeService         *models.UpgradeService
	IbcServices            *models.IbcService
	GravityService         *models.GravityService
	AkashService           *models.AkashService
	OracleService          *models.OracleService
}

type RpcServices struct {
	ValidatorService *models.ValidatorService
	ConsensusService *models.ConsensusService
}

func (rs RestServices) GetData(cfg *config.Config, rpcService RpcServices, blockHeight int64, blockData types.Blocks, denom string) *types.RESTData {
	// rpc
	var rpcData types.RPCData
	rpc := rpcData.New()

	// REST
	var restData types.RESTData
	constants.AccAddr = utils.GetAccAddrFromOperAddr(constants.OperAddr)

	rd := restData.New(blockHeight)

	rs.GetNodeInfo(cfg, rd)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		rpcService.GetRpcInfo(*cfg, rpc)
		rs.GetStakingInfo(*cfg, denom, rd)
		rs.GetSlashingInfo(*cfg, rd)
		rs.GetInflationInfo(*cfg, rd)
		rs.GetGovInfo(*cfg, rd)
		valInfo := rpcService.GetValidatorInfo(*cfg, blockHeight, rd)
		rs.GetBalanceInfo(*cfg, rd)
		rs.GetRewardsCommissionInfo(*cfg, rd)
		rs.GetSigningInfo(*cfg, valInfo[0], rd)

		consHexAddr := utils.Bech32AddrToHexAddr(valInfo[0])
		rs.GetCommitInfo(*cfg, rd, blockData, consHexAddr)
		zap.L().Info("", zap.Bool("Success", true), zap.String("Moniker", rd.Validator.Description.Moniker))
		zap.L().Info("", zap.Bool("Success", true), zap.String("VP", valInfo[1]))
		zap.L().Info("", zap.Bool("Success", true), zap.String("Precommit", fmt.Sprintf("%f", rd.Commit.ValidatorPrecommitStatus)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Balances", fmt.Sprint(rd.Balances)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Rewards", fmt.Sprint(rd.Rewards)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Commission", fmt.Sprint(rd.Commission)))
		rs.GetIbcChannelInfo(*cfg, rd)
		rs.GetIbcConnectionInfo(*cfg, rd)
		rs.GetTxnInfo(*cfg, blockHeight, rd)
		computerTPS(blockData, rd)
		rs.GetUpgradeInfo(*cfg, rd)
		// akash
		rs.GetAkashInfo(*cfg, rd)
		// gravity
		rs.GetGravityBridgeInfo(*cfg, rd)
		// oracle
		rs.GetOracleInfo(*cfg, rd)

		wg.Done()
	}()
	wg.Wait()
	return rd
}

func computerTPS(blockData types.Blocks, rd *types.RESTData) {
	lastTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", blockData.Block.Header.LastTimestamp)
	currentTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", blockData.Block.Header.Timestamp)
	interval := (currentTimestamp.UnixMilli() - lastTimestamp.UnixMilli()) / 1000 // ms -> s
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Block interval", strconv.Itoa(int(interval))))
	rd.BlockInterval = interval
	totalTxs, _ := strconv.Atoi(rd.TxInfo.Pagination.Total)
	tps := float64(totalTxs) / float64(interval)
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("TPS", fmt.Sprintf("%.2f", tps)))
	rd.TxInfo.TPS = tps
}
