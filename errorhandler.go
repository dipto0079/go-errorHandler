// Package errorhandler stores application errors in PostgreSQL.
//
//	h, err := errorhandler.New("user=postgres password=... host=localhost port=5432 dbname=mydb sslmode=disable")
//	id, err := h.Create(ctx, errorhandler.Error{ErrorCode: "E100", ErrorDetails: "boom", EnvType: "production", CreatedBy: "api"})
//
// The "errorhandler" table is created by the SQL in migrations/sql.
package errorhandler

import (
	"context"

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
	svc *errorHandler.ErrorCoreSvc
}

// New connects to PostgreSQL using a lib/pq connection string.
func New(dsn string) (*Handler, error) {
	st, err := postgres.NewStorage(dsn)
	if err != nil {
		return nil, err
	}
	return &Handler{svc: errorHandler.New(st)}, nil
}

// Create stores e and returns its id.
func (h *Handler) Create(ctx context.Context, e Error) (string, error) {
	return h.svc.CreateErrorHandler(ctx, e)
}

// Get returns the record with the given id.
func (h *Handler) Get(ctx context.Context, id string) (*Error, error) {
	return h.svc.GetErrorHandler(ctx, id)
}

// List returns all records.
func (h *Handler) List(ctx context.Context) ([]Error, error) {
	return h.svc.ListErrorHandler(ctx)
}

// Delete soft-deletes the record, recording the environment that deleted it.
func (h *Handler) Delete(ctx context.Context, id, env string) error {
	return h.svc.DeleteErrorHandler(ctx, id, env)
}
