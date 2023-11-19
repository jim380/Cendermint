-- +goose Up
-- +goose StatementBegin
CREATE TABLE absent_validators (
  block_height INT REFERENCES blocks(height),
  cons_pub_key TEXT REFERENCES validators(cons_pub_key),
  PRIMARY KEY (block_height, cons_pub_key)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE absent_validators;
-- +goose StatementEnd