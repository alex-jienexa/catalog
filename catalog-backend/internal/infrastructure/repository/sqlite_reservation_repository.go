package repository

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/domain/repository"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type sqliteReservationRepository struct {
	db *sqlx.DB
}

func NewSQLiteReservationRepository(db *sqlx.DB) repository.ReservationRepository {
	return &sqliteReservationRepository{db: db}
}

func (r *sqliteReservationRepository) Create(ctx context.Context, reservation *entity.ReservationCreate) (*entity.Reservation, error) {
	query := `
		INSERT INTO reservations (customer_id, product_id, status, created_at, updated_at)
		VALUES (?, ?, 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING *
	`
	var created entity.Reservation
	err := r.db.QueryRowxContext(ctx, query, reservation.CustomerID, reservation.ProductID).StructScan(&created)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}
	return &created, nil
}

func (r *sqliteReservationRepository) GetAll(ctx context.Context, limit, offset int) ([]entity.Reservation, error) {
	query := `SELECT * FROM reservations ORDER BY created_at DESC LIMIT ? OFFSET ?`
	var reservations []entity.Reservation
	err := r.db.SelectContext(ctx, &reservations, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}
	return reservations, nil
}

func (r *sqliteReservationRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM reservations")
	return count, err
}

func (r *sqliteReservationRepository) GetWithDetails(ctx context.Context, limit, offset int) ([]repository.ReservationWithDetails, error) {
	query := `
		SELECT 
			r.id, r.customer_id, r.product_id, r.status, r.created_at, r.updated_at,
			c.first_name as customer_first_name, c.last_name as customer_last_name, c.phone as customer_phone,
			p.name as product_name, p.price as product_price
		FROM reservations r
		JOIN customers c ON r.customer_id = c.id
		JOIN products p ON r.product_id = p.id
		ORDER BY r.created_at DESC
		LIMIT ? OFFSET ?
	`
	var results []repository.ReservationWithDetails
	err := r.db.SelectContext(ctx, &results, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations with details: %w", err)
	}
	return results, nil
}
