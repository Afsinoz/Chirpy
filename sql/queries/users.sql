-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES(
    $1,
    NOW(),
    $2,
    $3,
    $4
)
RETURNING *;

-- name: DeleteUsers :exec
TRUNCATE users CASCADE;

-- name: GetUser :one
SELECT * FROM users WHERE $1 = email;
