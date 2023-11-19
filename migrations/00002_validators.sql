-- +goose Up
-- +goose StatementBegin
CREATE TABLE validators (
  cons_pub_key TEXT PRIMARY KEY UNIQUE NOT NULL,
  cons_address TEXT NOT NULL,
  cons_address_hex TEXT NOT NULL,
  moniker TEXT NOT NULL,
  last_active TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE validators;
-- +goose StatementEnd