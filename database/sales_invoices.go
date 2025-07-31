package database

import (
	"fmt"
	"time"
)

// SalesInvoice operations
func (d *Database) CreateSalesInvoice(invoice *SalesInvoice) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Generate invoice number if not provided
	if invoice.InvoiceNumber == "" {
		invoice.InvoiceNumber = d.generateSalesInvoiceNumber()
	}

	// Insert sales invoice
	query := `
		INSERT INTO sales_invoices (invoice_number, customer_id, sales_category_id, issue_date, due_date, sub_total, vat_amount, total_amount, status, notes, notes_arabic, qr_code)
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

	// Insert sales invoice items
	for _, item := range invoice.Items {
		itemQuery := `
			INSERT INTO sales_invoice_items (invoice_id, product_id, quantity, unit_price, vat_rate, vat_amount, total_amount)
			VALUES (?, ?, ?, ?, ?, ?, ?)`

		_, err = tx.Exec(itemQuery, invoiceID, item.ProductID, item.Quantity, item.UnitPrice, item.VATRate, item.VATAmount, item.TotalAmount)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *Database) GetSalesInvoices() ([]SalesInvoice, error) {
	query := `
		SELECT 
			si.id, si.invoice_number, si.customer_id, si.sales_category_id, 
			si.issue_date, si.due_date, si.sub_total, si.vat_amount, si.total_amount, 
			si.status, si.notes, si.notes_arabic, si.qr_code, si.created_at, si.updated_at,
			c.id as customer_id, c.name as customer_name, c.email as customer_email, 
			c.phone as customer_phone, c.address as customer_address, c.city as customer_city, 
			c.country as customer_country, c.vat_number as customer_vat_number, 
			c.created_at as customer_created_at, c.updated_at as customer_updated_at
		FROM sales_invoices si
		LEFT JOIN customers c ON si.customer_id = c.id
		ORDER BY si.created_at DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []SalesInvoice
	for rows.Next() {
		var inv SalesInvoice
		var customer Customer
		var customerID, customerCreatedAt, customerUpdatedAt interface{}
		
		scanErr := rows.Scan(
			&inv.ID, &inv.InvoiceNumber, &inv.CustomerID, &inv.SalesCategoryID, 
			&inv.IssueDate, &inv.DueDate, &inv.SubTotal, &inv.VATAmount, &inv.TotalAmount, 
			&inv.Status, &inv.Notes, &inv.NotesArabic, &inv.QRCode, &inv.CreatedAt, &inv.UpdatedAt,
			&customerID, &customer.Name, &customer.Email, &customer.Phone, &customer.Address, 
			&customer.City, &customer.Country, &customer.VATNumber, 
			&customerCreatedAt, &customerUpdatedAt)
		if scanErr != nil {
			return nil, scanErr
		}
		
		// Only set customer if customer data exists
		if customerID != nil {
			customer.ID = int(customerID.(int64))
			if customerCreatedAt != nil {
				customer.CreatedAt = customerCreatedAt.(time.Time)
			}
			if customerUpdatedAt != nil {
				customer.UpdatedAt = customerUpdatedAt.(time.Time)
			}
			inv.Customer = &customer
		}
		
		invoices = append(invoices, inv)
	}
	return invoices, nil
}

