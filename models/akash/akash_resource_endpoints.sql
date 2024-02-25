CREATE TABLE akash_resource_endpoints (
  id SERIAL PRIMARY KEY,
  resource_id INT NOT NULL,
  kind TEXT NOT NULL,
  sequence_number INT NOT NULL,
  FOREIGN KEY (resource_id) REFERENCES akash_resources(id)
);