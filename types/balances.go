package types

type Balances struct {
	Balances []Coin
}

type Coin struct {
	Denom  string
	Amount string
}
