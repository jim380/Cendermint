package pages

import (
	"github.com/jim380/Cendermint/dashboard/components"
	"github.com/jim380/Cendermint/rest"
	"github.com/kyoto-framework/kyoto/v2"
)

type PIndexState struct {
	Block *kyoto.ComponentF[rest.Blocks]
	Node  *kyoto.ComponentF[rest.NodeInfo]
}

/*
Page
  - A page is a top-level component, which attaches components and
    defines rendering
*/
func PIndex(ctx *kyoto.Context) (state PIndexState) {
	// Define rendering
	kyoto.Template(ctx, "page.index.html")

	// Attach components
	state.Block = kyoto.Use(ctx, components.GetBlockInfo)
	state.Node = kyoto.Use(ctx, components.GetNodeInfo)

	return
}
