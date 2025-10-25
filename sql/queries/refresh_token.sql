-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (expires_at, token, user_id)
VALUES ($1, $2, $3)
RETURNING *;


-- name: GetRefreshTokenInfo :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens SET updated_at = $2, revoked_at = $3 WHERE token = $1;