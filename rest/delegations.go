package rest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jim380/Cendermint/config"
	"go.uber.org/zap"
)

type delegationsInfo struct {
	DelegationRes delegationRes `json:"delegation_responses"`
	Pagination    struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}
}

type delegationRes []struct {
	Delegation struct {
		DelegatorAddr string `json:"delegator_address"`
		ValidatorAddr string `json:"validator_address"`
		Shares        string `json:"shares"`
	}
	balance struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}
}

func (rd *RESTData) getDelegations(cfg config.Config) {
	var delInfo delegationsInfo
	var delRes map[string][]string = make(map[string][]string)

	route := getValidatorByAddressRoute(cfg)
	res, err := HttpQuery(RESTAddr + route + OperAddr + "/delegations" + "?pagination.limit=1000")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	
	// Unmarshal the JSON response and check for errors
	if err := json.Unmarshal(res, &delInfo); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Total delegations from range", fmt.Sprint(len(delInfo.DelegationRes))))
		zap.L().Info("", zap.Bool("Success", true), zap.String("Total delegations from pagination", delInfo.Pagination.Total))
	}

	for _, value := range delInfo.DelegationRes {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// precommit failure validator
				}
			}()
			delRes[value.Delegation.DelegatorAddr] = []string{value.balance.Amount}
		}()
	}

	rd.Delegations = delInfo
}
