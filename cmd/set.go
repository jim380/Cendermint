package cmd

import (
	common "github.com/jim380/Cendermint/rest"
)

// set rest and rpc addresses
func set_config() {
	common.Addr = restAddr
	common.OperAddr = operAddr
}
