package cmd

import (
	common "github.com/jim380/Cosmos-IE/rest/common"
)

func set_config() {
	common.Addr = restAddr
	common.OperAddr = operAddr
}
