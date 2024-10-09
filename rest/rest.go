package rest

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/types"

	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

var GetConspubMonikerMapWrapper = GetConspubMonikerMap

const activeSetSize = 180

func GetConspubMonikerMap() map[string]string {
	var v types.RpcValidators
	var vResult map[string]string = make(map[string]string)

	route := GetValidatorsRoute()
	res, err := utils.HttpQuery(constants.RESTAddr + route + "?status=BOND_STATUS_BONDED&pagination.limit=300")
	if err != nil {
		zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err:", err.Error()))
		return map[string]string{}
	}
	err = json.Unmarshal(res, &v)
	if err != nil {
		zap.L().Fatal("Failed to unmarshal response", zap.Bool("Success", false), zap.String("err:", err.Error()))
		return map[string]string{}
	}

	// Sort validators by tokens in descending order
	sort.Slice(v.Validators, func(i, j int) bool {
		tokensI, _ := strconv.ParseInt(v.Validators[i].Tokens, 10, 64)
		tokensJ, _ := strconv.ParseInt(v.Validators[j].Tokens, 10, 64)
		return tokensI > tokensJ
	})

	activeValidators := v.Validators
	if len(activeValidators) > activeSetSize {
		activeValidators = activeValidators[:activeSetSize]
	}

	for _, validator := range activeValidators {
		vResult[validator.ConsPubKey.Key] = validator.Description.Moniker
	}
	return vResult
}
