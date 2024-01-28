-- +goose Up
-- +goose StatementBegin
CREATE TABLE addresses (
  id SERIAL PRIMARY KEY,
  cons_pub_key TEXT NOT NULL UNIQUE,
  address TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE addresses;
-- +goose StatementEnd