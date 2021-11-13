package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	utils "github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

var (
	Addr     string
	OperAddr string
	AccAddr  string
)

type RESTData struct {
	BlockHeight   int64
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
}

func (rd RESTData) new(blockHeight int64) *RESTData {
	return &RESTData{
		BlockHeight:   blockHeight,
		Validatorsets: make(map[string][]string),
	}
}

func GetData(chain string, blockHeight int64, blockData Blocks, denom string) *RESTData {
	var restData RESTData
	AccAddr = utils.GetAccAddrFromOperAddr(OperAddr)

	rd := restData.new(blockHeight)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
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
		rd.getCommit(blockData, consHexAddr)
		zap.L().Info("", zap.Bool("Success", true), zap.String("Moniker:", rd.Validators.Description.Moniker))
		zap.L().Info("", zap.Bool("Success", true), zap.String("VP:", valMap[1]))
		zap.L().Info("", zap.Bool("Success", true), zap.String("Precommit:", fmt.Sprintf("%f", rd.Commit.ValidatorPrecommitStatus)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Balances", fmt.Sprint(rd.Balances)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Rewards:", fmt.Sprint(rd.Rewards)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Commission:", fmt.Sprint(rd.Commission)))
		rd.getIBCChannels()
		rd.getIBCConnections()
		rd.getNodeInfo()
		rd.getTxInfo(blockHeight)
		rd.getUpgradeInfo()
		wg.Done()
	}()
	wg.Wait()
	return rd
}

func GetDelegationsData(chain string, blockHeight int64, blockData Blocks, denom string) *RESTData {
	var restData RESTData
	AccAddr = utils.GetAccAddrFromOperAddr(OperAddr)

	rd := restData.new(blockHeight)
	rd.getDelegations()
	return rd
}

func RESTQuery(route string) ([]byte, error) {
	req, err := http.NewRequest("GET", Addr+route, nil)
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
