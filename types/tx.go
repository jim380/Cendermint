package types

type TxInfo struct {
	Txs        txs    `json:"txs"`
	TxResp     txResp `json:"tx_responses"`
	Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}
	Result TxResult
	TPS    float64
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

type txs []struct {
	AuthInfo struct {
		Fee []struct {
			Amount []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			}
			GasLimit string `json:"gas_limit"`
		}
	}
}

type txResp []struct {
	Hash string `json:"txhash"`
	Logs []struct {
		Events []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}
		}
	}
	GasWanted string `json:"gas_wanted"`
	GasUsd    string `json:"gas_used"`
}
