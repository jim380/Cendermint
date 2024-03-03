CREATE TABLE akash_resource_endpoints (
  id SERIAL PRIMARY KEY,
  group_dseq INT NOT NULL,
  kind TEXT NOT NULL,
  sequence_number INT NOT NULL,
  FOREIGN KEY (group_dseq) REFERENCES akash_resources(group_dseq)
);