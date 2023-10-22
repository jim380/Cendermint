package types

type OracleInfo struct {
	MissedCounterInfo
	PrevoteInfo
	VoteInfo
}

type MissedCounterInfo struct {
	MissedCounter struct {
		Validator string `json:"validator"`
		Counter   string `json:"counter"`
	} `json:"miss_counter"`
}

type PrevoteInfo struct {
	Prevote struct {
		Validator   string `json:"validator"`
		Hash        string `json:"hash"`
		SubmitBlock string `json:"submit_block"`
	} `json:"oracle_pre_vote"`
}

type VoteInfo struct {
	Vote struct {
		Validator   string `json:"validator"`
		ModuleVotes []struct {
			Module         string `json:"module"`
			NamespaceVotes []struct {
				Namespace string `json:"namespace"`
				Payload   string `json:"payload"`
			} `json:"namespace_votes"`
		} `json:"module_votes"`
	} `json:"oracle_vote"`
}
