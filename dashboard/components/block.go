package components

import (
	"encoding/json"
	"log"
	"os"

	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
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
func GetBlockInfo(ctx *kyoto.Context) (state types.Blocks) {
	route := "/cosmos/base/tendermint/v1beta1/blocks/latest" //TO-DO refactor this
	fetchBlockInfo := func() types.Blocks {
		var state types.Blocks
		resp, err := utils.HttpQuery(constants.RESTAddr + route)
		if err != nil {
			zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return types.Blocks{}
		}

		err = json.Unmarshal(resp, &state)
		if err != nil {
			zap.L().Fatal("Failed to unmarshal response", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return types.Blocks{}
		}

		// convert block hash from base64 to hex
		hashInHex, err := utils.Base64ToHex(state.BlockId.Hash)
		if err != nil {
			zap.L().Fatal("Failed to convert base64 to hex", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return types.Blocks{}
		}
		state.BlockId.Hash = hashInHex

		/*
			Find validators with missing signatures in the block
		*/
		var activeSet map[string][]string = make(map[string][]string)

		conspubMonikerMap := rest.GetConspubMonikerMap()

		for consPubKey, moniker := range conspubMonikerMap {
			prefix, err := utils.GetPrefix(os.Getenv("CHAIN"))
			if err != nil {
				zap.L().Fatal("Failed to get prefix", zap.Bool("Success", false), zap.String("err", err.Error()))
			}
			consAddrHex := utils.PubkeyToHexAddr(prefix, consPubKey)
			if consAddrHex == "" {
				zap.L().Fatal("Failed to convert public key to hex address", zap.Bool("Success", false))
			}
			// populate the map => [ConsAddr]{consPubKey, moniker}
			activeSet[consAddrHex] = []string{consPubKey, moniker}
		}

		/*
			- Create a map validatorConsAddrSignedMap using allSignaturesInBlock for quick lookup
			- validatorConsAddrSignedMap gives all validators who signed on this block
		*/
		allSignaturesInBlock := state.Block.LastCommit.Signatures
		validatorConsAddrSignedMap := make(map[string]bool)
		for _, signature := range allSignaturesInBlock {
			// Validator_address could be in hex or base64; hex is legacy so using base64 here
			validatorConsAddrSignedMap[signature.Validator_address] = true
		}

		// Check if validator.ConsAddr in activeSet exists in validatorConsAddrSignedMap
		for consAddrHex, props := range activeSet {
			// convert consAddrHex to base64
			consAddrInBase64, err := utils.HexToBase64(consAddrHex)
			if err != nil {
				zap.L().Fatal("HexToBase64", zap.Bool("Success", false), zap.String("err:", err.Error()))
				return types.Blocks{}
			}
			if _, exists := validatorConsAddrSignedMap[consAddrInBase64]; !exists {
				// If the Validator_address does not exist in allSignaturesInBlock, add it to MissingValidators
				state.MissingValidators = append(state.MissingValidators, struct {
					Moniker     string
					ConsHexAddr string
				}{
					Moniker:     props[1],
					ConsHexAddr: consAddrHex,
					// TO-DO add operator address
				})
			}
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
