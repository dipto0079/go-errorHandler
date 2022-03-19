package user

import (
	"context"
	"errorHandler/errorHandler/storage"
	tpu "errorHandler/gunk/v1/user"
)

type UserCoreStore interface {
	Create(context.Context, storage.User) (int64, error)
}

type Svc struct {
	tpu.UnimplementedUserRegServiceServer
	core UserCoreStore
}

func NewUserServer(c UserCoreStore) *Svc {
	return &Svc{
		core: c,
	}
}
