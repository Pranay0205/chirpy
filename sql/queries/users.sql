-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(), NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC', $1, $2)
RETURNING *;


-- name: DeleteAllUsers :exec
DELETE from users;


-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirp_red FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = $3
WHERE id = $4
RETURNING *;


-- name: UpgradeUserToChirpRed :exec
UPDATE users
SET is_chirp_red = true WHERE id = $1;