package services_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jim380/Cendermint/services"
	"github.com/stretchr/testify/require"
)

func TestValidatorService_Index(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Initialize the ValidatorService with the mock database
	vs := &services.ValidatorService{DB: db}

	// Define test inputs
	consPubKey := "e3BehnEIlGUAnJYn9V8gBXuMh4tXO8xxlxyXD1APGyk="
	consAddr := "cosmosvalcons1px0zkz2cxvc6lh34uhafveea9jnaagckmrlsye"
	consAddrHex := "099E2B09583331AFDE35E5FA96673D2CA7DEA316"
	moniker := "testMoniker"
	lastActive := time.Now()

	mock.ExpectExec(`INSERT INTO validators`).
		WithArgs(consPubKey, consAddr, consAddrHex, moniker, lastActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	validator, err := vs.Index(consPubKey, consAddr, consAddrHex, moniker, lastActive)

	require.NoError(t, err)
	require.NotNil(t, validator)
	require.Equal(t, consPubKey, validator.ConsPubKey)
	require.Equal(t, consAddr, validator.ConsAddr)
	require.Equal(t, consAddrHex, validator.ConsAddrHex)
	require.Equal(t, moniker, validator.Moniker)
	require.WithinDuration(t, lastActive, validator.LastActive, time.Second)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
