CREATE TABLE akash_escrow_accounts (
  id TEXT PRIMARY KEY,
  owner TEXT NOT NULL,
  state TEXT NOT NULL,
  balance_denom TEXT,
  balance_amount TEXT,
  transferred_denom TEXT,
  transferred_amount TEXT,
  settled_at INT,
  depositor TEXT,
  funds_denom TEXT,
  funds_amount TEXT,
  FOREIGN KEY (balance_denom) REFERENCES denoms(denom),
  FOREIGN KEY (transferred_denom) REFERENCES denoms(denom),
  FOREIGN KEY (funds_denom) REFERENCES denoms(denom)
);