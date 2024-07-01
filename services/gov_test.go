package services_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/services"
	"github.com/jim380/Cendermint/types"
	"github.com/stretchr/testify/require"
)

func TestGovService_GetInfo(t *testing.T) {
	constants.OperAddr = "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn"

	// Read proposal data
	proposalData, err := os.ReadFile("../testutil/json/proposal.json")
	require.NoError(t, err, "Failed to read proposal.json")

	var gov types.Gov
	err = json.Unmarshal(proposalData, &gov)
	require.NoError(t, err, "Failed to unmarshal proposal.json")

	// Read vote data
	voteData, err := os.ReadFile("../testutil/json/proposal_votes.json")
	require.NoError(t, err, "Failed to read proposal_votes.json")

	// Set up test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/votes/") {
			w.WriteHeader(http.StatusOK)
			w.Write(voteData)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(proposalData)
		}
	}))
	defer server.Close()

	// Override REST address
	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	// Initialize GovService and call GetInfo
	gs := &services.GovService{}
	cfg := config.Config{}
	rd := &types.RESTData{}

	gs.GetInfo(cfg, rd)

	// Assertions
	require.Equal(t, float64(len(gov.Proposals)), rd.Gov.TotalProposalCount, "TotalProposalCount mismatch")
	require.Equal(t, float64(1), rd.Gov.VotingProposalCount, "VotingProposalCount mismatch")         // Adjust as per your test data
	require.Equal(t, float64(1), rd.Gov.InVotingVotedCount, "InVotingVotedCount mismatch")           // Adjust as per your test data
	require.Equal(t, float64(0), rd.Gov.InVotingDidNotVoteCount, "InVotingDidNotVoteCount mismatch") // Adjust as per your test data
}
