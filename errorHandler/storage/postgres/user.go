package postgres

import (
	"context"
	"errorHandler/errorHandler/storage"
)

const insertUser = `
	INSERT INTO users(
		name,
		email,
		password,
		email_verified
	)VALUES(
		:name,
		:email,
		:password,
		:email_verified
	)RETURNING id;
`

func (s *Storage) CreateUser(ctx context.Context, t storage.User) (int64, error) {
	stmt, err := s.db.PrepareNamed(insertUser)
	if err != nil {
		return 0, err
	}
	var id int64
	if err := stmt.Get(&id, t); err != nil {
		return 0, err
	}
	return id, nil
}
