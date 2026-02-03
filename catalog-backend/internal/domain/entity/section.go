package entity

import "time"

type Section struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  *string   `json:"description,omitempty" db:"description"`
	ProductCount int       `json:"product_count" db:"product_count"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// SectionCreate представляет данные для создания раздела
type SectionCreate struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty"`
}

// SectionUpdate представляет данные для обновления раздела
type SectionUpdate struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
