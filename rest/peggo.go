package rest

import (
	"encoding/json"
	"os"

	utils "github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type peggoInfo struct {
	oracleEvent
	ValSetCount int
	ValActive   float64 // [0]: false, [1]: true
	erc20Price
	ethPrice
	BatchFees  float64
	BridgeFees float64
}

type erc20Price struct {
	contractAddr `json:"0xe54fbaecc50731afe54924c40dfd1274f718fe02"`
}

type contractAddr struct {
	ERC20Price float64 `json:"usd"`
}

type ethPrice struct {
	ETHUSD `json:"ethereum"`
}

type ETHUSD struct {
	ETHPrice float64 `json:"usd"`
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

type oracleEvent struct {
	LastClaimEvent lastClaimEvent `json:"last_claim_event"`
}

type lastClaimEvent struct {
	EventNonce  string `json:"ethereum_event_nonce"`
	EventHeight string `json:"ethereum_event_height"`
}

type valSetInfo struct {
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

func (rd *RESTData) getUmeePrice() {
	var p erc20Price

	contractAddr := os.Getenv("CONTRACT_ADDR")
	res, err := PeggoQuery("https://peggo-fakex-qhcqt.ondigitalocean.app/api/v3/simple/token_price/ethereum?contract_addresses=" + contractAddr + "&vs_currencies=usd")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &p)

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

	rd.getUmeePrice()
	feesTotal := rd.PeggoInfo.ERC20Price * (b.Fees / 1000000)
	rd.PeggoInfo.BatchFees = feesTotal
}

func (rd *RESTData) getBridgeFees() {
	var p ethPrice
	var bf float64

	res, err := PeggoQuery("https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &p)

	rd.PeggoInfo.ETHPrice = p.ETHPrice
	rd.getUmeePrice()
	bf = (0.00225 * rd.PeggoInfo.ETHPrice) / (100 * rd.PeggoInfo.ERC20Price)
	rd.PeggoInfo.BridgeFees = bf
}

func (rd *RESTData) getOracleEvent() {
	var evt oracleEvent

	orchAddr := os.Getenv("ORCH_ADDR")
	res, err := RESTQuery("/peggy/v1/oracle/event/" + orchAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &evt)

	rd.PeggoInfo.oracleEvent = evt
}

func (rd *RESTData) getValSet() {
	var vs valSetInfo

	var vsResult map[string]string = make(map[string]string)

	res, err := RESTQuery("/peggy/v1/valset/current")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &vs)

	for _, member := range vs.ValSet.Members {
		vsResult[member.ETHAddr] = member.Power
	}

	rd.PeggoInfo.ValSetCount = len(vs.ValSet.Members)

	_, found := vsResult[os.Getenv("ORCH_ADDR")]
	if found {
		rd.PeggoInfo.ValActive = 1.0
	}
}
