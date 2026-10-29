package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	
)


type Storage struct {
	db *sqlx.DB
}

func NewStorage(dbstring string) (*Storage, error) {
	db, err := sqlx.Connect("postgres", dbstring)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres '%s'", dbstring)
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)
	return &Storage{db: db}, nil
}

func NewTestStorage(dbstring string, migrationDir string) (*Storage, func()) {
	db, teardown := MustNewDevelopmentDB(dbstring, migrationDir)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(time.Hour)

	return &Storage{db: db}, teardown
}


// NewStorageFromDB wraps an existing database connection. The caller keeps
// ownership of db, but Close will still close it.
func NewStorageFromDB(db *sql.DB) *Storage {
	return &Storage{db: sqlx.NewDb(db, driver)}
}

// Close closes the underlying database connection.
func (s *Storage) Close() error {
	return s.db.Close()
}

// Exec runs a raw SQL statement, used for applying the embedded schema.
func (s *Storage) Exec(ctx context.Context, query string) error {
	_, err := s.db.ExecContext(ctx, query)
	return err
}
