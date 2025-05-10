CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)    NOT NULL,
    email VARCHAR(255)   UNIQUE NOT NULL,
    age INTEGER          NOT NULL,
    password_hash TEXT   NOT NULL
);
