CREATE TABLE akash_resources (
  id SERIAL PRIMARY KEY,
  group_dseq INT NOT NULL UNIQUE,
  cpu_units TEXT,
  memory_quantity TEXT,
  gpu_units TEXT,
  price_denom TEXT,
  price_amount TEXT,
);