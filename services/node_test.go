package services_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/services"
	"github.com/jim380/Cendermint/types"
	"github.com/stretchr/testify/require"
)

func TestNodeService_GetInfo(t *testing.T) {
	// Load mock data from JSON file
	mockData, err := os.ReadFile("../testutil/json/node.json")
	require.NoError(t, err)

	// Create a new HTTP server with a mock handler for node info
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockData)
	}))
	defer server.Close()

	// Override the RESTAddr
	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	// Initialize service
	ns := &services.NodeService{}

	// Prepare input and output
	cfg := &config.Config{}
	rd := &types.RESTData{}

	// Call method
	ns.GetInfo(cfg, rd)

	// Parse the mock data to get the expected result
	var expectedData types.NodeInfo
	err = json.Unmarshal(mockData, &expectedData)
	require.NoError(t, err)

	// Debug statements
	t.Logf("Mock Data: %s", string(mockData))
	t.Logf("Expected Node Info: %+v", expectedData)
	t.Logf("Actual Node Info: %+v", rd.NodeInfo)

	// Assertions
	require.NotNil(t, rd.NodeInfo)
	require.Equal(t, expectedData.Application.Version, rd.NodeInfo.Application.Version)
	require.Equal(t, expectedData.Application.SDKVersion, rd.NodeInfo.Application.SDKVersion)
	require.Equal(t, expectedData.Default.NodeID, rd.NodeInfo.Default.NodeID)
	require.Equal(t, expectedData.Default.TMVersion, rd.NodeInfo.Default.TMVersion)
	require.Equal(t, expectedData.Default.Moniker, rd.NodeInfo.Default.Moniker)
}
