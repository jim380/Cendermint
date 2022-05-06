package rest

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type RPCData struct {
	ConsensusState
	Validatorsets map[string][]string
}

type ConsensusState struct {
	Result struct {
		RoundState `json:"round_state"`
	} `json:"result"`
}

type RoundState struct {
	Height       string           `json:"height"`
	Round        int64            `json:"round"`
	Step         int64            `json:"step"`
	Validatorset rpcValidatorsets `json:"validators"`
	Votes        []struct {
		Prevotes           []string `json:"prevotes"`
		Precommits         []string `json:"precommits"`
		PrevotesBitArray   string   `json:"prevotes_bit_array"`
		PrecommitsBitArray string   `json:"precommits_bit_array"`
	} `json:"votes"`
}

type rpcValidatorsets struct {
	Validators []struct {
		ConsAddr string `json:"address"`
		//ConsPubKey       consPubKeyValSet `json:"pub_key"`
		ProposerPriority string `json:"proposer_priority"`
		VotingPower      string `json:"voting_power"`
	} `json:"validators"`
}

func (rpc *RPCData) getConsensusDump() {
	var cs ConsensusState
	var vSetsResult map[string][]string = make(map[string][]string)

	res, err := HttpQuery(RPCAddr + "/dump_consensus_state")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &cs)

	for _, validator := range cs.Result.Validatorset.Validators {
		// populate the map => [ConsAddr][]string{ConsAddr, VotingPower, ProposerPriority, prevote, precommit, commit}
		fmt.Println("key(" + validator.ConsAddr + ") " + "value(" + validator.VotingPower + ")")
		vSetsResult[validator.ConsAddr] = []string{validator.ConsAddr, validator.VotingPower, validator.ProposerPriority, "nil", "nil", "nil"}
	}

	// sort the map based on voting power
	vSetsResultSorted := SortRPC(vSetsResult)

	// temp check; remove after testing
	_, found := vSetsResultSorted["AC2D56057CD84765E6FBE318979093E8E44AA18F"]
	if !found {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Validator not found in the active set"))
	}

	// TO-DO: update (prevote, precommit, commit} in the map
	rpc.ConsensusState = cs
	rpc.Validatorsets = vSetsResultSorted

	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Consensus:", "height("+rpc.ConsensusState.Result.Height+") "+"round("+strconv.FormatInt(rpc.ConsensusState.Result.Round, 10)+") "+"step("+strconv.FormatInt(rpc.ConsensusState.Result.Step, 10)+")"))
	prevoteParsed := utils.ParseConsensusOutput(rpc.ConsensusState.Result.Votes[0].PrevotesBitArray, "\\= (.*)", 1)
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Prevote bit array:", fmt.Sprintf("%.2f", prevoteParsed)))
	precommitParsed := utils.ParseConsensusOutput(rpc.ConsensusState.Result.Votes[0].PrecommitsBitArray, "\\= (.*)", 1)
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Precommit bit array:", fmt.Sprintf("%.2f", precommitParsed)))
	zap.L().Info("", zap.Bool("Success", true), zap.String("# of validators from RPC: ", fmt.Sprint(len(rpc.Validatorsets))))
	// fmt.Println(rpc.Validatorsets)
}

func SortRPC(mapValue map[string][]string) map[string][]string {
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
		newMapValue[key][5] = strconv.Itoa(i + 1)
	}
	return newMapValue
}
