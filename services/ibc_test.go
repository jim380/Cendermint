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

func TestGetChannelInfo(t *testing.T) {
	// Load sample IBC channels data from ibc_channels.json
	data, err := os.ReadFile("../testutil/json/ibc_channels.json")
	require.NoError(t, err)

	var ibcChannelsData types.IbcChannelInfo
	err = json.Unmarshal(data, &ibcChannelsData)
	require.NoError(t, err)

	// Set up test server for IBC channels
	ibcChannelsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer ibcChannelsServer.Close()

	// Override REST address for IBC channels
	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = ibcChannelsServer.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	// Initialize IbcService
	is := &services.IbcService{}

	// Sample RESTData
	rd := &types.RESTData{}

	// Sample config
	cfg := config.Config{}

	// Call GetChannelInfo
	is.GetChannelInfo(cfg, rd)

	// Assertions
	require.Equal(t, 2, len(rd.IBC.IBCChannels), "IBC channels count mismatch")
	require.Equal(t, "STATE_OPEN", rd.IBC.IBCChannels["channel-882"][0], "First channel state mismatch")
	require.Equal(t, "ORDER_ORDERED", rd.IBC.IBCChannels["channel-882"][1], "First channel ordering mismatch")
	require.Equal(t, "channel-234", rd.IBC.IBCChannels["channel-882"][2], "First channel counterparty channel ID mismatch")
	require.Equal(t, "STATE_OPEN", rd.IBC.IBCChannels["channel-370"][0], "Second channel state mismatch")
	require.Equal(t, "ORDER_ORDERED", rd.IBC.IBCChannels["channel-370"][1], "Second channel ordering mismatch")
	require.Equal(t, "channel-2", rd.IBC.IBCChannels["channel-370"][2], "Second channel counterparty channel ID mismatch")
	require.Equal(t, 2, rd.IBC.IBCInfo.IbcChannelInfo.OpenChannels, "Open channels count mismatch")
}

func TestGetConnectionInfo(t *testing.T) {
	// Load sample IBC connections data from ibc_connections.json
	data, err := os.ReadFile("../testutil/json/ibc_connections.json")
	require.NoError(t, err)

	var ibcConnectionsData types.IbcConnectionInfo
	err = json.Unmarshal(data, &ibcConnectionsData)
	require.NoError(t, err)

	// Set up test server for IBC connections
	ibcConnectionsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer ibcConnectionsServer.Close()

	// Override REST address for IBC connections
	originalRESTAddr := constants.RESTAddr
	constants.RESTAddr = ibcConnectionsServer.URL
	defer func() { constants.RESTAddr = originalRESTAddr }()

	// Initialize IbcService
	is := &services.IbcService{}

	// Sample RESTData
	rd := &types.RESTData{}

	// Sample config
	cfg := config.Config{}

	// Call GetConnectionInfo
	is.GetConnectionInfo(cfg, rd)

	// Assertions
	require.Equal(t, 2, len(rd.IBC.IBCConnections), "IBC connections count mismatch")
	require.Equal(t, "STATE_TRYOPEN", rd.IBC.IBCConnections["connection-0"][0], "First connection state mismatch")
	require.Equal(t, "07-tendermint-1", rd.IBC.IBCConnections["connection-0"][1], "First connection client ID mismatch")
	require.Equal(t, "connection-0", rd.IBC.IBCConnections["connection-0"][2], "First connection counterparty connection ID mismatch")
	require.Equal(t, "07-tendermint-0", rd.IBC.IBCConnections["connection-0"][3], "First connection counterparty client ID mismatch")
	require.Equal(t, "STATE_OPEN", rd.IBC.IBCConnections["connection-1"][0], "Second connection state mismatch")
	require.Equal(t, "07-tendermint-1", rd.IBC.IBCConnections["connection-1"][1], "Second connection client ID mismatch")
	require.Equal(t, "connection-0", rd.IBC.IBCConnections["connection-1"][2], "Second connection counterparty connection ID mismatch")
	require.Equal(t, "07-tendermint-0", rd.IBC.IBCConnections["connection-1"][3], "Second connection counterparty client ID mismatch")
	require.Equal(t, 1, rd.IBC.IBCInfo.IbcConnectionInfo.OpenConnections, "Open connections count mismatch")
}
