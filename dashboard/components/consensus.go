package components

import (
	"encoding/json"
	"fmt"

	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"github.com/kyoto-framework/kyoto/v2"
	"go.uber.org/zap"
)

func GetConsensusInfo(ctx *kyoto.Context) (state types.RPCData) {
	fetchConsensusInfo := func() types.RPCData {
		var state types.RPCData
		var cs types.ConsensusState
		var vSetsResult map[string][]string = make(map[string][]string)

		resp, err := utils.HttpQuery(constants.RPCAddr + "/dump_consensus_state")
		if err != nil {
			zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return types.RPCData{}
		}

		err = json.Unmarshal(resp, &cs)
		if err != nil {
			zap.L().Fatal("Failed to unmarshal response", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return types.RPCData{}
		}

		conspubMonikerMap := rest.GetConspubMonikerMap()
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

			// populate the map => [ConsAddrHex][]string{ConsPubKey, VotingPower, ProposerPriority, prevote, precommit, Moniker}
			vSetsResult[validator.ConsAddrHex] = []string{validator.ConsPubKey.Key, validator.VotingPower, validator.ProposerPriority, prevote, precommit, validator.Moniker}
		}

		// cs.Result.RoundState.Votes[0].PrevotesBitArray = utils.ParseConsensusOutput(cs.Result.RoundState.Votes[0].PrevotesBitArray, "\\= (.*)", 1)
		cs.Result.Votes[0].PrecommitsBitArray = fmt.Sprintf("%.2f", utils.ParseConsensusOutput(cs.Result.Votes[0].PrecommitsBitArray, "\\= (.*)", 1))
		cs.Result.Votes[0].PrevotesBitArray = fmt.Sprintf("%.2f", utils.ParseConsensusOutput(cs.Result.Votes[0].PrevotesBitArray, "\\= (.*)", 1))
		state.ConsensusState = cs
		state.Validatorsets = utils.Sort(vSetsResult, 1) // sort by voting power
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
