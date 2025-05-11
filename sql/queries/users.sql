-- name: CreateUser :one
INSERT INTO users (
    email
) VALUES (
    $1
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: DeleteAllUsers :exec
DELETE FROM users;