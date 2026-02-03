package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(path string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Включаем foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		// Таблица разделов
		`CREATE TABLE IF NOT EXISTS sections (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			product_count INTEGER DEFAULT 0,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Таблица товаров
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			price REAL NOT NULL CHECK(price >= 0),
			section_id INTEGER NOT NULL,
			description TEXT,
			image_url TEXT,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE RESTRICT
		)`,

		// Таблица контактов
		`CREATE TABLE IF NOT EXISTS contacts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			platform TEXT NOT NULL,
			url TEXT NOT NULL,
			icon TEXT,
			is_active BOOLEAN DEFAULT TRUE,
			"order" INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Индексы
		`CREATE INDEX IF NOT EXISTS idx_products_section_id ON products(section_id)`,
		`CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active)`,
		`CREATE INDEX IF NOT EXISTS idx_sections_is_active ON sections(is_active)`,
		`CREATE INDEX IF NOT EXISTS idx_contacts_is_active_order ON contacts(is_active, "order")`,
	}

	for i, migration := range migrations {
		log.Printf("Running migration %d...", i+1)
		_, err := db.Exec(migration)
		if err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}
