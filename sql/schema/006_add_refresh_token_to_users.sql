-- +goose Up
ALTER TABLE users
ADD COLUMN refresh_token TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE users
DROP COLUMN refresh_token;
