-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (gen_random_uuid(), NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC', $1, $2)
RETURNING *;


-- name: DeleteAllChirps :exec
DELETE from chirps;


-- name: GetAllChirps :many
SELECT id, created_at, updated_at, body, user_id 
FROM chirps 
WHERE user_id = COALESCE($1, user_id)  -- If $1 is NULL, matches any user_id (user_id = user_id : always true thus returns all)
ORDER BY created_at ASC;


-- name: GetChirpByID :one
SELECT id, created_at, updated_at, body, user_id FROM chirps WHERE id = $1;

-- name: DeleteChirpByID :exec
DELETE FROM chirps WHERE id = $1;