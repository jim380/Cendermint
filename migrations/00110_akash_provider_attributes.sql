-- +goose Up
-- +goose StatementBegin
CREATE TABLE akash_provider_attributes (
  id SERIAL PRIMARY KEY,
  provider_owner TEXT NOT NULL,
  attribute_key TEXT NOT NULL,
  attribute_value TEXT NOT NULL,
  FOREIGN KEY (provider_owner) REFERENCES akash_providers(owner)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE akash_provider_attributes;
-- +goose StatementEnd