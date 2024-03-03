-- +goose Up
-- +goose StatementBegin
CREATE TABLE akash_resources (
  id SERIAL PRIMARY KEY,
  group_dseq INT NOT NULL UNIQUE,
  resource_id INT NOT NULL,
  cpu_units TEXT,
  memory_quantity TEXT,
  gpu_units TEXT,
  price_denom TEXT,
  price_amount TEXT,
  FOREIGN KEY (group_dseq) REFERENCES akash_groups(id),
  FOREIGN KEY (price_denom) REFERENCES denoms(denom)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE akash_resources;
-- +goose StatementEnd