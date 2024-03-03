CREATE TABLE akash_resources (
  id SERIAL PRIMARY KEY,
  group_dseq INT NOT NULL UNIQUE,
  resource_id INT NOT NULL,
  cpu_units TEXT,
  memory_quantity TEXT,
  gpu_units TEXT,
  price_denom TEXT,
  price_amount TEXT,
  FOREIGN KEY (group_dseq) REFERENCES akash_groups(dseq),
  FOREIGN KEY (price_denom) REFERENCES denoms(denom)
);