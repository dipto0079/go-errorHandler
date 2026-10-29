-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS errorhandler(
    id SERIAL NOT NULL,
    error_code TEXT NOT NULL,
    error_details TEXT NOT NULL,
    env_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL,
    deleted_at TIMESTAMP,
    delete_by_env_type TEXT,

    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS errorhandler;
-- +goose StatementEnd
