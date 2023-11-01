-- +goose Up
-- +goose StatementBegin
CREATE TABLE absent_validators (
  block_height INT REFERENCES blocks(height),
  cons_pub_address TEXT REFERENCES validators(cons_pub_address),
  PRIMARY KEY (block_height, cons_pub_address)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE absent_validators;
-- +goose StatementEnd