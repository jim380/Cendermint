package exporter

var labels []labelData = []labelData{
	{
		"labels_node_info",
		[]string{"chain_id", "node_moniker", "node_id", "tm_version", "app_name", "binary_name", "app_version", "git_commit", "go_version", "sdk_version"},
	},
	{
		"labels_addr",
		[]string{"operator_address", "account_address", "cons_address_hex"},
	},
	{
		"labels_upgrade",
		[]string{"upgrade_name", "upgrade_time", "upgrade_height", "upgrade_info"},
	},
}
