package terra

import (
	"encoding/json"
	"strings"

	utils "github.com/jim380/Cosmos-IE/utils"
	"go.uber.org/zap"
)

type Result struct {
	Height string `json:"height"`
	Result string `json:"result"`
}

func GetOracleMiss(restServer string, operAddr string, log *zap.Logger) float64 {

	var r Result

	// Oracle Missing
	res, _ := utils.RunRESTCommand(restServer, "/oracle/voters/"+operAddr+"/miss")
	json.Unmarshal(res, &r)
	// log
	if strings.Contains(string(res), "not found") {
		// handle error
		log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		log.Info("\t", zap.Bool("Success", true), zap.String("Terra Oracle Miss", r.Result))
	}

	return utils.StringToFloat64(r.Result)
}
