-- +goose Up
-- +goose StatementBegin
CREATE TABLE akash_deployments (
  owner TEXT NOT NULL,
  dseq INT NOT NULL,
  state TEXT NOT NULL,
  version TEXT NOT NULL,
  created_at INT NOT NULL,
  last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (owner, dseq),
  FOREIGN KEY (owner) REFERENCES akash_providers(owner)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE akash_deployments;
-- +goose StatementEnd