package entity

import "time"

type Reservation struct {
	ID         int       `json:"id" db:"id"`
	CustomerID int       `json:"customer_id" db:"customer_id"`
	ProductID  int       `json:"product_id" db:"product_id"`
	Status     string    `json:"status" db:"status"` // pending, contacted, completed
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type ReservationCreate struct {
	CustomerID int `json:"customer_id" validate:"required"`
	ProductID  int `json:"product_id" validate:"required"`
}
