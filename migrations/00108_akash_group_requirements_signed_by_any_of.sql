-- +goose Up
-- +goose StatementBegin
CREATE TABLE akash_group_requirements_signed_by_any_of (
  group_dseq INT NOT NULL UNIQUE PRIMARY KEY,
  signed_by_any_of TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE akash_group_requirements_signed_by_any_of;
-- +goose StatementEnd