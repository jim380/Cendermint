package rest

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jim380/Cendermint/config"
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
		ConsAddr   string `json:"address"`
		ConsPubKey struct {
			Type string `json:"type"`
			Key  string `json:"value"`
		} `json:"pub_key"`
		ProposerPriority string `json:"proposer_priority"`
		VotingPower      string `json:"voting_power"`
		Moniker          string
	} `json:"validators"`
}

type RpcValidators struct {
	Validators []struct {
		ConsPubKey struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"consensus_pubkey"`
		Description struct {
			Moniker string `json:"moniker"`
		} `json:"description"`
	} `json:"validators"`
}

func (rpc *RPCData) getConsensusDump(cfg config.Config) {
	var cs ConsensusState
	var vSetsResult map[string][]string = make(map[string][]string)

	res, err := HttpQuery(RPCAddr + "/dump_consensus_state")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	// Unmarshal the JSON response and check for errors
	if err := json.Unmarshal(res, &cs); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	conspubMonikerMap := rpc.getConspubMonikerMap(cfg)
	// cs.Result.Validatorset.Validators is already sorted based on voting power
	for index, validator := range cs.Result.Validatorset.Validators {
		var prevote, precommit string

		// get moniker
		validator.Moniker = conspubMonikerMap[validator.ConsPubKey.Key]
		if cs.Result.RoundState.Votes[0].Prevotes[index] != "nil-Vote" {
			prevote = "✅"
		} else {
			prevote = "❌"
		}

		if cs.Result.RoundState.Votes[0].Precommits[index] != "nil-Vote" {
			precommit = "✅"
		} else {
			precommit = "❌"
		}

		// populate the map => [ConsAddr][]string{ConsPubKey, VotingPower, ProposerPriority, prevote, precommit, moniker}
		vSetsResult[validator.ConsAddr] = []string{validator.ConsPubKey.Key, validator.VotingPower, validator.ProposerPriority, prevote, precommit, validator.Moniker}
	}

	rpc.ConsensusState = cs
	rpc.Validatorsets = vSetsResult

	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Consensus", "height("+rpc.ConsensusState.Result.Height+") "+"round("+strconv.FormatInt(rpc.ConsensusState.Result.Round, 10)+") "+"step("+strconv.FormatInt(rpc.ConsensusState.Result.Step, 10)+")"))
	prevoteParsed := utils.ParseConsensusOutput(rpc.ConsensusState.Result.Votes[0].PrevotesBitArray, "\\= (.*)", 1)
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Prevote bit array", fmt.Sprintf("%.2f", prevoteParsed)))
	precommitParsed := utils.ParseConsensusOutput(rpc.ConsensusState.Result.Votes[0].PrecommitsBitArray, "\\= (.*)", 1)
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Precommit bit array", fmt.Sprintf("%.2f", precommitParsed)))
	zap.L().Info("", zap.Bool("Success", true), zap.String("# of validators from RPC", fmt.Sprint(len(rpc.Validatorsets))))
}

func (rpc *RPCData) getConspubMonikerMap(cfg config.Config) map[string]string {
	var v RpcValidators
	var vResult map[string]string = make(map[string]string)

	route := GetValidatorsRoute()
	res, err := HttpQuery(RESTAddr + route + "?status=BOND_STATUS_BONDED&pagination.limit=300")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	// Unmarshal the JSON response and check for errors
	if err := json.Unmarshal(res, &v); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	for _, validator := range v.Validators {
		// populate the map => [conspub] -> (moniker)
		vResult[validator.ConsPubKey.Key] = validator.Description.Moniker
	}
	return vResult
}
