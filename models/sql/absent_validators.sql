CREATE TABLE absent_validators (
  block_height INT REFERENCES blocks(height),
  validator_cons_hex_addr TEXT REFERENCES validators(cons_hex_address),
  PRIMARY KEY (block_height, validator_cons_hex_addr)
);