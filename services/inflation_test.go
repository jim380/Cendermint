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
	"github.com/jim380/Cendermint/utils"
	"github.com/stretchr/testify/require"
)

func TestInflationService_GetInfo(t *testing.T) {
	// Read inflation data
	inflationData, err := os.ReadFile("../testutil/json/inflation.json")
	require.NoError(t, err, "Failed to read inflation.json")

	var inflation struct {
		Inflation string `json:"inflation"`
	}
	err = json.Unmarshal(inflationData, &inflation)
	require.NoError(t, err, "Failed to unmarshal inflation.json")

	// Set up test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(inflationData)
	}))
	defer server.Close()

	// Override REST address
	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	// Initialize InflationService and call GetInfo
	is := &services.InflationService{}
	cfg := config.Config{}
	rd := &types.RESTData{}

	is.GetInfo(cfg, rd)

	// Assertions
	expectedInflation := utils.StringToFloat64(inflation.Inflation)
	require.Equal(t, expectedInflation, rd.Inflation, "Inflation mismatch")
}
