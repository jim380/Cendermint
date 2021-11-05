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
	Result struct {
		GasWantedTotal       float64
		GasUsedTotal         float64
		EventsTotal          float64
		DelegateTotal        float64
		MessageTotal         float64
		TransferTotal        float64
		UnbondTotal          float64
		WithdrawRewardsTotal float64
		CreateValidatorTotal float64
	}
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
	var gasWantedTotal, gasUsedTotal, eventsTotal, delegateTotal, messageTotal, transferTotal, unbondTotal, withdrawRewardsTotal, createValidatorTotal int
	// var txRespResult map[string][]string = make(map[string][]string)

	res, err := RESTQuery("/cosmos/tx/v1beta1/txs?events=tx.height=" + fmt.Sprintf("%v", currentBlockHeight))
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

	for i, v := range txInfo.TxResp {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Tx #:", fmt.Sprintf("%v", i)))
		// txRespResult[v.Hash] = []string{v.GasWanted, v.GasUsd}
		gasWantedRes, _ := strconv.Atoi(v.GasWanted)
		gasWantedTotal = gasWantedTotal + gasWantedRes
		gasUsedRes, _ := strconv.Atoi(v.GasUsd)
		gasUsedTotal = gasUsedTotal + gasUsedRes
		for _, v := range v.Logs {
			for _, v := range v.Events {
				eventsTotal++
				switch v.Type {
				case "delegate":
					delegateTotal++
				case "message":
					messageTotal++
				case "transfer":
					transferTotal++
				case "unbond":
					unbondTotal++
				case "withdraw_rewards":
					withdrawRewardsTotal++
				case "create_validator":
					createValidatorTotal++
				}
				// for _, v := range v.Attributes {
				// 	attrKey = v.Key
				// 	attrValue = v.Value
				// }
			}
		}
	}
	zap.L().Info("", zap.Bool("Success", true), zap.String("GasUsed Total:", fmt.Sprintf("%v", gasWantedTotal)))
	zap.L().Info("", zap.Bool("Success", true), zap.String("Events Total:", fmt.Sprintf("%v", eventsTotal)))
	zap.L().Info("", zap.Bool("Success", true), zap.String("Transfer Total:", fmt.Sprintf("%v", transferTotal)))

	rd.TxInfo = txInfo
	rd.TxInfo.Result.GasUsedTotal = float64(gasWantedTotal)
	rd.TxInfo.Result.GasWantedTotal = float64(gasWantedTotal)
	rd.TxInfo.Result.EventsTotal = float64(eventsTotal)
	rd.TxInfo.Result.DelegateTotal = float64(delegateTotal)
	rd.TxInfo.Result.MessageTotal = float64(messageTotal)
	rd.TxInfo.Result.TransferTotal = float64(transferTotal)
	rd.TxInfo.Result.UnbondTotal = float64(unbondTotal)
	rd.TxInfo.Result.WithdrawRewardsTotal = float64(withdrawRewardsTotal)
	rd.TxInfo.Result.CreateValidatorTotal = float64(createValidatorTotal)
}
