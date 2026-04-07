package repository

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/domain/repository"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type sqliteCustomerRepository struct {
	db *sqlx.DB
}

func NewSQLiteCustomerRepository(db *sqlx.DB) repository.CustomerRepository {
	return &sqliteCustomerRepository{db: db}
}

func (r *sqliteCustomerRepository) Create(ctx context.Context, customer *entity.CustomerCreate) (*entity.Customer, error) {
	query := `
		INSERT INTO customers (first_name, last_name, phone, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING *
	`
	var last sql.NullString
	if customer.LastName != nil {
		last = sql.NullString{String: *customer.LastName, Valid: true}
	}
	var created entity.Customer
	err := r.db.QueryRowxContext(ctx, query, customer.FirstName, last, customer.Phone).StructScan(&created)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}
	return &created, nil
}

func (r *sqliteCustomerRepository) GetByPhone(ctx context.Context, phone string) (*entity.Customer, error) {
	query := `SELECT * FROM customers WHERE phone = ?`
	var customer entity.Customer
	err := r.db.GetContext(ctx, &customer, query, phone)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	return &customer, nil
}

func (r *sqliteCustomerRepository) GetByID(ctx context.Context, id int) (*entity.Customer, error) {
	query := `SELECT * FROM customers WHERE id = ?`
	var customer entity.Customer
	err := r.db.GetContext(ctx, &customer, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	return &customer, nil
}
