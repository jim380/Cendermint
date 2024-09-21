package types

type GravityInfo struct {
	GravityParams
	OracleEventNonce
	ValSetCount   int
	ValActive     float64 // [0]: false, [1]: true
	GravityActive float64 // [0]: false, [1]: true
	UmeePrice
	EthPrice
	BatchFees   float64
	BatchesFees float64
	BridgeFees  float64
}

type GravityParams struct {
	BridgeParams `json:"params"`
}

type BridgeParams struct {
	SignedValsetsWindow    string `json:"signed_valsets_window"`
	SignedBatchesWindow    string `json:"signed_batches_window"`
	TargetBatchTimeout     string `json:"target_batch_timeout"`
	SlashFractionValset    string `json:"slash_fraction_valset"`
	SlashFractionBatch     string `json:"slash_fraction_batch"`
	SlashFractionBadEthSig string `json:"slash_fraction_bad_eth_signature"`
	ValsetReward           `json:"valset_reward"`
	BridgeActive           bool `json:"bridge_active"`
}

type ValsetReward struct {
	Amount string `json:"amount"`
}

type EthPrice struct {
	ETHUSD `json:"ethereum"`
}

type ETHUSD struct {
	ETHPrice float64 `json:"usd"`
}

type UmeePrice struct {
	UMEEUSD `json:"umee"`
}

type UMEEUSD struct {
	UMEEPrice float64 `json:"usd"`
}

type Batches struct {
	Batches []batch `json:"batches"`
	Fees    float64
}

type batch struct {
	BatchNonce   string        `json:"batch_nonce"`
	BatchTimeout string        `json:"batch_timeout"`
	Transactions []transaction `json:"transactions"`
}

type transaction struct {
	ID         string     `json:"id"`
	Sender     string     `json:"sender"`
	DestAddr   string     `json:"dest_address"`
	ERC20Token erc20Token `json:"erc20_token"`
	ERC20Fee   erc20Fee   `json:"erc20_fee"`
}

type erc20Token struct {
	Contract string `json:"contract"`
	Amount   string `json:"amount"`
}

type erc20Fee struct {
	Contract string `json:"contract"`
	Amount   string `json:"amount"`
}

type BatchFees struct {
	BatchFees []batchFee `json:"batchFees"`
	Fees      float64
}

type batchFee struct {
	Token     string `json:"token"`
	TotalFees string `json:"total_fees"`
}

type OracleEventNonce struct {
	EventNonce string `json:"event_nonce"`
}

type ValSetInfo struct {
	ValSet   valSet `json:"valset"`
	ValCount int
}

type valSet struct {
	Nonce   string   `json:"nonce"`
	Members []member `json:"members"`
}

type member struct {
	Power   string `json:"power"`
	ETHAddr string `json:"ethereum_address"`
}
