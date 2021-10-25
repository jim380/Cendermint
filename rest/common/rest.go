package rest

import (
	"os/exec"

	"go.uber.org/zap"

	band "github.com/jim380/Cosmos-IE/rest/chains/band"
	terra "github.com/jim380/Cosmos-IE/rest/chains/terra"
	utils "github.com/jim380/Cosmos-IE/utils"
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
	Validator     validator
	Delegations   delegationInfo
	Balances      []Coin
	Rewards       []Coin
	Commission    []Coin
	Inflation     float64
	Gov           govInfo

	Oracle_terra float64
	Oracle_band  float64
}

func newRESTData(blockHeight int64) *RESTData {
	rd := &RESTData{
		BlockHeight:   blockHeight,
		Validatorsets: make(map[string][]string),
	}

	return rd
}

func GetData(chain string, blockHeight int64, blockData Blocks, denom string, log *zap.Logger) *RESTData {
	AccAddr = utils.GetAccAddrFromOperAddr(OperAddr, log)

	rd := newRESTData(blockHeight)
	rd.StakingPool = getStakingPool(denom, log)

	rd.Validatorsets = getValidatorsets(blockHeight, log)
	rd.Validator = getValidators(log)
	/* Block synchronization problem occurs
	   when using "/cosmos/staking/v1beta1/validators/{validator_addr}/delegations" in rest-server
	   after gaiad v4.2.0 */
	// if chain != "cosmos" {
	// 	rd.Delegations = getDelegations(log)
	// }

	rd.Balances = getBalances(AccAddr, log)
	rd.Rewards = getRewards(log)
	rd.Commission = getCommission(log)
	rd.Inflation = getInflation(chain, denom, log)
	rd.Gov = getGovInfo(log)

	consHexAddr := utils.Bech32AddrToHexAddr(rd.Validatorsets[rd.Validator.Consensus_pubkey.Key][0], log)
	rd.Commit = getCommit(blockData, consHexAddr)

	if chain == "band" {
		rd.Oracle_band = band.CheckOracleActive(Addr, OperAddr, log)
	}
	if chain == "terra" {
		rd.Oracle_terra = terra.GetOracleMiss(Addr, OperAddr, log)
	}

	return rd
}

func runRESTCommand(str string) ([]uint8, error) {
	cmd := "curl -s -XGET " + Addr + str + " -H \"accept:application/json\""
	// log.Println("runRestCommand: ", cmd);
	out, err := exec.Command("/bin/bash", "-c", cmd).Output()

	return out, err
}
