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

func TestSlashingService_GetSlashingParams(t *testing.T) {
	slashingParamsData, err := os.ReadFile("../testutil/json/slashing_params.json")
	require.NoError(t, err, "Failed to read slashing_params.json")

	var slashingInfo types.SlashingInfo
	err = json.Unmarshal(slashingParamsData, &slashingInfo)
	require.NoError(t, err, "Failed to unmarshal slashing_params.json")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(slashingParamsData)
	}))
	defer server.Close()

	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	ss := &services.SlashingService{}
	cfg := config.Config{}
	rd := &types.RESTData{}

	ss.GetSlashingParams(cfg, rd)

	require.Equal(t, slashingInfo.Params, rd.Slashing.Params, "Slashing Params mismatch")
}

func TestSlashingService_GetSigningInfo(t *testing.T) {
	signingInfoData, err := os.ReadFile("../testutil/json/slashing_signing_info.json")
	require.NoError(t, err, "Failed to read slashing_signing_info.json")

	var slashingInfo types.SlashingInfo
	err = json.Unmarshal(signingInfoData, &slashingInfo)
	require.NoError(t, err, "Failed to unmarshal slashing_signing_info.json")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(signingInfoData)
	}))
	defer server.Close()

	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	ss := &services.SlashingService{}
	cfg := config.Config{}
	rd := &types.RESTData{}
	consAddr := "cosmosvalcons1px0zkz2cxvc6lh34uhafveea9jnaagckmrlsye" // Replace with a valid consensus address

	ss.GetSigningInfo(cfg, consAddr, rd)

	require.Equal(t, slashingInfo.ValSigning, rd.Slashing.ValSigning, "Signing Info mismatch")
}
