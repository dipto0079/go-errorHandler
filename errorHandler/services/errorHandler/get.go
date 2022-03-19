package errorHandler

import (
	"context"

	erHan "errorHandler/gunk/v1/errorHandler"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	acts "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetErrorHandler(ctx context.Context, req *erHan.GetErrorHandlerRequest) (*erHan.GetErrorHandlerResponse, error) {

	res, err := s.erHan.GetErrorHandler(ctx, req.ID)
	if err != nil {

		return nil, status.Error(codes.NotFound, "GetErrorHandler doesn't exists")
	}

	getAnn := &erHan.GetErrorHandlerResponse{
		ErrorHandler: &erHan.ErrorHandler{
			ID:           res.ID,
			ErrorCode:    res.ErrorCode,
			ErrorDetails: res.ErrorDetails,
			EnvType:      res.EnvType,
			CreatedAt:    acts.New(res.CreatedAt),
			CreatedBy:    res.CreatedBy,
		},
	}
	return getAnn, nil

}
