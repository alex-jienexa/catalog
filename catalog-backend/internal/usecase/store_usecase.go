package usecase

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/domain/repository"
	"context"
	"errors"
)

type StoreUseCase struct {
	repo repository.StoreRepository
}

func NewStoreUseCase(repo repository.StoreRepository) *StoreUseCase {
	return &StoreUseCase{repo: repo}
}

func (uc *StoreUseCase) GetStoreInfo(ctx context.Context) (*entity.StoreInfo, error) {
	return uc.repo.Get()
}

func (uc *StoreUseCase) UpdateStoreInfo(ctx context.Context, info *entity.StoreInfo) (*entity.StoreInfo, error) {
	if info.Title == "" || info.Description == "" {
		return nil, errors.New("title and description are required")
	}
	if err := uc.repo.Update(info); err != nil {
		return nil, err
	}
	return info, nil
}
