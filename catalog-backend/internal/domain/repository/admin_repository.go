package repository

import (
	"catalog-backend/internal/domain/entity"
	"context"
)

type AdminRepository interface {
	Create(ctx context.Context, admin *entity.AdminCreate, passwordHash string) (*entity.Admin, error)
	GetByID(ctx context.Context, id int) (*entity.Admin, error)
	GetByUsername(ctx context.Context, username string) (*entity.Admin, error)
	GetAll(ctx context.Context) ([]entity.Admin, error)
	Update(ctx context.Context, id int, input *entity.AdminUpdate, passwordHash *string) (*entity.Admin, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}
