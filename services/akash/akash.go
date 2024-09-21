package akash

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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
	if err := json.Unmarshal(res, &deployments); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

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
	if err := json.Unmarshal(resActive, &activeDeployments); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

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
	if err := json.Unmarshal(res, &providers); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	// handle pagination
	nextKey := providers.Pagination.NextKey
	for nextKey != "" {
		res, err := utils.HttpQuery(constants.RESTAddr + route + "?pagination.key=" + nextKey)
		if err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}
		var nextPage akash.ProvidersResponse
		if err := json.Unmarshal(res, &nextPage); err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}

		// append to the response
		providers.Providers = append(providers.Providers, nextPage.Providers...)

		// update the next key
		nextKey = nextPage.Pagination.NextKey
	}

	return providers
}

func (as *AkashService) IndexProviders(cfg config.Config, providers akash.ProvidersResponse) error {
	if cfg.Chain.Name != "akash" {
		return nil
	}

	for _, provider := range providers.Providers {
		_, err := as.DB.Exec(`
		INSERT INTO akash_providers (owner, host_uri, email, website, last_updated)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		ON CONFLICT (owner)
		DO UPDATE SET host_uri = $2, email = $3, website = $4, last_updated = CURRENT_TIMESTAMP`,
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

func (as *AkashService) FindProviderOwnersPendingAuditorUpdate(period time.Duration) ([]string, error) {
	cutoff := time.Now().Add(-period)

	rows, err := as.DB.Query(`
		SELECT owner FROM (
			SELECT DISTINCT owner FROM akash_providers 
			WHERE owner NOT IN (
				SELECT DISTINCT provider_owner FROM akash_provider_auditors
			)
			UNION
			SELECT provider_owner FROM akash_provider_auditors 
			WHERE last_updated < $1
		) AS subquery
	`, cutoff)
	if err != nil {
		return nil, fmt.Errorf("error querying providers pending auditor update: %w", err)
	}
	defer rows.Close()

	var providerOwners []string
	for rows.Next() {
		var owner string
		err := rows.Scan(&owner)
		if err != nil {
			return nil, fmt.Errorf("error scanning provider owner: %w", err)
		}
		providerOwners = append(providerOwners, owner)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return providerOwners, nil
}

func (as *AkashService) IndexAuditorForProviderOwners(cfg config.Config, providerOwners []string) error {
	if cfg.Chain.Name != "akash" {
		return nil
	}

	for _, owner := range providerOwners {
		var auditors akash.AuditorsResponse

		route := rest.GetAuditorsForProviderOwnerRoute(owner)
		res, err := utils.HttpQuery(constants.RESTAddr + route)
		if err != nil {
			return fmt.Errorf("error querying auditors for provider owner: %w", err)
		}
		if err := json.Unmarshal(res, &auditors); err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}

		// handle pagination
		nextKey := auditors.Pagination.NextKey
		for nextKey != "" {
			res, err := utils.HttpQuery(constants.RESTAddr + route + "?pagination.key=" + nextKey)
			if err != nil {
				zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
			}
			var nextPage akash.AuditorsResponse
			if err := json.Unmarshal(res, &nextPage); err != nil {
				zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
			}

			// append to the response
			auditors.Providers = append(auditors.Providers, nextPage.Providers...)

			// update the next key
			nextKey = nextPage.Pagination.NextKey
		}

		for _, provider := range auditors.Providers {
			_, err = as.DB.Exec(`
			INSERT INTO akash_provider_auditors (provider_owner, auditor, last_updated) 
			VALUES ($1, $2, CURRENT_TIMESTAMP)
			ON CONFLICT (provider_owner, auditor)
			DO UPDATE SET last_updated = CURRENT_TIMESTAMP`,
				owner, provider.Auditor)
			if err != nil {
				return fmt.Errorf("error indexing provider auditor: %w", err)
			}
		}
	}

	return nil
}

func (as *AkashService) FindProviderOwnersPendingDeploymentUpdate(period time.Duration) ([]string, error) {
	cutoff := time.Now().Add(-period)

	rows, err := as.DB.Query(`
		SELECT owner FROM (
			SELECT DISTINCT owner FROM akash_providers 
			WHERE owner NOT IN (
				SELECT DISTINCT owner FROM akash_deployments
			)
			UNION
			SELECT owner FROM akash_deployments 
			WHERE last_updated < $1
		) AS subquery
	`, cutoff)
	if err != nil {
		return nil, fmt.Errorf("error querying providers pending deployment update: %w", err)
	}
	defer rows.Close()

	var providerOwners []string
	for rows.Next() {
		var owner string
		err := rows.Scan(&owner)
		if err != nil {
			return nil, fmt.Errorf("error scanning provider owner: %w", err)
		}
		providerOwners = append(providerOwners, owner)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return providerOwners, nil
}

func (as *AkashService) IndexDeploymentForProviderOwner(cfg config.Config, providerOwners []string) error {
	if cfg.Chain.Name != "akash" {
		return nil
	}

	for _, owner := range providerOwners {
		var deployments akash.DeploymentsResponse

		route := rest.GetDeploymentsForProviderOwnerRoute(owner)
		res, err := utils.HttpQuery(constants.RESTAddr + route)
		if err != nil {
			return fmt.Errorf("error querying deployments for provider owner: %w", err)
		}
		if err := json.Unmarshal(res, &deployments); err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}

		// handle pagination
		nextKey := deployments.Pagination.NextKey
		for nextKey != "" {
			res, err := utils.HttpQuery(constants.RESTAddr + route + "&pagination.key=" + nextKey)
			if err != nil {
				zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
			}
			var nextPage akash.DeploymentsResponse
			if err := json.Unmarshal(res, &nextPage); err != nil {
				zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
			}

			// append to the response
			deployments.Deployments = append(deployments.Deployments, nextPage.Deployments...)

			// update the next key
			nextKey = nextPage.Pagination.NextKey
		}

		for _, deployment := range deployments.Deployments {
			// index deployments
			_, err = as.DB.Exec(`
			INSERT INTO akash_deployments (owner, dseq, state, version, created_at, last_updated) 
			VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
			ON CONFLICT (owner, dseq)
			DO UPDATE SET state = $3, version = $4, created_at = $5, last_updated = CURRENT_TIMESTAMP`,
				deployment.DeploymentDetails.DeploymentId.Owner, deployment.DeploymentDetails.DeploymentId.Dseq, deployment.DeploymentDetails.State, deployment.DeploymentDetails.Version, deployment.DeploymentDetails.CreatedAt)
			if err != nil {
				return fmt.Errorf("error indexing deployment: %w", err)
			}

			// index groups
			for _, group := range deployment.Groups {
				_, err = as.DB.Exec(`
				INSERT INTO akash_groups (owner, dseq, gseq, state, name, created_at) 
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (owner, dseq, gseq)
				DO UPDATE SET state = $4, name = $5, created_at = $6`,
					deployment.DeploymentDetails.DeploymentId.Owner, deployment.DeploymentDetails.DeploymentId.Dseq, group.GroupId.Gseq, group.State, group.GroupSpec.Name, group.CreatedAt)
				if err != nil {
					return fmt.Errorf("error indexing group: %w", err)
				}

				// Index group requirements
				for _, signedByAnyOf := range group.GroupSpec.Requirements.SignedBy.AnyOf {
					_, err = as.DB.Exec(`
					INSERT INTO akash_group_requirements_signed_by_any_of (group_dseq, signed_by_any_of) 
					VALUES ($1, $2)
					ON CONFLICT (group_dseq)
					DO UPDATE SET signed_by_any_of = $2`,
						group.GroupId.Dseq, signedByAnyOf)
					if err != nil {
						return fmt.Errorf("error indexing group requirement signed by any of: %w", err)
					}
				}

				for _, signedByAllOf := range group.GroupSpec.Requirements.SignedBy.AllOf {
					_, err = as.DB.Exec(`
					INSERT INTO akash_group_requirements_signed_by_all_of (group_dseq, signed_by_all_of) 
					VALUES ($1, $2)
					ON CONFLICT (group_dseq)
					DO UPDATE SET signed_by_all_of = $2`,
						group.GroupId.Dseq, signedByAllOf)
					if err != nil {
						return fmt.Errorf("error indexing group requirement signed by all of: %w", err)
					}
				}

				for _, attribute := range group.GroupSpec.Requirements.Attributes {
					_, err = as.DB.Exec(`
						INSERT INTO akash_group_requirements_attributes (group_dseq, attribute_key, attribute_value) 
						VALUES ($1, $2, $3)
						ON CONFLICT (group_dseq)
						DO UPDATE SET attribute_key = $2, attribute_value = $3`,
						group.GroupId.Dseq, attribute.Key, attribute.Value)
					if err != nil {
						return fmt.Errorf("error indexing group requirement attribute: %w", err)
					}
				}

				// Index resources
				for _, resource := range group.GroupSpec.Resources {
					dseq, _ := strconv.Atoi(group.GroupId.Dseq)
					_, err = as.DB.Exec(`
						INSERT INTO akash_resources (group_dseq, cpu_units, memory_quantity, gpu_units, price_denom, price_amount) 
						VALUES ($1, $2, $3, $4, $5, $6)
						ON CONFLICT (group_dseq)
						DO UPDATE SET cpu_units = $2, memory_quantity = $3, gpu_units = $4, price_denom = $5, price_amount = $6`,
						dseq, resource.ResourceDetails.CPU.Units.Val, resource.ResourceDetails.Memory.Quantity.Val, resource.ResourceDetails.GPU.Units.Val, resource.Price.Denom, resource.Price.Amount)
					if err != nil {
						return fmt.Errorf("error indexing resource: %w", err)
					}

					// Index resource endpoints and storage
					for _, endpoint := range resource.ResourceDetails.Endpoints {
						_, err = as.DB.Exec(`
							INSERT INTO akash_resource_endpoints (group_dseq, kind, sequence_number) 
							VALUES ($1, $2, $3)`,
							dseq, endpoint.Kind, endpoint.Sequence_number)
						if err != nil {
							return fmt.Errorf("error indexing resource endpoint: %w", err)
						}
					}

					for _, storage := range resource.ResourceDetails.Storage {
						_, err = as.DB.Exec(`
							INSERT INTO akash_resource_storage (group_dseq, name, quantity) 
							VALUES ($1, $2, $3)`,
							dseq, storage.Name, storage.Quantity.Val)
						if err != nil {
							return fmt.Errorf("error indexing resource storage: %w", err)
						}
					}
				}
			}

			// index escrow account
			escrow := deployment.EscrowAccount
			_, err = as.DB.Exec(`
			INSERT INTO akash_escrow_accounts (id_scope, id_xid, owner, state, balance_denom, balance_amount, transferred_denom, transferred_amount, settled_at, depositor, funds_denom, funds_amount) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			ON CONFLICT (id_scope, id_xid)
			DO UPDATE SET owner = $3, state = $4, balance_denom = $5, balance_amount = $6, transferred_denom = $7, transferred_amount = $8, settled_at = $9, depositor = $10, funds_denom = $11, funds_amount = $12`,
				escrow.ID.Scope, escrow.ID.Xid, escrow.Owner, escrow.State, escrow.Balance.Denom, escrow.Balance.Amount, escrow.Transferred.Denom, escrow.Transferred.Amount, escrow.SettledAt, escrow.Depositor, escrow.Funds.Denom, escrow.Funds.Amount)
			if err != nil {
				return fmt.Errorf("error indexing escrow account: %w", err)
			}
		}
	}

	return nil
}
