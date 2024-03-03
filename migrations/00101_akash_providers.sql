-- +goose Up
-- +goose StatementBegin
CREATE TABLE akash_providers (
  owner TEXT PRIMARY KEY UNIQUE NOT NULL,
  host_uri TEXT NOT NULL,
  email TEXT,
  website TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE akash_providers;
-- +goose StatementEnd