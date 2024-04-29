CREATE TABLE users (
    id SERIAL,
    user_id CHAR(12) NOT NULL PRIMARY KEY,
    user_name VARCHAR(60) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(40) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);