package rpc

import (
	"fmt"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"go.uber.org/zap"
)

type CommitInfo struct {
	ChainId                  string
	VoteType                 float64 // [0]: false, [1]: prevote, [2]: precommit
	ValidatorPrecommitStatus float64 // [0]: false, [1]: true
	ValidatorProposingStatus float64 // [0]: false, [1]: true
}

func (rd *RPCData) getCommit(commitData *coretypes.ResultCommit, consHexAddr string) {
	var cInfo CommitInfo
	blockProposer := fmt.Sprint(commitData.SignedHeader.Header.ProposerAddress)

	cInfo.ChainId = commitData.SignedHeader.Header.ChainID
	cInfo.VoteType, cInfo.ValidatorPrecommitStatus, cInfo.ValidatorProposingStatus = 0.0, 0.0, 0.0

	for _, v := range commitData.SignedHeader.Commit.Signatures {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// panic("oops..something bad happened")
				}
			}()

			if consHexAddr == fmt.Sprint(v.ValidatorAddress) {
				cInfo.ValidatorPrecommitStatus = 1.0
				zap.L().Info("", zap.Bool("Success", true), zap.String("Precommit:", "signed"))
			} else if consHexAddr == blockProposer {
				cInfo.ValidatorProposingStatus = 1.0
			} else {
				// fmt.Println("Hex:", consHexAddr, "Validator Address:", v.Validator_address)
				zap.L().Error("", zap.Bool("Success", false), zap.String("Precommit:", "missed"))
			}
		}()

	}
	rd.Commit = cInfo
}
