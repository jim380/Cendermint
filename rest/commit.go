package rest

import (
	"strconv"

	"go.uber.org/zap"
)

type commitInfo struct {
	ChainId                  string
	ValidatorPrecommitStatus float64 // [0]: false, [1]: true
	ValidatorProposingStatus float64 // [0]: false, [1]: true
	MissedCount              int
	LastSigned               int
	MissThreshold            float64 // [0]: false, [1]: true
	MissConsecutive          float64 // [0]: false, [1]: true
}

func (rd *RESTData) getCommit(blockData Blocks, consHexAddr string) {
	var cInfo commitInfo
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
				if cInfo.MissedCount >= 2 {
					cInfo.MissThreshold = 1
				}
				// miss consecutively
				if currentHeight-cInfo.LastSigned == 2 {
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
