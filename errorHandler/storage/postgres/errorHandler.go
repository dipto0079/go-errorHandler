package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/dipto0079/go-errorHandler/errorHandler/storage"
)

const insertError = `
INSERT INTO errorhandler (
	error_code,
	error_details,
	env_type,
	created_at,
	created_by,
	stack_trace,
	fingerprint,
	last_seen_at,
	metadata
) VALUES (
	:error_code,
	:error_details,
	:env_type,
	:created_at,
	:created_by,
	:stack_trace,
	:fingerprint,
	:created_at,
	:metadata
)
ON CONFLICT (fingerprint) WHERE fingerprint <> '' AND status = 'open' AND deleted_at IS NULL
DO UPDATE SET
	occurrence_count = errorhandler.occurrence_count + 1,
	last_seen_at = EXCLUDED.last_seen_at
RETURNING
	id, occurrence_count
`

// SaveError stores an error. When e.Fingerprint is set and an open error with
// the same fingerprint exists, that row's occurrence count is incremented
// instead. It returns the row id and the resulting occurrence count.
func (s *Storage) SaveError(ctx context.Context, e storage.ErrorHandler) (string, int, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, insertError)
	if err != nil {
		return "", 0, err
	}
	defer stmt.Close()
	var res struct {
		ID    string `db:"id"`
		Count int    `db:"occurrence_count"`
	}
	if err := stmt.GetContext(ctx, &res, e); err != nil {
		return "", 0, err
	}
	return res.ID, res.Count, nil
}

func (s *Storage) CreateError(ctx context.Context, e storage.ErrorHandler) (string, error) {
	id, _, err := s.SaveError(ctx, e)
	return id, err
}

const errorColumns = `
	id,
	error_code,
	error_details,
	env_type,
	created_at,
	created_by,
	deleted_at,
	COALESCE(delete_by_env_type, '') AS delete_by_env_type,
	stack_trace,
	status,
	fingerprint,
	occurrence_count,
	last_seen_at,
	resolved_at,
	metadata
`

const getError = `SELECT ` + errorColumns + ` FROM errorhandler WHERE id = $1 AND deleted_at IS NULL`

func (s *Storage) GetError(ctx context.Context, id string) (*storage.ErrorHandler, error) {

	var errorHandler storage.ErrorHandler
	if err := s.db.GetContext(ctx, &errorHandler, getError, id); err != nil {
		if err == sql.ErrNoRows || isInvalidID(err) {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &errorHandler, nil
}

func (s *Storage) ListError(ctx context.Context) ([]storage.ErrorHandler, error) {
	return s.ListErrors(ctx, storage.ListFilter{})
}

// ListErrors returns errors, newest first, narrowed by f.
func (s *Storage) ListErrors(ctx context.Context, f storage.ListFilter) ([]storage.ErrorHandler, error) {
	var (
		where []string
		args  []interface{}
	)
	add := func(cond string, v interface{}) {
		args = append(args, v)
		where = append(where, fmt.Sprintf(cond, len(args)))
	}
	if !f.IncludeDeleted {
		where = append(where, "deleted_at IS NULL")
	}
	if f.EnvType != "" {
		add("env_type = $%d", f.EnvType)
	}
	if f.Status != "" {
		add("status = $%d", f.Status)
	}
	if f.ErrorCode != "" {
		add("error_code = $%d", f.ErrorCode)
	}
	if !f.Since.IsZero() {
		add("last_seen_at >= $%d", f.Since)
	}

	q := "SELECT " + errorColumns + " FROM errorhandler"
	if len(where) > 0 {
		q += " WHERE " + strings.Join(where, " AND ")
	}
	q += " ORDER BY id DESC"
	if f.Limit > 0 {
		args = append(args, f.Limit)
		q += fmt.Sprintf(" LIMIT $%d", len(args))
	}
	if f.Offset > 0 {
		args = append(args, f.Offset)
		q += fmt.Sprintf(" OFFSET $%d", len(args))
	}

	var errors []storage.ErrorHandler
	if err := s.db.SelectContext(ctx, &errors, q, args...); err != nil {
		return nil, err
	}
	return errors, nil
}

const resolveError = `
	UPDATE errorhandler
	SET status = 'resolved', resolved_at = now()
	WHERE id = $1 AND deleted_at IS NULL
`

// isInvalidID reports whether err is Postgres rejecting a non-numeric id.
func isInvalidID(err error) bool {
	var pe *pq.Error
	return errors.As(err, &pe) && pe.Code == "22P02"
}

// ResolveError marks an error as resolved. A later occurrence of the same
// error starts a new open row.
func (s *Storage) ResolveError(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, resolveError, id)
	if err != nil {
		if isInvalidID(err) {
			return storage.NotFound
		}
		return err
	}
	if n, err := res.RowsAffected(); err == nil && n == 0 {
		return storage.NotFound
	}
	return nil
}

const deleteError = `
	UPDATE
	errorhandler
	SET
		deleted_at = now(),
		delete_by_env_type = $1
	WHERE 
		id = $2 AND deleted_at IS NULL
`

func (s *Storage) DeleteError(ctx context.Context, id string, env string) error {

	res, err := s.db.ExecContext(ctx, deleteError, env, id)
	if err != nil {
		if isInvalidID(err) {
			return storage.NotFound
		}
		return err
	}
	if n, err := res.RowsAffected(); err == nil && n == 0 {
		return storage.NotFound
	}
	return nil
}

// PurgeErrors permanently removes errors last seen before the given time.
// When resolvedOnly is true, open errors are kept. It returns the number removed.
func (s *Storage) PurgeErrors(ctx context.Context, before time.Time, resolvedOnly bool) (int64, error) {
	q := "DELETE FROM errorhandler WHERE last_seen_at < $1"
	if resolvedOnly {
		q += " AND status = 'resolved'"
	}
	res, err := s.db.ExecContext(ctx, q, before)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// ErrorStats counts errors that are not deleted.
func (s *Storage) ErrorStats(ctx context.Context) (storage.Stats, error) {
	st := storage.Stats{ByEnv: map[string]int{}}
	var rows []struct {
		Env    string `db:"env_type"`
		Status string `db:"status"`
		N      int    `db:"n"`
		Recent int    `db:"recent"`
	}
	const q = `
		SELECT env_type, status, count(*) AS n,
		       count(*) FILTER (WHERE last_seen_at >= now() - interval '24 hours') AS recent
		FROM errorhandler
		WHERE deleted_at IS NULL
		GROUP BY env_type, status`
	if err := s.db.SelectContext(ctx, &rows, q); err != nil {
		return st, err
	}
	for _, r := range rows {
		st.Total += r.N
		st.Last24h += r.Recent
		st.ByEnv[r.Env] += r.N
		switch r.Status {
		case storage.StatusOpen:
			st.Open += r.N
		case storage.StatusResolved:
			st.Resolved += r.N
		}
	}
	return st, nil
}
