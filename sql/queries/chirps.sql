-- name: CreateChirp :one
INSERT INTO chirps (body, user_id)
VALUES ($1, $2)
RETURNING id, created_at, updated_at, body, user_id;

-- name: SelectAllChirps :many
SELECT *
FROM chirps
ORDER BY created_at ASC;

-- name: SelectOneChirp :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: SelectChirpByAuthor :many
SELECT *
FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: DeleteOneChirp :exec
DELETE FROM chirps
WHERE id = $1 AND user_id = $2;