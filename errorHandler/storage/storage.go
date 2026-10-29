package storage

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
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
	Metadata        Metadata     `db:"metadata"`
}

// Metadata is free-form key/value context stored as JSONB.
type Metadata map[string]string

func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return "{}", nil
	}
	b, err := json.Marshal(m)
	return string(b), err
}

func (m *Metadata) Scan(src interface{}) error {
	var b []byte
	switch v := src.(type) {
	case nil:
		*m = nil
		return nil
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("storage: cannot scan %T into Metadata", src)
	}
	return json.Unmarshal(b, m)
}

// Stats summarises stored errors that are not deleted.
type Stats struct {
	Total    int
	Open     int
	Resolved int
	Last24h  int // errors last seen in the past 24 hours
	ByEnv    map[string]int
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
