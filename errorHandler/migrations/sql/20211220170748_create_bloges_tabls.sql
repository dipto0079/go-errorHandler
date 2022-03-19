-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS blogs(
    id SERIAL NOT NULL,
    cat_id INTEGER NOT NULL,
    title TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    image TEXT NOT NULL,

    PRIMARY KEY(id),
    CONSTRAINT fk_category
      FOREIGN KEY(cat_id) 
	  REFERENCES categories(id)
	  ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blogs;
-- +goose StatementEnd
