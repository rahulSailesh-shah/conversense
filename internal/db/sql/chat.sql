-- name: CreateChatMessage :one
INSERT INTO meeting_chat_messages (
    meeting_id,
    user_id,
    role,
    content
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetChatMessages :many
SELECT * FROM meeting_chat_messages
WHERE meeting_id = $1
ORDER BY created_at ASC;

-- name: GetRecentChatMessages :many
SELECT * FROM meeting_chat_messages
WHERE meeting_id = $1
ORDER BY created_at DESC
LIMIT $2;
