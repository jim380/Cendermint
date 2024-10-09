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

	bondedValidators := v.Validators
	stakingParams := GetStakingParams()

	/*
	 *	This logic only applies to chains (e.g. the Cosmos Hub) where
	 *	a validator can be bonded but inactive (e.g. not signing blocks)
	 *  e.g. max_validators=200 but only the top 180 are active
	 *
	 *	For most other chains, if a validator is bonded, they are active
	 */

	// TO-DO replace MaxValidators here with active validators
	if len(bondedValidators) > stakingParams.Params.MaxValidators {
		bondedValidators = bondedValidators[:stakingParams.Params.MaxValidators]
	}

	for _, validator := range bondedValidators {
		vResult[validator.ConsPubKey.Key] = validator.Description.Moniker
	}
	return vResult
}

func GetStakingParams() types.StakingParams {
	route := GetStakingParamsRoute()
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err:", err.Error()))
		return types.StakingParams{}
	}

	var stakingParams types.StakingParams
	err = json.Unmarshal(res, &stakingParams)
	if err != nil {
		zap.L().Fatal("Failed to unmarshal response", zap.Bool("Success", false), zap.String("err:", err.Error()))
		return types.StakingParams{}
	}

	return stakingParams
}
