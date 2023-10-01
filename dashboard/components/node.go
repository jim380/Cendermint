package components

import (
	"encoding/json"

	"github.com/jim380/Cendermint/rest"
	"github.com/kyoto-framework/kyoto/v2"
	"go.uber.org/zap"
)

func GetNodeInfo(ctx *kyoto.Context) (state rest.NodeInfo) {
	route := rest.GetNodeInfoRoute()
	fetchNodeInfo := func() rest.NodeInfo {
		var state rest.NodeInfo
		resp, err := rest.HttpQuery(rest.RESTAddr + route)
		if err != nil {
			zap.L().Fatal("Connection to REST failed", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return rest.NodeInfo{}
		}

		err = json.Unmarshal(resp, &state)
		if err != nil {
			zap.L().Fatal("Failed to unmarshal response", zap.Bool("Success", false), zap.String("err:", err.Error()))
			return rest.NodeInfo{}
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
