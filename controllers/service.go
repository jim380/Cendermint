package controllers

import (
	"database/sql"

	"github.com/jim380/Cendermint/services"
)

type RestServices struct {
	BlockService           *services.BlockService
	TxnService             *services.TxnService
	AbsentValidatorService *services.AbsentValidatorService
	NodeService            *services.NodeService
	StakingService         *services.StakingService
	SlashingService        *services.SlashingService
	InflationService       *services.InflationService
	GovService             *services.GovService
	BankService            *services.BankService
	DelegationService      *services.DelegationService
	UpgradeService         *services.UpgradeService
	IbcServices            *services.IbcService
	GravityService         *services.GravityService
	AkashService           *services.AkashService
	OracleService          *services.OracleService
}

type RpcServices struct {
	ValidatorService *services.ValidatorService
	ConsensusService *services.ConsensusService
}

func InitializeRpcServices(db *sql.DB) RpcServices {
	rpcServices := []services.RpcServiceExecutor{
		&services.ValidatorService{},
		&services.ConsensusService{},
	}

	for _, service := range rpcServices {
		service.Init(db)
	}

	rpcServicesController := RpcServices{
		ValidatorService: rpcServices[0].(*services.ValidatorService),
		ConsensusService: rpcServices[1].(*services.ConsensusService),
	}

	return rpcServicesController
}

func InitializeRestServices(db *sql.DB) RestServices {
	restServices := []services.RestServiceExecutor{
		&services.BlockService{},
		&services.AbsentValidatorService{},
		&services.NodeService{},
		&services.StakingService{},
		&services.SlashingService{},
		&services.InflationService{},
		&services.GovService{},
		&services.BankService{},
		&services.DelegationService{},
		&services.UpgradeService{},
		&services.IbcService{},
		&services.GravityService{},
		&services.AkashService{},
		&services.OracleService{},
	}

	for _, service := range restServices {
		service.Init(db)
	}

	restServicesController := RestServices{
		BlockService:           restServices[0].(*services.BlockService),
		AbsentValidatorService: restServices[1].(*services.AbsentValidatorService),
		NodeService:            restServices[2].(*services.NodeService),
		StakingService:         restServices[3].(*services.StakingService),
		SlashingService:        restServices[4].(*services.SlashingService),
		InflationService:       restServices[5].(*services.InflationService),
		GovService:             restServices[6].(*services.GovService),
		BankService:            restServices[7].(*services.BankService),
		DelegationService:      restServices[8].(*services.DelegationService),
		UpgradeService:         restServices[9].(*services.UpgradeService),
		IbcServices:            restServices[10].(*services.IbcService),
		GravityService:         restServices[11].(*services.GravityService),
		AkashService:           restServices[12].(*services.AkashService),
		OracleService:          restServices[13].(*services.OracleService),
	}

	return restServicesController
}
