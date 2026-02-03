package repository

import (
	"catalog-backend/internal/domain/entity"
	"context"
)

type SectionRepository interface {
	Create(ctx context.Context, section *entity.SectionCreate) (*entity.Section, error)
	GetByID(ctx context.Context, id int) (*entity.Section, error)
	GetAll(ctx context.Context, activeOnly bool) ([]entity.Section, error)
	Update(ctx context.Context, id int, section *entity.SectionUpdate) (*entity.Section, error)
	Delete(ctx context.Context, id int) error
	UpdateProductCount(ctx context.Context, id int, delta int) error
}
