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

func TestDelegationService_GetInfo(t *testing.T) {
	mockData, err := os.ReadFile("../testutil/json/delegations.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockData)
	}))
	defer server.Close()

	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	ds := &services.DelegationService{}

	cfg := config.Config{}
	rd := &types.RESTData{}

	ds.GetInfo(cfg, rd)

	var expectedData types.DelegationsInfo
	err = json.Unmarshal(mockData, &expectedData)
	require.NoError(t, err)

	require.NotNil(t, rd.Delegations)
	require.Equal(t, expectedData, rd.Delegations)
	for i, delRes := range expectedData.DelegationRes {
		require.Equal(t, delRes.Balance.Amount, rd.Delegations.DelegationRes[i].Balance.Amount)
		require.Equal(t, delRes.Delegation.DelegatorAddr, rd.Delegations.DelegationRes[i].Delegation.DelegatorAddr)
		require.Equal(t, delRes.Delegation.ValidatorAddr, rd.Delegations.DelegationRes[i].Delegation.ValidatorAddr)
		require.Equal(t, delRes.Delegation.Shares, rd.Delegations.DelegationRes[i].Delegation.Shares)
	}
	require.Equal(t, expectedData.Pagination.Total, rd.Delegations.Pagination.Total)
	require.Equal(t, expectedData.Pagination.NextKey, rd.Delegations.Pagination.NextKey)
}
