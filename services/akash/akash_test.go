package akash_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/services/akash"
	"github.com/jim380/Cendermint/types"
	akash_types "github.com/jim380/Cendermint/types/akash"
	"github.com/stretchr/testify/require"
)

func TestGetAkashDeployments(t *testing.T) {
	// Read deployments data
	deploymentsData, err := os.ReadFile("../../testutil/json/akash_deployments.json")
	require.NoError(t, err, "Failed to read akash_deployments.json")

	// Read active deployments data
	activeDeploymentsData, err := os.ReadFile("../../testutil/json/akash_active_deployments.json")
	require.NoError(t, err, "Failed to read akash_active_deployments.json")

	// Set up test server for deployments
	deploymentsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("filters.state") == "active" {
			w.WriteHeader(http.StatusOK)
			w.Write(activeDeploymentsData)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(deploymentsData)
		}
	}))
	defer deploymentsServer.Close()

	// Override REST address for deployments
	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = deploymentsServer.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	// Initialize AkashService and call GetAkashDeployments
	as := &akash.AkashService{}
	cfg := config.Config{Chain: config.Chain{Name: "akash"}}
	data := &types.AsyncData{}

	deployments := as.GetAkashDeployments(cfg, data)

	// Assertions
	require.NotNil(t, deployments, "Deployments should not be nil")
	require.Equal(t, 306819, data.AkashInfo.TotalDeployments, "TotalDeployments mismatch")
	require.Equal(t, 7074, data.AkashInfo.ActiveDeployments, "ActiveDeployments mismatch")
	require.Equal(t, 306819-7074, data.AkashInfo.ClosedDeployments, "ClosedDeployments mismatch")
}

