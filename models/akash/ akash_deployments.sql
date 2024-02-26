CREATE TABLE akash_deployments (
  owner TEXT NOT NULL,
  dseq INT NOT NULL,
  state TEXT NOT NULL,
  version TEXT NOT NULL,
  created_at INT NOT NULL,
  PRIMARY KEY (owner, dseq),
  FOREIGN KEY (owner) REFERENCES akash_providers(owner)
);