package dashboard

import (
	"fmt"
	"os"

	"github.com/jim380/Cendermint/dashboard/components"
	"github.com/jim380/Cendermint/dashboard/pages"
	"github.com/jim380/Cendermint/models"
	"github.com/kyoto-framework/kyoto/v2"
	"go.uber.org/zap"
)

func StartDashboard() {
	// Setup a database connection
	cfg := models.DefaultPostgresConfig()
	fmt.Println(cfg.String())
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Database connection", "OK"))
	}
	defer db.Close()

	// Setup our model services
	// userService := models.UserService{
	// 	DB: db,
	// }
	// sessionService := models.SessionService{
	// 	DB: db,
	// }

	go func() {
		port := os.Getenv("DASHBOARD_PORT")
		// Register page
		kyoto.HandlePage("/", pages.PIndex)
		// Client
		kyoto.HandleAction(components.GetBlockInfo)
		kyoto.HandleAction(components.GetNodeInfo)
		kyoto.HandleAction(components.GetConsensusInfo)
		// Serve
		kyoto.Serve(":" + port)
	}()
}
