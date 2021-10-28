package rest

import (
	"fmt"
	"os/exec"

	utils "github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

var (
	Addr     string
	OperAddr string
	AccAddr  string
)

type RESTData struct {
	BlockHeight int64
	Commit      commitInfo
	StakingPool stakingPool

	Validatorsets map[string][]string
	Validators    validator
	Delegations   delegationInfo
	Balances      []Coin
	Rewards       []Coin
	Commission    []Coin
	Inflation     float64
	Gov           govInfo
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
	rd.getStakingPool(denom)
	rd.getInflation(chain, denom)
	rd.getGovInfo()
	rd.getValidatorsets(blockHeight)
	rd.getValidator()
	rd.getBalances()
	rd.getRewardsCommission()

	consHexAddr := utils.Bech32AddrToHexAddr(rd.Validatorsets[rd.Validators.ConsPubKey.Value][0])
	rd.getCommit(blockData, consHexAddr)
	zap.L().Info("", zap.Bool("Success", true), zap.String("Moniker:", rd.Validators.Description.Moniker))
	zap.L().Info("", zap.Bool("Success", true), zap.String("VP:", rd.Validatorsets[rd.Validators.ConsPubKey.Value][1]))
	zap.L().Info("", zap.Bool("Success", true), zap.String("Precommit:", "signed"))
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Balances", fmt.Sprint(rd.Balances)))
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Rewards:", fmt.Sprint(rd.Rewards)))
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Commission:", fmt.Sprint(rd.Commission)))

	return rd
}

func runRESTCommand(str string) ([]uint8, error) {
	cmd := "curl -s -XGET " + Addr + str + " -H \"accept:application/json\""
	out, err := exec.Command("/bin/bash", "-c", cmd).Output()

	return out, err
}
