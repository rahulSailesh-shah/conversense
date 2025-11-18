-- name: CreateAgent :one
INSERT INTO agent (name, user_id, instructions)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAgentByID :one
SELECT * FROM agent WHERE id = $1;

-- name: GetAgents :many
SELECT id, name, instructions, user_id, created_at, updated_at, COUNT(*) OVER() as total_count
FROM agent
WHERE user_id = $1
    AND (CASE WHEN $2::text != '' THEN name ILIKE '%' || $2 || '%' ELSE TRUE END)
ORDER BY updated_at DESC
LIMIT $3 OFFSET $4;

-- name: GetAgent :one
SELECT * FROM agent WHERE id = $1 AND user_id = $2;

-- name: UpdateAgent :one
UPDATE agent
SET name = $2, instructions = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAgent :exec
DELETE FROM agent WHERE id = $1;

-- name: DeleteAgentsByUserID :exec
DELETE FROM agent WHERE user_id = $1;
