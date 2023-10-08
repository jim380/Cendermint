package rest

import (
	"encoding/json"
	"strings"

	"github.com/jim380/Cendermint/config"
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

func (rd *RESTData) getSlashingParams(cfg config.Config) {
	var d slashingInfo

	route := getSlashingParamsRoute(cfg)
	res, err := HttpQuery(RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if err := json.Unmarshal(res, &d); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Bool("Success", false), zap.String("err", err.Error()))
		return
	}
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}
	rd.Slashing.Params = d.Params
}

func (rd *RESTData) getSigningInfo(cfg config.Config, consAddr string) {
	var d slashingInfo

	route := getSigningInfoByAddressRoute(cfg)
	res, err := HttpQuery(RESTAddr + route + consAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if err := json.Unmarshal(res, &d); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Bool("Success", false), zap.String("err", err.Error()))
		return
	}
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.Slashing.ValSigning = d.ValSigning
}
