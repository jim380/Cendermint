CREATE TABLE transaction_signers (
  transaction_hash TEXT NOT NULL,
  address_id INT NOT NULL,
  PRIMARY KEY (transaction_hash, address_id),
  FOREIGN KEY (transaction_hash) REFERENCES transactions(hash),
  FOREIGN KEY (address_id) REFERENCES addresses(id)
);