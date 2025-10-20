-- +goose Up
ALTER TABLE users
ADD COLUMN token TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE users
DROP COLUMN token;
