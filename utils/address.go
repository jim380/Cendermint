package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/bech32"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetPrefix(chain string) (string, error) {
	chains := []struct {
		Chain      string `json:"chain"`
		AddrPrefix string `json:"addr_prefix"`
	}{}

	jsonFile, err := os.Open("chains.json")
	if err != nil {
		return "", fmt.Errorf("failed to open chains.json: %w", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return "", fmt.Errorf("failed to read chains.json: %w", err)
	}

	err = json.Unmarshal(byteValue, &chains)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal chains.json: %w", err)
	}

	for _, c := range chains {
		if c.Chain == chain {
			return c.AddrPrefix, nil
		}
	}

	return "", fmt.Errorf("prefix not found for chain: %s", chain)
}

// valcons -> hex
func Bech32AddrToHexAddr(bech32str string) (string, error) {
	_, bz, err := bech32.DecodeAndConvert(bech32str)
	if err != nil {
		return "", fmt.Errorf("failed to decode and convert bech32 address: %w", err)
	}

	return fmt.Sprintf("%X", bz), nil
}

func GetAccAddrFromOperAddr(operAddr string) (string, error) {
	hexAddr, err := sdk.ValAddressFromBech32(operAddr)
	if err != nil {
		return "", fmt.Errorf("failed to convert operator address to hex: %w", err)
	}

	accAddr, err := sdk.AccAddressFromHexUnsafe(fmt.Sprint(hexAddr))
	if err != nil {
		return "", fmt.Errorf("failed to convert hex address to account address: %w", err)
	}

	return accAddr.String(), nil
}

// bech32Prefixes format: [osmo, osmovaloper]
func GetAccAddrFromOperAddrWithLocalPrefix(operAddr string, bech32Prefixes []string) (string, error) {
	bz, err := sdk.GetFromBech32(operAddr, bech32Prefixes[1])
	if err != nil {
		return "", fmt.Errorf("failed to get from bech32: %w", err)
	}

	accAddr, err := bech32.ConvertAndEncode(bech32Prefixes[0], bz)
	if err != nil {
		return "", fmt.Errorf("failed to convert and encode bech32: %w", err)
	}

	return accAddr, nil
}

func Base64ToHex(base64String string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

func HexToBase64(hexAddr string) (string, error) {
	bytes, err := decodeHex([]byte(hexAddr))
	if err != nil {
		return "", fmt.Errorf("failed to decode hex: %w", err)
	}
	return string(base64Encode(bytes)), nil
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

func PubkeyToHexAddr(prefix, pubkey string) string {
	// decode the base64 pubkey
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubkey)
	if err != nil {
		log.Println("Error decoding base64 pubKey:", err)
		return ""
	}

	// hash using SHA-256
	hash := sha256.Sum256(pubKeyBytes)

	// keep only the first 20 bytes
	addressBytes := hash[:20]

	// valcons -> hex
	bech32Addr, err := bech32.ConvertAndEncode(prefix+"valcons", addressBytes)
	if err != nil {
		log.Println("Error encoding to Bech32:", err)
		return ""
	}

	_, bz, err := bech32.DecodeAndConvert(bech32Addr)
	if err != nil {
		log.Println("PubkeyToHexAddr failed:", err)
		return ""
	}

	return fmt.Sprintf("%X", bz)
}
