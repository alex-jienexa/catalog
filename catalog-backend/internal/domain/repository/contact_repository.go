package repository

import (
	"catalog-backend/internal/domain/entity"
	"context"
)

type ContactRepository interface {
	Create(ctx context.Context, contact *entity.ContactCreate) (*entity.Contact, error)
	GetByID(ctx context.Context, id int) (*entity.Contact, error)
	GetAll(ctx context.Context, activeOnly bool) ([]entity.Contact, error)
	Update(ctx context.Context, id int, contact *entity.ContactUpdate) (*entity.Contact, error)
	Delete(ctx context.Context, id int) error
}
