package repository

import (
	"catalog-backend/internal/domain/entity"
)

type StoreRepository interface {
	Get() (*entity.StoreInfo, error)
	Update(info *entity.StoreInfo) error
	EnsureDefault(defaultInfo *entity.StoreInfo) error
}
