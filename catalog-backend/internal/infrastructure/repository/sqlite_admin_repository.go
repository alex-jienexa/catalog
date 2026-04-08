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

type sqliteAdminRepository struct {
	db *sqlx.DB
}

func NewSQLiteAdminRepository(db *sqlx.DB) repository.AdminRepository {
	return &sqliteAdminRepository{db: db}
}

func (r *sqliteAdminRepository) Create(ctx context.Context, input *entity.AdminCreate, passwordHash string) (*entity.Admin, error) {
	query := `
		INSERT INTO admins (name, username, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING *
	`
	var admin entity.Admin
	err := r.db.QueryRowxContext(ctx, query, input.Name, input.Username, passwordHash).StructScan(&admin)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin: %w", err)
	}
	return &admin, nil
}

func (r *sqliteAdminRepository) GetByID(ctx context.Context, id int) (*entity.Admin, error) {
	query := `SELECT * FROM admins WHERE id = ?`
	var admin entity.Admin
	err := r.db.GetContext(ctx, &admin, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}
	return &admin, nil
}

func (r *sqliteAdminRepository) GetByUsername(ctx context.Context, username string) (*entity.Admin, error) {
	query := `SELECT * FROM admins WHERE username = ?`
	var admin entity.Admin
	err := r.db.GetContext(ctx, &admin, query, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get admin by username: %w", err)
	}
	return &admin, nil
}

func (r *sqliteAdminRepository) GetAll(ctx context.Context) ([]entity.Admin, error) {
	query := `SELECT * FROM admins ORDER BY created_at ASC`
	var admins []entity.Admin
	err := r.db.SelectContext(ctx, &admins, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get admins: %w", err)
	}
	return admins, nil
}

func (r *sqliteAdminRepository) Update(ctx context.Context, id int, input *entity.AdminUpdate, passwordHash *string) (*entity.Admin, error) {
	updates := []string{}
	args := []interface{}{}

	if input.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *input.Name)
	}
	if input.Username != nil {
		updates = append(updates, "username = ?")
		args = append(args, *input.Username)
	}
	if passwordHash != nil {
		updates = append(updates, "password_hash = ?")
		args = append(args, *passwordHash)
	}

	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}

	updates = append(updates, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, id)

	query := `UPDATE admins SET ` + strings.Join(updates, ", ") + ` WHERE id = ? RETURNING *`

	var admin entity.Admin
	err := r.db.QueryRowxContext(ctx, query, args...).StructScan(&admin)
	if err != nil {
		return nil, fmt.Errorf("failed to update admin: %w", err)
	}
	return &admin, nil
}

func (r *sqliteAdminRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM admins WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete admin: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("admin not found")
	}
	return nil
}

func (r *sqliteAdminRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM admins`)
	return count, err
}
