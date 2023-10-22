package exporter

var gaugesNamespaceList = [...]string{
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
	"gov_voting_proposal_voted_count",
	"gov_voting_proposal_did_not_vote_count",

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

	// gravity
	"gravity_signed_valsets_window",
	"gravity_signed_batches_window",
	"gravity_target_batch_timeout",
	"gravity_slash_fraction_valset",
	"gravity_slash_fraction_batch",
	"gravity_slash_fraction_bad_eth_sig",
	"gravity_valset_reward_amount",
	"gravity_bridge_active",
	"gravity_valset_count",
	"gravity_valset_active",
	"gravity_event_nonce",
	// "gravity_last_claim_height",
	"gravity_erc20_price",
	"gravity_batch_fees",
	"gravity_batches_fees",
	"gravity_bridge_fees",

	// akash
	"akash_total_deployments",
	"akash_active_deployments",
	"akash_closed_deployments",

	// oracle
	"oracle_validator_missed_blocks_counter",
}
