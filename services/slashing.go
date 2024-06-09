package services

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

func (sls *SlashingService) Init(db *sql.DB) {
	sls.DB = db
}

type MissingValidators []struct {
	Moniker     string
	ConsPubAddr string
}

func (ss *SlashingService) GetSlashingParams(cfg config.Config, rd *types.RESTData) {
	var d types.SlashingInfo

	route := rest.GetSlashingParamsRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &d)
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
	json.Unmarshal(res, &d)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.Slashing.ValSigning = d.ValSigning
}

func (ss *SlashingService) GetCommitInfo(cfg config.Config, rd *types.RESTData, blockData types.Blocks, consHexAddr string) MissingValidators {
	var cInfo types.CommitInfo
	missed := true

	blockProposer := blockData.Block.Header.ProposerAddress
	cInfo.ChainId = blockData.Block.Header.ChainID
	cInfo.ValidatorPrecommitStatus, cInfo.ValidatorProposingStatus, cInfo.MissThreshold, cInfo.MissConsecutive = 0.0, 0.0, 0.0, 0.0
	currentHeight, _ := strconv.Atoi(blockData.Block.Header.Height)

	/*
		Find validators with missing signatures in the block
	*/
	// var cs types.ConsensusState
	var activeSet map[string][]string = make(map[string][]string)
	var missingValidators MissingValidators

	conspubMonikerMap := rest.GetConspubMonikerMap()
	for consPubKey, moniker := range conspubMonikerMap {
		prefix, err := utils.GetPrefix(cfg.Chain.Name)
		if err != nil {
			zap.L().Fatal("Failed to get prefix", zap.Bool("Success", false), zap.String("err", err.Error()))
		}
		consAddrInHex := utils.PubkeyToHexAddr(prefix, consPubKey)
		if consAddrInHex == "" {
			zap.L().Fatal("Failed to convert public key to hex address", zap.Bool("Success", false))
		}
		// populate the map => [ConsAddr]{consPubKey, moniker}
		activeSet[consAddrInHex] = []string{consPubKey, moniker}
	}

	/*
		- Create a map validatorConsAddrSignedMap using allSignaturesInBlock for quick lookup
		- validatorConsAddrSignedMap gives all validators who signed on this block
	*/
	allSignaturesInBlock := blockData.Block.LastCommit.Signatures
	validatorConsAddrSignedMap := make(map[string]bool)
	for _, signature := range allSignaturesInBlock {
		if consHexAddr == blockProposer {
			cInfo.ValidatorProposingStatus = 1.0
			zap.L().Info("", zap.Bool("Success", true), zap.String("Proposer:", "true"))
		}

		// Validator_address could be in hex or base64; hex is legacy so using base64 here
		validatorConsAddrSignedMap[signature.Validator_address] = true
	}

	// Check if consAddrInHex in activeSet exists in validatorConsAddrSignedMap
	for consAddrInHex, props := range activeSet {
		// convert consAddrInHex to base64
		base64ConsAddr, err := utils.HexToBase64(consAddrInHex)
		if err != nil {
			zap.L().Fatal("Failed to convert hex to base64", zap.Bool("Success", false), zap.String("err:", err.Error()))
		}

		if _, exists := validatorConsAddrSignedMap[base64ConsAddr]; !exists {
			// If the Validator_address does not exist in allSignaturesInBlock, add it to MissingValidators
			missingValidators = append(missingValidators, struct {
				Moniker     string
				ConsPubAddr string
			}{
				Moniker:     props[1],
				ConsPubAddr: props[0],
				// TO-DO add operator address
			})
		}
	}

	base64Addr, err := utils.HexToBase64(consHexAddr)
	if err != nil {
		zap.L().Fatal("Failed to convert hex to base64", zap.Bool("Success", false), zap.String("err:", err.Error()))
	}

	if _, exists := validatorConsAddrSignedMap[base64Addr]; exists {
		// If exists, then the validator signed this block
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

	if missed {
		cInfo.MissedCount += 1
	}

	rd.Commit = cInfo

	return missingValidators
}
