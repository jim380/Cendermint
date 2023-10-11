package models

import (
	"database/sql"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type SlashingService struct {
	DB *sql.DB
}

func (ss *SlashingService) GetSlashingParams(cfg config.Config, rd *types.RESTData) {
	var d types.SlashingInfo

	route := rest.GetSlashingParamsRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &d); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}
	rd.Slashing.Params = d.Params
}

func (ss *SlashingService) GetSigningInfo(cfg config.Config, consAddr string, rd *types.RESTData) {
	var d types.SlashingInfo

	route := rest.GetSigningInfoByAddressRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route + consAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &d); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.Slashing.ValSigning = d.ValSigning
}

func (ss *SlashingService) GetCommitInfo(rd *types.RESTData, blockData types.Blocks, consHexAddr string) {
	var cInfo types.CommitInfo
	missed := true

	blockProposer := blockData.Block.Header.Proposer_address
	cInfo.ChainId = blockData.Block.Header.ChainID
	cInfo.ValidatorPrecommitStatus, cInfo.ValidatorProposingStatus, cInfo.MissThreshold, cInfo.MissConsecutive = 0.0, 0.0, 0.0, 0.0
	currentHeight, _ := strconv.Atoi(blockData.Block.Header.Height)

	for _, v := range blockData.Block.LastCommit.Signatures {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// precommit failure validator
				}
			}()

			if consHexAddr == blockProposer {
				cInfo.ValidatorProposingStatus = 1.0
				zap.L().Info("", zap.Bool("Success", true), zap.String("Proposer:", "true"))
			}

			if consHexAddr == v.Validator_address {
				missed = false
				cInfo.LastSigned = currentHeight
				cInfo.ValidatorPrecommitStatus = 1.0
				// if missed more than threshold
				threshold, _ := strconv.Atoi(os.Getenv("MISS_THRESHOLD"))
				if cInfo.MissedCount >= threshold {
					zap.L().Warn("Missed >= threshold", zap.Bool("Success", true), zap.String("MissedCount", strconv.Itoa(cInfo.MissedCount)))
					zap.L().Warn("Missed >= threshold", zap.Bool("Success", true), zap.String("Threshold", os.Getenv("MISS_THRESHOLD")))
					cInfo.MissThreshold = 1
				}
				// miss consecutively
				consecutive, _ := strconv.Atoi(os.Getenv("MISS_CONSECUTIVE"))
				if currentHeight-cInfo.LastSigned == consecutive {
					zap.L().Warn("MissConsecutive >= threshold", zap.Bool("Success", true), zap.String("MissedCount", strconv.Itoa(currentHeight-cInfo.LastSigned)))
					zap.L().Warn("MissConsecutive >= threshold", zap.Bool("Success", true), zap.String("Threshold", os.Getenv("MISS_CONSECUTIVE")))
					cInfo.MissConsecutive = 1
				}
				// MissedCount resets when the validator signs again
				cInfo.MissedCount = 0
			}
		}()
	}
	if missed {
		cInfo.MissedCount += 1
	}
	rd.Commit = cInfo
}
