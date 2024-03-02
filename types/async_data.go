package types

type AsyncData struct {
	AkashInfo AkashInfo
}

func (ad AsyncData) New() *AsyncData {
	return &AsyncData{}
}
