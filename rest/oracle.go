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
	MissesCount string `json:"miss_counter"`
}

func (rd *RESTData) getOracleMissesCount() {
	var valInfo oracleValidatorsInfo
	res, err := HttpQuery(RESTAddr + "/umee/oracle/v1/validators/" + os.Getenv("OPERATOR_ADDR") + "/miss")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &valInfo)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.OracleInfo.MissesCount = valInfo.MissesCount
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Oracle misses", rd.OracleInfo.MissesCount))
}
