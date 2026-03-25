-- +goose Up
-- +goose StatementBegin

-- Add indexes for common lookup patterns (no schema type changes to avoid breaking sqlc generated code)
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_listens_user_id ON listens (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_listens_user_id;
DROP INDEX IF EXISTS idx_sessions_user_id;

-- +goose StatementEnd
