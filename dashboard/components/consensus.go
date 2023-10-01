package components

import (
	"encoding/json"
	"fmt"

	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/utils"
	"github.com/kyoto-framework/kyoto/v2"
	"go.uber.org/zap"
)

func GetConsensusInfo(ctx *kyoto.Context) (state rest.RPCData) {
	fetchConsensusInfo := func() rest.RPCData {
		var state rest.RPCData
		var cs rest.ConsensusState
		var vSetsResult map[string][]string = make(map[string][]string)

		resp, err := rest.HttpQuery(rest.RPCAddr + "/dump_consensus_state")
		if err != nil {
			zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return rest.RPCData{}
		}

		err = json.Unmarshal(resp, &cs)
		if err != nil {
			zap.L().Fatal("Failed to unmarshal response", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return rest.RPCData{}
		}

		conspubMonikerMap := getConspubMonikerMap()
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

			// populate the map => [ConsAddr][]string{ConsAddr, VotingPower, ProposerPriority, prevote, precommit, Moniker}
			vSetsResult[validator.ConsAddr] = []string{validator.ConsPubKey.Key, validator.VotingPower, validator.ProposerPriority, prevote, precommit, validator.Moniker}
		}

		// cs.Result.RoundState.Votes[0].PrevotesBitArray = utils.ParseConsensusOutput(cs.Result.RoundState.Votes[0].PrevotesBitArray, "\\= (.*)", 1)
		cs.Result.Votes[0].PrecommitsBitArray = fmt.Sprintf("%.2f", utils.ParseConsensusOutput(cs.Result.Votes[0].PrecommitsBitArray, "\\= (.*)", 1))
		cs.Result.Votes[0].PrevotesBitArray = fmt.Sprintf("%.2f", utils.ParseConsensusOutput(cs.Result.Votes[0].PrevotesBitArray, "\\= (.*)", 1))
		state.ConsensusState = cs
		state.Validatorsets = rest.Sort(vSetsResult, 1) // sort by voting power
		return state
	}

	handled := kyoto.Action(ctx, "Reload Consensus", func(args ...any) {
		state = fetchConsensusInfo()
	})

	if handled {
		return
	}

	state = fetchConsensusInfo()

	return
}

func getConspubMonikerMap() map[string]string {
	var v rest.RpcValidators
	var vResult map[string]string = make(map[string]string)

	route := rest.GetValidatorsRoute()
	res, err := rest.HttpQuery(rest.RESTAddr + route + "?status=BOND_STATUS_BONDED&pagination.limit=300")
	if err != nil {
		zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err:", err.Error()))
		return map[string]string{}
	}
	err = json.Unmarshal(res, &v)
	if err != nil {
		zap.L().Fatal("Failed to unmarshal response", zap.Bool("Success", false), zap.String("err:", err.Error()))
		return map[string]string{}
	}

	for _, validator := range v.Validators {
		// populate the map => [conspub] -> (moniker)
		vResult[validator.ConsPubKey.Key] = validator.Description.Moniker
	}
	return vResult
}
