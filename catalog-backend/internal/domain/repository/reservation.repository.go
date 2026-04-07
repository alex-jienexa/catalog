package repository

import (
	"catalog-backend/internal/domain/entity"
	"context"
)

type ReservationRepository interface {
	Create(ctx context.Context, reservation *entity.ReservationCreate) (*entity.Reservation, error)
	GetAll(ctx context.Context, limit, offset int) ([]entity.Reservation, error)
	Count(ctx context.Context) (int, error)
	GetWithDetails(ctx context.Context, limit, offset int) ([]ReservationWithDetails, error)
}

type ReservationWithDetails struct {
	ReservationID     int     `json:"reservation_id"`
	CustomerFirstName string  `json:"customer_first_name"`
	CustomerLastName  *string `json:"customer_last_name"`
	CustomerPhone     string  `json:"customer_phone"`
	ProductName       string  `json:"product_name"`
	ProductPrice      float64 `json:"product_price"`
}
