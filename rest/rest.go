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
		rd.getBalances()
		rd.getRewardsCommission()
		rd.getSigningInfo(rd.Validatorsets[rd.Validators.ConsPubKey.Key][0])

		consHexAddr := utils.Bech32AddrToHexAddr(rd.Validatorsets[rd.Validators.ConsPubKey.Key][0])
		rd.getCommit(blockData, consHexAddr)
		zap.L().Info("", zap.Bool("Success", true), zap.String("Moniker:", rd.Validators.Description.Moniker))
		zap.L().Info("", zap.Bool("Success", true), zap.String("VP:", rd.Validatorsets[rd.Validators.ConsPubKey.Key][1]))
		zap.L().Info("", zap.Bool("Success", true), zap.String("Precommit:", fmt.Sprintf("%f", rd.Commit.ValidatorPrecommitStatus)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Balances", fmt.Sprint(rd.Balances)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Rewards:", fmt.Sprint(rd.Rewards)))
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Commission:", fmt.Sprint(rd.Commission)))
		rd.getIBCChannels()
		rd.getIBCConnections()
		// takes ~5-6 blocks to return results per request
		// might halt the node. Caution !!!
		// rd.getDelegations()
		wg.Done()
	}()
	wg.Wait()
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
