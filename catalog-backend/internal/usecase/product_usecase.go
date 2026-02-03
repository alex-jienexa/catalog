package usecase

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/domain/repository"
	"context"
	"errors"
)

type ProductUseCase struct {
	productRepo repository.ProductRepository
	sectionRepo repository.SectionRepository
}

func NewProductUseCase(productRepo repository.ProductRepository, sectionRepo repository.SectionRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
		sectionRepo: sectionRepo,
	}
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, input entity.ProductCreate) (*entity.Product, error) {
	// Проверяем существование раздела
	section, err := uc.sectionRepo.GetByID(ctx, input.SectionID)
	if err != nil {
		return nil, err
	}
	if section == nil {
		return nil, ErrSectionNotFound
	}

	// Создаем продукт
	product, err := uc.productRepo.Create(ctx, &input)
	if err != nil {
		return nil, err
	}

	// Обновляем счетчик товаров в разделе
	err = uc.sectionRepo.UpdateProductCount(ctx, input.SectionID, 1)
	if err != nil {
		// В случае ошибки откатываем создание продукта
		_ = uc.productRepo.Delete(ctx, product.ID)
		return nil, err
	}

	return product, nil
}

func (uc *ProductUseCase) GetProduct(ctx context.Context, id int) (*entity.Product, error) {
	return uc.productRepo.GetByID(ctx, id)
}

func (uc *ProductUseCase) GetProducts(ctx context.Context, filters repository.ProductFilters) ([]entity.Product, int, error) {
	products, err := uc.productRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.productRepo.Count(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, id int, input entity.ProductUpdate) (*entity.Product, error) {
	// Получаем текущий продукт
	currentProduct, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if currentProduct == nil {
		return nil, ErrProductNotFound
	}

	// Если меняется раздел, обновляем счетчики
	if input.SectionID != nil && *input.SectionID != currentProduct.SectionID {
		// Проверяем новый раздел
		newSection, err := uc.sectionRepo.GetByID(ctx, *input.SectionID)
		if err != nil {
			return nil, err
		}
		if newSection == nil {
			return nil, ErrSectionNotFound
		}

		// Обновляем счетчики
		err = uc.sectionRepo.UpdateProductCount(ctx, currentProduct.SectionID, -1)
		if err != nil {
			return nil, err
		}
		err = uc.sectionRepo.UpdateProductCount(ctx, *input.SectionID, 1)
		if err != nil {
			// Откатываем изменения
			_ = uc.sectionRepo.UpdateProductCount(ctx, currentProduct.SectionID, 1)
			return nil, err
		}
	}

	// Обновляем продукт
	return uc.productRepo.Update(ctx, id, &input)
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id int) error {
	// Получаем продукт для получения section_id
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrProductNotFound
	}

	// Удаляем продукт
	err = uc.productRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Обновляем счетчик товаров в разделе
	return uc.sectionRepo.UpdateProductCount(ctx, product.SectionID, -1)
}

var (
	ErrProductNotFound = errors.New("product not found")
	ErrSectionNotFound = errors.New("section not found")
)
