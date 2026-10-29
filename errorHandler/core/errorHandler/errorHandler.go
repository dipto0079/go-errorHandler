package errorHandler

import (
	"context"

	"github.com/dipto0079/go-errorHandler/errorHandler/storage"
	"github.com/dipto0079/go-errorHandler/errorHandler/storage/postgres"
)

type ErrorCoreSvc struct {
	st *postgres.Storage
}

func New(a *postgres.Storage) *ErrorCoreSvc {
	return &ErrorCoreSvc{
		st: a,
	}
}

func (cs ErrorCoreSvc) CreateErrorHandler(ctx context.Context, errorHandler storage.ErrorHandler) (string, error) {

	id, err := cs.st.CreateError(ctx, errorHandler)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (cs ErrorCoreSvc) GetErrorHandler(ctx context.Context, id string) (*storage.ErrorHandler, error) {

	errorhan, err := cs.st.GetError(ctx, id)
	if err != nil {
		return nil, err
	}
	return errorhan, nil
}

func (cs ErrorCoreSvc) DeleteErrorHandler(ctx context.Context, id string, env string) error {

	return cs.st.DeleteError(ctx, id, env)
}

func (cs ErrorCoreSvc) ListErrorHandler(ctx context.Context) ([]storage.ErrorHandler, error) {

	res, err := cs.st.ListError(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
