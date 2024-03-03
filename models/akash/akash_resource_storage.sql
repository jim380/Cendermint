CREATE TABLE akash_resource_storage (
  id SERIAL PRIMARY KEY,
  group_dseq INT NOT NULL,
  name TEXT NOT NULL,
  quantity TEXT NOT NULL,
  FOREIGN KEY (group_dseq) REFERENCES akash_resources(group_dseq)
);