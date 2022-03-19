package user

import (
	"context"
	"errorHandler/errorHandler/storage"
	"errorHandler/errorHandler/storage/postgres"
)

type CoreSve struct {
	store *postgres.Storage
}

func NewCoreSve(b *postgres.Storage) *CoreSve {
	return &CoreSve{
		store: b,
	}
}
func (cs CoreSve) Create(ctx context.Context, t storage.User) (int64, error) {
	return cs.store.CreateUser(ctx, t)
}
