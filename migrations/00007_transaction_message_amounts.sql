-- +goose Up
-- +goose StatementBegin
CREATE TABLE transaction_message_amounts (
  id SERIAL PRIMARY KEY,
  amount INT NOT NULL,
  message_id INT NOT NULL,
  denom_id INT NOT NULL,
  FOREIGN KEY (message_id) REFERENCES messages(id),
  FOREIGN KEY (denom_id) REFERENCES denoms(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transaction_message_amounts;
-- +goose StatementEnd