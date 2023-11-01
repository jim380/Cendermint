-- +goose Up
-- +goose StatementBegin
CREATE TABLE validators (
  cons_hex_address TEXT PRIMARY KEY UNIQUE NOT NULL,
  moniker TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE validators;
-- +goose StatementEnd