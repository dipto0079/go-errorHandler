package category

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

func (cs CoreSve) Create_ser(ctx context.Context, t storage.Category) (int64, error) {
	return cs.store.Create_sto(ctx, t)
}

func (cs CoreSve) Get_AllData_ser(ctx context.Context) ([]storage.Category, error) {
	return cs.store.Get_all_Data(ctx)
}

func (cs CoreSve) Delete(ctx context.Context, id int64) error {
	return cs.store.Delete(ctx, id)
}

func (cs CoreSve) Get_single_ser(ctx context.Context, id int64) (storage.Category, error) {
	return cs.store.Get_sto(ctx, id)
}

func (cs CoreSve) Update(ctx context.Context, c storage.Category) error {
	return cs.store.Update(ctx, c)
}
