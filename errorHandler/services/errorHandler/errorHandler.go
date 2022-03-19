package errorHandler

import (
	"context"

	"errorHandler/errorHandler/storage"

	erHan "errorHandler/gunk/v1/errorHandler"
)

type ErrorSVC interface {
	CreateErrorHandler(context.Context, storage.ErrorHandler) (string, error)
	GetErrorHandler(context.Context, string) (*storage.ErrorHandler, error)
	DeleteErrorHandler(context.Context, string, string) error
	ListErrorHandler(context.Context) ([]storage.ErrorHandler, error)
}

type Svc struct {
	erHan.UnimplementedErrorHandlerServiceServer
	erHan ErrorSVC
}

func New(erHan ErrorSVC) *Svc {
	return &Svc{
		erHan: erHan,
	}
}
