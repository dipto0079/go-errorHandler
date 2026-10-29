-- +goose Up
-- +goose StatementBegin
ALTER TABLE errorhandler
    ADD COLUMN IF NOT EXISTS stack_trace TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'open',
    ADD COLUMN IF NOT EXISTS fingerprint TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS occurrence_count INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS last_seen_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS resolved_at TIMESTAMP;

CREATE UNIQUE INDEX IF NOT EXISTS errorhandler_open_fingerprint
    ON errorhandler (fingerprint)
    WHERE fingerprint <> '' AND status = 'open' AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS errorhandler_env_type ON errorhandler (env_type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS errorhandler_env_type;
DROP INDEX IF EXISTS errorhandler_open_fingerprint;
ALTER TABLE errorhandler
    DROP COLUMN IF EXISTS resolved_at,
    DROP COLUMN IF EXISTS last_seen_at,
    DROP COLUMN IF EXISTS occurrence_count,
    DROP COLUMN IF EXISTS fingerprint,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS stack_trace;
-- +goose StatementEnd
