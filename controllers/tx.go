package controllers

import (
	"strconv"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
	"go.uber.org/zap"
)

func (rs RestServices) GetTxnInfo(cfg config.Config, currentBlockHeight int64, rd *types.RESTData) {
	txsInBlock, err := rs.TxnService.GetTxnsInBlock(cfg, currentBlockHeight)
	if err != nil {
		zap.L().Error("Error getting txns in block", zap.String("Height", strconv.Itoa(int(currentBlockHeight))))
		return
	}
	rs.TxnService.PopulateRestData(rd, txsInBlock)
	rs.IndexTxnsInBlock(cfg, currentBlockHeight, txsInBlock)
}

func (rs RestServices) IndexTxnsInBlock(cfg config.Config, height int64, txsInBlock types.TxInfo) {
	err := rs.TxnService.Index(cfg, height, txsInBlock)
	if err != nil {
		zap.L().Error("Error indexing txns for block "+strconv.FormatInt(height, 10), zap.String("Error", err.Error()))
		return
	} else {
		zap.L().Debug("Txns successfully indexed", zap.String("Height", strconv.FormatInt(height, 10)))
	}
}
