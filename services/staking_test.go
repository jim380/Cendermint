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
	"go.uber.org/zap"
)

type MockStakingService struct {
	services.StakingService
	MockGetTotalSupply func(cfg config.Config, denom string, log *zap.Logger) float64
}

func TestStakingService_GetInfo(t *testing.T) {
	stakingPoolData, err := os.ReadFile("../testutil/json/staking_pool.json")
	require.NoError(t, err)

	coinSupplyData, err := os.ReadFile("../testutil/json/coin_supply.json")
	require.NoError(t, err)

	stakingPoolServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(stakingPoolData)
	}))
	defer stakingPoolServer.Close()

	coinSupplyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(coinSupplyData)
	}))
	defer coinSupplyServer.Close()

	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = stakingPoolServer.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	mss := &MockStakingService{
		MockGetTotalSupply: func(cfg config.Config, denom string, log *zap.Logger) float64 {
			res, err := utils.HttpQuery(coinSupplyServer.URL + "/supply/" + denom)
			if err != nil {
				log.Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
			}
			var ts types.Supply
			json.Unmarshal(res, &ts)
			return utils.StringToFloat64(ts.Amount.Amount)
		},
	}

	cfg := config.Config{}
	rd := &types.RESTData{}

	mss.GetInfo(cfg, "uatom", rd)

	var expectedStakingPool types.StakingPool
	err = json.Unmarshal(stakingPoolData, &expectedStakingPool)
	require.NoError(t, err)

	var expectedCoinSupply types.Supply
	err = json.Unmarshal(coinSupplyData, &expectedCoinSupply)
	require.NoError(t, err)

	require.NotNil(t, rd.StakingPool)
	require.Equal(t, expectedStakingPool.Pool.Bonded_tokens, rd.StakingPool.Pool.Bonded_tokens)
	require.Equal(t, expectedStakingPool.Pool.Not_bonded_tokens, rd.StakingPool.Pool.Not_bonded_tokens)
	require.Equal(t, utils.StringToFloat64(expectedCoinSupply.Amount.Amount), rd.StakingPool.Pool.Total_supply)
}
