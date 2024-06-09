package services_test

import (
	"database/sql"
	"database/sql/driver"
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
	testtypes "github.com/jim380/Cendermint/testutil/types"
	"github.com/jim380/Cendermint/types"
	"github.com/stretchr/testify/require"
)

func TestBlockService_GetInfo(t *testing.T) {
	data, err := os.ReadFile("../testutil/json/block.json")
	require.NoError(t, err)

	var expectedBlock types.Blocks
	err = json.Unmarshal(data, &expectedBlock)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	bs := &services.BlockService{}
	cfg := config.Config{}

	result := bs.GetInfo(cfg)

	require.Equal(t, expectedBlock, result)
}

func TestBlockService_GetLastBlockTimestamp(t *testing.T) {
	data, err := os.ReadFile("../testutil/json/block.json")
	require.NoError(t, err)

	var lastBlock types.LastBlock
	err = json.Unmarshal(data, &lastBlock)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = server.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	bs := &services.BlockService{
		Block: &types.Blocks{},
	}
	cfg := config.Config{}
	currentHeight := int64(20783592)

	result := bs.GetLastBlockTimestamp(cfg, currentHeight)

	require.Equal(t, lastBlock.Block.Header.Timestamp, result.Block.Header.LastTimestamp)
}

func TestBlockService_Index(t *testing.T) {
	fixedTime := time.Date(2024, time.June, 8, 17, 38, 23, 0, time.UTC)

	tests := []struct {
		name        string
		data        testtypes.TestDataBlock
		mock        testtypes.MockData
		expectError bool
	}{
		{
			name: "Valid Insert",
			data: testtypes.TestDataBlock{
				Height:    12345,
				Hash:      "somehash",
				Timestamp: fixedTime,
				Proposer:  "someproposer",
				TxnCount:  10,
			},
			mock: testtypes.MockData{
				Query: `INSERT INTO blocks \(height, block_hash, timestamp, proposer_address, txn_count\) VALUES \(\$1, \$2, \$3, \$4, \$5\) ON CONFLICT \(height\) DO NOTHING RETURNING block_hash`,
				Args:  []interface{}{12345, "somehash", fixedTime, "someproposer", 10},
				Rows:  sqlmock.NewRows([]string{"block_hash"}).AddRow("somehash"),
				Err:   nil,
			},
			expectError: false,
		},
		{
			name: "Insert Error",
			data: testtypes.TestDataBlock{
				Height:    12345,
				Hash:      "somehash",
				Timestamp: fixedTime,
				Proposer:  "someproposer",
				TxnCount:  10,
			},
			mock: testtypes.MockData{
				Query: `INSERT INTO blocks \(height, block_hash, timestamp, proposer_address, txn_count\) VALUES \(\$1, \$2, \$3, \$4, \$5\) ON CONFLICT \(height\) DO NOTHING RETURNING block_hash`,
				Args:  []interface{}{12345, "somehash", fixedTime, "someproposer", 10},
				Rows:  nil,
				Err:   sql.ErrConnDone,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			bs := &services.BlockService{DB: db}

			args := make([]driver.Value, len(tt.mock.Args))
			for i, arg := range tt.mock.Args {
				args[i] = arg
			}

			query := mock.ExpectQuery(tt.mock.Query).WithArgs(args...)
			if tt.mock.Err != nil {
				query.WillReturnError(tt.mock.Err)
			} else {
				query.WillReturnRows(tt.mock.Rows)
			}

			result, err := bs.Index(tt.data.Height, tt.data.Hash, tt.data.Timestamp, tt.data.Proposer, tt.data.TxnCount)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tt.data.Height, result.Height)
				require.Equal(t, tt.data.Hash, result.BlockHash)
				require.Equal(t, tt.data.Timestamp, result.Timestamp)
				require.Equal(t, tt.data.Proposer, result.Proposer)
				require.Equal(t, tt.data.TxnCount, result.TxnCount)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
