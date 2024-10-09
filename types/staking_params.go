package types

type StakingParams struct {
	Params struct {
		UnbondingTime     string `json:"unbonding_time"`
		MaxValidators     int    `json:"max_validators"`
		MaxEntries        int    `json:"max_entries"`
		HistoricalEntries int    `json:"historical_entries"`
		BondDenom         string `json:"bond_denom"`
		MinCommissionRate string `json:"min_commission_rate"`
	} `json:"params"`
}
