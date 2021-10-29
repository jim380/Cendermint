package cmd

import (
	"fmt"
	"log"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/jim380/Cendermint/utils"
)

func SetSDKConfig(chain string) {

	// Bech32MainPrefix is the common prefix of all prefixes
	Bech32MainPrefix := utils.GetPrefix(chain)
	// Bech32PrefixAccAddr is the prefix of account addresses
	Bech32PrefixAccAddr := Bech32MainPrefix
	// Bech32PrefixAccPub is the prefix of account public keys
	Bech32PrefixAccPub := Bech32MainPrefix + sdktypes.PrefixPublic
	// Bech32PrefixValAddr is the prefix of validator operator addresses
	Bech32PrefixValAddr := Bech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixOperator
	// Bech32PrefixValPub is the prefix of validator operator public keys
	Bech32PrefixValPub := Bech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixOperator + sdktypes.PrefixPublic
	// Bech32PrefixConsAddr is the prefix of consensus node addresses
	Bech32PrefixConsAddr := Bech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixConsensus
	// Bech32PrefixConsPub is the prefix of consensus node public keys
	Bech32PrefixConsPub := Bech32MainPrefix + sdktypes.PrefixValidator + sdktypes.PrefixConsensus + sdktypes.PrefixPublic
	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.Seal()
}

func CheckInputs(chain, operAddr, restAddr, listenPort string, chainList []string) {
	if chain == "" {
		log.Fatal("Chain was not provided.")
	} else {
		valid := false
		for _, c := range chainList {
			if chain == c {
				valid = true
			}
		}
		if !valid {
			log.Fatal(fmt.Sprintf("%s is not supported", chain) + fmt.Sprint("\nList of supported chains: ", chainList))
		}
	}

	// TODO add more robust checks
	if operAddr == "" {
		log.Fatal("Operator address was not provided.")
	}

	if restAddr == "" {
		log.Fatal("REST address was not provided.")
	}

	if listenPort == "" {
		log.Fatal("Listening port was not provided.")
	}
}
