package rest

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type txInfo struct {
	Txs        txs    `json:"txs"`
	TxResp     txResp `json:"tx_responses"`
	Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}
	Result txResult
	TPS    float64
}

type txResult struct {
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

func (rd *RESTData) getTxInfo(currentBlockHeight int64) {
	var txInfo txInfo
	var txRes txResult

	res, err := HttpQuery(RESTAddr + "/cosmos/tx/v1beta1/txs?events=tx.height=" + fmt.Sprintf("%v", currentBlockHeight))
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &txInfo)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Total txs in this block:", txInfo.Pagination.Total))
	}

	for _, v := range txInfo.TxResp {
		// zap.L().Info("", zap.Bool("Success", true), zap.String("Tx #:", fmt.Sprintf("%v", i)))
		gasWantedRes, _ := strconv.ParseFloat(v.GasWanted, 64)
		txRes.GasWantedTotal = txRes.GasWantedTotal + gasWantedRes
		gasUsedRes, _ := strconv.ParseFloat(v.GasUsd, 64)
		txRes.GasUsedTotal = txRes.GasUsedTotal + gasUsedRes
		for _, v := range v.Logs {
			for _, v := range v.Events {
				txRes.Default.EventsTotal++
				switch v.Type {
				case "delegate":
					txRes.Default.DelegateTotal++
				case "message":
					txRes.Default.MessageTotal++
				case "transfer":
					txRes.Default.TransferTotal++
				case "unbond":
					txRes.Default.UnbondTotal++
				case "withdraw_rewards":
					txRes.Default.WithdrawRewardsTotal++
				case "create_validator":
					txRes.Default.CreateValidatorTotal++
				case "proposal_vote":
					txRes.Default.ProposalVote++
				case "fungible_token_packet":
					txRes.IBC.FungibleTokenPacketTotal++
				case "ibc_transfer":
					txRes.IBC.IbcTransferTotal++
				case "send_packet":
					txRes.IBC.SendPacketTotal++
				case "recv_packet":
					txRes.IBC.RecvPacketTotal++
				case "redelegate":
					txRes.Default.RedelegateTotal++
				case "update_client":
					txRes.IBC.UpdateClientTotal++
				case "acknowledge_packet":
					txRes.IBC.AckPacketTotal++
				case "write_acknowledgement":
					txRes.IBC.WriteAckTotal++
				case "timeout":
					txRes.IBC.TimeoutTotal++
				case "timeout_packet":
					txRes.IBC.TimeoutPacketTotal++
				case "denomination_trace":
					txRes.IBC.DenomTraceTotal++
				case "swap_within_batch":
					txRes.Swap.SwapWithinBatchTotal++
				case "withdraw_within_batch":
					txRes.Swap.WithdrawWithinBatchTotal++
				case "deposit_within_batch":
					txRes.Swap.DepositWithinBatchTotal++
				default:
					txRes.OthersTotal++
				}
				// if v.Type == "message" {
				// 	for _, v := range v.Attributes {
				// 		if v.Key == "action" {
				// 			actionsTotal++
				// 			switch v.Value {
				// 			case "send":
				// 				sendTotal++
				// 			case "delegate":
				// 				delegateActionTotal++
				// 			case "begin_unbonding":
				// 				beginUnbondingTotal++
				// 			case "withdraw_delegator_reward":
				// 				withdrawDelegatorRewardTotal++
				// 			}
				// 		}
				// 	}
				// }
			}
		}
	}
	zap.L().Info("", zap.Bool("Success", true), zap.String("Events Total:", fmt.Sprintf("%v", txRes.Default.EventsTotal)))
	if txRes.OthersTotal != 0 {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Others Total:", fmt.Sprintf("%v", txRes.OthersTotal)))
	}

	rd.TxInfo = txInfo
	rd.TxInfo.Result = txRes
}
