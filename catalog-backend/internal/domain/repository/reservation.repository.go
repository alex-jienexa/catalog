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
	UpdateStatus(ctx context.Context, id int, status string) (*entity.Reservation, error) // новый
}

type ReservationWithDetails struct {
	entity.Reservation
	CustomerFirstName string  `json:"customer_first_name" db:"customer_first_name"`
	CustomerLastName  *string `json:"customer_last_name" db:"customer_last_name"`
	CustomerPhone     string  `json:"customer_phone" db:"customer_phone"`
	ProductName       string  `json:"product_name" db:"product_name"`
	ProductPrice      float64 `json:"product_price" db:"product_price"`
}
