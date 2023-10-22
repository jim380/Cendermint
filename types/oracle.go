package types

type OracleInfo struct {
	MissedCounterInfo
}

type MissedCounterInfo struct {
	MissedCounter struct {
		Validator string `json:"validator"`
		Counter   string `json:"counter"`
	} `json:"miss_counter"`
}
