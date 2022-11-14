package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	utils "github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

var (
	RESTAddr, RESTAddrSputnik, RESTAddrApollo string
	RPCAddr, RPCAddrSputnik, RPCAddrApollo    string
	OperAddr                                  string
	AccAddr                                   string
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

func GetData(chain string, heightProvider, heightSputnik, heightApollo int64, blockData Blocks, denom string) *RESTData {
	// validator set comparison
	var sputnikValSetExistsInProvider, apolloValSetExistsInProvider bool = true, true
	var missingValsInSputnik, missingValsInApollo []string

	// rpc
	var rpcData RPCData
	rpc := rpcData.new()

	// REST
	var restData, restDataSputnik, restDataApollo RESTData
	AccAddr = utils.GetAccAddrFromOperAddr(OperAddr)

	rd := restData.new(heightProvider)
	rdSputnik := restDataSputnik.new(heightSputnik)
	rdApollo := restDataApollo.new(heightApollo)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		rpc.getConsensusDump()
		rd.getStakingPool(denom)
		rd.getSlashingParams()
		rd.getInflation(chain, denom)
		rd.getGovInfo()
		rd.getValidatorsets(heightProvider)
		rdSputnik.getValidatorsets(heightSputnik)
		rdApollo.getValidatorsets(heightApollo)
		// compare validator sets in provider and sputnik
		for kSputnik, valSputnik := range rdSputnik.Validatorsets {
			if _, found := rd.Validatorsets[kSputnik]; !found {
				sputnikValSetExistsInProvider = false
				missingValsInSputnik = append(missingValsInSputnik, valSputnik[0])
			}
		}
		if sputnikValSetExistsInProvider {
			zap.L().Info("", zap.Bool("Success", false), zap.String("------ Validator set in Sputnik exists in Provider -----", ""))
			zap.L().Info("", zap.Bool("Success", false), zap.String("Provider height: ", strconv.Itoa(int(heightProvider))))
			zap.L().Info("", zap.Bool("Success", false), zap.String("Sputnik height: ", strconv.Itoa(int(heightSputnik))))
		} else {
			zap.L().Info("", zap.Bool("Success", false), zap.String("------ Validator set in Sputnik does NOT exist in Provider -----", ""))
			zap.L().Info("", zap.Bool("Success", false), zap.String("Provider height: ", strconv.Itoa(int(heightProvider))))
			zap.L().Info("", zap.Bool("Success", false), zap.String("Sputnik height: ", strconv.Itoa(int(heightSputnik))))
			zap.L().Info("", zap.Bool("Success", false), zap.String("Validators not found in Provider: ", strings.Join(missingValsInSputnik, " ")))
		}
		// compare validator sets in provider and apollo
		for kApollo, valApollo := range rdApollo.Validatorsets {
			if _, found := rd.Validatorsets[kApollo]; !found {
				apolloValSetExistsInProvider = false
				missingValsInApollo = append(missingValsInApollo, valApollo[0])
			}
		}
		if apolloValSetExistsInProvider {
			zap.L().Info("", zap.Bool("Success", false), zap.String("------ Validator set in Apollo exists in Provider -----", ""))
			zap.L().Info("", zap.Bool("Success", false), zap.String("Provider height: ", strconv.Itoa(int(heightProvider))))
			zap.L().Info("", zap.Bool("Success", false), zap.String("Apollo height: ", strconv.Itoa(int(heightApollo))))
		} else {
			zap.L().Info("", zap.Bool("Success", false), zap.String("------ Validator set in Apollo does NOT exist in Provider -----", ""))
			zap.L().Info("", zap.Bool("Success", false), zap.String("Provider height: ", strconv.Itoa(int(heightProvider))))
			zap.L().Info("", zap.Bool("Success", false), zap.String("Apollo height: ", strconv.Itoa(int(heightApollo))))
			zap.L().Info("", zap.Bool("Success", false), zap.String("Validators not found in Provider: ", strings.Join(missingValsInApollo, " ")))
		}

		// ----------------------
		rd.getValidator()
		valMap, found := rd.Validatorsets[rd.Validators.ConsPubKey.Key]
		if !found {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Validator not found in the active set"))
		}
		rd.getBalances()
		rd.getRewardsCommission()
		rd.getSigningInfo(valMap[0])

		consHexAddr := utils.Bech32AddrToHexAddr(valMap[0])
		rd.getCommit(blockData, consHexAddr)
		zap.L().Info("", zap.Bool("Success", true), zap.String("Moniker", rd.Validators.Description.Moniker))
		zap.L().Info("", zap.Bool("Success", true), zap.String("VP", valMap[1]))
		zap.L().Info("", zap.Bool("Success", true), zap.String("Precommit", fmt.Sprintf("%f", rd.Commit.ValidatorPrecommitStatus)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Balances", fmt.Sprint(rd.Balances)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Rewards", fmt.Sprint(rd.Rewards)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Commission", fmt.Sprint(rd.Commission)))
		rd.getIBCChannels()
		rd.getIBCConnections()
		rd.getNodeInfo()
		rd.getTxInfo(heightProvider)
		rd.computerTPS(blockData)
		rd.getUpgradeInfo()
		// gravity
		rd.getBridgeParams()
		rd.getValSet()
		rd.getOracleEventNonce()
		rd.getBatchFees()
		rd.getBatchesFees()
		rd.getBridgeFees()
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

func GetDelegationsData(chain string, blockHeight int64, blockData Blocks, denom string) *RESTData {
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
