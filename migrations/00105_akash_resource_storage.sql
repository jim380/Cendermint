-- +goose Up
-- +goose StatementBegin
CREATE TABLE akash_resource_storage (
  id SERIAL PRIMARY KEY,
  resource_id INT NOT NULL,
  name TEXT NOT NULL,
  quantity TEXT NOT NULL,
  FOREIGN KEY (resource_id) REFERENCES akash_resources(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE akash_resource_storage;
-- +goose StatementEnd