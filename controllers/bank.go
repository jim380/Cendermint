package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetBalanceInfo(cfg config.Config, rd *types.RESTData) {
	rs.BankService.GetBalanceInfo(cfg, rd)
}

func (rs RestServices) GetRewardsCommissionInfo(cfg config.Config, rd *types.RESTData) {
	rs.BankService.GetRewardsCommissionInfo(cfg, rd)
}
