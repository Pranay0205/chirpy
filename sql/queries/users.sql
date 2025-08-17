-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(), NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC', $1, $2)
RETURNING *;


-- name: DeleteAllUsers :exec
DELETE from users;


-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password FROM users WHERE email = $1;

