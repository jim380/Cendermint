package controllers

import (
	"time"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/types/akash"
	"go.uber.org/zap"
)

func (rs RestServices) GetAkashInfo(cfg config.Config, data *types.AsyncData) {
	rs.AkashService.GetAkashDeployments(cfg, data)
	providers := rs.AkashService.GetAkashProviders(cfg, data)
	// index
	rs.IndexAkashProviders(cfg, providers)
}

func (rs RestServices) IndexAkashProviders(cfg config.Config, providers akash.ProvidersResponse) {
	err := rs.AkashService.IndexProviders(cfg, providers)
	if err != nil {
		zap.L().Error("Error indexing akash providers", zap.String("Error", err.Error()))
		return
	} else {
		zap.L().Info("Akash providers successfully indexed", zap.Int("Amount: ", len(providers.Providers)))
	}
}

func (rs RestServices) IndexAkashAuditors(cfg config.Config) {
	providerOwnersPendingUpdate, err := rs.AkashService.FindProviderOwnersPendingAuditorUpdate(time.Duration(constants.LastUpdatedMoreThan) * time.Second)
	if err != nil {
		zap.L().Error("Error finding akash providers pending auditor update", zap.String("IndexAkashAuditors", err.Error()))
		return
	}

	if len(providerOwnersPendingUpdate) == 0 {
		zap.L().Info("No akash providers pending auditor update", zap.String("IndexAkashAuditors", ""))
		return
	}

	err = rs.AkashService.IndexAuditorForProviderOwners(cfg, providerOwnersPendingUpdate)
	if err != nil {
		zap.L().Error("Error indexing akash auditors", zap.String("IndexAkashAuditors", err.Error()))
		return
	} else {
		zap.L().Info("Akash auditors successfully indexed", zap.Int("Amount: ", len(providerOwnersPendingUpdate)))
	}
}

func (rs RestServices) IndexAkashDeployments(cfg config.Config) {
	providerOwnersPendingUpdate, err := rs.AkashService.FindProviderOwnersPendingDeploymentUpdate(time.Duration(constants.LastUpdatedMoreThan) * time.Second)
	if err != nil {
		zap.L().Error("Error finding akash providers pending deployment update", zap.String("IndexAkashDeployments", err.Error()))
		return
	}

	if len(providerOwnersPendingUpdate) == 0 {
		zap.L().Info("No akash providers pending deployment update", zap.String("IndexAkashDeployments", ""))
		return
	}

	err = rs.AkashService.IndexDeploymentForProviderOwner(cfg, providerOwnersPendingUpdate)
	if err != nil {
		zap.L().Error("Error indexing akash deployments", zap.String("IndexAkashDeployments", err.Error()))
		return
	} else {
		zap.L().Info("Akash deployments successfully indexed", zap.Int("Amount: ", len(providerOwnersPendingUpdate)))
	}

}
