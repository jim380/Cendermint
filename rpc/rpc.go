package rpc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	tmhttp "github.com/tendermint/tendermint/rpc/client/http"
)

type RPCData struct {
	Commit CommitInfo
}

var (
	Addr        string
	ConsHexAddr string
	Client      *tmhttp.HTTP
)

func newRPCData() *RPCData {
	rd := &RPCData{}
	return rd
}

func connect() {
	client, err := tmhttp.NewWithClient("tcp://"+Addr, "/ws", http.DefaultClient)
	if err != nil {
		zap.L().Fatal("RPC-Server", zap.Bool("Success", false), zap.String("err", fmt.Sprintf("%s", err)))
	} else {
		zap.L().Info("RPC-Server", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("Connected to", "tcp://"+Addr+"/websocket"))
	}
	defer client.Stop()
}

func GetData(blockHeight int64, consHexAddr string) *RPCData {
	var commitHeight int64 = blockHeight - 1
	cxtTimeout := 3 * time.Second
	rd := newRPCData()

	ctx, cancel := context.WithTimeout(context.Background(), cxtTimeout)
	defer cancel()

	commitData, err := Client.Commit(ctx, &commitHeight)
	if err != nil {
		zap.L().Fatal("RPC-Server", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	} else {
		zap.L().Info("RPC-Server", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("Get Data", "Commit Data"))
	}

	rd.getCommit(commitData, consHexAddr)

	return rd
}
