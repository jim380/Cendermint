package types

type UpgradeInfo struct {
	Planned bool
	Plan    struct {
		Name   string `json:"name"`
		Time   string `json:"time"`
		Height string `json:"height"`
		Info   string `json:"info"`
	}
}
