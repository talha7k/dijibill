package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	database := &Database{db: db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

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
			customer_id INTEGER,
			sales_category_id INTEGER,
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

// insertDefaultSettings inserts default settings data
func (d *Database) insertDefaultSettings() error {
	// Insert default payment types
	var paymentCount int
	err := d.db.QueryRow("SELECT COUNT(*) FROM payment_types").Scan(&paymentCount)
	if err != nil {
		return err
	}

	if paymentCount == 0 {
		paymentTypes := []PaymentType{
			{Name: "Cash", NameArabic: "نقدي", Code: "cash", Description: "Cash payment", IsDefault: true, IsActive: true},
			{Name: "Credit Card", NameArabic: "بطاقة ائتمان", Code: "card", Description: "Credit card payment", IsDefault: false, IsActive: true},
			{Name: "Bank Transfer", NameArabic: "تحويل بنكي", Code: "bank_transfer", Description: "Bank transfer payment", IsDefault: false, IsActive: true},
		}

		for _, pt := range paymentTypes {
			if err := d.CreatePaymentType(&pt); err != nil {
				log.Printf("Warning: Could not insert payment type %s: %v", pt.Name, err)
			}
		}
	}

	// Insert default sales categories
	var salesCount int
	err = d.db.QueryRow("SELECT COUNT(*) FROM sales_categories").Scan(&salesCount)
	if err != nil {
		return err
	}

	if salesCount == 0 {
		salesCategories := []SalesCategory{
			{Name: "Retail", NameArabic: "تجزئة", Code: "retail", Description: "Retail sales", DescriptionArabic: "مبيعات التجزئة", IsDefault: true, IsActive: true},
			{Name: "Wholesale", NameArabic: "جملة", Code: "wholesale", Description: "Wholesale sales", DescriptionArabic: "مبيعات الجملة", IsDefault: false, IsActive: true},
			{Name: "Service", NameArabic: "خدمة", Code: "service", Description: "Service sales", DescriptionArabic: "مبيعات الخدمات", IsDefault: false, IsActive: true},
		}

		for _, sc := range salesCategories {
			if err := d.CreateSalesCategory(&sc); err != nil {
				log.Printf("Warning: Could not insert sales category %s: %v", sc.Name, err)
			}
		}
	}

	// Insert default tax rates
	var taxCount int
	err = d.db.QueryRow("SELECT COUNT(*) FROM tax_rates").Scan(&taxCount)
	if err != nil {
		return err
	}

	if taxCount == 0 {
		taxRates := []TaxRate{
			{Name: "Standard VAT", NameArabic: "ضريبة القيمة المضافة", Rate: 15.0, Description: "Standard VAT rate", IsDefault: true, IsActive: true},
			{Name: "Zero VAT", NameArabic: "معفى من الضريبة", Rate: 0.0, Description: "Zero VAT rate", IsDefault: false, IsActive: true},
			{Name: "Exempt", NameArabic: "معفى", Rate: 0.0, Description: "VAT exempt", IsDefault: false, IsActive: true},
		}

		for _, tr := range taxRates {
			if err := d.CreateTaxRate(&tr); err != nil {
				log.Printf("Warning: Could not insert tax rate %s: %v", tr.Name, err)
			}
		}
	}

	// Insert default units of measurement
	var unitCount int
	err = d.db.QueryRow("SELECT COUNT(*) FROM units_of_measurement").Scan(&unitCount)
	if err != nil {
		return err
	}

	if unitCount == 0 {
		units := []UnitOfMeasurement{
			{Value: "pcs", Label: "Pieces", Arabic: "قطعة", IsDefault: true, IsActive: true},
			{Value: "kg", Label: "Kilograms", Arabic: "كيلوغرام", IsDefault: false, IsActive: true},
			{Value: "m", Label: "Meters", Arabic: "متر", IsDefault: false, IsActive: true},
			{Value: "l", Label: "Liters", Arabic: "لتر", IsDefault: false, IsActive: true},
			{Value: "box", Label: "Box", Arabic: "صندوق", IsDefault: false, IsActive: true},
			{Value: "pack", Label: "Pack", Arabic: "عبوة", IsDefault: false, IsActive: true},
		}

		for _, u := range units {
			if err := d.CreateUnitOfMeasurement(&u); err != nil {
				log.Printf("Warning: Could not insert unit %s: %v", u.Label, err)
			}
		}
	}

	// Insert default product settings
	var defaultSettingsCount int
	err = d.db.QueryRow("SELECT COUNT(*) FROM default_product_settings").Scan(&defaultSettingsCount)
	if err != nil {
		return err
	}

	if defaultSettingsCount == 0 {
		// Get default IDs for foreign keys
		var defaultTaxRateID, defaultUnitID, defaultSalesCategoryID int
		
		d.db.QueryRow("SELECT id FROM tax_rates WHERE is_default = 1 LIMIT 1").Scan(&defaultTaxRateID)
		d.db.QueryRow("SELECT id FROM units_of_measurement WHERE is_default = 1 LIMIT 1").Scan(&defaultUnitID)
		d.db.QueryRow("SELECT id FROM sales_categories WHERE is_default = 1 LIMIT 1").Scan(&defaultSalesCategoryID)

		_, err = d.db.Exec(`
			INSERT INTO default_product_settings (
				default_stock, default_tax_rate_id, default_unit_id, 
				default_sales_category_id, default_product_type, default_product_status, 
				default_markup, default_price_includes_tax, default_price_change_allowed
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			1, defaultTaxRateID, defaultUnitID, defaultSalesCategoryID,
			"product", true, 0.0, false, true)
		
		if err != nil {
			log.Printf("Warning: Could not insert default product settings: %v", err)
		}
	}

	return nil
}

// insertDefaultCompany inserts a default company record
func (d *Database) insertDefaultCompany() error {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM companies").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = d.db.Exec(`
			INSERT INTO companies (name, name_arabic, vat_number, cr_number, email, phone, address, address_arabic, city, city_arabic, country, country_arabic)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			"Your Company Name", "اسم شركتك", "123456789012345", "1234567890",
			"info@company.com", "+966501234567",
			"123 Business Street", "شارع الأعمال ١٢٣",
			"Riyadh", "الرياض", "Saudi Arabia", "المملكة العربية السعودية")
		return err
	}
	return nil
}

// Customer operations
func (d *Database) CreateCustomer(customer *Customer) error {
	query := `
		INSERT INTO customers (name, name_arabic, vat_number, email, phone, address, address_arabic, city, city_arabic, country, country_arabic)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := d.db.Exec(query, customer.Name, customer.NameArabic, customer.VATNumber,
		customer.Email, customer.Phone, customer.Address, customer.AddressArabic,
		customer.City, customer.CityArabic, customer.Country, customer.CountryArabic)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	customer.ID = int(id)
	return nil
}

func (d *Database) GetCustomers() ([]Customer, error) {
	query := `SELECT id, name, name_arabic, vat_number, email, phone, address, address_arabic, 
			  city, city_arabic, country, country_arabic, created_at, updated_at FROM customers ORDER BY name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.ID, &c.Name, &c.NameArabic, &c.VATNumber, &c.Email, &c.Phone,
			&c.Address, &c.AddressArabic, &c.City, &c.CityArabic, &c.Country, &c.CountryArabic,
			&c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (d *Database) GetCustomerByID(id int) (*Customer, error) {
	query := `SELECT id, name, name_arabic, vat_number, email, phone, address, address_arabic, 
			  city, city_arabic, country, country_arabic, created_at, updated_at FROM customers WHERE id = ?`

	var c Customer
	err := d.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.NameArabic, &c.VATNumber, &c.Email, &c.Phone,
		&c.Address, &c.AddressArabic, &c.City, &c.CityArabic, &c.Country, &c.CountryArabic,
		&c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *Database) UpdateCustomer(customer *Customer) error {
	query := `
		UPDATE customers SET name = ?, name_arabic = ?, vat_number = ?, email = ?, phone = ?, 
		address = ?, address_arabic = ?, city = ?, city_arabic = ?, country = ?, country_arabic = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := d.db.Exec(query, customer.Name, customer.NameArabic, customer.VATNumber,
		customer.Email, customer.Phone, customer.Address, customer.AddressArabic,
		customer.City, customer.CityArabic, customer.Country, customer.CountryArabic, customer.ID)
	return err
}

func (d *Database) DeleteCustomer(id int) error {
	_, err := d.db.Exec("DELETE FROM customers WHERE id = ?", id)
	return err
}

// ProductCategory operations
func (d *Database) CreateProductCategory(category *ProductCategory) error {
	query := `INSERT INTO product_categories (name, name_arabic, description, description_arabic) VALUES (?, ?, ?, ?)`

	result, err := d.db.Exec(query, category.Name, category.NameArabic, category.Description, category.DescriptionArabic)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	category.ID = int(id)
	return nil
}

func (d *Database) GetProductCategories() ([]ProductCategory, error) {
	query := `SELECT id, name, name_arabic, description, description_arabic, created_at, updated_at FROM product_categories ORDER BY name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []ProductCategory
	for rows.Next() {
		var c ProductCategory
		err := rows.Scan(&c.ID, &c.Name, &c.NameArabic, &c.Description, &c.DescriptionArabic, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// Product operations
func (d *Database) CreateProduct(product *Product) error {
	query := `
		INSERT INTO products (name, name_arabic, description, description_arabic, category_id, unit_price, vat_rate, unit, unit_arabic, sku, barcode, stock, min_stock, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := d.db.Exec(query, product.Name, product.NameArabic, product.Description, product.DescriptionArabic,
		product.CategoryID, product.UnitPrice, product.VATRate, product.Unit, product.UnitArabic, 
		product.SKU, product.Barcode, product.Stock, product.MinStock, product.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = int(id)
	return nil
}

func (d *Database) UpdateProduct(product *Product) error {
	query := `
		UPDATE products SET name = ?, name_arabic = ?, description = ?, description_arabic = ?, 
		category_id = ?, unit_price = ?, vat_rate = ?, unit = ?, unit_arabic = ?, 
		sku = ?, barcode = ?, stock = ?, min_stock = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := d.db.Exec(query, product.Name, product.NameArabic, product.Description, product.DescriptionArabic,
		product.CategoryID, product.UnitPrice, product.VATRate, product.Unit, product.UnitArabic,
		product.SKU, product.Barcode, product.Stock, product.MinStock, product.IsActive, product.ID)
	return err
}

func (d *Database) GetProducts() ([]Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.name_arabic, p.description, p.description_arabic, 
			p.category_id, COALESCE(pc.name, '') as category_name,
			p.unit_price, p.vat_rate, p.unit, p.unit_arabic, 
			p.sku, p.barcode, p.stock, p.min_stock, p.is_active, 
			p.created_at, p.updated_at
		FROM products p
		LEFT JOIN product_categories pc ON p.category_id = pc.id
		ORDER BY p.name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.NameArabic, &product.Description, &product.DescriptionArabic,
			&product.CategoryID, &product.CategoryName, &product.UnitPrice, &product.VATRate, &product.Unit, &product.UnitArabic,
			&product.SKU, &product.Barcode, &product.Stock, &product.MinStock, &product.IsActive, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		
		// Set default values for new fields that don't exist in the database yet
		product.Color = ""
		product.ImageURL = ""
		product.ServiceNotUsingStock = false
		
		products = append(products, product)
	}
	return products, nil
}

// Invoice operations
func (d *Database) CreateInvoice(invoice *Invoice) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Generate invoice number if not provided
	if invoice.InvoiceNumber == "" {
		invoice.InvoiceNumber = d.generateInvoiceNumber()
	}

	// Insert invoice
	query := `
		INSERT INTO invoices (invoice_number, customer_id, sales_category_id, issue_date, due_date, sub_total, vat_amount, total_amount, status, notes, notes_arabic, qr_code)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := tx.Exec(query, invoice.InvoiceNumber, invoice.CustomerID, invoice.SalesCategoryID, invoice.IssueDate, invoice.DueDate,
		invoice.SubTotal, invoice.VATAmount, invoice.TotalAmount, invoice.Status, invoice.Notes, invoice.NotesArabic, invoice.QRCode)
	if err != nil {
		return err
	}

	invoiceID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	invoice.ID = int(invoiceID)

	// Insert invoice items
	for _, item := range invoice.Items {
		itemQuery := `
			INSERT INTO invoice_items (invoice_id, product_id, quantity, unit_price, vat_rate, vat_amount, total_amount)
			VALUES (?, ?, ?, ?, ?, ?, ?)`

		_, err = tx.Exec(itemQuery, invoiceID, item.ProductID, item.Quantity, item.UnitPrice, item.VATRate, item.VATAmount, item.TotalAmount)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *Database) GetInvoices() ([]Invoice, error) {
	query := `SELECT id, invoice_number, customer_id, sales_category_id, issue_date, due_date, sub_total, vat_amount, total_amount, status, notes, notes_arabic, qr_code, created_at, updated_at FROM invoices ORDER BY created_at DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var inv Invoice
		err := rows.Scan(&inv.ID, &inv.InvoiceNumber, &inv.CustomerID, &inv.SalesCategoryID, &inv.IssueDate, &inv.DueDate,
			&inv.SubTotal, &inv.VATAmount, &inv.TotalAmount, &inv.Status, &inv.Notes, &inv.NotesArabic,
			&inv.QRCode, &inv.CreatedAt, &inv.UpdatedAt)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, inv)
	}
	return invoices, nil
}

func (d *Database) GetInvoiceByID(id int) (*Invoice, error) {
	query := `SELECT id, invoice_number, customer_id, sales_category_id, issue_date, due_date, sub_total, vat_amount, total_amount, status, notes, notes_arabic, qr_code, created_at, updated_at FROM invoices WHERE id = ?`

	var inv Invoice
	err := d.db.QueryRow(query, id).Scan(&inv.ID, &inv.InvoiceNumber, &inv.CustomerID, &inv.SalesCategoryID, &inv.IssueDate, &inv.DueDate,
		&inv.SubTotal, &inv.VATAmount, &inv.TotalAmount, &inv.Status, &inv.Notes, &inv.NotesArabic,
		&inv.QRCode, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Get customer
	if inv.CustomerID > 0 {
		customer, err := d.GetCustomerByID(inv.CustomerID)
		if err == nil {
			inv.Customer = customer
		}
	}

	// Get sales category
	if inv.SalesCategoryID > 0 {
		salesCategory, err := d.GetSalesCategoryByID(inv.SalesCategoryID)
		if err == nil {
			inv.SalesCategory = salesCategory
		}
	}

	// Get invoice items
	items, err := d.GetInvoiceItems(inv.ID)
	if err == nil {
		inv.Items = items
	}

	return &inv, nil
}

func (d *Database) GetInvoiceItems(invoiceID int) ([]InvoiceItem, error) {
	query := `SELECT id, invoice_id, product_id, quantity, unit_price, vat_rate, vat_amount, total_amount, created_at FROM invoice_items WHERE invoice_id = ?`

	rows, err := d.db.Query(query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []InvoiceItem
	for rows.Next() {
		var item InvoiceItem
		err := rows.Scan(&item.ID, &item.InvoiceID, &item.ProductID, &item.Quantity, &item.UnitPrice,
			&item.VATRate, &item.VATAmount, &item.TotalAmount, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (d *Database) GetCompany() (*Company, error) {
	query := `SELECT id, name, name_arabic, vat_number, cr_number, email, phone, address, address_arabic, city, city_arabic, country, country_arabic, COALESCE(logo, '') as logo FROM companies LIMIT 1`

	var c Company
	err := d.db.QueryRow(query).Scan(&c.ID, &c.Name, &c.NameArabic, &c.VATNumber, &c.CRNumber, &c.Email, &c.Phone,
		&c.Address, &c.AddressArabic, &c.City, &c.CityArabic, &c.Country, &c.CountryArabic, &c.Logo)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *Database) UpdateCompany(company *Company) error {
	query := `
		UPDATE companies SET name = ?, name_arabic = ?, vat_number = ?, cr_number = ?, email = ?, phone = ?, 
		address = ?, address_arabic = ?, city = ?, city_arabic = ?, country = ?, country_arabic = ?, logo = NULLIF(?, '')
		WHERE id = ?`

	_, err := d.db.Exec(query, company.Name, company.NameArabic, company.VATNumber, company.CRNumber,
		company.Email, company.Phone, company.Address, company.AddressArabic, company.City, company.CityArabic,
		company.Country, company.CountryArabic, company.Logo, company.ID)
	return err
}



func (d *Database) generateInvoiceNumber() string {
	var count int
	d.db.QueryRow("SELECT COUNT(*) FROM invoices").Scan(&count)
	return fmt.Sprintf("INV-%06d", count+1)
}

// PaymentType operations
func (d *Database) CreatePaymentType(paymentType *PaymentType) error {
	query := `INSERT INTO payment_types (name, name_arabic, code, description, is_default, is_active) VALUES (?, ?, ?, ?, ?, ?)`

	result, err := d.db.Exec(query, paymentType.Name, paymentType.NameArabic, paymentType.Code, paymentType.Description, paymentType.IsDefault, paymentType.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	paymentType.ID = int(id)
	return nil
}

func (d *Database) GetPaymentTypes() ([]PaymentType, error) {
	query := `SELECT id, name, name_arabic, code, description, is_default, is_active, created_at, updated_at FROM payment_types ORDER BY name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paymentTypes []PaymentType
	for rows.Next() {
		var pt PaymentType
		err := rows.Scan(&pt.ID, &pt.Name, &pt.NameArabic, &pt.Code, &pt.Description, &pt.IsDefault, &pt.IsActive, &pt.CreatedAt, &pt.UpdatedAt)
		if err != nil {
			return nil, err
		}
		paymentTypes = append(paymentTypes, pt)
	}
	return paymentTypes, nil
}

func (d *Database) UpdatePaymentType(paymentType *PaymentType) error {
	query := `UPDATE payment_types SET name = ?, name_arabic = ?, code = ?, description = ?, is_default = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	_, err := d.db.Exec(query, paymentType.Name, paymentType.NameArabic, paymentType.Code, paymentType.Description, paymentType.IsDefault, paymentType.IsActive, paymentType.ID)
	return err
}

func (d *Database) DeletePaymentType(id int) error {
	_, err := d.db.Exec("DELETE FROM payment_types WHERE id = ?", id)
	return err
}

// SalesCategory operations
func (d *Database) CreateSalesCategory(salesCategory *SalesCategory) error {
	query := `INSERT INTO sales_categories (name, name_arabic, code, description, description_arabic, is_default, is_active) VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := d.db.Exec(query, salesCategory.Name, salesCategory.NameArabic, salesCategory.Code, salesCategory.Description, salesCategory.DescriptionArabic, salesCategory.IsDefault, salesCategory.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	salesCategory.ID = int(id)
	return nil
}

func (d *Database) GetSalesCategories() ([]SalesCategory, error) {
	query := `SELECT id, name, name_arabic, code, description, description_arabic, is_default, is_active, created_at, updated_at FROM sales_categories ORDER BY name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var salesCategories []SalesCategory
	for rows.Next() {
		var sc SalesCategory
		err := rows.Scan(&sc.ID, &sc.Name, &sc.NameArabic, &sc.Code, &sc.Description, &sc.DescriptionArabic, &sc.IsDefault, &sc.IsActive, &sc.CreatedAt, &sc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		salesCategories = append(salesCategories, sc)
	}
	return salesCategories, nil
}

func (d *Database) GetSalesCategoryByID(id int) (*SalesCategory, error) {
	query := `SELECT id, name, name_arabic, code, description, description_arabic, is_default, is_active, created_at, updated_at FROM sales_categories WHERE id = ?`

	var sc SalesCategory
	err := d.db.QueryRow(query, id).Scan(&sc.ID, &sc.Name, &sc.NameArabic, &sc.Code, &sc.Description, &sc.DescriptionArabic, &sc.IsDefault, &sc.IsActive, &sc.CreatedAt, &sc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

func (d *Database) UpdateSalesCategory(salesCategory *SalesCategory) error {
	query := `UPDATE sales_categories SET name = ?, name_arabic = ?, code = ?, description = ?, description_arabic = ?, is_default = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	_, err := d.db.Exec(query, salesCategory.Name, salesCategory.NameArabic, salesCategory.Code, salesCategory.Description, salesCategory.DescriptionArabic, salesCategory.IsDefault, salesCategory.IsActive, salesCategory.ID)
	return err
}

func (d *Database) DeleteSalesCategory(id int) error {
	_, err := d.db.Exec("DELETE FROM sales_categories WHERE id = ?", id)
	return err
}

// TaxRate operations
func (d *Database) CreateTaxRate(taxRate *TaxRate) error {
	query := `INSERT INTO tax_rates (name, name_arabic, rate, description, is_default, is_active) VALUES (?, ?, ?, ?, ?, ?)`

	result, err := d.db.Exec(query, taxRate.Name, taxRate.NameArabic, taxRate.Rate, taxRate.Description, taxRate.IsDefault, taxRate.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	taxRate.ID = int(id)
	return nil
}

func (d *Database) GetTaxRates() ([]TaxRate, error) {
	query := `SELECT id, name, name_arabic, rate, description, is_default, is_active, created_at, updated_at FROM tax_rates ORDER BY name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taxRates []TaxRate
	for rows.Next() {
		var tr TaxRate
		err := rows.Scan(&tr.ID, &tr.Name, &tr.NameArabic, &tr.Rate, &tr.Description, &tr.IsDefault, &tr.IsActive, &tr.CreatedAt, &tr.UpdatedAt)
		if err != nil {
			return nil, err
		}
		taxRates = append(taxRates, tr)
	}
	return taxRates, nil
}

func (d *Database) UpdateTaxRate(taxRate *TaxRate) error {
	query := `UPDATE tax_rates SET name = ?, name_arabic = ?, rate = ?, description = ?, is_default = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	_, err := d.db.Exec(query, taxRate.Name, taxRate.NameArabic, taxRate.Rate, taxRate.Description, taxRate.IsDefault, taxRate.IsActive, taxRate.ID)
	return err
}

func (d *Database) DeleteTaxRate(id int) error {
	_, err := d.db.Exec("DELETE FROM tax_rates WHERE id = ?", id)
	return err
}

// UnitOfMeasurement operations
func (d *Database) CreateUnitOfMeasurement(unit *UnitOfMeasurement) error {
	query := `INSERT INTO units_of_measurement (value, label, arabic, is_default, is_active) VALUES (?, ?, ?, ?, ?)`

	result, err := d.db.Exec(query, unit.Value, unit.Label, unit.Arabic, unit.IsDefault, unit.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	unit.ID = int(id)
	return nil
}

func (d *Database) GetUnitsOfMeasurement() ([]UnitOfMeasurement, error) {
	query := `SELECT id, value, label, arabic, is_default, is_active, created_at, updated_at FROM units_of_measurement ORDER BY label`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []UnitOfMeasurement
	for rows.Next() {
		var u UnitOfMeasurement
		err := rows.Scan(&u.ID, &u.Value, &u.Label, &u.Arabic, &u.IsDefault, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	return units, nil
}

func (d *Database) UpdateUnitOfMeasurement(unit *UnitOfMeasurement) error {
	query := `UPDATE units_of_measurement SET value = ?, label = ?, arabic = ?, is_default = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	_, err := d.db.Exec(query, unit.Value, unit.Label, unit.Arabic, unit.IsDefault, unit.IsActive, unit.ID)
	return err
}

func (d *Database) DeleteUnitOfMeasurement(id int) error {
	_, err := d.db.Exec("DELETE FROM units_of_measurement WHERE id = ?", id)
	return err
}

// DefaultProductSettings operations
func (d *Database) GetDefaultProductSettings() (*DefaultProductSettings, error) {
	query := `SELECT id, default_stock, default_tax_rate_id, default_unit_id, 
			  default_product_type, default_product_status, default_markup, 
			  default_price_includes_tax, default_price_change_allowed, created_at, updated_at 
			  FROM default_product_settings LIMIT 1`

	var dps DefaultProductSettings
	err := d.db.QueryRow(query).Scan(&dps.ID, &dps.DefaultStock, &dps.DefaultTaxRateID, &dps.DefaultUnitID,
		&dps.DefaultProductType, &dps.DefaultProductStatus,
		&dps.DefaultMarkup, &dps.DefaultPriceIncludesTax, &dps.DefaultPriceChangeAllowed, &dps.CreatedAt, &dps.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &dps, nil
}

func (d *Database) UpdateDefaultProductSettings(settings *DefaultProductSettings) error {
	query := `UPDATE default_product_settings SET 
			  default_stock = ?, default_tax_rate_id = ?, default_unit_id = ?, 
			  default_product_type = ?, default_product_status = ?, 
			  default_markup = ?, default_price_includes_tax = ?, default_price_change_allowed = ?, 
			  updated_at = CURRENT_TIMESTAMP 
			  WHERE id = ?`

	_, err := d.db.Exec(query, settings.DefaultStock, settings.DefaultTaxRateID, settings.DefaultUnitID,
		settings.DefaultProductType,
		settings.DefaultProductStatus, settings.DefaultMarkup, settings.DefaultPriceIncludesTax,
		settings.DefaultPriceChangeAllowed, settings.ID)
	return err
}
