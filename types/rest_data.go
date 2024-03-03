package types

type RESTData struct {
	BlockHeight   int64
	BlockInterval int64
	Commit        CommitInfo
	NodeInfo      NodeInfo
	TxInfo        TxInfo
	StakingPool   StakingPool
	Slashing      SlashingInfo
	Validatorsets map[string][]string
	Validator     Validator
	Delegations   DelegationsInfo
	Balances      []Coin
	Rewards       []Coin
	Commission    []Coin
	Inflation     float64
	Gov           GovInfo
	IBC           struct {
		IBCChannels    map[string][]string
		IBCConnections map[string][]string
		IBCInfo        IbcInfo
	}
	UpgradeInfo UpgradeInfo
	GravityInfo GravityInfo
	OracleInfo  OracleInfo
}

func (rd RESTData) New(blockHeight int64) *RESTData {
	return &RESTData{
		BlockHeight:   blockHeight,
		Validatorsets: make(map[string][]string),
	}
}
