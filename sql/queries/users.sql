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

-- name: UpdateUser :exec
UPDATE users SET email=$1, hashed_password=$2, updated_at=$3 WHERE id=$4; 
