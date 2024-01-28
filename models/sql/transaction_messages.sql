CREATE TABLE transaction_messages (
  id SERIAL PRIMARY KEY,
  type TEXT NOT NULL,
  from_address TEXT NOT NULL,
  to_address TEXT NOT NULL,
  transaction_hash TEXT NOT NULL,
  FOREIGN KEY (transaction_hash) REFERENCES transactions(hash)
);