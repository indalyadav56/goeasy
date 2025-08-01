package application

import (
	"context"
	"log/slog"

	"github.com/test/authtest-fixed/internal/product/domain/repository"
	"github.com/test/authtest-fixed/internal/product/domain/entity"
)

type ProductService interface {
	Create(item *entity.Product) (*entity.Product, error)
	GetByID(id string) (*entity.Product, error)
	Update(id string, item *entity.Product) (*entity.Product, error)
	Delete(id string) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) Create(item *entity.Product) (*entity.Product, error) {
	slog.Info("creating user")
	err := s.repo.Insert(context.Background(), item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *productService) GetByID(id string) (*entity.Product, error) {
	return s.repo.FindByID(context.Background(), id)
}

func (s *productService) Update(id string, item *entity.Product) (*entity.Product, error) {
	err := s.repo.Update(context.Background(), item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *productService) Delete(id string) error {
	return s.repo.Delete(context.Background(), id)
}
