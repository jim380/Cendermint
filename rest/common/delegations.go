package rest

import (
	"encoding/json"
	"fmt"
	"strings"

	utils "github.com/jim380/Cosmos-IE/utils"
	"go.uber.org/zap"
)

type delegationInfo struct {
	DelegationCount float64
	SelfDelegation  float64
}

type delegations struct {
	Delegation_responses []struct {
		Delegation delegation
	}

	Pagination struct {
		Total string
	}
}

type selfDelegation struct {
	Delegation_response struct {
		Delegation delegation
	}
}

type delegation struct {
	Delegator_address string `json:"delegator_address"`
	Validator_address string `json:"validator_address"`
	Shares            string `json:"shares"`
}

var (
	dInfo delegationInfo
)

func getDelegations(log *zap.Logger) delegationInfo {
	var d delegations

	res, _ := runRESTCommand("/cosmos/staking/v1beta1/validators/" + OperAddr + "/delegations?pagination.limit=10000")
	json.Unmarshal(res, &d)
	if strings.Contains(string(res), "not found") {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		log.Info("", zap.Bool("Success", true), zap.String("Delegation Count", fmt.Sprint(len(d.Delegation_responses))))
	}

	dInfo.DelegationCount = float64(len(d.Delegation_responses))

	for _, value := range d.Delegation_responses {
		if AccAddr == value.Delegation.Delegator_address {
			dInfo.SelfDelegation = utils.StringToFloat64(value.Delegation.Shares)
		}
	}

	return dInfo
}
