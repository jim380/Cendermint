package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type IbcService struct {
	DB *sql.DB
}

func (ibcs *IbcService) Init(db *sql.DB) {
	ibcs.DB = db
}

func (is *IbcService) GetChannelInfo(cfg config.Config, rd *types.RESTData) {
	var ibcInfo types.IbcChannelInfo
	var ibcChannels map[string][]string = make(map[string][]string)

	ibcInfo.OpenChannels = 0
	route := rest.GetIBCChannelsRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route + "?pagination.limit=1000000")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &ibcInfo)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Active IBC channels", fmt.Sprint(len(ibcInfo.IBCChannels))))
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
	zap.L().Info("", zap.Bool("Success", true), zap.String("Open IBC channels", fmt.Sprint(ibcInfo.OpenChannels)))

	rd.IBC.IBCChannels = ibcChannels
	rd.IBC.IBCInfo.IbcChannelInfo = ibcInfo
}

func (is *IbcService) GetConnectionInfo(cfg config.Config, rd *types.RESTData) {
	var ibcInfo types.IbcConnectionInfo
	var ibcConnections map[string][]string = make(map[string][]string)

	ibcInfo.OpenConnections = 0
	route := rest.GetIBCConnectionsRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route + "?pagination.limit=100000")
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
	zap.L().Info("", zap.Bool("Success", true), zap.String("Open IBC connections", fmt.Sprint(ibcInfo.OpenConnections)))

	rd.IBC.IBCConnections = ibcConnections
	rd.IBC.IBCInfo.IbcConnectionInfo = ibcInfo
}
