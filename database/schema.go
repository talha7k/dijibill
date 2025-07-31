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
		`CREATE TABLE IF NOT EXISTS suppliers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			company_name TEXT NOT NULL,
			company_name_arabic TEXT,
			contact_person TEXT NOT NULL,
			contact_person_arabic TEXT,
			vat_number TEXT,
			email TEXT,
			phone TEXT,
			address TEXT,
			address_arabic TEXT,
			city TEXT,
			city_arabic TEXT,
			country TEXT,
			country_arabic TEXT,
			payment_terms TEXT DEFAULT 'net_30',
			active BOOLEAN DEFAULT 1,
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
		`CREATE TABLE IF NOT EXISTS sales_invoices (
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
		`CREATE TABLE IF NOT EXISTS sales_invoice_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			invoice_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity REAL NOT NULL,
			unit_price REAL NOT NULL,
			vat_rate REAL NOT NULL,
			vat_amount REAL NOT NULL,
			total_amount REAL NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (invoice_id) REFERENCES sales_invoices(id) ON DELETE CASCADE,
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`,
		`CREATE TABLE IF NOT EXISTS purchase_invoices (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			invoice_number TEXT UNIQUE NOT NULL,
			supplier_id INTEGER NOT NULL,
			issue_date DATETIME NOT NULL,
			due_date DATETIME,
			sub_total REAL NOT NULL,
			vat_amount REAL NOT NULL,
			vat_rate REAL DEFAULT 15.0,
			vat_inclusive BOOLEAN DEFAULT 0,
			total_amount REAL NOT NULL,
			status TEXT DEFAULT 'draft',
			notes TEXT,
			notes_arabic TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (supplier_id) REFERENCES suppliers(id)
		)`,
		`CREATE TABLE IF NOT EXISTS purchase_invoice_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			invoice_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity REAL NOT NULL,
			unit_price REAL NOT NULL,
			vat_rate REAL NOT NULL,
			vat_amount REAL NOT NULL,
			total_amount REAL NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (invoice_id) REFERENCES purchase_invoices(id) ON DELETE CASCADE,
			FOREIGN KEY (product_id) REFERENCES products(id)
		)`,
		`CREATE VIEW IF NOT EXISTS invoices AS SELECT * FROM sales_invoices`,
		`CREATE VIEW IF NOT EXISTS invoice_items AS SELECT * FROM sales_invoice_items`,
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
		`CREATE TABLE IF NOT EXISTS purchase_product_categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			name_arabic TEXT,
			description TEXT,
			description_arabic TEXT,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS purchase_products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			name_arabic TEXT,
			description TEXT,
			description_arabic TEXT,
			category_id INTEGER,
			unit_price REAL DEFAULT 0.0,
			vat_rate REAL DEFAULT 15.0,
			unit TEXT DEFAULT 'pcs',
			unit_arabic TEXT DEFAULT 'قطعة',
			sku TEXT,
			barcode TEXT,
			is_active BOOLEAN DEFAULT 1,
			notes TEXT,
			notes_arabic TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (category_id) REFERENCES purchase_product_categories(id)
		)`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}

	// Run migrations for existing databases
	if err := d.runMigrations(); err != nil {
		log.Printf("Warning: Could not run migrations: %v", err)
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

// runMigrations handles database schema migrations for existing databases
func (d *Database) runMigrations() error {
	// Check if vat_rate column exists in purchase_invoices table
	var columnExists bool
	err := d.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('purchase_invoices') WHERE name='vat_rate'").Scan(&columnExists)
	if err != nil {
		return err
	}

	// Add vat_rate column if it doesn't exist
	if !columnExists {
		_, err = d.db.Exec("ALTER TABLE purchase_invoices ADD COLUMN vat_rate REAL DEFAULT 15.0")
		if err != nil {
			return fmt.Errorf("error adding vat_rate column: %v", err)
		}
		log.Println("Added vat_rate column to purchase_invoices table")
	}

	// Check if vat_inclusive column exists in purchase_invoices table
	columnExists = false
	err = d.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('purchase_invoices') WHERE name='vat_inclusive'").Scan(&columnExists)
	if err != nil {
		return err
	}

	// Add vat_inclusive column if it doesn't exist
	if !columnExists {
		_, err = d.db.Exec("ALTER TABLE purchase_invoices ADD COLUMN vat_inclusive BOOLEAN DEFAULT 0")
		if err != nil {
			return fmt.Errorf("error adding vat_inclusive column: %v", err)
		}
		log.Println("Added vat_inclusive column to purchase_invoices table")
	}

	return nil
}