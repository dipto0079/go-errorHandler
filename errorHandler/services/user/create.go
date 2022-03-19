package user

import (
	"context"
	"errorHandler/errorHandler/storage"
	tpu "errorHandler/gunk/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) Create(ctx context.Context, req *tpu.CreateUserRequest) (*tpu.CreateUserResponse, error) {

	user := storage.User{
		Name:          req.User.Name,
		Email:         req.User.Email,
		Password:      req.User.Password,
		EmailVerified: req.User.EmailVerified,
	}

	id, err := s.core.Create(context.Background(), user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create category: %s", err)
	}
	return &tpu.CreateUserResponse{
		ID: id,
	}, nil

}
