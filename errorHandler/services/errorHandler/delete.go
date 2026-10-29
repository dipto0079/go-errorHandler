package errorHandler

import (
	"context"
	"errors"

	"github.com/dipto0079/go-errorHandler/errorHandler/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	erHan "github.com/dipto0079/go-errorHandler/gunk/v1/errorHandler"
)

func (s *Svc) DeleteErrorHandler(ctx context.Context, req *erHan.DeleteErrorHandlerRequest) (*erHan.DeleteErrorHandlerResponse, error) {
	
	if err := s.erHan.DeleteErrorHandler(ctx, req.ID, req.DeleteByEnvType); err != nil {
		
		if errors.Is(err, storage.NotFound) {
			return nil, status.Error(codes.NotFound, "DeleteErrorHandler doesn't exists")
		}
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return &erHan.DeleteErrorHandlerResponse{}, nil
}
