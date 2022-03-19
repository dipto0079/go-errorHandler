package blog

import (
	"context"
	"errorHandler/errorHandler/storage"
	tpb "errorHandler/gunk/v1/blog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) GetBlog(ctx context.Context, req *tpb.GetBlogRequest) (*tpb.GetBlogResponse, error) {

	var blo storage.Blog

	blo, err := s.core.GetBlog(context.Background(), req.GetID())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get Blog.")
	}

	return &tpb.GetBlogResponse{
		Blog: &tpb.Blog{
			ID:          blo.ID,
			CatID:       blo.CatID,
			Title:       blo.Title,
			Description: blo.Description,
			Image:       blo.Image,
		},
	}, nil

}
