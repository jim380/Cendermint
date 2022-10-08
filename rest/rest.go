package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

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
	Validators    validator
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
	OracleInfo  OracleInfo
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

func GetData(chain string, blockHeight int64, blockData SDKBlook, denom string) *RESTData {
	// rpc
	var rpcData RPCData
	rpc := rpcData.new()

	// REST
	var restData RESTData
	AccAddr = utils.GetAccAddrFromOperAddr(OperAddr)

	rd := restData.new(blockHeight)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		rpc.getConsensusDump()
		rd.getStakingPool(denom)
		rd.getSlashingParams()
		rd.getInflation(chain, denom)
		rd.getGovInfo()
		rd.getValidatorsets(blockHeight)
		rd.getValidator()
		valMap, found := rd.Validatorsets[rd.Validators.ConsPubKey.Key]
		if !found {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Validator not found in the active set"))
		}
		rd.getBalances()
		rd.getRewardsCommission()
		rd.getSigningInfo(valMap[0])

		consHexAddr := utils.Bech32AddrToHexAddr(valMap[0])
		valAddr := utils.HexToBase64(consHexAddr)
		operatorAddr := os.Getenv("OPERATOR_ADDR")
		rd.getCommit(blockData, operatorAddr, valAddr)
		zap.L().Info("", zap.Bool("Success", true), zap.String("Moniker", rd.Validators.Description.Moniker))
		zap.L().Info("", zap.Bool("Success", true), zap.String("VP", valMap[1]))
		zap.L().Info("", zap.Bool("Success", true), zap.String("Precommit", fmt.Sprintf("%f", rd.Commit.ValidatorPrecommitStatus)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Balances", fmt.Sprint(rd.Balances)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Rewards", fmt.Sprint(rd.Rewards)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Commission", fmt.Sprint(rd.Commission)))
		rd.getIBCChannels()
		rd.getIBCConnections()
		rd.getNodeInfo()
		rd.getTxInfo(blockHeight)
		rd.computerTPS(blockData)
		rd.getUpgradeInfo()
		// gravity
		rd.getBridgeParams()
		rd.getValSet()
		rd.getOracleEventNonce()
		rd.getBatchFees()
		rd.getBatchesFees()
		rd.getBridgeFees()

		// oracle
		rd.getOracleMissesCount()
		rd.getOracleSubmitBlock()
		rd.getOracleFeederDelegate()
		wg.Done()
	}()
	wg.Wait()
	return rd
}

func (rd *RESTData) computerTPS(blockData SDKBlook) {
	// lastTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", blockData.Block.Header.LastTimestamp)
	// currentTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", blockData.Block.Header.Timestamp)
	// interval := (currentTimestamp.UnixMilli() - lastTimestamp.UnixMilli()) / 1000 // ms -> s
	// zap.L().Info("\t", zap.Bool("Success", true), zap.String("Block interval", strconv.Itoa(int(interval))))
	// rd.BlockInterval = interval
	// totalTxs, _ := strconv.Atoi(rd.TxInfo.Pagination.Total)
	// tps := float64(totalTxs) / float64(interval)
	// zap.L().Info("\t", zap.Bool("Success", true), zap.String("TPS", fmt.Sprintf("%.2f", tps)))
	// rd.TxInfo.TPS = tps
}

func GetDelegationsData(chain string, blockHeight int64, blockData SDKBlook, denom string) *RESTData {
	var restData RESTData
	AccAddr = utils.GetAccAddrFromOperAddr(OperAddr)

	rd := restData.new(blockHeight)
	rd.getDelegations()
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
