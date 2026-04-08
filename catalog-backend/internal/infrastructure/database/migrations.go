package database

import (
	"database/sql"
	"fmt"
	"log"
)

// RunMigrations выполняет все миграции
func RunMigrations(db *sql.DB) error {
	// Таблица admins (не имеет зависимостей)
	if err := createAdminsTable(db); err != nil {
		return err
	}
	// Таблица customers (не имеет зависимостей)
	if err := createCustomersTable(db); err != nil {
		return err
	}
	// Таблица sections (не имеет зависимостей)
	if err := createSectionsTable(db); err != nil {
		return err
	}
	// Таблица products (зависит от sections)
	if err := createProductsTable(db); err != nil {
		return err
	}
	// Таблица reservations (зависит от customers и products)
	if err := createReservationsTable(db); err != nil {
		return err
	}
	// Таблица contacts (не имеет зависимостей)
	if err := createContactsTable(db); err != nil {
		return err
	}
	// Индексы
	if err := createIndexes(db); err != nil {
		return err
	}
	log.Println("All migrations completed successfully")
	return nil
}

func createAdminsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS admins (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("admins table: %w", err)
	}
	return nil
}

func createCustomersTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS customers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT,
			phone TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("customers table: %w", err)
	}
	return nil
}

func createReservationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS reservations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			customer_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			status TEXT DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("reservations table: %w", err)
	}
	return nil
}

func createSectionsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS sections (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			product_count INTEGER DEFAULT 0,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("sections table: %w", err)
	}
	return nil
}

func createProductsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS products (
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
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("products table: %w", err)
	}
	return nil
}

func createContactsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS contacts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			platform TEXT NOT NULL,
			url TEXT NOT NULL,
			icon TEXT,
			is_active BOOLEAN DEFAULT TRUE,
			"order" INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("contacts table: %w", err)
	}
	return nil
}

func createIndexes(db *sql.DB) error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_products_section_id ON products(section_id)",
		"CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active)",
		"CREATE INDEX IF NOT EXISTS idx_sections_is_active ON sections(is_active)",
		"CREATE INDEX IF NOT EXISTS idx_contacts_is_active_order ON contacts(is_active, \"order\")",
		"CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(phone)",
		"CREATE INDEX IF NOT EXISTS idx_reservations_customer_id ON reservations(customer_id)",
		"CREATE INDEX IF NOT EXISTS idx_reservations_product_id ON reservations(product_id)",
		"CREATE INDEX IF NOT EXISTS idx_reservations_created_at ON reservations(created_at)",
	}
	for _, idx := range indexes {
		if _, err := db.Exec(idx); err != nil {
			return fmt.Errorf("index %s: %w", idx, err)
		}
	}
	return nil
}
