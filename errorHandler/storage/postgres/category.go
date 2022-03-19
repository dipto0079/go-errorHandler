package postgres

import (
	"context"
	"errorHandler/errorHandler/storage"
)

const insertCategory = `
	INSERT INTO categories(
		title
	)VALUES(
		:title
	)RETURNING id;
`

func (s *Storage) Create_sto(ctx context.Context, t storage.Category) (int64, error) {
	stmt, err := s.db.PrepareNamed(insertCategory)
	if err != nil {
		return 0, err
	}
	var id int64
	if err := stmt.Get(&id, t); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) Get_sto(ctx context.Context, id int64) (storage.Category, error) {
	var c storage.Category

	if err := s.db.Get(&c, "SELECT * FROM categories WHERE id=$1", id); err != nil {
		return c, err
	}
	return c, nil
}

const updateCategory = `
	UPDATE categories 
	SET
		title =:title
		
	WHERE
	id =:id
	RETURNING *;
`

func (s *Storage) Update(ctx context.Context, t storage.Category) error {
	stmt, err := s.db.PrepareNamed(updateCategory)
	if err != nil {
		return err
	}
	var cat storage.Category
	if err := stmt.Get(&cat, t); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Delete(ctx context.Context, id int64) error {
	var b storage.Category
	if err := s.db.Get(&b, "DELETE FROM categories WHERE id=$1 RETURNING * ", id); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Get_all_Data(ctx context.Context) ([]storage.Category, error) {

	var c []storage.Category

	if err := s.db.Select(&c, "SELECT * FROM categories order by id desc"); err != nil {
		return c, err
	}
	return c, nil
}
