CREATE TABLE  akash_providers (
  owner TEXT PRIMARY KEY UNIQUE NOT NULL,
  host_uri TEXT NOT NULL,
  email TEXT,
  website TEXT,
  last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);