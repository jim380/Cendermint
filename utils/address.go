package utils

import (
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
	case "nym":
		return "punk"
	default:
		return "cosmos"
	}
}

// Bech32 Addr -> Hex Addr
func Bech32AddrToHexAddr(bech32str string) string {
	_, bz, err := bech32.DecodeAndConvert(bech32str)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("Convert Address", "Bech32Addr To HexAddr"))
	}

	return fmt.Sprintf("%X", bz)
}

func GetAccAddrFromOperAddr(operAddr string) string {
	hexAddr, err := sdk.ValAddressFromBech32(operAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("Convert Address", "OperAddr To HexAddr"))
	}

	accAddr, err := sdk.AccAddressFromHex(fmt.Sprint(hexAddr))
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("Convert Address", "HexAddr To AccAddr"))
	}

	return accAddr.String()
}

func GetAccAddrFromOperAddr_localPrefixes(operAddr string, bech32Prefixes []string) string {
	bz, err := sdk.GetFromBech32(operAddr, bech32Prefixes[2])
	if err != nil {

		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("Convert Address", "OperAddr To HexAddr"))
	}

	accAddr, err := bech32.ConvertAndEncode(bech32Prefixes[0], bz)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", fmt.Sprint(err)))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("Convert Address", "HexAddr To AccAddr"))
	}

	return accAddr
}