func TestIndexProviders(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Initialize AkashService with mock database
	as := &akash.AkashService{DB: db}
	cfg := config.Config{Chain: config.Chain{Name: "akash"}}

	// Sample providers data
	providers := akash_types.ProvidersResponse{
		Providers: []akash_types.Provider{
			{
				Owner:   "akash1qphp30grcwp369uf5x5jsen56q2vtgc5l40prd",
				HostURI: "https://provider.spheron.wiki:8443",
				Attributes: []akash_types.Attribute{
					{Key: "region", Value: "eu-east"},
					{Key: "host", Value: "akash"},
					{Key: "tier", Value: "community"},
					{Key: "organization", Value: "spheron"},
					{Key: "capabilities/storage/1/class", Value: "beta3"},
					{Key: "capabilities/storage/1/persistent", Value: "true"},
				},
				Info: akash_types.Info{Email: "", Website: ""},
			},
			{
				Owner:   "akash1qzkvrygsazvwyn7jm0r3efs26qkxde635cgvls",
				HostURI: "https://provider-test-14.testcoders.com:8443",
				Attributes: []akash_types.Attribute{
					{Key: "host", Value: "Testcoders Test 14"},
				},
				Info: akash_types.Info{Email: "", Website: ""},
			},
		},
	}

	// Expect database interactions
	for _, provider := range providers.Providers {
		mock.ExpectExec(`INSERT INTO akash_providers`).
			WithArgs(provider.Owner, provider.HostURI, provider.Info.Email, provider.Info.Website).
			WillReturnResult(sqlmock.NewResult(1, 1))

		for _, attribute := range provider.Attributes {
			mock.ExpectExec(`INSERT INTO akash_provider_attributes`).
				WithArgs(provider.Owner, attribute.Key, attribute.Value).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
	}

	// Call IndexProviders
	err = as.IndexProviders(cfg, providers)
	require.NoError(t, err)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestFindProviderOwnersPendingAuditorUpdate(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Initialize AkashService with mock database
	as := &akash.AkashService{DB: db}

	// Define the period
	period := 30 * 24 * time.Hour // 30 days

	// Sample data to return from the query
	rows := sqlmock.NewRows([]string{"owner"}).
		AddRow("akash1qphp30grcwp369uf5x5jsen56q2vtgc5l40prd").
		AddRow("akash1qzkvrygsazvwyn7jm0r3efs26qkxde635cgvls")

	// Expect the query to be executed
	mock.ExpectQuery(`SELECT owner FROM \(
			SELECT DISTINCT owner FROM akash_providers 
			WHERE owner NOT IN \(
				SELECT DISTINCT provider_owner FROM akash_provider_auditors
			\)
			UNION
			SELECT provider_owner FROM akash_provider_auditors 
			WHERE last_updated < \$1
		\) AS subquery`).
		WithArgs(sqlmock.AnyArg()). // Use AnyArg to match any time argument
		WillReturnRows(rows)

	// Call FindProviderOwnersPendingAuditorUpdate
	providerOwners, err := as.FindProviderOwnersPendingAuditorUpdate(period)
	require.NoError(t, err)

	// Assertions
	require.Equal(t, 2, len(providerOwners), "Provider owners count mismatch")
	require.Equal(t, "akash1qphp30grcwp369uf5x5jsen56q2vtgc5l40prd", providerOwners[0], "First provider owner mismatch")
	require.Equal(t, "akash1qzkvrygsazvwyn7jm0r3efs26qkxde635cgvls", providerOwners[1], "Second provider owner mismatch")

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestIndexAuditorForProviderOwners(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Initialize AkashService with mock database
	as := &akash.AkashService{DB: db}
	cfg := config.Config{Chain: config.Chain{Name: "akash"}}

	// Sample provider owners
	providerOwners := []string{
		"akash1qphp30grcwp369uf5x5jsen56q2vtgc5l40prd",
		"akash1qzkvrygsazvwyn7jm0r3efs26qkxde635cgvls",
	}

	// Sample auditors data
	auditorsData := akash_types.AuditorsResponse{
		Providers: []struct {
			Auditor string `json:"auditor"`
		}{
			{Auditor: "auditor1"},
			{Auditor: "auditor2"},
		},
		Pagination: akash_types.Pagination{NextKey: "", Total: "2"},
	}
	auditorsDataJSON, err := json.Marshal(auditorsData)
	require.NoError(t, err, "Failed to marshal auditors data")

	// Set up test server for auditors
	auditorsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(auditorsDataJSON)
	}))
	defer auditorsServer.Close()

	// Override REST address for auditors
	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = auditorsServer.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	// Expect database interactions
	for _, owner := range providerOwners {
		for _, auditor := range auditorsData.Providers {
			mock.ExpectExec(`INSERT INTO akash_provider_auditors`).
				WithArgs(owner, auditor.Auditor).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
	}

	// Call IndexAuditorForProviderOwners
	err = as.IndexAuditorForProviderOwners(cfg, providerOwners)
	require.NoError(t, err)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestFindProviderOwnersPendingDeploymentUpdate(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Initialize AkashService with mock database
	as := &akash.AkashService{DB: db}

	// Define the period
	period := 30 * 24 * time.Hour // 30 days

	// Sample data to return from the query
	rows := sqlmock.NewRows([]string{"owner"}).
		AddRow("akash1qphp30grcwp369uf5x5jsen56q2vtgc5l40prd").
		AddRow("akash1qzkvrygsazvwyn7jm0r3efs26qkxde635cgvls")

	// Expect the query to be executed
	mock.ExpectQuery(`SELECT owner FROM \(
			SELECT DISTINCT owner FROM akash_providers 
			WHERE owner NOT IN \(
				SELECT DISTINCT owner FROM akash_deployments
			\)
			UNION
			SELECT owner FROM akash_deployments 
			WHERE last_updated < \$1
		\) AS subquery`).
		WithArgs(sqlmock.AnyArg()). // Use AnyArg to match any time argument
		WillReturnRows(rows)

	// Call FindProviderOwnersPendingDeploymentUpdate
	providerOwners, err := as.FindProviderOwnersPendingDeploymentUpdate(period)
	require.NoError(t, err)

	// Assertions
	require.Equal(t, 2, len(providerOwners), "Provider owners count mismatch")
	require.Equal(t, "akash1qphp30grcwp369uf5x5jsen56q2vtgc5l40prd", providerOwners[0], "First provider owner mismatch")
	require.Equal(t, "akash1qzkvrygsazvwyn7jm0r3efs26qkxde635cgvls", providerOwners[1], "Second provider owner mismatch")

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
