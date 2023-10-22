package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetOracleInfo(cfg config.Config, rd *types.RESTData) {
	rs.OracleService.GetMissedCounterInfoByValidator(cfg, rd)
	rs.OracleService.GetPrevoteInfoByValidator(cfg, rd)
	rs.OracleService.GetVoteInfoByValidator(cfg, rd)
}
