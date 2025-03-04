-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES(
    $1,
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: DeleteUsers :exec
TRUNCATE users CASCADE;
