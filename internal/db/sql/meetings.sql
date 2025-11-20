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
    m.transcript_url,
    m.recording_url,
    m.summary,
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
    name = COALESCE($3, name),
    agent_id = COALESCE($4, agent_id),
    status = COALESCE($5, status),
    start_time = COALESCE($6, start_time),
    end_time = COALESCE($7, end_time),
    transcript_url = COALESCE($8, transcript_url),
    recording_url = COALESCE($9, recording_url),
    summary = COALESCE($10, summary),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteMeeting :exec
DELETE FROM meeting WHERE id = $1;

-- name: DeleteMeetingsByUserID :exec
DELETE FROM meeting WHERE user_id = $1;
