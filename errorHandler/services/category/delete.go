package category

import (
	"context"
	tpc "errorHandler/gunk/v1/category"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) Delete(ctx context.Context, req *tpc.DeleteCategoryRequest) (*tpc.DeleteCategoryResponse, error) {

	err := s.core.Delete(context.Background(), req.GetID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create category: %s", err)
	}
	return &tpc.DeleteCategoryResponse{}, nil

}
