// Package errorhandler stores application errors in PostgreSQL.
//
//	h, err := errorhandler.New("user=postgres password=... host=localhost port=5432 dbname=mydb sslmode=disable")
//	id, err := h.Create(ctx, errorhandler.Error{ErrorCode: "E100", ErrorDetails: "boom", EnvType: "production", CreatedBy: "api"})
//
// Call Migrate once to create the "errorhandler" table.
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

// Handler creates, reads, lists and deletes error records.
type Handler struct {
	st  *postgres.Storage
	svc *errorHandler.ErrorCoreSvc
}

// New connects to PostgreSQL using a lib/pq connection string.
func New(dsn string) (*Handler, error) {
	st, err := postgres.NewStorage(dsn)
	if err != nil {
		return nil, err
	}
	return &Handler{st: st, svc: errorHandler.New(st)}, nil
}

// NewWithDB uses an existing connection, which must be a PostgreSQL database.
// Close will close db as well.
func NewWithDB(db *sql.DB) *Handler {
	st := postgres.NewStorageFromDB(db)
	return &Handler{st: st, svc: errorHandler.New(st)}
}

// Close closes the database connection.
func (h *Handler) Close() error {
	return h.st.Close()
}

// Migrate creates the errorhandler table if it does not exist. It is safe to
// call on every start.
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

// Create stores e and returns its id.
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

// List returns all records.
func (h *Handler) List(ctx context.Context) ([]Error, error) {
	return h.svc.ListErrorHandler(ctx)
}

// Delete soft-deletes the record, recording the environment that deleted it.
// It returns ErrNotFound if the id does not exist.
func (h *Handler) Delete(ctx context.Context, id, env string) error {
	return h.svc.DeleteErrorHandler(ctx, id, env)
}
