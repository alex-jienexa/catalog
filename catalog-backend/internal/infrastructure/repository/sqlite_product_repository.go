package repository

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/domain/repository"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type sqliteProductRepository struct {
	db *sqlx.DB
}

func NewSQLiteProductRepository(db *sqlx.DB) repository.ProductRepository {
	return &sqliteProductRepository{db: db}
}

func (r *sqliteProductRepository) Create(ctx context.Context, product *entity.ProductCreate) (*entity.Product, error) {
	// Вместо именованных параметров используем позиционные
	query := `
        INSERT INTO products (name, price, section_id, description, image_url, is_active, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    `

	result, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Price,
		product.SectionID,
		product.Description,
		product.ImageURL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Получаем созданный продукт
	return r.GetByID(ctx, int(id))
}

func (r *sqliteProductRepository) GetByID(ctx context.Context, id int) (*entity.Product, error) {
	query := `SELECT * FROM products WHERE id = ?`

	var product entity.Product
	err := r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

func (r *sqliteProductRepository) GetAll(ctx context.Context, filters repository.ProductFilters) ([]entity.Product, error) {
	query := `SELECT * FROM products WHERE 1=1`
	args := []interface{}{}

	if filters.SectionID != nil {
		query += ` AND section_id = ?`
		args = append(args, *filters.SectionID)
	}

	if filters.IsActive != nil {
		query += ` AND is_active = ?`
		args = append(args, *filters.IsActive)
	}

	query += ` ORDER BY created_at DESC`

	if filters.Limit > 0 {
		query += ` LIMIT ?`
		args = append(args, filters.Limit)
	}

	if filters.Offset > 0 {
		query += ` OFFSET ?`
		args = append(args, filters.Offset)
	}

	var products []entity.Product
	err := r.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

func (r *sqliteProductRepository) Update(ctx context.Context, id int, product *entity.ProductUpdate) (*entity.Product, error) {
	// Начинаем строить запрос
	query := `UPDATE products SET `
	args := []interface{}{}
	updates := []string{}

	if product.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *product.Name)
	}

	if product.Price != nil {
		updates = append(updates, "price = ?")
		args = append(args, *product.Price)
	}

	if product.SectionID != nil {
		updates = append(updates, "section_id = ?")
		args = append(args, *product.SectionID)
	}

	if product.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *product.Description)
	}

	if product.ImageURL != nil {
		updates = append(updates, "image_url = ?")
		args = append(args, *product.ImageURL)
	}

	if product.IsActive != nil {
		updates = append(updates, "is_active = ?")
		args = append(args, *product.IsActive)
	}

	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}

	updates = append(updates, "updated_at = ?")
	args = append(args, time.Now())

	query += strings.Join(updates, ", ")
	query += ` WHERE id = ? RETURNING *`
	args = append(args, id)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	defer rows.Close()

	var updatedProduct entity.Product
	if rows.Next() {
		if err := rows.StructScan(&updatedProduct); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
	}

	return &updatedProduct, nil
}

func (r *sqliteProductRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *sqliteProductRepository) Count(ctx context.Context, filters repository.ProductFilters) (int, error) {
	query := `SELECT COUNT(*) FROM products WHERE 1=1`
	args := []interface{}{}

	if filters.SectionID != nil {
		query += ` AND section_id = ?`
		args = append(args, *filters.SectionID)
	}

	if filters.IsActive != nil {
		query += ` AND is_active = ?`
		args = append(args, *filters.IsActive)
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}
