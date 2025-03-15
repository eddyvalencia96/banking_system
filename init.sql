\c banking_system;

CREATE TABLE users (
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    initial_balance INTEGER NOT NULL,
    balance INTEGER NOT NULL
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    transaction_date TIMESTAMP NOT NULL,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

INSERT INTO users ("username", "password")
VALUES ('edvalencia', '960919edval');

CREATE USER "bs_ed" WITH PASSWORD '960919';
GRANT ALL PRIVILEGES ON DATABASE banking_system TO "bs_ed";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE accounts TO "bs_ed";
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE transactions TO bs_ed;
GRANT USAGE, SELECT ON SEQUENCE accounts_id_seq TO bs_ed;
GRANT USAGE, SELECT ON SEQUENCE transactions_id_seq TO bs_ed;