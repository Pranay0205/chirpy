-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($1, $2, $3, $4, $5, NULL)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = $1 WHERE token = $2;

-- name: GetUserFromRefreshToken :one
SELECT u.id, u.created_at, u.updated_at, u.email
FROM users u
WHERE u.id = (
    SELECT user_id 
    FROM refresh_tokens 
    WHERE token = $1
);