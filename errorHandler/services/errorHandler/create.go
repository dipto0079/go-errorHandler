package errorHandler

import (
	"context"
	"time"

	"github.com/dipto0079/go-errorHandler/errorHandler/storage"
	erHan "github.com/dipto0079/go-errorHandler/gunk/v1/errorHandler"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CreateErrorHandler(ctx context.Context, req *erHan.CreateErrorHandlerRequest) (*erHan.CreateErrorHandlerResponse, error) {
	
	if req.ErrorHandler == nil {
		return nil, status.Error(codes.InvalidArgument, "errorHandler is required")
	}
	createdAt := time.Now()
	if req.ErrorHandler.CreatedAt != nil {
		createdAt = req.ErrorHandler.CreatedAt.AsTime()
	}
	dbPrm := storage.ErrorHandler{
		ErrorCode:    req.ErrorHandler.ErrorCode,
		ErrorDetails: req.ErrorHandler.ErrorDetails,
		EnvType:      req.ErrorHandler.EnvType,
		CreatedBy:    req.ErrorHandler.CreatedBy,
		CreatedAt:    createdAt,
	}
	res, err := s.erHan.CreateErrorHandler(ctx, dbPrm)
	if err != nil {
		
		return nil, status.Error(codes.Internal, "failed to create CreateErrorHandler")
	}

	return &erHan.CreateErrorHandlerResponse{
		ID: res,
	}, nil

}
