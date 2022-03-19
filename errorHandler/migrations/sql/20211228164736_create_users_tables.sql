-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id SERIAL NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email_verified TEXT NOT NULL,
    is_verify BOOLEAN DEFAULT false,

    PRIMARY KEY(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
