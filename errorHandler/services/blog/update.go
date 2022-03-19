package blog

import (
	"context"
	"errorHandler/errorHandler/storage"
	tpb "errorHandler/gunk/v1/blog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdateBlog(ctx context.Context, req *tpb.UpdateBlogRequest) (*tpb.UpdateBlogResponse, error) {
	blog := storage.Blog{
		ID:          req.GetBlog().ID,
		CatID:       req.GetBlog().CatID,
		Title:       req.GetBlog().Title,
		Description: req.GetBlog().Description,
		Image:       req.GetBlog().Image,
	}

	err := s.core.UpdateBlog(context.Background(), blog)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to Update category.")
	}
	return &tpb.UpdateBlogResponse{}, nil
}
