-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users WHERE refresh_token = $1;

-- name: GetUserFromAccessToken :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserInformation :one
UPDATE users 
SET updated_at = $2,
    email = $3, 
    hashed_password = $4 
WHERE id = $1 
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;
