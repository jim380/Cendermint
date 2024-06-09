package types

type Validatorsets struct {
	Block_Height string `json:"block_height"`
	Validators   []struct {
		ConsAddr         string           `json:"address"`
		ConsPubKey       consPubKeyValSet `json:"pub_key"`
		ProposerPriority string           `json:"proposer_priority"`
		VotingPower      string           `json:"voting_power"`
	} `json:"validators"`
}

type consPubKeyValSet struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}
