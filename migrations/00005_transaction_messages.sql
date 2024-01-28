-- +goose Up
-- +goose StatementBegin
CREATE TABLE transaction_messages (
  id SERIAL PRIMARY KEY,
  type TEXT NOT NULL,
  transaction_hash TEXT NOT NULL,
  FOREIGN KEY (transaction_hash) REFERENCES transactions(hash)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transaction_messages;
-- +goose StatementEnd