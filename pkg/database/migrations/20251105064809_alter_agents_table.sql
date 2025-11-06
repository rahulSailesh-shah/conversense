-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE agent ALTER COLUMN user_id TYPE VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE agent ALTER COLUMN user_id TYPE UUID USING user_id::UUID;
-- +goose StatementEnd
