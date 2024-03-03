package types

import "github.com/jim380/Cendermint/types/akash"

type AsyncData struct {
	AkashInfo akash.AkashInfo
}

func (ad AsyncData) New() *AsyncData {
	return &AsyncData{}
}
