package rest

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type ibcChannelInfo struct {
	Open        int
	IBCChannels ibcChannels `json:"channels"`
	Pagination  struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}
	Height struct {
		RevisionNumber string `json:"revision_number"`
		RevisionHeight string `json:"revision_height"`
	}
}

type ibcChannels []struct {
	State          string       `json:"state"`
	Ordering       string       `json:"ordering"`
	Counterparty   counterparty `json:"counterparty"`
	ConnectionHops []struct {
		string
	}
	Version   string `json:"version"`
	PortID    string `json:"port_id"`
	ChannelID string `json:"channel_id"`
}

type counterparty struct {
	PortID    string `json:"port_id"`
	ChannelID string `json:"channel_id"`
}

func (rd *RESTData) getIBCChannels() {
	var ibc ibcChannelInfo
	var ibcChannels map[string][]string = make(map[string][]string)
	ibc.Open = 0
	res, err := RESTQuery("/ibc/core/channel/v1beta1/channels" + "?pagination.limit=1000")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &ibc)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Active IBC Channels", fmt.Sprint(len(ibc.IBCChannels))))
	}

	for _, value := range ibc.IBCChannels {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// precommit failure validator
				}
			}()
			// populate the ibc channel map => [ChanelID][]string{State, Ordering, Counterparty_ChannelID}
			ibcChannels[value.ChannelID] = []string{value.State, value.Ordering, value.Counterparty.ChannelID}
			if value.State == "STATE_OPEN" {
				ibc.Open++
			}
		}()
	}
	zap.L().Info("", zap.Bool("Success", true), zap.String("Open IBC Channels", fmt.Sprint(ibc.Open)))

	rd.IBC.IBCChannels = ibcChannels
	rd.IBC.IBCInfo = ibc
}
