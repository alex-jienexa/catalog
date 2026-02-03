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

type sqliteContactRepository struct {
	db *sqlx.DB
}

func NewSQLiteContactRepository(db *sqlx.DB) repository.ContactRepository {
	return &sqliteContactRepository{db: db}
}

func (r *sqliteContactRepository) Create(ctx context.Context, contact *entity.ContactCreate) (*entity.Contact, error) {
	query := `
		INSERT INTO contacts (platform, url, icon, "order", created_at, updated_at)
		VALUES (:platform, :url, :icon, :order, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING *
	`

	rows, err := r.db.NamedQueryContext(ctx, query, contact)
	if err != nil {
		return nil, fmt.Errorf("failed to create contact: %w", err)
	}
	defer rows.Close()

	var createdContact entity.Contact
	if rows.Next() {
		if err := rows.StructScan(&createdContact); err != nil {
			return nil, fmt.Errorf("failed to scan contact: %w", err)
		}
	}

	return &createdContact, nil
}

func (r *sqliteContactRepository) GetByID(ctx context.Context, id int) (*entity.Contact, error) {
	query := `SELECT * FROM contacts WHERE id = ?`

	var contact entity.Contact
	err := r.db.GetContext(ctx, &contact, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	return &contact, nil
}

func (r *sqliteContactRepository) GetAll(ctx context.Context, activeOnly bool) ([]entity.Contact, error) {
	query := `SELECT * FROM contacts`

	if activeOnly {
		query += ` WHERE is_active = TRUE`
	}

	query += ` ORDER BY "order", platform`

	var contacts []entity.Contact
	err := r.db.SelectContext(ctx, &contacts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get contacts: %w", err)
	}

	return contacts, nil
}

func (r *sqliteContactRepository) Update(ctx context.Context, id int, contact *entity.ContactUpdate) (*entity.Contact, error) {
	query := `UPDATE contacts SET `
	args := []interface{}{}
	updates := []string{}

	if contact.Platform != nil {
		updates = append(updates, "platform = ?")
		args = append(args, *contact.Platform)
	}

	if contact.URL != nil {
		updates = append(updates, "url = ?")
		args = append(args, *contact.URL)
	}

	if contact.Icon != nil {
		updates = append(updates, "icon = ?")
		args = append(args, *contact.Icon)
	}

	if contact.IsActive != nil {
		updates = append(updates, "is_active = ?")
		args = append(args, *contact.IsActive)
	}

	if contact.Order != nil {
		updates = append(updates, `"order" = ?`)
		args = append(args, *contact.Order)
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
		return nil, fmt.Errorf("failed to update contact: %w", err)
	}
	defer rows.Close()

	var updatedContact entity.Contact
	if rows.Next() {
		if err := rows.StructScan(&updatedContact); err != nil {
			return nil, fmt.Errorf("failed to scan contact: %w", err)
		}
	}

	return &updatedContact, nil
}

func (r *sqliteContactRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM contacts WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("contact not found")
	}

	return nil
}
