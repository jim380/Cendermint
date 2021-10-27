package rest

import (
	"os/exec"

	"github.com/jim380/Cosmos-IE/rpc"
	utils "github.com/jim380/Cosmos-IE/utils"
)

var (
	Addr     string
	OperAddr string
	AccAddr  string
)

type RESTData struct {
	BlockHeight int64
	Commit      rpc.CommitInfo
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
	rd.getValidatorsets(blockHeight)
	rd.getValidators()
	/* Block synchronization problem occurs
	   when using "/cosmos/staking/v1beta1/validators/{validator_addr}/delegations" in rest-server
	   after gaiad v4.2.0 */
	// if chain != "cosmos" {
	// 	rd.Delegations = getDelegations(log)
	// }
	rd.getBalances()
	rd.getRewards()
	rd.getCommission()
	rd.getInflation(chain, denom)
	rd.getGovInfo()

	return rd
}

func runRESTCommand(str string) ([]uint8, error) {
	cmd := "curl -s -XGET " + Addr + str + " -H \"accept:application/json\""
	// log.Println("runRestCommand: ", cmd);
	out, err := exec.Command("/bin/bash", "-c", cmd).Output()

	return out, err
}
