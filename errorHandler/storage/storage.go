package storage

import (
	"database/sql"
	"errors"
	"time"
)

type ErrorHandler struct {
	ID              string       `db:"id"`
	ErrorCode       string       `db:"error_code"`
	ErrorDetails    string       `db:"error_details"`
	EnvType         string       `db:"env_type"`
	CreatedAt       time.Time    `db:"created_at,omitempty"`
	CreatedBy       string       `db:"created_by"`
	DeletedAt       sql.NullTime `db:"deleted_at,omitempty"`
	DeleteByEnvType string       `db:"delete_by_env_type"`
}

var (
	// NotFound is returned when the requested resource does not exist.
	NotFound = errors.New("not found")
)
