package postgres

import (
	"context"
	"errorHandler/errorHandler/storage"
)

const insertBlog = `
	INSERT INTO blogs(
		cat_id,
		title,
		description,
		image
	)VALUES(
		:cat_id,
		:title,
		:description,
		:image
	)RETURNING id;
`

func (s *Storage) Create(ctx context.Context, t storage.Blog) (int64, error) {
	stmt, err := s.db.PrepareNamed(insertBlog)
	if err != nil {
		return 0, err
	}
	var id int64
	if err := stmt.Get(&id, t); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) ListBlog(ctx context.Context) ([]storage.Blog, error) {

	var b []storage.Blog

	if err := s.db.Select(&b, "SELECT blogs.id, cat_id, blogs.title, description, image, categories.title as catname FROM blogs LEFT JOIN categories ON categories.id = blogs.cat_id  order by blogs.id desc"); err != nil {
		return b, err
	}
	return b, nil
}

func (s *Storage) GetBlog(ctx context.Context, id int64) (storage.Blog, error) {
	var b storage.Blog

	if err := s.db.Get(&b, "SELECT * FROM blogs WHERE id=$1", id); err != nil {
		return b, err
	}
	return b, nil
}

const UpdateBlog = `
	UPDATE blogs 
	SET
		cat_id =:cat_id,
		title =:title,
		description =:description,
		image =:image
		
	WHERE
	id =:id
	RETURNING *;
`

func (s *Storage) UpdateBlog(ctx context.Context, t storage.Blog) error {
	stmt, err := s.db.PrepareNamed(UpdateBlog)
	if err != nil {
		return err
	}
	var blog storage.Blog
	if err := stmt.Get(&blog, t); err != nil {
		return err
	}
	return nil
}

func (s *Storage) BlogDelete(ctx context.Context, id int64) error {
	var b storage.Blog
	if err := s.db.Get(&b, "DELETE FROM blogs WHERE id=$1 RETURNING * ", id); err != nil {
		return err
	}
	return nil
}
