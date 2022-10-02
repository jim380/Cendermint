package rest

import (
	"encoding/json"
	"os"
	"strings"

	"go.uber.org/zap"
)

type OracleInfo struct {
	oracleValidatorsInfo
}

type oracleValidatorsInfo struct {
	missesCount
	prevote
	feederDelegate
}

type missesCount struct {
	MissesCount string `json:"miss_counter"`
}

type prevote struct {
	AggregatePrevote struct {
		SubmitBlock string `json:"submit_block"`
	} `json:"aggregate_prevote"`
}

type feederDelegate struct {
	Address string `json:"feeder_addr"`
}

func (rd *RESTData) getOracleMissesCount() {
	var count missesCount
	res, err := HttpQuery(RESTAddr + "/umee/oracle/v1/validators/" + os.Getenv("OPERATOR_ADDR") + "/miss")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &count)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.OracleInfo.MissesCount = count.MissesCount
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Oracle misses", rd.OracleInfo.MissesCount))
}

func (rd *RESTData) getOracleSubmitBlock() {
	var vote prevote
	res, err := HttpQuery(RESTAddr + "/umee/oracle/v1/validators/" + os.Getenv("OPERATOR_ADDR") + "/aggregate_prevote")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &vote)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.OracleInfo.AggregatePrevote.SubmitBlock = vote.AggregatePrevote.SubmitBlock
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Oracle submit block", rd.OracleInfo.AggregatePrevote.SubmitBlock))
}

func (rd *RESTData) getOracleFeederDelegate() {
	var fd feederDelegate
	res, err := HttpQuery(RESTAddr + "/umee/oracle/v1/validators/" + os.Getenv("OPERATOR_ADDR") + "/feeder")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &fd)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.OracleInfo.feederDelegate.Address = fd.Address
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Oracle feeder delegate", rd.OracleInfo.feederDelegate.Address))
}
