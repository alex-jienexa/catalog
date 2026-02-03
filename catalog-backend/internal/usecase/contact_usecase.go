package usecase

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/domain/repository"
	"context"
)

type ContactUseCase struct {
	contactRepo repository.ContactRepository
}

func NewContactUseCase(contactRepo repository.ContactRepository) *ContactUseCase {
	return &ContactUseCase{
		contactRepo: contactRepo,
	}
}

func (uc *ContactUseCase) CreateContact(ctx context.Context, input entity.ContactCreate) (*entity.Contact, error) {
	return uc.contactRepo.Create(ctx, &input)
}

func (uc *ContactUseCase) GetContact(ctx context.Context, id int) (*entity.Contact, error) {
	return uc.contactRepo.GetByID(ctx, id)
}

func (uc *ContactUseCase) GetContacts(ctx context.Context, activeOnly bool) ([]entity.Contact, error) {
	return uc.contactRepo.GetAll(ctx, activeOnly)
}

func (uc *ContactUseCase) UpdateContact(ctx context.Context, id int, input entity.ContactUpdate) (*entity.Contact, error) {
	return uc.contactRepo.Update(ctx, id, &input)
}

func (uc *ContactUseCase) DeleteContact(ctx context.Context, id int) error {
	return uc.contactRepo.Delete(ctx, id)
}
