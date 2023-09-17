package utils

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
)

func StringToFloat64(str string) float64 {
	var result float64
	result, _ = strconv.ParseFloat(str, 64)
	return result
}

func BoolToFloat64(b bool) float64 {
	var result float64
	if b {
		result = 1
	} else {
		result = 0
	}

	return result
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
