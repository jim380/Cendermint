package models

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type OracleService struct {
	DB *sql.DB
}

func (os *OracleService) Init(db *sql.DB) {
	os.DB = db
}

func (os *OracleService) GetMissedCounterInfoByValidator(cfg config.Config, rd *types.RESTData) {
	var ms types.MissedCounterInfo

	route := rest.GetMissedCounterRoute()

	res, err := utils.HttpQuery(constants.RESTAddr + route + "/" + cfg.OperatorAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &ms)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.OracleInfo.MissedCounterInfo = ms
}

func (os *OracleService) GetPrevoteInfoByValidator(cfg config.Config, rd *types.RESTData) {
	var pv types.PrevoteInfo

	route := rest.GetPrevoteRoute()

	res, err := utils.HttpQuery(constants.RESTAddr + route + "/" + cfg.OperatorAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &pv)
	if strings.Contains(string(res), "not found") {
		zap.L().Warn("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Warn("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.OracleInfo.PrevoteInfo = pv
}

func (os *OracleService) GetVoteInfoByValidator(cfg config.Config, rd *types.RESTData) {
	var v types.VoteInfo

	route := rest.GetVoteRoute()

	res, err := utils.HttpQuery(constants.RESTAddr + route + "/" + cfg.OperatorAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &v)
	if strings.Contains(string(res), "not found") {
		zap.L().Warn("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Warn("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.OracleInfo.VoteInfo = v
}
