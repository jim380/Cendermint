CREATE TABLE akash_provider_attributes (
  id SERIAL PRIMARY KEY,
  provider_owner TEXT NOT NULL,
  attribute_key TEXT NOT NULL,
  attribute_value TEXT NOT NULL,
  FOREIGN KEY (provider_owner) REFERENCES providers(owner)
);