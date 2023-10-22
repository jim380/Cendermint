package types

type GovInfo struct {
	TotalProposalCount      float64
	VotingProposalCount     float64
	InVotingVotedCount      float64
	InVotingDidNotVoteCount float64
}

type Gov struct {
	Proposals  []proposal `json:"proposals"`
	Pagination struct {
		Total string `json:"total"`
	} `json:"pagination"`
}

type proposal struct {
	ProposalID string `json:"proposal_id"`
	Status     string `json:"status"`
}

type Vote struct {
	Votes struct {
		Option string `json:"option"`
	} `json:"vote"`
}
