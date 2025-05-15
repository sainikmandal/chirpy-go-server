-- name: CreateChirp :one
INSERT INTO chirps (
    body,
    user_id
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE id = $1 LIMIT 1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;
