package types

type LastBlock struct {
	Block struct {
		Header lastBlockHeader `json:"header"`
	}
}

type lastBlockHeader struct {
	Timestamp string `json:"time"`
}

type Blocks struct {
	BlockId struct {
		Hash string `json:"hash"`
	} `json:"block_id"`
	Block struct {
		Header     header     `json:"header"`
		LastCommit lastCommit `json:"last_commit"`
	} `json:"block"`
	MissingValidators []struct {
		Moniker     string
		ConsHexAddr string
	} // not part of the response so no json tag
}

type header struct {
	ChainID          string `json:"chain_id"`
	Height           string `json:"height"`
	Proposer_address string `json:"proposer_address"`
	Timestamp        string `json:"time"`
	LastTimestamp    string // not part of the response so no json tag
}

type lastCommit struct {
	Signatures []struct {
		Validator_address string `json:"validator_address"`
		Signature         string `json:"signature"`
	} `json:"signatures"`
}
