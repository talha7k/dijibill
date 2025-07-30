package database

import (
	"fmt"
)

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

	// Get customer only if CustomerID is not 0
	if inv.CustomerID > 0 {
		customer, err := d.GetCustomerByID(inv.CustomerID)
		if err == nil {
			inv.Customer = customer
		}
		// If customer retrieval fails, inv.Customer remains nil which is handled in PDF generation
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



func (d *Database) generateInvoiceNumber() string {
	var count int
	d.db.QueryRow("SELECT COUNT(*) FROM invoices").Scan(&count)
	return fmt.Sprintf("INV-%06d", count+1)
}