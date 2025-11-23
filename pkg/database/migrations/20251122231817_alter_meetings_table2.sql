-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE meeting
    ADD CONSTRAINT meeting_agent_id_fkey
    FOREIGN KEY (agent_id)
    REFERENCES agent(id)
    ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE meeting
    DROP CONSTRAINT IF EXISTS meeting_agent_id_fkey;
-- +goose StatementEnd
