-- +goose Up
-- +goose StatementBegin
CREATE TABLE transaction_fee_amounts (
  id SERIAL PRIMARY KEY,
  amount INT NOT NULL,
  transaction_hash TEXT NOT NULL,
  denom_id INT NOT NULL,
  FOREIGN KEY (transaction_hash) REFERENCES transactions(hash),
  FOREIGN KEY (denom_id) REFERENCES denoms(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transaction_fee_amounts;
-- +goose StatementEnd