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
	mockData, err := os.ReadFile("../testutil/json/node.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockData)
	}))
	defer server.Close()

	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	ns := &services.NodeService{}

	cfg := &config.Config{}
	rd := &types.RESTData{}

	ns.GetInfo(cfg, rd)

	var expectedData types.NodeInfo
	err = json.Unmarshal(mockData, &expectedData)
	require.NoError(t, err)

	require.NotNil(t, rd.NodeInfo)
	require.Equal(t, expectedData.Application.Version, rd.NodeInfo.Application.Version)
	require.Equal(t, expectedData.Application.SDKVersion, rd.NodeInfo.Application.SDKVersion)
	require.Equal(t, expectedData.Default.NodeID, rd.NodeInfo.Default.NodeID)
	require.Equal(t, expectedData.Default.TMVersion, rd.NodeInfo.Default.TMVersion)
	require.Equal(t, expectedData.Default.Moniker, rd.NodeInfo.Default.Moniker)
}
