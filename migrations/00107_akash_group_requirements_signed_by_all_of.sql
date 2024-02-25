-- +goose Up
-- +goose StatementBegin
CREATE TABLE akash_group_requirements_signed_by_all_of (
  id SERIAL PRIMARY KEY,
  group_id INT NOT NULL,
  signed_by_all_of TEXT,
  FOREIGN KEY (group_id) REFERENCES akash_groups(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE akash_group_requirements_signed_by_all_of;
-- +goose StatementEnd