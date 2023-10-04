package rest

import (
	"encoding/json"
	"strings"

	"github.com/jim380/Cendermint/config"
	"go.uber.org/zap"
)

type validators struct {
	Validator validator `json:"validator"`
}

type validator struct {
	OperAddr        string        `json:"operator_address"`
	ConsPubKey      consPubKeyVal `json:"consensus_pubkey"`
	Jailed          bool          `json:"jailed"`
	Status          int           `json:"status"`
	Tokens          string        `json:"tokens"`
	DelegatorShares string        `json:"delegator_shares"`
	Description     struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	}
	UnbondingHeight string `json:"unbonding_height"`
	UnbondingTime   string `json:"unbonding_time"`
	Commission      struct {
		Commission commission_rates `json:"commission_rates"`
		UpdateTime string           `json:"update_time"`
	}
	MinSelfDelegation string `json:"min_self_delegation"`
}

type consPubKeyVal struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type commission_rates struct {
	Rate            string `json:"rate"`
	Max_rate        string `json:"max_rate"`
	Max_change_rate string `json:"max_change_rate"`
}

func (rd *RESTData) getValidator(cfg config.Config) {
	var v validators

	route := getValidatorByAddressRoute(cfg)
	res, err := HttpQuery(RESTAddr + route + OperAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if err := json.Unmarshal(res, &v); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.Validator = v.Validator
}
