package rest

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type validatorsets struct {
	Validators []struct {
		Address string
		Pub_key struct {
			Type string `json:"@type"`
			Key  string
		}
		Voting_power string
	}
}

func getValidatorsets(currentBlockHeight int64, log *zap.Logger) map[string][]string {
	var vSets validatorsets
	var vSetsResult map[string][]string = make(map[string][]string)

	res, _ := runRESTCommand("/cosmos/base/tendermint/v1beta1/validatorsets/" + fmt.Sprint(currentBlockHeight) + "?pagination.limit=1000")
	json.Unmarshal(res, &vSets)
	if strings.Contains(string(res), "not found") {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		log.Info("", zap.Bool("Success", true), zap.String("Number of loaded validators", fmt.Sprint(len(vSets.Validators))))
	}

	for _, value := range vSets.Validators {
		// populate the validator set map
		vSetsResult[value.Pub_key.Key] = []string{value.Address, value.Voting_power}
	}

	return vSetsResult
}
