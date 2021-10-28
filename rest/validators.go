package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"
)

type validators struct {
	Height string `json:"height"`
	Result validator
}

type validator struct {
	OperAddr        string `json:"operator_address"`
	ConsPubKey      string `json:"consensus_pubkey"`
	Jailed          bool   `json:"jailed"`
	Status          int    `json:"status"`
	Tokens          string `json:"tokens"`
	DelegatorShares string `json:"delegator_shares"`
	Description     struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	}
	UnbondingHeight string `json:"unbonding_height"`
	UnbondingTime   string `json:"unbonding_time"`
	Commission      struct {
		Commission_rates struct {
			Rate            string `json:"rate"`
			Max_rate        string `json:"max_rate"`
			Max_change_rate string `json:"max_change_rate"`
		}
		UpdateTime string `json:"update_time"`
	}
	MinSelfDelegation string `json:"min_self_delegation"`
}

func (rd *RESTData) getValidators() {
	var v validators

	res, err := runRESTCommand("/staking/validators/" + OperAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Failed to connect to REST-Server"))
	}
	json.Unmarshal(res, &v)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Validator Moniker", v.Result.Description.Moniker))
	}

	rd.Validators = v.Result
}
