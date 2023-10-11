package types

type SlashingInfo struct {
	Params     Params      `json:"params"`
	ValSigning SigningInfo `json:"val_signing_info"`
}

type Params struct {
	SignedBlocksWindow      string `json:"signed_blocks_window"`
	MinSignedPerWindow      string `json:"min_signed_per_window"`
	DowntimeJailDuration    string `json:"downtime_jail_duration"`
	SlashFractionDoubleSign string `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   string `json:"slash_fraction_downtime"`
}

type SigningInfo struct {
	StartHeight         string `json:"start_height"`
	IndexOffset         string `json:"index_offset"`
	JailedUntil         string `json:"jailed_until"`
	Tombstoned          bool   `json:"tombstoned"`
	MissedBlocksCounter string `json:"missed_blocks_counter"`
}
