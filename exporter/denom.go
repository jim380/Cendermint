package exporter

import (
	"go.uber.org/zap"
)

func getDenomList(chain string) []string {

	var dList []string

	// Add a staking denom to index 0
	switch chain {
	case "cosmos":
		dList = []string{"uatom"}
	case "iris":
		dList = []string{"uiris"}
	case "umee":
		dList = []string{"uumee"}
	case "osmosis":
		dList = []string{"uosmo"}
	case "juno":
		dList = []string{"ujuno"}
	case "akash":
		dList = []string{"uakt"}
	case "regen":
		dList = []string{"uregen"}
	case "microtick":
		dList = []string{"utick"}
	case "nyx":
		dList = []string{"unyx"}
	case "evmos":
		dList = []string{"aevmos"}
	case "assetMantle":
		dList = []string{"aphoton"}
	case "rizon":
		dList = []string{"uatolo"}
	case "stargaze":
		dList = []string{"ustars"}
	case "chihuahua":
		dList = []string{"uhuahua"}
	case "gravity":
		dList = []string{"ugraviton"}
	case "lum":
		dList = []string{"ulum"}
	case "provenance":
		dList = []string{"nhash"}
	case "crescent":
		dList = []string{"ucre"}
	case "sifchain":
		dList = []string{"urowan"}
	default:
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "denom not supported"))
	}

	return dList
}
