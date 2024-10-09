package types

type RPCData struct {
	ConsensusState
	Validatorsets map[string][]string
}

type ConsensusState struct {
	Result struct {
		RoundState `json:"round_state"`
	} `json:"result"`
}

type RoundState struct {
	Height       string           `json:"height"`
	Round        int64            `json:"round"`
	Step         int64            `json:"step"`
	Validatorset rpcValidatorsets `json:"validators"`
	Votes        []struct {
		Prevotes           []string `json:"prevotes"`
		Precommits         []string `json:"precommits"`
		PrevotesBitArray   string   `json:"prevotes_bit_array"`
		PrecommitsBitArray string   `json:"precommits_bit_array"`
	} `json:"votes"`
}

type rpcValidatorsets struct {
	Validators []struct {
		ConsAddrHex string `json:"address"`
		ConsPubKey  struct {
			Type string `json:"type"`
			Key  string `json:"value"`
		} `json:"pub_key"`
		ProposerPriority string `json:"proposer_priority"`
		VotingPower      string `json:"voting_power"`
		Moniker          string
	} `json:"validators"`
}

type RpcValidators struct {
	Validators []struct {
		ConsPubKey struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"consensus_pubkey"`
		Tokens      string `json:"tokens"`
		Description struct {
			Moniker string `json:"moniker"`
		} `json:"description"`
	} `json:"validators"`
}

func (rpc RPCData) New() *RPCData {
	return &RPCData{Validatorsets: make(map[string][]string)}
}
