package entity

import "time"

type Contact struct {
	ID        int       `json:"id" db:"id"`
	Platform  string    `json:"platform" db:"platform"`
	URL       string    `json:"url" db:"url"`
	Icon      *string   `json:"icon,omitempty" db:"icon"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	Order     int       `json:"order" db:"order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ContactCreate представляет данные для создания контакта
type ContactCreate struct {
	Platform string  `json:"platform" validate:"required,min=1,max=50"`
	URL      string  `json:"url" validate:"required,url"`
	Icon     *string `json:"icon,omitempty"`
	Order    int     `json:"order" validate:"min=0"`
}

// ContactUpdate представляет данные для обновления контакта
type ContactUpdate struct {
	Platform *string `json:"platform,omitempty" validate:"omitempty,min=1,max=50"`
	URL      *string `json:"url,omitempty" validate:"omitempty,url"`
	Icon     *string `json:"icon,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	Order    *int    `json:"order,omitempty" validate:"omitempty,min=0"`
}
