-- +goose Up
CREATE TABLE users
(
    id         UUID PRIMARY KEY,
    email      TEXT UNIQUE NOT NULL,
    password   TEXT        NOT NULL,
    role       TEXT        NOT NULL,
    created_at TIMESTAMP   NOT NULL,
    updated_at TIMESTAMP   NOT NULL
);

CREATE TABLE refresh_tokens
(
    id         UUID PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token      TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP   NOT NULL
);

-- +goose Down
DROP TABLE users;
DROP TABLE refresh_tokens;