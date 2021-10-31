package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"
)

type slashingInfo struct {
	Params Params `json:"params"`
}

type Params struct {
	SignedBlocksWindow      string `json:"signed_blocks_window"`
	MinSignedPerWindow      string `json:"min_signed_per_window"`
	DowntimeJailDuration    string `json:"downtime_jail_duration"`
	SlashFractionDoubleSign string `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   string `json:"slash_fraction_downtime"`
}

func (rd *RESTData) getSlashingParams() {
	var d slashingInfo

	res, err := RESTQuery("/cosmos/slashing/v1beta1/params")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &d)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.Slashing.Params.SignedBlocksWindow = d.Params.SignedBlocksWindow
	rd.Slashing.Params.MinSignedPerWindow = d.Params.MinSignedPerWindow
	rd.Slashing.Params.DowntimeJailDuration = d.Params.DowntimeJailDuration
	rd.Slashing.Params.SlashFractionDoubleSign = d.Params.SlashFractionDoubleSign
	rd.Slashing.Params.SlashFractionDowntime = d.Params.SlashFractionDowntime
}
