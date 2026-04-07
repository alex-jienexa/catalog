package usecase

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/domain/repository"
	"context"
	"errors"
)

type ReservationUseCase struct {
	reservationRepo repository.ReservationRepository
	productRepo     repository.ProductRepository
}

func NewReservationUseCase(reservationRepo repository.ReservationRepository, productRepo repository.ProductRepository) *ReservationUseCase {
	return &ReservationUseCase{
		reservationRepo: reservationRepo,
		productRepo:     productRepo,
	}
}

func (uc *ReservationUseCase) CreateReservation(ctx context.Context, customerID, productID int) (*entity.Reservation, error) {
	// Проверяем существование товара
	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	// Создаём бронирование
	input := &entity.ReservationCreate{
		CustomerID: customerID,
		ProductID:  productID,
	}
	return uc.reservationRepo.Create(ctx, input)
}

func (uc *ReservationUseCase) GetReservations(ctx context.Context, page, limit int) ([]repository.ReservationWithDetails, int, error) {
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit
	reservations, err := uc.reservationRepo.GetWithDetails(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := uc.reservationRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return reservations, total, nil
}
