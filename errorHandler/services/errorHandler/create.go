package errorHandler

import (
	"context"

	"errorHandler/errorHandler/storage"
	erHan "errorHandler/gunk/v1/errorHandler"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CreateErrorHandler(ctx context.Context, req *erHan.CreateErrorHandlerRequest) (*erHan.CreateErrorHandlerResponse, error) {
	
	dbPrm := storage.ErrorHandler{
		ErrorCode:    req.ErrorHandler.ErrorCode,
		ErrorDetails: req.ErrorHandler.ErrorDetails,
		EnvType:      req.ErrorHandler.EnvType,
		CreatedBy:    req.ErrorHandler.CreatedBy,
		CreatedAt:    req.ErrorHandler.CreatedAt.AsTime(),
	}
	res, err := s.erHan.CreateErrorHandler(ctx, dbPrm)
	if err != nil {
		
		return nil, status.Error(codes.Internal, "failed to create CreateErrorHandler")
	}

	return &erHan.CreateErrorHandlerResponse{
		ID: res,
	}, nil

}
