package entity

import "time"

type Admin struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type AdminCreate struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

type AdminUpdate struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
}

// AdminResponse — публичное представление без хеша пароля
type AdminResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Admin) ToResponse() AdminResponse {
	return AdminResponse{
		ID:        a.ID,
		Name:      a.Name,
		Username:  a.Username,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
