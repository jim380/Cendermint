package services_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/services"
	"github.com/jim380/Cendermint/types"
	"github.com/stretchr/testify/require"
)

func TestTxnService_Index(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	ts := &services.TxnService{DB: db}

	cfg := config.Config{}
	height := int64(12345)
	timestamp := time.Now().Format(time.RFC3339)
	txInfo := types.TxInfo{
		Txs: []types.Txs{
			{
				Body: struct {
					Messages []struct {
						Type string `json:"@type"`
					} `json:"messages"`
					Memo string `json:"memo"`
				}{
					Messages: []struct {
						Type string `json:"@type"`
					}{
						{Type: "test message type"},
					},
					Memo: "test memo",
				},
				AuthInfo: struct {
					Fee struct {
						Amount []struct {
							Denom  string `json:"denom"`
							Amount string `json:"amount"`
						} `json:"amount"`
						GasLimit string `json:"gas_limit"`
						Payer    string `json:"payer"`
						Granter  string `json:"granter"`
					} `json:"fee"`
				}{
					Fee: struct {
						Amount []struct {
							Denom  string `json:"denom"`
							Amount string `json:"amount"`
						} `json:"amount"`
						GasLimit string `json:"gas_limit"`
						Payer    string `json:"payer"`
						Granter  string `json:"granter"`
					}{
						Amount: []struct {
							Denom  string `json:"denom"`
							Amount string `json:"amount"`
						}{
							{Denom: "test denom", Amount: "100"},
						},
						Payer:   "test payer",
						Granter: "test granter",
					},
				},
			},
		},
		TxResp: []types.TxResp{
			{
				Hash:      "test hash",
				Timestamp: timestamp,
				Tx: struct {
					Type string `json:"@type"`
				}{
					Type: "test tx type",
				},
				GasWanted: "200",
				GasUsed:   "150",
			},
		},
	}

	mock.ExpectBegin()

	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs("test hash", height, timestamp, "test tx type", "200", "150", "test memo", "test payer", "test granter").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO transaction_messages`).
		WithArgs("test message type", "test hash").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(`INSERT INTO denoms`).
		WithArgs("test denom").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(`INSERT INTO transaction_fee_amounts`).
		WithArgs("100", "test hash", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = ts.Index(cfg, height, txInfo)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestTxnService_GetTxnsInBlock(t *testing.T) {
	mockData, err := os.ReadFile("../testutil/json/tx.json")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockData)
	}))
	defer server.Close()

	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	cfg := config.Config{}
	ts := &services.TxnService{}

	txInfo, err := ts.GetTxnsInBlock(cfg, 20796813)
	require.NoError(t, err)

	var expectedTxInfo types.TxInfo
	err = json.Unmarshal(mockData, &expectedTxInfo)
	require.NoError(t, err)

	require.NotNil(t, txInfo)
	require.Equal(t, expectedTxInfo, txInfo)
}
