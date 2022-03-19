package errorHandler

import (
	"context"

	"errorHandler/errorHandler/storage"
	"errorHandler/errorHandler/storage/postgres"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

		return "", status.Error(codes.Internal, "processing failed")
	}

	return id, nil
}

func (cs ErrorCoreSvc) GetErrorHandler(ctx context.Context, id string) (*storage.ErrorHandler, error) {

	errorhan, err := cs.st.GetError(ctx, id)
	if err != nil && err != storage.NotFound {

		return nil, status.Error(codes.Internal, "processing failed")
	}
	return errorhan, nil
}

func (cs ErrorCoreSvc) DeleteErrorHandler(ctx context.Context, id string, env string) error {

	if err := cs.st.DeleteError(ctx, id, env); err != nil {

		return status.Error(codes.Internal, "processing failed")
	}

	return nil
}

func (cs ErrorCoreSvc) ListErrorHandler(ctx context.Context) ([]storage.ErrorHandler, error) {

	res, err := cs.st.ListError(ctx)
	if err != nil {

		return res, status.Error(codes.Internal, "processing failed")
	}

	return res, nil
}
