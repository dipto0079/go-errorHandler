package postgres

import (
	"context"
	"database/sql"

	"github.com/dipto0079/go-errorHandler/errorHandler/storage"
)

const insertError = `
INSERT INTO errorhandler (
	error_code,
	error_details,
	env_type,
	created_at,
	created_by
) VALUES (
	:error_code,
	:error_details,
	:env_type,
	:created_at,
	:created_by
) RETURNING
	id
`

func (s *Storage) CreateError(ctx context.Context, errors storage.ErrorHandler) (string, error) {

	stmt, err := s.db.PrepareNamed(insertError)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var id string
	if err := stmt.GetContext(ctx, &id, errors); err != nil {
		return "", err
	}

	return id, nil
}

const getError = `
	SELECT 
		id,
		error_code,
		error_details,
		env_type,
		created_by
	FROM errorhandler
	WHERE id = $1
`

func (s *Storage) GetError(ctx context.Context, id string) (*storage.ErrorHandler, error) {

	var errorHandler storage.ErrorHandler
	if err := s.db.GetContext(ctx, &errorHandler, getError, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &errorHandler, nil
}

func (s *Storage) ListError(ctx context.Context) ([]storage.ErrorHandler, error) {

	var errors []storage.ErrorHandler

	if err := s.db.SelectContext(ctx, &errors, "SELECT * FROM errorhandler ORDER BY id DESC"); err != nil {

		return errors, err
	}
	return errors, nil
}

const deleteError = `
	UPDATE
	errorhandler
	SET
		deleted_at = now(),
		delete_by_env_type = $1
	WHERE 
		id = $2
`

func (s *Storage) DeleteError(ctx context.Context, id string, env string) error {

	res, err := s.db.ExecContext(ctx, deleteError, env, id)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err == nil && n == 0 {
		return storage.NotFound
	}
	return nil
}
