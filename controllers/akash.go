package controllers

import (
	"github.com/jim380/Cendermint/config"
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
	providersWithNullAuditor, err := rs.AkashService.FindProvidersWithNullAuditor()
	if err != nil {
		zap.L().Error("Error finding akash providers with null auditor", zap.String("Error", err.Error()))
		return
	}

	if len(providersWithNullAuditor) == 0 {
		zap.L().Info("No akash providers with null auditor found")
		return
	}

	err = rs.AkashService.IndexAuditorForProviderOwners(cfg, providersWithNullAuditor)
	if err != nil {
		zap.L().Error("Error indexing akash auditors", zap.String("Error", err.Error()))
		return
	} else {
		zap.L().Info("Akash auditors successfully indexed", zap.Int("Amount: ", len(providersWithNullAuditor)))
	}
}
