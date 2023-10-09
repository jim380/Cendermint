package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetTxnInfo(cfg config.Config, currentBlockHeight int64, rd *types.RESTData) {
	rs.TxnService.GetInfo(cfg, currentBlockHeight, rd)
}
