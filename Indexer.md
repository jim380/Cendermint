## Indexer

### Tables

| Table                     | Columns                                                                                     |
| ------------------------- | ------------------------------------------------------------------------------------------- |
| `blocks`                  | `height`, `block_hash`, `proposer_address`, `txn_count`, `timestamp`                        |
| `validators`              | `cons_pub_key`, `cons_address`, `cons_address_hex`, `moniker`, `last_active`                |
| `absent_validators`       | `block_height`, `cons_pub_key`                                                              |
| `transactions`            | `hash`, `height`, `timestamp`, `type`, `gas_wanted`, `gas_used`, `memo`, `payer`, `granter` |
| `transaction_fee_amounts` | `id`, `amount`, `transaction_hash`, `denom_id`                                              |
| `transaction_messages`    | `id`, `type`, `transaction_hash`                                                            |
| `denoms`                  | `id`, `denom`                                                                               |
