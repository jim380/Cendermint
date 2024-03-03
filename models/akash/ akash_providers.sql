CREATE TABLE  akash_providers (
  owner TEXT PRIMARY KEY UNIQUE NOT NULL,
  host_uri TEXT NOT NULL,
  email TEXT,
  website TEXT
);