-- name: CreateAgent :one
INSERT INTO agent (name, user_id, instructions)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAgentByID :one
SELECT * FROM agent WHERE id = $1;

-- name: GetAgentsByUserID :many
SELECT * FROM agent WHERE user_id = $1;

-- name: GetAgentByUserID :one
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
