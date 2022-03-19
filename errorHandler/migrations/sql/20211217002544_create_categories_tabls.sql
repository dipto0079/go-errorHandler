-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories(
    id SERIAL NOT NULL,
    title TEXT NOT NULL,
    is_completed BOOLEAN DEFAULT false,

    PRIMARY KEY(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
