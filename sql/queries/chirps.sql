-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (gen_random_uuid(), NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC', $1, $2)
RETURNING *;


-- name: DeleteAllChirps :exec
DELETE from chirps;