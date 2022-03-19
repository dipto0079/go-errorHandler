package errorHandler

import (
	"context"

	erHan "errorHandler/gunk/v1/errorHandler"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) ListErrorHandler(ctx context.Context, req *erHan.ListErrorHandlerRequest) (*erHan.ListErrorHandlerResponse, error) {

	ids, err := s.erHan.ListErrorHandler(context.Background())
	if err != nil {

		return nil, status.Error(codes.NotFound, "ListErrorHandler doesn't exists")
	}
	var errhan []*erHan.ErrorHandler
	for _, v := range ids {
		errhan = append(errhan, &erHan.ErrorHandler{
			ID:              v.ID,
			ErrorCode:       v.ErrorCode,
			ErrorDetails:    v.ErrorDetails,
			EnvType:         v.EnvType,
			CreatedBy:       v.CreatedBy,
			DeleteByEnvType: v.DeleteByEnvType,
		})
	}
	return &erHan.ListErrorHandlerResponse{
		ErrorHandler: errhan,
	}, nil

}
