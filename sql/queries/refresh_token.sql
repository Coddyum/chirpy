-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (expires_at, token)
VALUES ($1, $2)
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at;


-- name: GetRefreshTokenInfo :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens SET updated_at = $2, revoked_at = $3 WHERE token = $1;