package exporter

import (
	rest "github.com/jim380/Cosmos-IE/rest/common"
)

var (
	previousBlockHeight int64

	gaugesNamespaceList = [...]string{"blockHeight",
		"notBondedTokens",
		"bondedTokens",
		"totalSupply",
		"bondedRatio",
		"totalProposalCount",
		"votingProposalCount",
		"votingPower",
		"minSelfDelegation",
		"jailStatus",
		"delegationShares",
		"delegationRatio",
		"delegatorCount",
		"delegationSelf",
		"commissionRate",
		"commissionMaxRate",
		"commissionMaxChangeRate",
		"precommitStatus",
		"inflation",
		"actualInflation",
	}

	metricData metric
)

type metric struct {
	Network struct {
		ChainID       string
		BlockHeight   int64
		PrecommitRate float64

		Staking struct {
			NotBondedTokens float64
			BondedTokens    float64
			TotalSupply     float64
			BondedRatio     float64
		}

		Minting struct {
			Inflation       float64
			ActualInflation float64
		}

		Gov struct {
			TotalProposalCount  float64
			VotingProposalCount float64
		}
	}

	Validator struct {
		Moniker           string
		VotingPower       float64
		MinSelfDelegation float64
		JailStatus        float64

		Address struct {
			Account      string
			Operator     string
			ConsensusHex string
		}

		// Proposer struct {
		// 	Ranking float64
		// 	Status  float64
		// }

		Delegation struct {
			Shares         float64
			Ratio          float64
			DelegatorCount float64
			Self           float64
		}

		Commission struct {
			Rate          float64
			MaxRate       float64
			MaxChangeRate float64
		}

		Account struct {
			Balances   []rest.Coin
			Commission []rest.Coin
			Rewards    []rest.Coin
		}

		Commit struct {
			// VoteType        float64
			PrecommitStatus float64
		}
	}
}

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
	}

	return dList
}
