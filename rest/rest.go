package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jim380/Cendermint/config"
	utils "github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

var (
	RESTAddr string
	RPCAddr  string
	OperAddr string
	AccAddr  string
)

type RESTData struct {
	BlockHeight   int64
	BlockInterval int64
	Commit        commitInfo
	NodeInfo      nodeInfo
	TxInfo        txInfo
	StakingPool   stakingPool
	Slashing      slashingInfo
	Validatorsets map[string][]string
	Validator     validator
	Delegations   delegationsInfo
	Balances      []Coin
	Rewards       []Coin
	Commission    []Coin
	Inflation     float64
	Gov           govInfo
	IBC           struct {
		IBCChannels    map[string][]string
		IBCConnections map[string][]string
		IBCInfo        ibcInfo
	}
	UpgradeInfo upgradeInfo
	GravityInfo gravityInfo
	AkashInfo   akashInfo
}

func (rd RESTData) new(blockHeight int64) *RESTData {
	return &RESTData{
		BlockHeight:   blockHeight,
		Validatorsets: make(map[string][]string),
	}
}

func (rpc RPCData) new() *RPCData {
	return &RPCData{Validatorsets: make(map[string][]string)}
}

func GetData(cfg *config.Config, blockHeight int64, blockData Blocks, denom string) *RESTData {
	// rpc
	var rpcData RPCData
	rpc := rpcData.new()

	// REST
	var restData RESTData
	AccAddr = utils.GetAccAddrFromOperAddr(OperAddr)

	rd := restData.new(blockHeight)
	rd.getNodeInfo(cfg) // get node info first

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		rpc.getConsensusDump(*cfg)
		rd.getStakingPool(*cfg, denom)
		rd.getSlashingParams(*cfg)
		rd.getInflation(*cfg, denom)
		rd.getGovInfo(*cfg)
		rd.getValidatorsets(*cfg, blockHeight)
		rd.getValidator(*cfg)
		// TO-DO if consumer chain, use cosmoshub's ConsPubKey
		valMap, found := rd.Validatorsets[rd.Validator.ConsPubKey.Key]
		if !found {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Validator not found in the active set"))
		}

		rd.getBalances(*cfg)
		rd.getRewardsCommission(*cfg)
		rd.getSigningInfo(*cfg, valMap[0])

		consHexAddr := utils.Bech32AddrToHexAddr(valMap[0])
		rd.getCommit(blockData, consHexAddr)
		zap.L().Info("", zap.Bool("Success", true), zap.String("Moniker", rd.Validator.Description.Moniker))
		zap.L().Info("", zap.Bool("Success", true), zap.String("VP", valMap[1]))
		zap.L().Info("", zap.Bool("Success", true), zap.String("Precommit", fmt.Sprintf("%f", rd.Commit.ValidatorPrecommitStatus)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Balances", fmt.Sprint(rd.Balances)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Rewards", fmt.Sprint(rd.Rewards)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Commission", fmt.Sprint(rd.Commission)))
		rd.getIBCChannels(*cfg)
		rd.getIBCConnections(*cfg)
		rd.getTxInfo(*cfg, blockHeight)
		rd.computerTPS(blockData)
		rd.getUpgradeInfo(*cfg)
		rd.getAkashDeployments(*cfg)
		// gravity
		rd.getBridgeParams(*cfg)
		rd.getValSet(*cfg)
		rd.getOracleEventNonce(*cfg)
		rd.getBatchFees(*cfg)
		rd.getBatchesFees(*cfg)
		rd.getBridgeFees(*cfg)
		wg.Done()
	}()
	wg.Wait()
	return rd
}

func (rd *RESTData) computerTPS(blockData Blocks) {
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

func GetDelegationsData(cfg config.Config, chain string, blockHeight int64, blockData Blocks, denom string) *RESTData {
	var restData RESTData
	AccAddr = utils.GetAccAddrFromOperAddr(OperAddr)

	rd := restData.new(blockHeight)
	rd.getDelegations(cfg)
	return rd
}

func HttpQuery(route string) ([]byte, error) {
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return body, err
}
