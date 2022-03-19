package category

import (
	"context"
	tpc "errorHandler/gunk/v1/category"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) GetAllData(ctx context.Context, req *tpc.GetAllDataCategoryRequest) (*tpc.GetAllDataCategoryResponse, error) {
	ids, err := s.core.Get_AllData_ser(context.Background())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create category: %s", err)
	}
	var ctl []*tpc.Category
	for _, v := range ids {
		ctl = append(ctl, &tpc.Category{
			ID:         v.ID,
			Title:      v.Title,
			IsComplete: v.IsComplete,
		})
	}
	return &tpc.GetAllDataCategoryResponse{
		Category: ctl,
	}, nil
}
