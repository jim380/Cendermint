-- +goose Up
-- +goose StatementBegin
CREATE TABLE denoms (
  id SERIAL PRIMARY KEY,
  denom TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE denoms;
-- +goose StatementEnd