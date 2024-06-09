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

func TestUpgradeService_GetInfo(t *testing.T) {
	mockData, err := os.ReadFile("../testutil/json/upgrade.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockData)
	}))
	defer server.Close()

	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	us := &services.UpgradeService{}

	cfg := config.Config{}
	rd := &types.RESTData{}

	us.GetInfo(cfg, rd)

	var expectedData types.UpgradeInfo
	err = json.Unmarshal(mockData, &expectedData)
	require.NoError(t, err)

	require.NotNil(t, rd.UpgradeInfo)
	require.Equal(t, expectedData.Plan, rd.UpgradeInfo.Plan)
	require.Equal(t, expectedData.Planned, rd.UpgradeInfo.Planned)
}
