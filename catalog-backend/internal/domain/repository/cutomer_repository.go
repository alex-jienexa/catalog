package repository

import (
	"catalog-backend/internal/domain/entity"
	"context"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.CustomerCreate) (*entity.Customer, error)
	GetByPhone(ctx context.Context, phone string) (*entity.Customer, error)
	GetByID(ctx context.Context, id int) (*entity.Customer, error)
}
