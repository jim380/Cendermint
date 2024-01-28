-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
  hash TEXT PRIMARY KEY UNIQUE NOT NULL,
  height INT NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  type TEXT NOT NULL,
  gas_wanted INT NOT NULL,
  gas_used INT NOT NULL,
  memo TEXT,
  payer TEXT,
  granter TEXT,
  FOREIGN KEY (height) REFERENCES blocks(height)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd