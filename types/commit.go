package types

type CommitInfo struct {
	ChainId                  string
	ValidatorPrecommitStatus float64 // [0]: false, [1]: true
	ValidatorProposingStatus float64 // [0]: false, [1]: true
	MissedCount              int
	LastSigned               int
	MissThreshold            float64
	MissConsecutive          float64
}
