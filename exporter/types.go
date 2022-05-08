package exporter

import (
	"github.com/jim380/Cendermint/rest"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	defaultGauges       []prometheus.Gauge
	gaugesDenom         []prometheus.Gauge
	previousBlockHeight int64
	metricData          metric
)

type metric struct {
	Network struct {
		ChainID       string
		BlockHeight   int64
		BlockInterval int64

		NodeInfo struct {
			NodeID     string
			TMVersion  string
			Moniker    string
			Name       string
			AppName    string
			Version    string
			GitCommit  string
			GoVersion  string
			SDKVersion string
		}

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
			StartHeight             float64
			IndexOffset             float64
			JailedUntil             float64
			Tombstoned              float64
			MissedBlocksCounter     float64
		}

		Minting struct {
			Inflation       float64
			ActualInflation float64
		}

		Gov struct {
			TotalProposalCount      float64
			VotingProposalCount     float64
			InVotingVotedCount      float64
			InVotingDidNotVoteCount float64
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
			PrecommitStatus float64
			MissedCount     float64
			LastSigned      float64
			MissThreshold   float64
			MissConsecutive float64
		}
	}

	Upgrade struct {
		Planned float64
		Name    string
		Time    string
		Height  string
		Info    string
	}

	IBC struct {
		IBCChannels struct {
			Total float64
			Open  float64
		}
		IBCConnections struct {
			Total float64
			Open  float64
		}
	}

	Tx struct {
		TPS            float64
		GasWantedTotal float64
		GasUsedTotal   float64
		// default
		EventsTotal          float64
		DelegateTotal        float64
		MessageTotal         float64
		TransferTotal        float64
		UnbondTotal          float64
		WithdrawRewardsTotal float64
		CreateValidatorTotal float64
		RedelegateTotal      float64
		ProposalVote         float64
		// IBC
		FungibleTokenPacketTotal float64
		IbcTransferTotal         float64
		UpdateClientTotal        float64
		AckPacketTotal           float64
		WriteAckTotal            float64
		SendPacketTotal          float64
		RecvPacketTotal          float64
		TimeoutTotal             float64
		TimeoutPacketTotal       float64
		DenomTraceTotal          float64
		// swap
		SwapWithinBatchTotal     float64
		WithdrawWithinBatchTotal float64
		DepositWithinBatchTotal  float64

		OthersTotal float64
	}

	Gravity struct {
		gravityParams struct {
			SignedValsetsWindow    float64
			SignedBatchesWindow    float64
			TargetBatchTimeout     float64
			SlashFractionValset    float64
			SlashFractionBatch     float64
			SlashFractionBadEthSig float64
			ValsetReward           struct {
				Amount float64
			}
		}
		ValSetCount   float64
		ValSetActive  float64
		GravityActive float64
		EventNonce    float64
		// LastClaimHeight float64
		Erc20Price  float64
		BatchFees   float64
		BatchesFees float64
		BridgeFees  float64
	}
}
