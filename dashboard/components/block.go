package components

import (
	"encoding/json"
	"log"

	"github.com/jim380/Cendermint/rest"
	"github.com/kyoto-framework/kyoto/v2"
	"go.uber.org/zap"
)

/*
Component
  - Each component is a context receiver, which returns its state
  - Each component becomes a part of the page or top-level component,
    which executes component asynchronously and gets a state future object
  - Context holds common objects like http.ResponseWriter, *http.Request, etc
*/
func GetBlockInfo(ctx *kyoto.Context) (state rest.Blocks) {
	route := "/cosmos/base/tendermint/v1beta1/blocks/latest" //TO-DO refactor this
	fetchBlockInfo := func() rest.Blocks {
		var state rest.Blocks
		resp, err := rest.HttpQuery(rest.RESTAddr + route)
		if err != nil {
			zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return rest.Blocks{}
		}

		err = json.Unmarshal(resp, &state)
		if err != nil {
			zap.L().Fatal("Failed to unmarshal response", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return rest.Blocks{}
		}

		return state
	}

	/*
		Handle Actions
			- To call an action of parent component, use $ prefix in action name
			- To call an action of component by id, use <id:action> as an action name
		    - To push multiple component UI updates during a single action call,
		        call kyoto.ActionFlush(ctx, state) to initiate an update
	*/
	handled := kyoto.Action(ctx, "Reload Block", func(args ...any) {
		// add logic here
		state = fetchBlockInfo()
		log.Println("New block info fetched on block", state.Block.Header.Height)
	})
	// Prevent further execution if action handled
	if handled {
		return
	}
	// Default loading behavior if not handled
	state = fetchBlockInfo()

	return
}
