package dashboard

import (
	"os"

	"github.com/jim380/Cendermint/dashboard/components"
	"github.com/jim380/Cendermint/dashboard/pages"
	"github.com/kyoto-framework/kyoto/v2"
)

func StartDashboard() {
	port := os.Getenv("DASHBOARD_PORT")
	// Register page
	kyoto.HandlePage("/", pages.PIndex)
	// Client
	kyoto.HandleAction(components.GetBlockInfo)
	kyoto.HandleAction(components.GetNodeInfo)
	kyoto.HandleAction(components.GetConsensusInfo)
	// Serve
	kyoto.Serve(":" + port)
}
