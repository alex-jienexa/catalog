package usecase

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/domain/repository"
	"context"
)

type CustomerUseCase struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerUseCase(customerRepo repository.CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{customerRepo: customerRepo}
}

func (uc *CustomerUseCase) GetOrCreateCustomer(ctx context.Context, input entity.CustomerCreate) (*entity.Customer, bool, error) {
	// Проверяем существование клиента по телефону
	existing, err := uc.customerRepo.GetByPhone(ctx, input.Phone)
	if err != nil {
		return nil, false, err
	}
	if existing != nil {
		return existing, true, nil // существующий клиент
	}
	// Создаём нового
	created, err := uc.customerRepo.Create(ctx, &input)
	return created, false, err
}
