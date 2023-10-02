package utils

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"go.uber.org/zap"

	"github.com/cosmos/cosmos-sdk/types/bech32"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetPrefix(chain string) string {
	switch chain {
	case "cosmos":
		return "cosmos"
	case "umee":
		return "umee"
	case "osmosis":
		return "osmo"
	case "juno":
		return "juno"
	case "akash":
		return "akash"
	case "regen":
		return "regen"
	case "stargaze":
		return "stars"
	case "evmos":
		return "evmos"
	case "rizon":
		return "rizon"
	case "gravity":
		return "gravity"
	case "lum":
		return "lum"
	case "provenance":
		return "pb"
	case "crescent":
		return "cre"
	case "assetMantle":
		return "mantle"
	case "sifchain":
		return "sif"
	case "passage":
		return "pas"
	case "stride":
		return "stride"
	case "canto":
		return "canto"
	case "teritori":
		return "tori"
	case "nym":
		return "n"
	default:
		return "cosmos"
	}
}

// Bech32 Addr -> Hex Addr
func Bech32AddrToHexAddr(bech32str string) string {
	_, bz, err := bech32.DecodeAndConvert(bech32str)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	}

	return fmt.Sprintf("%X", bz)
}

func GetAccAddrFromOperAddr(operAddr string) string {
	hexAddr, err := sdk.ValAddressFromBech32(operAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	}

	accAddr, err := sdk.AccAddressFromHex(fmt.Sprint(hexAddr))
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	}

	return accAddr.String()
}

func GetAccAddrFromOperAddr_localPrefixes(operAddr string, bech32Prefixes []string) string {
	bz, err := sdk.GetFromBech32(operAddr, bech32Prefixes[2])
	if err != nil {

		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	}

	accAddr, err := bech32.ConvertAndEncode(bech32Prefixes[0], bz)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	}

	return accAddr
}

func Base64ToHex(base64String string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HexToBase64(hexAddr string) string {
	bytes, err := decodeHex([]byte(hexAddr))
	if err != nil {
		fmt.Println(err)
	}
	return string(base64Encode(bytes))
}

func decodeHex(input []byte) ([]byte, error) {
	db := make([]byte, hex.DecodedLen(len(input)))
	_, err := hex.Decode(db, input)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func base64Encode(input []byte) []byte {
	eb := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
	base64.StdEncoding.Encode(eb, input)

	return eb
}
