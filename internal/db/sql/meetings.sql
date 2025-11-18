-- name: CreateMeeting :one
INSERT INTO meeting (name, user_id, agent_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetMeetingByID :one
SELECT * FROM meeting WHERE id = $1;

-- name: GetMeetings :many
SELECT
    m.id,
    m.name,
    m.user_id,
    m.agent_id,
    m.start_time,
    m.end_time,
    m.status,
    m.created_at,
    m.updated_at,
    COUNT(*) OVER() as total_count,
    a.name AS agent_name,
    a.instructions AS agent_instructions
FROM meeting AS m
JOIN agent AS a
    ON m.agent_id = a.id
WHERE m.user_id = $1
    AND (
        CASE
            WHEN $2::text != '' THEN m.name ILIKE '%' || $2 || '%'
            ELSE TRUE
        END
    )
ORDER BY m.updated_at DESC
LIMIT $3 OFFSET $4;


-- name: GetMeeting :one
SELECT
    m.id,
    m.name,
    m.user_id,
    m.agent_id,
    m.start_time,
    m.end_time,
    m.status,
    m.created_at,
    m.updated_at,
    a.name AS agent_name,
    a.instructions AS agent_instructions
FROM meeting AS m
JOIN agent AS a
    ON m.agent_id = a.id
WHERE m.id = $1
    AND m.user_id = $2;

-- name: UpdateMeeting :one
UPDATE meeting
SET
    name = COALESCE($2, name),
    agent_id = COALESCE($3, agent_id),
    status = COALESCE($4, status),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteMeeting :exec
DELETE FROM meeting WHERE id = $1;

-- name: DeleteMeetingsByUserID :exec
DELETE FROM meeting WHERE user_id = $1;
