-- +goose Up
-- +goose StatementBegin
ALTER TABLE meeting
ALTER COLUMN agent_id TYPE UUID USING agent_id::uuid;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE meeting
ALTER COLUMN agent_id TYPE VARCHAR(255);
-- +goose StatementEnd
