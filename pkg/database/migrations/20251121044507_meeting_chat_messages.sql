-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meeting_id UUID NOT NULL REFERENCES meeting(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL, -- "ai" for bot, or user ID for user
    role TEXT NOT NULL, -- "user" or "ai"
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE meeting_chat_messages;
-- +goose StatementEnd
