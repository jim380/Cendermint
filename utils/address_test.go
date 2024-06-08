package utils_test

import (
	"testing"

	"github.com/jim380/Cendermint/utils"
)

func TestGetPrefix(t *testing.T) {
	tests := []struct {
		chain       string
		want        string
		errExpected bool
	}{
		{"cosmos", "cosmos", false},
		{"stargaze", "stars", false},
		{"test", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.chain, func(t *testing.T) {
			got, err := utils.GetPrefix(tt.chain)
			if (err != nil) != tt.errExpected {
				t.Errorf("GetPrefix() error = %v, wantErr %v", err, tt.errExpected)
				return
			}
			if got != tt.want {
				t.Errorf("GetPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBech32AddrToHexAddr(t *testing.T) {
	tests := []struct {
		bech32Addr string
		want       string
		wantErr    bool
	}{
		{"cosmosvalcons1px0zkz2cxvc6lh34uhafveea9jnaagckmrlsye", "099E2B09583331AFDE35E5FA96673D2CA7DEA316", false},
		{"osmovalcons1eddx8wg73a8w3kunt9pvhcjhy33kg70qqrwjct", "CB5A63B91E8F4EE8DB935942CBE25724636479E0", false},
		{"osmovalconsdde1x8wg73a8w3kunt9pvhcjhy33kg70qqrwjct", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.bech32Addr, func(t *testing.T) {
			got, err := utils.Bech32AddrToHexAddr(tt.bech32Addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bech32AddrToHexAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Bech32AddrToHexAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccAddrFromOperAddr(t *testing.T) {
	tests := []struct {
		operAddr string
		want     string
		wantErr  bool
	}{
		{"cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", "cosmos1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9q", false},
		{"cosmosvaloperlc1pqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.operAddr, func(t *testing.T) {
			got, err := utils.GetAccAddrFromOperAddr(tt.operAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccAddrFromOperAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAccAddrFromOperAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccAddrFromOperAddrWithLocalPrefix(t *testing.T) {
	bech32Prefixes := []string{"osmo", "osmovaloper"}

	tests := []struct {
		operAddr       string
		bech32Prefixes []string
		want           string
		wantErr        bool
	}{
		{"osmovaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4ep88n0y4", bech32Prefixes, "osmo1clpqr4nrk4khgkxj78fcwwh6dl3uw4epasmvnj", false},
		{"invalidOperAddr", bech32Prefixes, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.operAddr, func(t *testing.T) {
			got, err := utils.GetAccAddrFromOperAddrWithLocalPrefix(tt.operAddr, tt.bech32Prefixes)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccAddrFromOperAddrWithLocalPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAccAddrFromOperAddrWithLocalPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase64ToHex(t *testing.T) {
	tests := []struct {
		base64String string
		want         string
		wantErr      bool
	}{
		{"14upFOwoe0ORClUvpd6+Puvjufg=", "D78BA914EC287B43910A552FA5DEBE3EEBE3B9F8", false},
		{"CujluSpeFo0zFV83keTQhv2frWYkkdVOmDDw+S120sk=", "0AE8E5B92A5E168D33155F3791E4D086FD9FAD662491D54E9830F0F92D76D2C9", false},
		{"invalidBase64", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.base64String, func(t *testing.T) {
			got, err := utils.Base64ToHex(tt.base64String)
			if (err != nil) != tt.wantErr {
				t.Errorf("Base64ToHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Base64ToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHexToBase64(t *testing.T) {
	tests := []struct {
		hexAddr string
		want    string
		wantErr bool
	}{
		{"0AE8E5B92A5E168D33155F3791E4D086FD9FAD662491D54E9830F0F92D76D2C9", "CujluSpeFo0zFV83keTQhv2frWYkkdVOmDDw+S120sk=", false},
		{"D78BA914EC287B43910A552FA5DEBE3EEBE3B9F8", "14upFOwoe0ORClUvpd6+Puvjufg=", false},
		{"invalidHex", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.hexAddr, func(t *testing.T) {
			got, err := utils.HexToBase64(tt.hexAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexToBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HexToBase64() = %v, want %v", got, tt.want)
			}
		})
	}
}
