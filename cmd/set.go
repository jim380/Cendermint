package cmd

import (
	common "github.com/jim380/Cosmos-IE/rest/common"
	"github.com/jim380/Cosmos-IE/rpc"
)

// set rest and rpc addresses
func set_config() {
	common.Addr = restAddr
	common.OperAddr = operAddr
	rpc.Addr = rpcAddr
}
