package rest

type commitInfo struct {
	ChainId                  string
	ValidatorPrecommitStatus float64 // [0]: false, [1]: true
}

func getCommit(blockData Blocks, consHexAddr string) commitInfo {
	var cInfo commitInfo

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
			}
		}()

	}

	return cInfo
}
