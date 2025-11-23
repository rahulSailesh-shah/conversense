-- name: CreateAgent :one
INSERT INTO agent (name, user_id, instructions)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAgentByID :one
SELECT * FROM agent WHERE id = $1;

-- name: GetAgents :many
SELECT
    a.id,
    a.name,
    a.instructions,
    a.user_id,
    a.created_at,
    a.updated_at,
    COUNT(m.id) AS meeting_count,
    COUNT(*) OVER() AS total_count
FROM agent a
LEFT JOIN meeting m ON a.id = m.agent_id
WHERE a.user_id = $1
  AND ($2::text = '' OR a.name ILIKE '%' || $2 || '%')
GROUP BY
    a.id,
    a.name,
    a.instructions,
    a.user_id,
    a.created_at,
    a.updated_at
ORDER BY a.updated_at DESC
LIMIT $3 OFFSET $4;


-- name: GetAgent :one
SELECT
 a.*,
 COALESCE(m.meeting_count, 0) AS meeting_count
FROM agent a
LEFT JOIN (
    SELECT agent_id, COUNT(*) AS meeting_count
    FROM meeting
    WHERE agent_id = $1
    GROUP BY agent_id
) m ON a.id = m.agent_id
WHERE a.id = $1 AND a.user_id = $2;

-- name: UpdateAgent :one
UPDATE agent
SET name = $2, instructions = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAgent :exec
DELETE FROM agent WHERE id = $1 AND user_id = $2 ;
