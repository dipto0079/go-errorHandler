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
	StackTrace      string       `db:"stack_trace"`
	Status          string       `db:"status"`
	Fingerprint     string       `db:"fingerprint"`
	OccurrenceCount int          `db:"occurrence_count"`
	LastSeenAt      time.Time    `db:"last_seen_at"`
	ResolvedAt      sql.NullTime `db:"resolved_at"`
}

const (
	StatusOpen     = "open"
	StatusResolved = "resolved"
)

// ListFilter narrows ListErrors. Zero values mean "no filter".
type ListFilter struct {
	EnvType        string
	Status         string
	ErrorCode      string
	Since          time.Time
	Limit          int
	Offset         int
	IncludeDeleted bool
}

var (
	// NotFound is returned when the requested resource does not exist.
	NotFound = errors.New("not found")
)
