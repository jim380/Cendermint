package services_test

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jim380/Cendermint/services"
	testtypes "github.com/jim380/Cendermint/testutil/types"
	"github.com/stretchr/testify/require"
)

func TestAbsentValidatorService_Index(t *testing.T) {
	tests := []struct {
		name        string
		data        testtypes.TestDataAbsentValidator
		mock        testtypes.MockData
		expectError bool
	}{
		{
			name: "Valid Insert",
			data: testtypes.TestDataAbsentValidator{
				Height:         12345,
				ConsAddrBase64: "testConsAddrBase64",
			},
			mock: testtypes.MockData{
				Query: `INSERT INTO absent_validators \(block_height, cons_pub_key\) VALUES \(\$1, \$2\) RETURNING cons_pub_key`,
				Args:  []interface{}{12345, "testConsAddrBase64"},
				Rows:  sqlmock.NewRows([]string{"cons_pub_key"}).AddRow("testConsAddrBase64"),
				Err:   nil,
			},
			expectError: false,
		},
		{
			name: "Insert Error",
			data: testtypes.TestDataAbsentValidator{
				Height:         12345,
				ConsAddrBase64: "testConsAddrBase64",
			},
			mock: testtypes.MockData{
				Query: `INSERT INTO absent_validators \(block_height, cons_pub_key\) VALUES \(\$1, \$2\) RETURNING cons_pub_key`,
				Args:  []interface{}{12345, "testConsAddrBase64"},
				Rows:  nil,
				Err:   sqlmock.ErrCancelled,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Initialize the service
			service := &services.AbsentValidatorService{}
			service.Init(db)

			// Convert mockArgs to []driver.Value
			args := make([]driver.Value, len(tt.mock.Args))
			for i, arg := range tt.mock.Args {
				args[i] = arg
			}

			// Set up the expected query and result
			query := mock.ExpectQuery(tt.mock.Query).WithArgs(args...)
			if tt.mock.Err != nil {
				query.WillReturnError(tt.mock.Err)
			} else {
				query.WillReturnRows(tt.mock.Rows)
			}

			// Call the method
			result, err := service.Index(tt.data.Height, tt.data.ConsAddrBase64)

			// Assertions
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tt.data.Height, result.BlockHeight)
				require.Equal(t, tt.data.ConsAddrBase64, result.ConsAddrBase64)
			}

			// Ensure all expectations were met
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
