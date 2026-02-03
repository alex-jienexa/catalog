package entity

import "time"

type Product struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Price       float64   `json:"price" db:"price"`
	SectionID   int       `json:"section_id" db:"section_id"`
	Description *string   `json:"description,omitempty" db:"description"`
	ImageURL    *string   `json:"image_url,omitempty" db:"image_url"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ProductCreate представляет данные для создания продукта
type ProductCreate struct {
	Name        string  `json:"name" validate:"required,min=1,max=200"`
	Price       float64 `json:"price" validate:"required,min=0"`
	SectionID   int     `json:"section_id" validate:"required,min=1"`
	Description *string `json:"description,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
}

// ProductUpdate представляет данные для обновления продукта
type ProductUpdate struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	SectionID   *int     `json:"section_id,omitempty" validate:"omitempty,min=1"`
	Description *string  `json:"description,omitempty"`
	ImageURL    *string  `json:"image_url,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}
