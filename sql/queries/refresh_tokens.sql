-- name: CreateRefToken :one

INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, revoked_at, user_id) 
VALUES (
    $1,
    NOW(),
    NOW(),
    NOW() + interval '60 days',
    NULL,
    $2
)
RETURNING *;


-- name: GetUserFromRefreshToken :one
SELECT user_id, expires_at, revoked_at FROM refresh_tokens 
WHERE token=$1; 


-- name: RevokeRefToken :exec
UPDATE refresh_tokens
SET revoked_at=$1, updated_at=$2
WHERE token=$3;


