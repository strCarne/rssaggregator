-- +goose Up
ALTER TABLE users ADD COLUMN api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (

    -- generating random bytes
    -- casting them into byte array
    -- hashing the byte array with sha256, so it will have fixed length
    -- encoding it into hexadecimal
    encode(sha256(random()::text::bytea ), 'hex')
);

-- +goose Down
ALTER TABLE users DROP COLUMN api_key;