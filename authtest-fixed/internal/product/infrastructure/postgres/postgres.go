package postgres

import (
	"context"
	"database/sql"
	
	"github.com/test/authtest-fixed/internal/product/domain/entity"
)

type PostgresRepository interface{}

type postgresRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *postgresRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Insert(ctx context.Context, entity *entity.Product) error {
	return nil
}

func (r *postgresRepository) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	return nil, nil
}

func (r *postgresRepository) Update(ctx context.Context, entity *entity.Product) error {
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	return nil
}