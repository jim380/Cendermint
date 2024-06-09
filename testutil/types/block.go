package types

import "time"

type TestDataBlock struct {
	Height    int
	Hash      string
	Timestamp time.Time
	Proposer  string
	TxnCount  int
}
