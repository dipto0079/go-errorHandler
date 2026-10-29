// Package errorhandler stores application errors in PostgreSQL.
//
//	h, err := errorhandler.New(dsn, errorhandler.WithEnv("production"), errorhandler.WithService("api"))
//	defer h.Close()
//	h.Migrate(ctx)
//
//	h.Capture(ctx, err)                          // record an error, de-duplicated
//	http.ListenAndServe(addr, h.Middleware(mux)) // record panics
//
//	h.List(ctx, errorhandler.ListOptions{Status: errorhandler.StatusOpen, Limit: 20})
//	h.Resolve(ctx, id)
package errorhandler

import (
	"context"
	"database/sql"
	"time"

	"github.com/dipto0079/go-errorHandler/errorHandler/core/errorHandler"
	"github.com/dipto0079/go-errorHandler/errorHandler/storage"
	"github.com/dipto0079/go-errorHandler/errorHandler/storage/postgres"

	_ "github.com/lib/pq"
)

// Error is a stored error record.
type Error = storage.ErrorHandler

// ErrNotFound is returned when an error record does not exist.
var ErrNotFound = storage.NotFound

// Error statuses.
const (
	StatusOpen     = storage.StatusOpen
	StatusResolved = storage.StatusResolved
)

// ListOptions narrows List. Zero values mean "no filter". Limit 0 returns everything.
type ListOptions struct {
	Env            string
	Status         string
	Code           string
	Since          time.Time // only errors last seen at or after this time
	Limit          int
	Offset         int
	IncludeDeleted bool
}

// Handler creates, reads, lists, resolves and deletes error records.
type Handler struct {
	st       *postgres.Storage
	svc      *errorHandler.ErrorCoreSvc
	env      string
	service  string
	notifier []Notifier
}

// New connects to PostgreSQL using a lib/pq connection string.
func New(dsn string, opts ...Option) (*Handler, error) {
	st, err := postgres.NewStorage(dsn)
	if err != nil {
		return nil, err
	}
	return newHandler(st, opts), nil
}

// NewWithDB uses an existing connection, which must be a PostgreSQL database.
// Close will close db as well.
func NewWithDB(db *sql.DB, opts ...Option) *Handler {
	return newHandler(postgres.NewStorageFromDB(db), opts)
}

func newHandler(st *postgres.Storage, opts []Option) *Handler {
	h := &Handler{st: st, svc: errorHandler.New(st), env: "production", service: "app"}
	for _, o := range opts {
		o(h)
	}
	return h
}

// Close closes the database connection.
func (h *Handler) Close() error {
	return h.st.Close()
}

// Migrate creates or upgrades the errorhandler table. It is safe to call on
// every start.
func (h *Handler) Migrate(ctx context.Context) error {
	stmts, err := upSQL()
	if err != nil {
		return err
	}
	for _, s := range stmts {
		if err := h.st.Exec(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

// Create stores e as a new record and returns its id. Use Capture to record
// a Go error with de-duplication and notifications.
func (h *Handler) Create(ctx context.Context, e Error) (string, error) {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	return h.svc.CreateErrorHandler(ctx, e)
}

// Get returns the record with the given id, or ErrNotFound.
func (h *Handler) Get(ctx context.Context, id string) (*Error, error) {
	return h.svc.GetErrorHandler(ctx, id)
}

// List returns records, newest first. Deleted records are skipped unless
// IncludeDeleted is set.
func (h *Handler) List(ctx context.Context, opts ...ListOptions) ([]Error, error) {
	var o ListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	return h.st.ListErrors(ctx, storage.ListFilter{
		EnvType:        o.Env,
		Status:         o.Status,
		ErrorCode:      o.Code,
		Since:          o.Since,
		Limit:          o.Limit,
		Offset:         o.Offset,
		IncludeDeleted: o.IncludeDeleted,
	})
}

// Resolve marks the record as resolved, or returns ErrNotFound. If the same
// error happens again it is recorded as a new open record.
func (h *Handler) Resolve(ctx context.Context, id string) error {
	return h.st.ResolveError(ctx, id)
}

// Delete soft-deletes the record, recording the environment that deleted it.
// It returns ErrNotFound if the id does not exist.
func (h *Handler) Delete(ctx context.Context, id, env string) error {
	return h.svc.DeleteErrorHandler(ctx, id, env)
}
