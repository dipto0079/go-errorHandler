package storage

import (
	"database/sql"
	"errors"
	"time"
)

type Category struct {
	ID         int64  `db:"id"`
	Title      string `db:"title"`
	IsComplete bool   `db:"is_completed"`
}

type Blog struct {
	ID          int64  `db:"id"`
	CatID       int64  `db:"cat_id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Image       string `db:"image"`
	CatName     string `db:"catname"`
}

type User struct {
	ID            int64  `db:"id"`
	Name          string `db:"name"`
	Email         string `db:"email"`
	Password      string `db:"password"`
	EmailVerified string `db:"email_verified"`
	IsVerify      bool   `db:"is_verify"`
}
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
