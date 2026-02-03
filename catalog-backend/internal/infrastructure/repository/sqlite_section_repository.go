package repository

import (
	"catalog-backend/internal/domain/entity"
	"catalog-backend/internal/domain/repository"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type sqliteSectionRepository struct {
	db *sqlx.DB
}

func NewSQLiteSectionRepository(db *sqlx.DB) repository.SectionRepository {
	return &sqliteSectionRepository{db: db}
}

func (r *sqliteSectionRepository) Create(ctx context.Context, section *entity.SectionCreate) (*entity.Section, error) {
	query := `
		INSERT INTO sections (name, description, created_at, updated_at)
		VALUES (:name, :description, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING *
	`

	rows, err := r.db.NamedQueryContext(ctx, query, section)
	if err != nil {
		return nil, fmt.Errorf("failed to create section: %w", err)
	}
	defer rows.Close()

	var createdSection entity.Section
	if rows.Next() {
		if err := rows.StructScan(&createdSection); err != nil {
			return nil, fmt.Errorf("failed to scan section: %w", err)
		}
	}

	return &createdSection, nil
}

func (r *sqliteSectionRepository) GetByID(ctx context.Context, id int) (*entity.Section, error) {
	query := `SELECT * FROM sections WHERE id = ?`

	var section entity.Section
	err := r.db.GetContext(ctx, &section, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get section: %w", err)
	}

	return &section, nil
}

func (r *sqliteSectionRepository) GetAll(ctx context.Context, activeOnly bool) ([]entity.Section, error) {
	query := `SELECT * FROM sections`

	if activeOnly {
		query += ` WHERE is_active = TRUE`
	}

	query += ` ORDER BY name`

	var sections []entity.Section
	err := r.db.SelectContext(ctx, &sections, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sections: %w", err)
	}

	return sections, nil
}

func (r *sqliteSectionRepository) Update(ctx context.Context, id int, section *entity.SectionUpdate) (*entity.Section, error) {
	query := `UPDATE sections SET `
	args := []interface{}{}
	updates := []string{}

	if section.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *section.Name)
	}

	if section.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *section.Description)
	}

	if section.IsActive != nil {
		updates = append(updates, "is_active = ?")
		args = append(args, *section.IsActive)
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
		return nil, fmt.Errorf("failed to update section: %w", err)
	}
	defer rows.Close()

	var updatedSection entity.Section
	if rows.Next() {
		if err := rows.StructScan(&updatedSection); err != nil {
			return nil, fmt.Errorf("failed to scan section: %w", err)
		}
	}

	return &updatedSection, nil
}

func (r *sqliteSectionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM sections WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete section: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("section not found")
	}

	return nil
}

func (r *sqliteSectionRepository) UpdateProductCount(ctx context.Context, id int, delta int) error {
	query := `UPDATE sections SET product_count = product_count + ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, delta, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update product count: %w", err)
	}

	return nil
}
