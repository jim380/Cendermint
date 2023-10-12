package types

type Validators struct {
	Validator Validator `json:"validator"`
}

type Validator struct {
	OperAddr        string        `json:"operator_address"`
	ConsPubKey      consPubKeyVal `json:"consensus_pubkey"`
	Jailed          bool          `json:"jailed"`
	Status          string        `json:"status"`
	Tokens          string        `json:"tokens"`
	DelegatorShares string        `json:"delegator_shares"`
	Description     struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	}
	UnbondingHeight string `json:"unbonding_height"`
	UnbondingTime   string `json:"unbonding_time"`
	Commission      struct {
		Commission commissionRates `json:"commission_rates"`
		UpdateTime string          `json:"update_time"`
	}
	MinSelfDelegation string `json:"min_self_delegation"`
}

type consPubKeyVal struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type commissionRates struct {
	Rate          string `json:"rate"`
	MaxRate       string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
}
