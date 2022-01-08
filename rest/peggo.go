package rest

import (
	"encoding/json"
	"os"
	"strings"

	utils "github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type peggoInfo struct {
	erc20Price
	BatchFees float64
}

type erc20Price struct {
	contractAddr `json:"0xe54fbaecc50731afe54924c40dfd1274f718fe02"`
}

type contractAddr struct {
	ERC20Price float64 `json:"usd"`
}

type batches struct {
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

func (rd *RESTData) getPrice() {
	var p erc20Price

	contractAddr := os.Getenv("CONTRACT_ADDR")
	res, err := PeggoQuery("https://peggo-fakex-qhcqt.ondigitalocean.app/api/v3/simple/token_price/ethereum?contract_addresses=" + contractAddr + "&vs_currencies=usd")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &p)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		// zap.L().Info("\t", zap.Bool("Success", true), zap.String("ETH Price:", fmt.Sprintf("%f", p.contractAddr.ERC20Price)))
	}

	rd.PeggoInfo.ERC20Price = p.ERC20Price
}

func (rd *RESTData) getBatchFees() {
	var b batches

	res, err := RESTQuery("/peggy/v1/batch/outgoingtx")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &b)

	for _, batch := range b.Batches {
		for _, tx := range batch.Transactions {
			b.Fees += utils.StringToFloat64(tx.ERC20Fee.Amount)
		}
	}

	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		//zap.L().Info("\t", zap.Bool("Success", true), zap.String("Total batch fees:", fmt.Sprintf("%f", b.Batches)))
	}
	rd.getPrice()
	feesTotal := rd.PeggoInfo.ERC20Price * (b.Fees / 1000000)
	rd.PeggoInfo.BatchFees = feesTotal
}
