package akash

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/types/akash"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type AkashService struct {
	DB *sql.DB
}

func (as *AkashService) Init(db *sql.DB) {
	as.DB = db
}

func (as *AkashService) GetAkashDeployments(cfg config.Config, data *types.AsyncData) akash.Deployments {
	if cfg.Chain.Name != "akash" {
		return akash.Deployments{}
	}
	var deployments, activeDeployments akash.DeploymentsResponse

	route := rest.GetDeploymentsRoute()
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &deployments)

	// get total deployments count
	totalDeploymentsCount, err := strconv.Atoi(deployments.Pagination.Total)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	data.AkashInfo.TotalDeployments = totalDeploymentsCount

	// get active deployments count
	resActive, err := utils.HttpQuery(constants.RESTAddr + route + "?filters.state=active")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(resActive, &activeDeployments)

	activeDeploymentsCount, err := strconv.Atoi(activeDeployments.Pagination.Total)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	data.AkashInfo.ActiveDeployments = activeDeploymentsCount

	// get closed deployments count
	data.AkashInfo.ClosedDeployments = totalDeploymentsCount - activeDeploymentsCount

	return data.AkashInfo.Deployments
}

func (as *AkashService) GetAkashProviders(cfg config.Config, data *types.AsyncData) akash.ProvidersResponse {
	if cfg.Chain.Name != "akash" {
		return akash.ProvidersResponse{}
	}
	var providers akash.ProvidersResponse

	route := rest.GetProvidersRoute()
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &providers)

	// handle pagination
	nextKey := providers.Pagination.NextKey
	for nextKey != "" {
		res, err := utils.HttpQuery(constants.RESTAddr + route + "?pagination.key=" + nextKey)
		if err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}
		var nextPage akash.ProvidersResponse
		json.Unmarshal(res, &nextPage)

		// append to the response
		providers.Providers = append(providers.Providers, nextPage.Providers...)

		// update the next key
		nextKey = nextPage.Pagination.NextKey
	}

	zap.L().Info("", zap.Bool("Success", true), zap.String("Providers", fmt.Sprint(providers.Providers)))

	return providers
}

func (as *AkashService) IndexProviders(cfg config.Config, providers akash.ProvidersResponse) error {
	if cfg.Chain.Name != "akash" {
		return nil
	}

	for _, provider := range providers.Providers {
		_, err := as.DB.Exec(`
		INSERT INTO akash_providers (owner, host_uri, email, website)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (owner)
		DO UPDATE SET host_uri = $2, email = $3, website = $4`,
			provider.Owner, provider.HostURI, provider.Info.Email, provider.Info.Website)
		if err != nil {
			return fmt.Errorf("error indexing provider: %w", err)
		}

		for _, attribute := range provider.Attributes {
			_, err := as.DB.Exec(`
			INSERT INTO akash_provider_attributes (provider_owner, attribute_key, attribute_value)
			VALUES ($1, $2, $3)`,
				provider.Owner, attribute.Key, attribute.Value)
			if err != nil {
				return fmt.Errorf("error indexing provider attribute: %w", err)
			}
		}
	}

	return nil
}
