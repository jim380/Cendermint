package services_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/services"
	"github.com/jim380/Cendermint/types"
	"github.com/stretchr/testify/require"
)

func TestConsensusService_GetConsensusDump(t *testing.T) {
	data, err := os.ReadFile("../testutil/json/dump_consensus_state.json")
	require.NoError(t, err)

	var consensusState types.ConsensusState
	err = json.Unmarshal(data, &consensusState)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	originalRPCAddr := constants.RPCAddr
	constants.RPCAddr = server.URL
	defer func() { constants.RPCAddr = originalRPCAddr }()

	originalGetConspubMonikerMapWrapper := rest.GetConspubMonikerMapWrapper
	rest.GetConspubMonikerMapWrapper = func() map[string]string {
		return map[string]string{
			"6Nz09YGHzwWxjczG0IhK4Iv0qY2IcX0P/5KitvRXTUc=": "moniker1",
		}
	}
	defer func() { rest.GetConspubMonikerMapWrapper = originalGetConspubMonikerMapWrapper }()

	css := &services.ConsensusService{}
	cfg := config.Config{}
	rpc := &types.RPCData{}

	result := css.GetConsensusDump(cfg, rpc)

	t.Logf("Result: %+v", result)
	t.Logf("ConsensusState: %+v", consensusState)

	require.NotNil(t, result)
	require.Equal(t, "moniker1", result["CB5A63B91E8F4EE8DB935942CBE25724636479E0"][5])
	require.Equal(t, "❌", result["CB5A63B91E8F4EE8DB935942CBE25724636479E0"][3])
	require.Equal(t, "❌", result["CB5A63B91E8F4EE8DB935942CBE25724636479E0"][4])
	require.Equal(t, consensusState, rpc.ConsensusState)
	require.Equal(t, result, rpc.Validatorsets)
}
