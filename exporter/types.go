package exporter

import "github.com/jim380/Cendermint/rest"

var (
	previousBlockHeight int64

	gaugesNamespaceList = [...]string{

		// chain
		"chain_blockHeight",
		"chain_block_interval",

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
		"slashing_start_Height",
		"slashing_index_offset",
		"slashing_jailed_until",
		"slashing_tombstoned",
		"slashing_missed_blocks_counter",

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
		"validator_last_signed_height",
		"validator_miss_consecutive",
		"validator_miss_threshold",
		"validator_miss_count",

		// upgrade
		"upgrade_planned",

		// ibc
		"ibc_channels_total",
		"ibc_channels_open",
		"ibc_connections_total",
		"ibc_connections_open",

		// tx
		"tx_tps",
		"tx_gas_wanted_total",
		"tx_gas_used_total",
		"tx_events_total",
		"tx_delegate_total",
		"tx_message_total",
		"tx_transfer_total",
		"tx_unbond_total",
		"tx_withdraw_rewards_total",
		"tx_create_validator_total",
		"tx_redelegate_total",
		"tx_proposal_vote_total",
		"tx_ibc_fungible_token_packet_total",
		"tx_ibc_transfer_total",
		"tx_ibc_update_client_total",
		"tx_ibc_ack_packet_total",
		"tx_ibc_send_packet_yotal",
		"tx_ibc_recv_packet_total",
		"tx_ibc_timeout_total",
		"tx_ibc_timeout_packet_total",
		"tx_ibc_denom_trace_total",
		"tx_swap_swap_within_batch_total",
		"tx_swap_withdraw_within_batch_total",
		"tx_swap_deposit_within_batch_total",
		"tx_others_total",

		// peggo
		"peggo_valset_count",
		"peggo_valset_active",
		"peggo_last_claim_nonce",
		"peggo_last_claim_height",
		"peggo_erc20_price",
		"peggo_batch_fees",
		"peggo_batches_fees",
		"peggo_bridge_fees",
	}

	metricData metric
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

	Peggo struct {
		ValSetCount     float64
		ValSetActive    float64
		LastClaimNonce  float64
		LastClaimHeight float64
		Erc20Price      float64
		BatchFees       float64
		BatchesFees     float64
		BridgeFees      float64
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

	case "nym":
		dList = []string{"upunk"}
	case "evmos":
		dList = []string{"aphoton"}
	}
	return dList
}
