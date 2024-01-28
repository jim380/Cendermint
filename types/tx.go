package types

type TxInfo struct {
	Txs        []Txs    `json:"txs"`
	TxResp     []TxResp `json:"tx_responses"`
	Pagination struct {
		NextKey string `json:"next_key"`
	} `json:"pagination"`
	Total  string `json:"total"`
	Result TxResult
	TPS    float64
}

type Txs struct {
	Body struct {
		Messages []struct {
			Type string `json:"@type"`
		} `json:"messages"`
		Memo string `json:"memo"`
	} `json:"body"`
	AuthInfo struct {
		Fee struct {
			Amount []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"amount"`
			GasLimit string `json:"gas_limit"`
			Payer    string `json:"payer"`
			Granter  string `json:"granter"`
		} `json:"fee"`
	} `json:"auth_info"`
}

type TxResp struct {
	Hash      string `json:"txhash"`
	Timestamp string `json:"timestamp"`
	Tx        struct {
		Type string `json:"@type"`
	} `json:"tx"`
	Logs []struct {
		Events []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"attributes"`
		} `json:"events"`
	} `json:"logs"`
	GasWanted string `json:"gas_wanted"`
	GasUsed   string `json:"gas_used"`
}

type TxResult struct {
	GasWantedTotal float64
	GasUsedTotal   float64
	Default        struct {
		EventsTotal          float64
		DelegateTotal        float64
		MessageTotal         float64
		TransferTotal        float64
		UnbondTotal          float64
		RedelegateTotal      float64
		WithdrawRewardsTotal float64
		CreateValidatorTotal float64
		ProposalVote         float64
	}
	IBC struct {
		FungibleTokenPacketTotal float64
		IbcTransferTotal         float64
		UpdateClientTotal        float64
		AckPacketTotal           float64
		WriteAckTotal            float64
		SendPacketTotal          float64
		RecvPacketTotal          float64
		TimeoutTotal             float64
		TimeoutPacketTotal       float64
		DenomTraceTotal          float64
	}
	Swap struct {
		SwapWithinBatchTotal     float64
		WithdrawWithinBatchTotal float64
		DepositWithinBatchTotal  float64
	}
	OthersTotal float64
	// ActionsTotal                 float64
	// SendTotal                    float64
	// DelegateActionTotal          float64
	// BeginUnbondingTotal          float64
	// WithdrawDelegatorRewardTotal float64
}
