package rest

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type ibcInfo struct {
	ibcChannelInfo
	ibcConnectionInfo
}

type ibcChannelInfo struct {
	OpenChannels int
	IBCChannels  ibcChannels `json:"channels"`
	Pagination   struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}
	Height struct {
		RevisionNumber string `json:"revision_number"`
		RevisionHeight string `json:"revision_height"`
	}
}

type ibcConnectionInfo struct {
	OpenConnections int
	IBConnections   ibcConnections `json:"connections"`
	Pagination      struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}
	Height struct {
		RevisionNumber string `json:"revision_number"`
		RevisionHeight string `json:"revision_height"`
	}
}

type ibcChannels []struct {
	State          string              `json:"state"`
	Ordering       string              `json:"ordering"`
	Counterparty   counterpartyChannel `json:"counterparty"`
	ConnectionHops []struct {
		string
	}
	Version   string `json:"version"`
	PortID    string `json:"port_id"`
	ChannelID string `json:"channel_id"`
}

type ibcConnections []struct {
	ID           string                 `json:"id"`
	ClientID     string                 `json:"client_id"`
	Versions     connectionVersions     `json:"versions"`
	State        string                 `json:"state"`
	Counterparty counterpartyConnection `json:"counterparty"`
	DelayPeriod  string                 `json:"delay_period"`
}

type connectionVersions []struct {
	Identifier string `json:"identifier"`
	Features   []struct {
		string
	}
}

type counterpartyChannel struct {
	PortID    string `json:"port_id"`
	ChannelID string `json:"channel_id"`
}

type counterpartyConnection struct {
	ClientID     string `json:"client_id"`
	ConnectionID string `json:"connection_id"`
	Prefix       struct {
		KeyPrefix string `json:"key_prefix"`
	}
}

func (rd *RESTData) getIBCChannels() {
	var ibcInfo ibcChannelInfo
	var ibcChannels map[string][]string = make(map[string][]string)
	ibcInfo.OpenChannels = 0
	res, err := RESTQuery("/ibc/core/channel/v1beta1/channels" + "?pagination.limit=1000")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &ibcInfo)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Active IBC Channels", fmt.Sprint(len(ibcInfo.IBCChannels))))
	}

	for _, value := range ibcInfo.IBCChannels {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// precommit failure validator
				}
			}()
			ibcChannels[value.ChannelID] = []string{value.State, value.Ordering, value.Counterparty.ChannelID}
			if value.State == "STATE_OPEN" {
				ibcInfo.OpenChannels++
			}
		}()
	}
	zap.L().Info("", zap.Bool("Success", true), zap.String("Open  IBC Channels", fmt.Sprint(ibcInfo.OpenChannels)))

	rd.IBC.IBCChannels = ibcChannels
	rd.IBC.IBCInfo.ibcChannelInfo = ibcInfo
}

func (rd *RESTData) getIBCConnections() {
	var ibcInfo ibcConnectionInfo
	var ibcConnections map[string][]string = make(map[string][]string)
	ibcInfo.OpenConnections = 0
	res, err := RESTQuery("/ibc/core/connection/v1beta1/connections" + "?pagination.limit=1000")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &ibcInfo)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Active IBC connections", fmt.Sprint(len(ibcInfo.IBConnections))))
	}

	for _, value := range ibcInfo.IBConnections {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// precommit failure validator
				}
			}()
			ibcConnections[value.ID] = []string{value.State, value.ClientID, value.Counterparty.ConnectionID, value.Counterparty.ClientID}
			if value.State == "STATE_OPEN" {
				ibcInfo.OpenConnections++
			}
		}()
	}
	zap.L().Info("", zap.Bool("Success", true), zap.String("Open IBC Channels", fmt.Sprint(ibcInfo.OpenConnections)))

	rd.IBC.IBCConnections = ibcConnections
	rd.IBC.IBCInfo.ibcConnectionInfo = ibcInfo
}
