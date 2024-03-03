CREATE TABLE akash_provider_auditors (
  id SERIAL PRIMARY KEY,
  provider_owner TEXT NOT NULL,
  auditor TEXT NOT NULL,
  last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (provider_owner) REFERENCES akash_providers(owner),
  UNIQUE (provider_owner, auditor)
);