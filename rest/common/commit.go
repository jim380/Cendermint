package rest

import (
	utils "github.com/jim380/Cosmos-IE/utils"
	"go.uber.org/zap"
)

type commitInfo struct {
	ChainId                  string
	ValidatorPrecommitStatus float64 // [0]: false, [1]: true
}

func (rd *RESTData) getCommit(blockData Blocks) {
	var cInfo commitInfo
	consHexAddr := utils.Bech32AddrToHexAddr(rd.Validatorsets[rd.Validator.Consensus_pubkey.Key][0])
	cInfo.ChainId = blockData.Block.Header.ChainID
	cInfo.ValidatorPrecommitStatus = 0.0

	for _, v := range blockData.Block.Last_commit.Signatures {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// panic("oops..something bad happened")
				}
			}()

			if consHexAddr == v.Validator_address {
				cInfo.ValidatorPrecommitStatus = 1.0
				zap.L().Info("", zap.Bool("Success", true), zap.String("Precommit:", "signed"))
			} else {
				zap.L().Fatal("", zap.Bool("Success", false), zap.String("Precommit:", "missed"))
			}
		}()

	}

	rd.Commit = cInfo
}
