CREATE TABLE akash_groups (
  id SERIAL,
  owner TEXT NOT NULL,
  dseq INT NOT NULL,
  gseq INT NOT NULL,
  state TEXT NOT NULL,
  name TEXT NOT NULL,
  created_at INT NOT NULL,
  PRIMARY KEY (id),
  UNIQUE (owner, dseq, gseq),
  FOREIGN KEY (owner, dseq) REFERENCES akash_deployments(owner, dseq)
);