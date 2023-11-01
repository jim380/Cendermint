package rest

import (
	"encoding/json"

	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

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

	for _, validator := range v.Validators {
		// populate the map => [conspub]moniker
		vResult[validator.ConsPubKey.Key] = validator.Description.Moniker
	}
	return vResult
}
