-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'errorhandler'
          AND column_name = 'created_at'
          AND data_type = 'timestamp without time zone'
    ) THEN
        ALTER TABLE errorhandler
            ALTER COLUMN created_at TYPE TIMESTAMPTZ,
            ALTER COLUMN deleted_at TYPE TIMESTAMPTZ,
            ALTER COLUMN last_seen_at TYPE TIMESTAMPTZ,
            ALTER COLUMN resolved_at TYPE TIMESTAMPTZ;
    END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd
