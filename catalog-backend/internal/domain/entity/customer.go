package entity

import "time"

type Customer struct {
	ID        int       `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  *string   `json:"last_name,omitempty" db:"last_name"`
	Phone     string    `json:"phone" db:"phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CustomerCreate struct {
	FirstName string  `json:"first_name" validate:"required,min=1,max=100"`
	LastName  *string `json:"last_name,omitempty"`
	Phone     string  `json:"phone" validate:"required,min=10,max=20"`
}
