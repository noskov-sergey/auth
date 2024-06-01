-- +goose Up
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    role integer NOT NULL,
    created_at timestamp NOT NULL default now(),
    updated_at timestamp
);

-- +goose Down
DROP TABLE users;
