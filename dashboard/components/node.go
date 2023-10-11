package components

import (
	"encoding/json"

	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"github.com/kyoto-framework/kyoto/v2"
	"go.uber.org/zap"
)

func GetNodeInfo(ctx *kyoto.Context) (state types.NodeInfo) {
	route := rest.GetNodeInfoRoute()
	fetchNodeInfo := func() types.NodeInfo {
		var state types.NodeInfo
		resp, err := utils.HttpQuery(constants.RESTAddr + route)
		if err != nil {
			zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return types.NodeInfo{}
		}

		err = json.Unmarshal(resp, &state)
		if err != nil {
			zap.L().Fatal("Failed to unmarshal response", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return types.NodeInfo{}
		}

		return state
	}

	handled := kyoto.Action(ctx, "Reload Node", func(args ...any) {
		state = fetchNodeInfo()
	})

	if handled {
		return
	}

	state = fetchNodeInfo()

	return
}
