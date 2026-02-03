package repository

import (
	"catalog-backend/internal/domain/entity"
	"context"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.ProductCreate) (*entity.Product, error)
	GetByID(ctx context.Context, id int) (*entity.Product, error)
	GetAll(ctx context.Context, filters ProductFilters) ([]entity.Product, error)
	Update(ctx context.Context, id int, product *entity.ProductUpdate) (*entity.Product, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context, filters ProductFilters) (int, error)
}

type ProductFilters struct {
	SectionID *int
	IsActive  *bool
	Limit     int
	Offset    int
}
