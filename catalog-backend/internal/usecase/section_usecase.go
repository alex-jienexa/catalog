package usecase

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/domain/repository"
	"context"
	"errors"
)

type SectionUseCase struct {
	sectionRepo repository.SectionRepository
	productRepo repository.ProductRepository
}

func NewSectionUseCase(sectionRepo repository.SectionRepository, productRepo repository.ProductRepository) *SectionUseCase {
	return &SectionUseCase{
		sectionRepo: sectionRepo,
		productRepo: productRepo,
	}
}

func (uc *SectionUseCase) CreateSection(ctx context.Context, input entity.SectionCreate) (*entity.Section, error) {
	return uc.sectionRepo.Create(ctx, &input)
}

func (uc *SectionUseCase) GetSection(ctx context.Context, id int) (*entity.Section, error) {
	return uc.sectionRepo.GetByID(ctx, id)
}

func (uc *SectionUseCase) GetSections(ctx context.Context, activeOnly bool) ([]entity.Section, error) {
	return uc.sectionRepo.GetAll(ctx, activeOnly)
}

func (uc *SectionUseCase) UpdateSection(ctx context.Context, id int, input entity.SectionUpdate) (*entity.Section, error) {
	return uc.sectionRepo.Update(ctx, id, &input)
}

func (uc *SectionUseCase) DeleteSection(ctx context.Context, id int) error {
	// Проверяем, есть ли товары в разделе
	filters := repository.ProductFilters{
		SectionID: &id,
		IsActive:  boolPtr(true),
		Limit:     1,
		Offset:    0,
	}

	products, err := uc.productRepo.GetAll(ctx, filters)
	if err != nil {
		return err
	}

	if len(products) > 0 {
		return ErrSectionHasProducts
	}

	return uc.sectionRepo.Delete(ctx, id)
}

func boolPtr(b bool) *bool {
	return &b
}

var ErrSectionHasProducts = errors.New("cannot delete section with products")
