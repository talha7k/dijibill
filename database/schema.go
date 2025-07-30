package database

import (
	"fmt"
	"log"
)

// createTables creates all necessary tables
func (d *Database) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS companies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			name_arabic TEXT,
			vat_number TEXT NOT NULL,
			cr_number TEXT,
			email TEXT,
			phone TEXT,
			address TEXT,
			address_arabic TEXT,
			city TEXT,
			city_arabic TEXT,
			country TEXT,
			country_arabic TEXT,
			logo TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS customers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			name_arabic TEXT,
			vat_number TEXT,
			email TEXT,
			phone TEXT,
			address TEXT,
			address_arabic TEXT,
			city TEXT,
			city_arabic TEXT,
			country TEXT,
			country_arabic TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS product_categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			name_arabic TEXT,
			description TEXT,
			description_arabic TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			name_arabic TEXT,
			description TEXT,
			description_arabic TEXT,
			category_id INTEGER,
			unit_price REAL NOT NULL,
			vat_rate REAL DEFAULT 15.0,
			unit TEXT DEFAULT 'pcs',
			unit_arabic TEXT DEFAULT 'قطعة',
			sku TEXT,
			barcode TEXT,
			stock INTEGER DEFAULT 0,
			min_stock INTEGER DEFAULT 0,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (category_id) REFERENCES product_categories(id)
		)`,
		`CREATE TABLE IF NOT EXISTS invoices (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			invoice_number TEXT UNIQUE NOT NULL,
			customer_id INTEGER NOT NULL,
			sales_category_id INTEGER NOT NULL,
			issue_date DATETIME NOT NULL,
			due_date DATETIME,
			sub_total REAL NOT NULL,
			vat_amount REAL NOT NULL,
			total_amount REAL NOT NULL,
			status TEXT DEFAULT 'draft',
			notes TEXT,
			notes_arabic TEXT,
			qr_code TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (customer_id) REFERENCES customers(id),
			FOREIGN KEY (sales_category_id) REFERENCES sales_categories(id)
		)`,
		`CREATE TABLE IF NOT EXISTS invoice_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			invoice_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity REAL NOT NULL,
			unit_price REAL NOT NULL,
			vat_rate REAL NOT NULL,
			vat_amount REAL NOT NULL,
			total_amount REAL NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`,
		`CREATE TABLE IF NOT EXISTS payment_types (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			name_arabic TEXT,
			code TEXT UNIQUE,
			description TEXT,
			is_default BOOLEAN DEFAULT 0,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS payments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			invoice_id INTEGER NOT NULL,
			payment_type_id INTEGER NOT NULL,
			amount REAL NOT NULL,
			payment_date DATETIME NOT NULL,
			reference TEXT,
			notes TEXT,
			notes_arabic TEXT,
			status TEXT DEFAULT 'completed',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
			FOREIGN KEY (payment_type_id) REFERENCES payment_types(id)
		)`,
		`CREATE TABLE IF NOT EXISTS sales_categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			name_arabic TEXT,
			code TEXT UNIQUE,
			description TEXT,
			description_arabic TEXT,
			is_default BOOLEAN DEFAULT 0,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS tax_rates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			name_arabic TEXT,
			rate REAL NOT NULL,
			description TEXT,
			is_default BOOLEAN DEFAULT 0,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS units_of_measurement (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			value TEXT NOT NULL,
			label TEXT NOT NULL,
			arabic TEXT,
			is_default BOOLEAN DEFAULT 0,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS default_product_settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			default_stock INTEGER DEFAULT 0,
			default_tax_rate_id INTEGER,
			default_unit_id INTEGER,
			default_product_type TEXT DEFAULT 'product',
			default_product_status BOOLEAN DEFAULT 1,
			default_markup REAL DEFAULT 0.0,
			default_price_includes_tax BOOLEAN DEFAULT 0,
			default_price_change_allowed BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (default_tax_rate_id) REFERENCES tax_rates(id),
			FOREIGN KEY (default_unit_id) REFERENCES units_of_measurement(id)
		)`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}

	// Insert default company if not exists
	if err := d.insertDefaultCompany(); err != nil {
		log.Printf("Warning: Could not insert default company: %v", err)
	}

	// Insert default settings data
	if err := d.insertDefaultSettings(); err != nil {
		log.Printf("Warning: Could not insert default settings: %v", err)
	}

	return nil
}