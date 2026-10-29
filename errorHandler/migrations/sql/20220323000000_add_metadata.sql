-- +goose Up
-- +goose StatementBegin
ALTER TABLE errorhandler
    ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE errorhandler DROP COLUMN IF EXISTS metadata;
-- +goose StatementEnd
