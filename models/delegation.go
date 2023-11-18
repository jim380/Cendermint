package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type DelegationService struct {
	DB *sql.DB
}

func (ds *DelegationService) Init(db *sql.DB) {
	ds.DB = db
}

func (ds *DelegationService) GetInfo(cfg config.Config, rd *types.RESTData) {
	var delInfo types.DelegationsInfo
	var delRes map[string][]string = make(map[string][]string)

	route := rest.GetValidatorByAddressRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route + constants.OperAddr + "/delegations" + "?pagination.limit=1000")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &delInfo)
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
			delRes[value.Delegation.DelegatorAddr] = []string{value.Balance.Amount}
		}()
	}

	rd.Delegations = delInfo
}
