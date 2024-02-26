-- +goose Up
-- +goose StatementBegin
CREATE TABLE akash_group_requirements_attributes (
  id SERIAL PRIMARY KEY,
  group_id INT NOT NULL,
  attribute_key TEXT NOT NULL,
  attribute_value TEXT NOT NULL,
  FOREIGN KEY (group_id) REFERENCES akash_groups(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE akash_group_requirements_attributes;
-- +goose StatementEnd