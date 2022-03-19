package blog

import (
	"context"
	"errorHandler/errorHandler/storage"
	tpb "errorHandler/gunk/v1/blog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CreateBlog(ctx context.Context, req *tpb.CreateBlogRequest) (*tpb.CreateBlogResponse, error) {
	blog := storage.Blog{
		CatID:       req.Blog.CatID,
		Title:       req.Blog.Title,
		Description: req.Blog.Description,
		Image:       req.Blog.Image,
	}

	id, err := s.core.Create(context.Background(), blog)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create Blog: %s", err)
	}
	return &tpb.CreateBlogResponse{
		ID: id,
	}, nil

}
