package category

import (
	"context"
	"errorHandler/errorHandler/storage"
	tpc "errorHandler/gunk/v1/category"
)

type categoryCoreStore interface {
	Create_ser(context.Context, storage.Category) (int64, error)
	Get_AllData_ser(context.Context) ([]storage.Category, error)
	Get_single_ser(context.Context, int64) (storage.Category, error)
	Delete(context.Context, int64) error
	Update(context.Context, storage.Category) error
}

type Svc struct {
	tpc.UnimplementedCategoryServiceServer
	core categoryCoreStore
}

func NewCategoryServer(c categoryCoreStore) *Svc {
	return &Svc{
		core: c,
	}
}
