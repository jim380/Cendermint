package types

import "github.com/DATA-DOG/go-sqlmock"

// TestData holds the input data for the tests
type TestData struct {
	Height         int
	ConsAddrBase64 string
}

// MockData holds the mock query, arguments, rows, and error
type MockData struct {
	Query string
	Args  []interface{}
	Rows  *sqlmock.Rows
	Err   error
}
