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
			ConsAddr         string           `json:"address"`
			ConsPubKey       consPubKeyValSet `json:"pub_key"`
			ProposerPriority string           `json:"proposer_priority"`
			VotingPower      string           `json:"voting_power"`
		}
	}
}

type consPubKeyValSet struct {
	Type string `json:"type"`
	Key  string `json:"value"`
}

func (rd *RESTData) getValidatorsets(currentBlockHeight int64) {
	var vSets, vSets2, vSets3 validatorsets
	var vSetsResult map[string][]string = make(map[string][]string)
	var vSetsResult2 map[string][]string = make(map[string][]string)
	var vSetsResult3 map[string][]string = make(map[string][]string)

	runPages(currentBlockHeight, &vSets, vSetsResult, 1)
	runPages(currentBlockHeight, &vSets2, vSetsResult2, 2)
	runPages(currentBlockHeight, &vSets2, vSetsResult2, 3)

	for _, value := range vSets.Result.Validators {
		// populate the validatorset map => [ConsPubKey][]string{ConsAddr, VotingPower, ProposerPriority}
		vSetsResult[value.ConsPubKey.Key] = []string{value.ConsAddr, value.VotingPower, value.ProposerPriority, "0"}
	}

	for _, value := range vSets2.Result.Validators {
		// populate the validatorset map => [ConsPubKey][]string{ConsAddr, VotingPower, ProposerPriority}
		vSetsResult2[value.ConsPubKey.Key] = []string{value.ConsAddr, value.VotingPower, value.ProposerPriority, "0"}
	}

	for _, value := range vSets3.Result.Validators {
		// populate the validatorset map => [ConsPubKey][]string{ConsAddr, VotingPower, ProposerPriority}
		vSetsResult2[value.ConsPubKey.Key] = []string{value.ConsAddr, value.VotingPower, value.ProposerPriority, "0"}
	}
	vSetsResultTemp := mergeMap(vSetsResult, vSetsResult2)
	vSetsResultFinal := mergeMap(vSetsResultTemp, vSetsResult3)
	rd.Validatorsets = Sort(vSetsResultFinal)
	zap.L().Info("", zap.Bool("Success", true), zap.String("Active validators", fmt.Sprint(len(vSets.Result.Validators)+len(vSets2.Result.Validators)+len(vSets3.Result.Validators))))
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

func mergeMap(a map[string][]string, b map[string][]string) map[string][]string {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func runPages(currentBlockHeight int64, vSets *validatorsets, vSetsResult map[string][]string, pages int) {
	res, err := HttpQuery(RESTAddr + "/validatorsets/" + fmt.Sprint(currentBlockHeight) + "?page=" + strconv.Itoa(pages) + "&limit=100")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	json.Unmarshal(res, &vSets)

	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	for _, value := range vSets.Result.Validators {
		// populate the validatorset map => [ConsPubKey][]string{ConsAddr, VotingPower, ProposerPriority}
		vSetsResult[value.ConsPubKey.Key] = []string{value.ConsAddr, value.VotingPower, value.ProposerPriority, "0"}
	}
}
