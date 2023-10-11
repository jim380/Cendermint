package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type TxnService struct {
	DB *sql.DB
}

func (ts *TxnService) GetInfo(cfg config.Config, currentBlockHeight int64, rd *types.RESTData) {
	var txInfo types.TxInfo
	var txRes types.TxResult

	route := rest.GetTxByHeightRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route + fmt.Sprintf("%v", currentBlockHeight))
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if !json.Valid(res) {
		zap.L().Error("Response is not valid JSON")
		return
	}
	if err := json.Unmarshal(res, &txInfo); err != nil {
		zap.L().Error("Failed to unmarshal JSON response", zap.Error(err))
		return
	}
	switch {
	case strings.Contains(string(res), "not found"):
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	case strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":"):
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	default:
		zap.L().Info("", zap.Bool("Success", true), zap.String("Total txs in this block", txInfo.Pagination.Total))
	}

	for _, v := range txInfo.TxResp {
		// zap.L().Info("", zap.Bool("Success", true), zap.String("Tx #:", fmt.Sprintf("%v", i)))
		gasWantedRes, _ := strconv.ParseFloat(v.GasWanted, 64)
		txRes.GasWantedTotal += gasWantedRes
		gasUsedRes, _ := strconv.ParseFloat(v.GasUsd, 64)
		txRes.GasUsedTotal += gasUsedRes
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
	zap.L().Info("", zap.Bool("Success", true), zap.String("Events Total", fmt.Sprintf("%v", txRes.Default.EventsTotal)))
	if txRes.OthersTotal != 0 {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Others Total", fmt.Sprintf("%v", txRes.OthersTotal)))
	}

	rd.TxInfo = txInfo
	rd.TxInfo.Result = txRes
}
