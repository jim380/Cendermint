package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"
)

type slashingInfo struct {
	Params     Params      `json:"params"`
	ValSigning SigningInfo `json:"val_signing_info"`
}

type Params struct {
	SignedBlocksWindow      string `json:"signed_blocks_window"`
	MinSignedPerWindow      string `json:"min_signed_per_window"`
	DowntimeJailDuration    string `json:"downtime_jail_duration"`
	SlashFractionDoubleSign string `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   string `json:"slash_fraction_downtime"`
}

type SigningInfo struct {
	StartHeight         string `json:"start_height"`
	IndexOffset         string `json:"index_offset"`
	JailedUntil         string `json:"jailed_until"`
	Tombstoned          bool   `json:"tombstoned"`
	MissedBlocksCounter string `json:"missed_blocks_counter"`
}

func (rd *RESTData) getSlashingParams() {
	var d slashingInfo

	res, err := HttpQuery(RESTAddr + "/cosmos/slashing/v1beta1/params")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &d)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}
	rd.Slashing.Params = d.Params
}

func (rd *RESTData) getSigningInfo(consAddr string) {
	var d slashingInfo

	res, err := HttpQuery(RESTAddr + "/cosmos/slashing/v1beta1/signing_infos/" + consAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &d)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.Slashing.ValSigning = d.ValSigning
}
