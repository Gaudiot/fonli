-- +goose Up
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    canonical_username VARCHAR(255) NOT NULL UNIQUE,
    lifestyle VARCHAR(255),
    lifestyle_topics TEXT
);

-- +goose Down
DROP TABLE IF EXISTS users;
