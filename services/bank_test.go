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

func TestBankService_GetBalanceInfo(t *testing.T) {
	data, err := os.ReadFile("../testutil/json/bank.json")
	require.NoError(t, err)

	var balances types.Balances
	err = json.Unmarshal(data, &balances)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	bs := &services.BankService{}
	cfg := config.Config{}
	rd := &types.RESTData{}

	bs.GetBalanceInfo(cfg, rd)

	require.Equal(t, balances.Balances, rd.Balances)
}

func TestBankService_GetRewardsCommissionInfo(t *testing.T) {
	data, err := os.ReadFile("../testutil/json/distribution.json")
	require.NoError(t, err)

	var rewardsAndCommission struct {
		Height string `json:"height"`
		Result struct {
			Operator_Address string       `json:"operator_address"`
			Selfbond_Rewards []types.Coin `json:"self_bond_rewards"`
			Commission       struct {
				Commission []types.Coin `json:"commission"`
			} `json:"val_commission"`
		} `json:"result"`
	}
	err = json.Unmarshal(data, &rewardsAndCommission)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	bs := &services.BankService{}
	cfg := config.Config{}
	rd := &types.RESTData{}

	bs.GetRewardsCommissionInfo(cfg, rd)

	require.Equal(t, rewardsAndCommission.Result.Selfbond_Rewards, rd.Rewards)
	require.Equal(t, rewardsAndCommission.Result.Commission.Commission, rd.Commission)
}
