package services

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

func (os *TxnService) Init(db *sql.DB) {
	os.DB = db
}

func (ts *TxnService) PopulateRestData(rd *types.RESTData, txInfo types.TxInfo) {
	var txRes types.TxResult

	for _, v := range txInfo.TxResp {
		gasWantedRes, _ := strconv.ParseFloat(v.GasWanted, 64)
		txRes.GasWantedTotal += gasWantedRes
		gasUsedRes, _ := strconv.ParseFloat(v.GasUsed, 64)
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

func (ts *TxnService) GetTxnsInBlock(cfg config.Config, height int64) (types.TxInfo, error) {
	var txInfo types.TxInfo

	route := rest.GetTxByHeightRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route + fmt.Sprintf("%v", height))
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if err := json.Unmarshal(res, &txInfo); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("", zap.Bool("Success", true), zap.String("Total txs in this block", txInfo.Total))
	}

	return txInfo, nil
}

func (ts *TxnService) Index(cfg config.Config, height int64, txsInBlock types.TxInfo) error {
	// start a new transaction
	tx, err := ts.DB.Begin()
	if err != nil {
		return err
	}

	// create a slice of TransactionData
	var txData []types.TransactionData
	for i := range txsInBlock.Txs {
		txData = append(txData, types.TransactionData{TxsData: txsInBlock.Txs[i], TxRespData: txsInBlock.TxResp[i]})
	}

	for _, txnData := range txData {
		// Insert into transactions table
		_, err = tx.Exec(`
	INSERT INTO transactions (hash, height, timestamp, type, gas_wanted, gas_used, memo, payer, granter)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (hash) DO NOTHING`,
			txnData.TxRespData.Hash, height, txnData.TxRespData.Timestamp, txnData.TxRespData.Tx.Type, txnData.TxRespData.GasWanted, txnData.TxRespData.GasUsed, txnData.TxsData.Body.Memo, txnData.TxsData.AuthInfo.Fee.Payer, txnData.TxsData.AuthInfo.Fee.Granter)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				zap.L().Error("failed to rollback transaction", zap.Error(rbErr))
			}
			return err
		}

		// insert into transaction_messages table
		for _, msg := range txnData.TxsData.Body.Messages {
			_, err = tx.Exec(`
			INSERT INTO transaction_messages (type, transaction_hash)
			VALUES ($1, $2)`,
				msg.Type, txnData.TxRespData.Hash)
			if err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					zap.L().Error("failed to rollback transaction", zap.Error(rbErr))
				}
				return err
			}
		}

		// insert into the denoms table & transaction_fee_amounts table
		for _, v := range txnData.TxsData.AuthInfo.Fee.Amount {
			var denomId int
			err = tx.QueryRow(`
				INSERT INTO denoms (denom)
				VALUES ($1) ON CONFLICT (denom) DO UPDATE SET denom = $1 RETURNING id`,
				v.Denom).Scan(&denomId)
			if err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					zap.L().Error("failed to rollback transaction", zap.Error(rbErr))
				}
				return err
			}

			_, err = tx.Exec(`
				INSERT INTO transaction_fee_amounts (amount, transaction_hash, denom_id)
				VALUES ($1, $2, $3)`,
				v.Amount, txnData.TxRespData.Hash, denomId)
			if err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					zap.L().Error("failed to rollback transaction", zap.Error(rbErr))
				}
				return err
			}
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			zap.L().Error("failed to rollback transaction", zap.Error(rbErr))
		}
		return err
	}

	return nil
}
