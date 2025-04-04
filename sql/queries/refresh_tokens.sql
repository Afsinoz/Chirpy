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



