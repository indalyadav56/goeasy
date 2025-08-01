package repository

import (
	"context"
	"github.com/test/authtest-fixed/internal/product/domain/entity"
)

type ProductRepository interface {
	Insert(ctx context.Context, entity *entity.Product) error
	FindByID(ctx context.Context, id string) (*entity.Product, error)
	Update(ctx context.Context, entity *entity.Product) error
	Delete(ctx context.Context, id string) error
}
