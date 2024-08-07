package types

type GovInfo struct {
	TotalProposalCount      float64
	VotingProposalCount     float64
	InVotingVotedCount      float64
	InVotingDidNotVoteCount float64
}

type Gov struct {
	Proposals  []Proposal `json:"proposals"`
	Pagination struct {
		Total string `json:"total"`
	} `json:"pagination"`
}

type Proposal struct {
	ProposalID string `json:"id"`
	Status     string `json:"status"`
}

type Vote struct {
	Vote struct {
		ProposalID string `json:"proposal_id"`
		Voter      string `json:"voter"`
		Options    []struct {
			Option string `json:"option"`
			Weight string `json:"weight"`
		} `json:"options"`
		Metadata string `json:"metadata"`
	} `json:"vote"`
}
