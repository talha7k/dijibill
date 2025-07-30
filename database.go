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
		`CREATE TABLE IF NOT EXISTS item_categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			name_arabic TEXT,
			description TEXT,
			description_arabic TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sale_items (
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
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (category_id) REFERENCES item_categories(id)
		)`,
		`CREATE TABLE IF NOT EXISTS invoices (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			invoice_number TEXT UNIQUE NOT NULL,
			customer_id INTEGER NOT NULL,
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
			FOREIGN KEY (customer_id) REFERENCES customers(id)
		)`,
		`CREATE TABLE IF NOT EXISTS invoice_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			invoice_id INTEGER NOT NULL,
			sale_item_id INTEGER NOT NULL,
			quantity REAL NOT NULL,
			unit_price REAL NOT NULL,
			vat_rate REAL NOT NULL,
			vat_amount REAL NOT NULL,
			total_amount REAL NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
			FOREIGN KEY (sale_item_id) REFERENCES sale_items(id)
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

// ItemCategory operations
func (d *Database) CreateItemCategory(category *ItemCategory) error {
	query := `INSERT INTO item_categories (name, name_arabic, description, description_arabic) VALUES (?, ?, ?, ?)`

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

func (d *Database) GetItemCategories() ([]ItemCategory, error) {
	query := `SELECT id, name, name_arabic, description, description_arabic, created_at, updated_at FROM item_categories ORDER BY name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []ItemCategory
	for rows.Next() {
		var c ItemCategory
		err := rows.Scan(&c.ID, &c.Name, &c.NameArabic, &c.Description, &c.DescriptionArabic, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// SaleItem operations
func (d *Database) CreateSaleItem(item *SaleItem) error {
	query := `
		INSERT INTO sale_items (name, name_arabic, description, description_arabic, category_id, unit_price, vat_rate, unit, unit_arabic)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := d.db.Exec(query, item.Name, item.NameArabic, item.Description, item.DescriptionArabic,
		item.CategoryID, item.UnitPrice, item.VATRate, item.Unit, item.UnitArabic)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	item.ID = int(id)
	return nil
}

func (d *Database) GetSaleItems() ([]SaleItem, error) {
	query := `SELECT id, name, name_arabic, description, description_arabic, category_id, unit_price, vat_rate, unit, unit_arabic, created_at, updated_at FROM sale_items ORDER BY name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []SaleItem
	for rows.Next() {
		var item SaleItem
		err := rows.Scan(&item.ID, &item.Name, &item.NameArabic, &item.Description, &item.DescriptionArabic,
			&item.CategoryID, &item.UnitPrice, &item.VATRate, &item.Unit, &item.UnitArabic, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
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
		INSERT INTO invoices (invoice_number, customer_id, issue_date, due_date, sub_total, vat_amount, total_amount, status, notes, notes_arabic, qr_code)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := tx.Exec(query, invoice.InvoiceNumber, invoice.CustomerID, invoice.IssueDate, invoice.DueDate,
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
			INSERT INTO invoice_items (invoice_id, sale_item_id, quantity, unit_price, vat_rate, vat_amount, total_amount)
			VALUES (?, ?, ?, ?, ?, ?, ?)`

		_, err = tx.Exec(itemQuery, invoiceID, item.SaleItemID, item.Quantity, item.UnitPrice, item.VATRate, item.VATAmount, item.TotalAmount)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *Database) GetInvoices() ([]Invoice, error) {
	query := `SELECT id, invoice_number, customer_id, issue_date, due_date, sub_total, vat_amount, total_amount, status, notes, notes_arabic, qr_code, created_at, updated_at FROM invoices ORDER BY created_at DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var inv Invoice
		err := rows.Scan(&inv.ID, &inv.InvoiceNumber, &inv.CustomerID, &inv.IssueDate, &inv.DueDate,
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
	query := `SELECT id, invoice_number, customer_id, issue_date, due_date, sub_total, vat_amount, total_amount, status, notes, notes_arabic, qr_code, created_at, updated_at FROM invoices WHERE id = ?`

	var inv Invoice
	err := d.db.QueryRow(query, id).Scan(&inv.ID, &inv.InvoiceNumber, &inv.CustomerID, &inv.IssueDate, &inv.DueDate,
		&inv.SubTotal, &inv.VATAmount, &inv.TotalAmount, &inv.Status, &inv.Notes, &inv.NotesArabic,
		&inv.QRCode, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Get customer
	customer, err := d.GetCustomerByID(inv.CustomerID)
	if err == nil {
		inv.Customer = customer
	}

	// Get invoice items
	items, err := d.GetInvoiceItems(inv.ID)
	if err == nil {
		inv.Items = items
	}

	return &inv, nil
}

func (d *Database) GetInvoiceItems(invoiceID int) ([]InvoiceItem, error) {
	query := `SELECT id, invoice_id, sale_item_id, quantity, unit_price, vat_rate, vat_amount, total_amount, created_at FROM invoice_items WHERE invoice_id = ?`

	rows, err := d.db.Query(query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []InvoiceItem
	for rows.Next() {
		var item InvoiceItem
		err := rows.Scan(&item.ID, &item.InvoiceID, &item.SaleItemID, &item.Quantity, &item.UnitPrice,
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

func (d *Database) UpdateInvoiceQRCode(invoiceID int, qrCode string) error {
	query := `UPDATE invoices SET qr_code = ? WHERE id = ?`
	_, err := d.db.Exec(query, qrCode, invoiceID)
	return err
}

func (d *Database) generateInvoiceNumber() string {
	var count int
	d.db.QueryRow("SELECT COUNT(*) FROM invoices").Scan(&count)
	return fmt.Sprintf("INV-%06d", count+1)
}
