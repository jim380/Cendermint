package exporter

import "github.com/jim380/Cendermint/rest"

var (
	previousBlockHeight int64

	gaugesNamespaceList = [...]string{

		// chain
		"chain_blockHeight",

		// minting
		"minting_inflation",
		"minting_actual_inflation",

		// staking
		"staking_not_bonded_tokens",
		"staking_bonded_tokens",
		"staking_total_supply",
		"staking_bonded_ratio",

		// slashing
		"slashing_signed_blocks_window",
		"slashing_min_signed_per_window",
		"slashing_downtime_jail_duration",
		"slashing_slash_fraction_double_sign",
		"slashing_slash_fraction_downtime",

		// gov
		"gov_total_proposal_count",
		"gov_voting_proposal_count",

		// validator
		"validator_voting_power",
		"validator_min_self_delegation",
		"validator_jail_status",
		// vadalidator_delegation
		"validator_delegation_shares",
		"validator_delegation_ratio",
		// vadalidator_commission
		"validator_commission_rate",
		"validator_commission_max_rate",
		"validator_commission_max_change_rate",
		// vadalidator_signing
		"validator_precommit_status",
		"validator_proposer_status",

		// ibc
		"ibc_channels_total",
		"ibc_channels_open",
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

		Slashing struct {
			SignedBlocksWindow      float64
			MinSignedPerWindow      float64
			DowntimeJailDuration    float64
			SlashFractionDoubleSign float64
			SlashFractionDowntime   float64
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

		Proposer struct {
			Ranking float64
			Status  float64
		}

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
	IBC struct {
		IBCChannels struct {
			Total float64
			Open  float64
		}
	}
}

func getDenomList(chain string) []string {

	var dList []string

	// Add a staking denom to index 0
	switch chain {
	case "cosmos":
		dList = []string{"uatom"}
	case "umee":
		dList = []string{"uumee"}
	case "nym":
		dList = []string{"upunk"}
	}
	return dList
}
