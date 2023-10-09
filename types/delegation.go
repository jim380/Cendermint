package types

type DelegationsInfo struct {
	DelegationRes delegationRes `json:"delegation_responses"`
	Pagination    struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}
}

type delegationRes []struct {
	Delegation struct {
		DelegatorAddr string `json:"delegator_address"`
		ValidatorAddr string `json:"validator_address"`
		Shares        string `json:"shares"`
	}
	Balance struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}
}