func (d *Database) GetSalesInvoiceByID(id int) (*SalesInvoice, error) {
	query := `SELECT id, invoice_number, customer_id, sales_category_id, issue_date, due_date, sub_total, vat_amount, total_amount, status, notes, notes_arabic, qr_code, created_at, updated_at FROM sales_invoices WHERE id = ?`

	var inv SalesInvoice
	err := d.db.QueryRow(query, id).Scan(&inv.ID, &inv.InvoiceNumber, &inv.CustomerID, &inv.SalesCategoryID, &inv.IssueDate, &inv.DueDate,
		&inv.SubTotal, &inv.VATAmount, &inv.TotalAmount, &inv.Status, &inv.Notes, &inv.NotesArabic,
		&inv.QRCode, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Get customer only if CustomerID is not 0
	if inv.CustomerID > 0 {
		customer, customerErr := d.GetCustomerByID(inv.CustomerID)
		if customerErr == nil {
			inv.Customer = customer
		}
		// If customer retrieval fails, inv.Customer remains nil which is handled in PDF generation
	}

	// Get sales category
	if inv.SalesCategoryID > 0 {
		salesCategory, categoryErr := d.GetSalesCategoryByID(inv.SalesCategoryID)
		if categoryErr == nil {
			inv.SalesCategory = salesCategory
		}
	}

	// Get sales invoice items
	items, itemsErr := d.GetSalesInvoiceItems(inv.ID)
	if itemsErr == nil {
		inv.Items = items
	}

	return &inv, nil
}

func (d *Database) UpdateSalesInvoice(invoice *SalesInvoice) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update sales invoice
	query := `
		UPDATE sales_invoices 
		SET invoice_number = ?, customer_id = ?, sales_category_id = ?, issue_date = ?, due_date = ?, 
		    sub_total = ?, vat_amount = ?, total_amount = ?, status = ?, notes = ?, notes_arabic = ?, qr_code = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err = tx.Exec(query, invoice.InvoiceNumber, invoice.CustomerID, invoice.SalesCategoryID, invoice.IssueDate, invoice.DueDate,
		invoice.SubTotal, invoice.VATAmount, invoice.TotalAmount, invoice.Status, invoice.Notes, invoice.NotesArabic, invoice.QRCode, invoice.ID)
	if err != nil {
		return err
	}

	// Delete existing items
	_, err = tx.Exec("DELETE FROM sales_invoice_items WHERE invoice_id = ?", invoice.ID)
	if err != nil {
		return err
	}

	// Insert updated items
	for _, item := range invoice.Items {
		itemQuery := `
			INSERT INTO sales_invoice_items (invoice_id, product_id, quantity, unit_price, vat_rate, vat_amount, total_amount)
			VALUES (?, ?, ?, ?, ?, ?, ?)`

		_, execErr := tx.Exec(itemQuery, invoice.ID, item.ProductID, item.Quantity, item.UnitPrice, item.VATRate, item.VATAmount, item.TotalAmount)
		if execErr != nil {
			return execErr
		}
	}

	return tx.Commit()
}

func (d *Database) DeleteSalesInvoice(id int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete sales invoice items first (due to foreign key constraint)
	_, err = tx.Exec("DELETE FROM sales_invoice_items WHERE invoice_id = ?", id)
	if err != nil {
		return err
	}

	// Delete sales invoice
	_, err = tx.Exec("DELETE FROM sales_invoices WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *Database) GetSalesInvoiceItems(invoiceID int) ([]SalesInvoiceItem, error) {
	query := `SELECT id, invoice_id, product_id, quantity, unit_price, vat_rate, vat_amount, total_amount, created_at FROM sales_invoice_items WHERE invoice_id = ?`

	rows, err := d.db.Query(query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []SalesInvoiceItem
	for rows.Next() {
		var item SalesInvoiceItem
		scanErr := rows.Scan(&item.ID, &item.InvoiceID, &item.ProductID, &item.Quantity, &item.UnitPrice,
			&item.VATRate, &item.VATAmount, &item.TotalAmount, &item.CreatedAt)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	return items, nil
}

func (d *Database) generateSalesInvoiceNumber() string {
	var count int
	d.db.QueryRow("SELECT COUNT(*) FROM sales_invoices").Scan(&count)
	return fmt.Sprintf("SI-%06d", count+1)
}

// Legacy functions for backward compatibility
func (d *Database) CreateInvoice(invoice *Invoice) error {
	return d.CreateSalesInvoice((*SalesInvoice)(invoice))
}

func (d *Database) GetInvoices() ([]Invoice, error) {
	salesInvoices, err := d.GetSalesInvoices()
	if err != nil {
		return nil, err
	}
	
	invoices := make([]Invoice, len(salesInvoices))
	for i, si := range salesInvoices {
		invoices[i] = Invoice(si)
	}
	return invoices, nil
}

func (d *Database) GetInvoiceByID(id int) (*Invoice, error) {
	salesInvoice, err := d.GetSalesInvoiceByID(id)
	if err != nil {
		return nil, err
	}
	return (*Invoice)(salesInvoice), nil
}

func (d *Database) GetInvoiceItems(invoiceID int) ([]InvoiceItem, error) {
	salesItems, err := d.GetSalesInvoiceItems(invoiceID)
	if err != nil {
		return nil, err
	}
	
	items := make([]InvoiceItem, len(salesItems))
	for i, si := range salesItems {
		items[i] = InvoiceItem(si)
	}
	return items, nil
}