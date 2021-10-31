package rest

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type validatorsets struct {
	Height string `json:"height"`

	Result struct {
		Block_Height string `json:"block_height"`
		Validators   []struct {
			ConsAddr         string     `json:"address"`
			ConsPubKey       consPubKey `json:"pub_key"`
			ProposerPriority string     `json:"proposer_priority"`
			VotingPower      string     `json:"voting_power"`
		}
	}
}

func (rd *RESTData) getValidatorsets(currentBlockHeight int64) {
	var vSets validatorsets
	var vSetsResult map[string][]string = make(map[string][]string)

	res, err := RESTQuery("/validatorsets/" + fmt.Sprint(currentBlockHeight) + "?pagination.limit=1000")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &vSets)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Active validators", fmt.Sprint(len(vSets.Result.Validators))))
	}

	for _, value := range vSets.Result.Validators {
		// populate the validatorset map => [ConsPubKey][]string{ConsAddr, VotingPower, ProposerPriority}
		vSetsResult[value.ConsPubKey.Value] = []string{value.ConsAddr, value.VotingPower, value.ProposerPriority, "0"}
	}

	rd.Validatorsets = Sort(vSetsResult)
}

func Sort(mapValue map[string][]string) map[string][]string {

	keys := []string{}
	newMapValue := mapValue

	for key := range mapValue {
		keys = append(keys, key)
	}

	// Sort by proposer_priority
	sort.Slice(keys, func(i, j int) bool {
		a, _ := strconv.Atoi(mapValue[keys[i]][2])
		b, _ := strconv.Atoi(mapValue[keys[j]][2])
		return a > b
	})

	for i, key := range keys {
		// proposer_ranking
		newMapValue[key][3] = strconv.Itoa(i + 1)
	}
	return newMapValue
}
