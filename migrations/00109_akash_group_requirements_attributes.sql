-- +goose Up
-- +goose StatementBegin
CREATE TABLE akash_group_requirements_attributes (
  id SERIAL PRIMARY KEY,
  group_dseq INT NOT NULL UNIQUE,
  attribute_key TEXT NOT NULL,
  attribute_value TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE akash_group_requirements_attributes;
-- +goose StatementEnd