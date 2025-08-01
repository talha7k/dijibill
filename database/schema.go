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
		`CREATE TABLE IF NOT EXISTS system_settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			currency TEXT DEFAULT 'SAR',
			language TEXT DEFAULT 'en',
			timezone TEXT DEFAULT 'Asia/Riyadh',
			date_format TEXT DEFAULT 'DD/MM/YYYY',
			invoice_language TEXT DEFAULT 'english',
			zatca_enabled BOOLEAN DEFAULT 1,
			auto_backup BOOLEAN DEFAULT 1,
			backup_frequency TEXT DEFAULT 'daily',
			last_backup_time DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
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

// runMultiTenantMigration applies the multi-tenant migration
func (d *Database) runMultiTenantMigration() error {
	// Create users table
	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		is_active BOOLEAN NOT NULL DEFAULT 1,
		company_id INTEGER NOT NULL,
		last_login DATETIME,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
	)`

	if _, err := d.db.Exec(createUsersTable); err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	// Add company_id columns to existing tables
	tables := []string{
		"customers", "suppliers", "product_categories", "products", 
		"sales_invoices", "purchase_invoices", "payment_types", "payments",
		"sales_categories", "tax_rates", "units_of_measurement", 
		"default_product_settings", "purchase_product_categories", 
		"purchase_products", "system_settings",
	}

	for _, table := range tables {
		// Check if company_id column already exists
		var columnExists bool
		err := d.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info(?) WHERE name='company_id'", table).Scan(&columnExists)
		if err != nil {
			return fmt.Errorf("error checking company_id column in %s: %v", table, err)
		}

		if !columnExists {
			query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN company_id INTEGER DEFAULT 1", table)
			if _, err := d.db.Exec(query); err != nil {
				return fmt.Errorf("error adding company_id to %s: %v", table, err)
			}
		}
	}

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)",
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_customers_company_id ON customers(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_suppliers_company_id ON suppliers(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_product_categories_company_id ON product_categories(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_products_company_id ON products(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_sales_invoices_company_id ON sales_invoices(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_purchase_invoices_company_id ON purchase_invoices(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_payment_types_company_id ON payment_types(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_payments_company_id ON payments(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_sales_categories_company_id ON sales_categories(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_tax_rates_company_id ON tax_rates(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_units_of_measurement_company_id ON units_of_measurement(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_default_product_settings_company_id ON default_product_settings(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_purchase_product_categories_company_id ON purchase_product_categories(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_purchase_products_company_id ON purchase_products(company_id)",
		"CREATE INDEX IF NOT EXISTS idx_system_settings_company_id ON system_settings(company_id)",
	}

	for _, indexQuery := range indexes {
		if _, err := d.db.Exec(indexQuery); err != nil {
			return fmt.Errorf("error creating index: %v", err)
		}
	}

	// Create a default admin user for the first company
	createAdminUser := `INSERT OR IGNORE INTO users (
		username, email, password, first_name, last_name, role, is_active, company_id
	) VALUES (
		'admin', 'admin@company.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 
		'Admin', 'User', 'admin', 1, 1
	)`

	if _, err := d.db.Exec(createAdminUser); err != nil {
		return fmt.Errorf("error creating default admin user: %v", err)
	}

	// Update existing data to belong to company 1
	updateQueries := []string{
		"UPDATE customers SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE suppliers SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE product_categories SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE products SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE sales_invoices SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE purchase_invoices SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE payment_types SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE payments SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE sales_categories SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE tax_rates SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE units_of_measurement SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE default_product_settings SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE purchase_product_categories SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE purchase_products SET company_id = 1 WHERE company_id IS NULL",
		"UPDATE system_settings SET company_id = 1 WHERE company_id IS NULL",
	}

	for _, updateQuery := range updateQueries {
		if _, err := d.db.Exec(updateQuery); err != nil {
			return fmt.Errorf("error updating existing data: %v", err)
		}
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

	// Check if table_number column exists in sales_invoices table
	columnExists = false
	err = d.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sales_invoices') WHERE name='table_number'").Scan(&columnExists)
	if err != nil {
		return err
	}

	// Add table_number column if it doesn't exist
	if !columnExists {
		_, err = d.db.Exec("ALTER TABLE sales_invoices ADD COLUMN table_number TEXT")
		if err != nil {
			return fmt.Errorf("error adding table_number column: %v", err)
		}
		log.Println("Added table_number column to sales_invoices table")
	}

	// Multi-tenant migration: Check if users table exists
	var tableExists bool
	err = d.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableExists)
	if err != nil {
		return err
	}

	// Run multi-tenant migration if users table doesn't exist
	if !tableExists {
		if err := d.runMultiTenantMigration(); err != nil {
			return fmt.Errorf("error running multi-tenant migration: %v", err)
		}
		log.Println("Successfully applied multi-tenant migration")
	}

	// Check if intro_viewed column exists in users table
	columnExists = false
	err = d.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='intro_viewed'").Scan(&columnExists)
	if err != nil {
		return err
	}

	// Add intro_viewed column if it doesn't exist
	if !columnExists {
		_, err = d.db.Exec("ALTER TABLE users ADD COLUMN intro_viewed BOOLEAN DEFAULT 0")
		if err != nil {
			return fmt.Errorf("error adding intro_viewed column: %v", err)
		}
		log.Println("Added intro_viewed column to users table")
	}

	// Check if created_by column exists in sales_invoices table
	columnExists = false
	err = d.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sales_invoices') WHERE name='created_by'").Scan(&columnExists)
	if err != nil {
		return err
	}

	// Add created_by and updated_by columns to sales_invoices if they don't exist
	if !columnExists {
		_, err = d.db.Exec("ALTER TABLE sales_invoices ADD COLUMN created_by INTEGER")
		if err != nil {
			return fmt.Errorf("error adding created_by column to sales_invoices: %v", err)
		}
		_, err = d.db.Exec("ALTER TABLE sales_invoices ADD COLUMN updated_by INTEGER")
		if err != nil {
			return fmt.Errorf("error adding updated_by column to sales_invoices: %v", err)
		}
		log.Println("Added created_by and updated_by columns to sales_invoices table")
	}

	// Check if created_by column exists in purchase_invoices table
	columnExists = false
	err = d.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('purchase_invoices') WHERE name='created_by'").Scan(&columnExists)
	if err != nil {
		return err
	}

	// Add created_by and updated_by columns to purchase_invoices if they don't exist
	if !columnExists {
		_, err = d.db.Exec("ALTER TABLE purchase_invoices ADD COLUMN created_by INTEGER")
		if err != nil {
			return fmt.Errorf("error adding created_by column to purchase_invoices: %v", err)
		}
		_, err = d.db.Exec("ALTER TABLE purchase_invoices ADD COLUMN updated_by INTEGER")
		if err != nil {
			return fmt.Errorf("error adding updated_by column to purchase_invoices: %v", err)
		}
		log.Println("Added created_by and updated_by columns to purchase_invoices table")
	}

	// Create indexes for the new user tracking columns
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_sales_invoices_created_by ON sales_invoices(created_by)",
		"CREATE INDEX IF NOT EXISTS idx_sales_invoices_updated_by ON sales_invoices(updated_by)",
		"CREATE INDEX IF NOT EXISTS idx_purchase_invoices_created_by ON purchase_invoices(created_by)",
		"CREATE INDEX IF NOT EXISTS idx_purchase_invoices_updated_by ON purchase_invoices(updated_by)",
	}

	for _, indexQuery := range indexes {
		if _, err := d.db.Exec(indexQuery); err != nil {
			log.Printf("Warning: Could not create index: %v", err)
		}
	}



	// Check if logo_file_id column exists in companies table
	columnExists = false
	err = d.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('companies') WHERE name='logo_file_id'").Scan(&columnExists)
	if err != nil {
		return err
	}

	// Add logo_file_id column if it doesn't exist
	if !columnExists {
		_, err = d.db.Exec("ALTER TABLE companies ADD COLUMN logo_file_id INTEGER")
		if err != nil {
			return fmt.Errorf("error adding logo_file_id column: %v", err)
		}
		
		// Create index for logo_file_id
		_, err = d.db.Exec("CREATE INDEX IF NOT EXISTS idx_companies_logo_file_id ON companies(logo_file_id)")
		if err != nil {
			return fmt.Errorf("error creating index for logo_file_id: %v", err)
		}
		
		log.Println("Added logo_file_id column to companies table")
	}

	return nil
}