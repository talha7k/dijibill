package database

import (
	"fmt"
	"time"
)

// PurchaseInvoice operations
func (d *Database) CreatePurchaseInvoice(invoice *PurchaseInvoice) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Generate invoice number if not provided
	if invoice.InvoiceNumber == "" {
		invoice.InvoiceNumber = d.generatePurchaseInvoiceNumber()
	}

	// Insert purchase invoice
	query := `
		INSERT INTO purchase_invoices (invoice_number, supplier_id, issue_date, due_date, sub_total, vat_amount, vat_rate, vat_inclusive, total_amount, status, notes, notes_arabic, created_by, updated_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := tx.Exec(query, invoice.InvoiceNumber, invoice.SupplierID, invoice.IssueDate.Time, invoice.DueDate.Time,
		invoice.SubTotal, invoice.VATAmount, invoice.VATRate, invoice.VATInclusive, invoice.TotalAmount, invoice.Status, invoice.Notes, invoice.NotesArabic, invoice.CreatedBy, invoice.CreatedBy)
	if err != nil {
		return err
	}

	invoiceID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	invoice.ID = int(invoiceID)

	// Insert purchase invoice items
	for _, item := range invoice.Items {
		itemQuery := `
			INSERT INTO purchase_invoice_items (invoice_id, product_id, quantity, unit_price, vat_rate, vat_amount, total_amount)
			VALUES (?, ?, ?, ?, ?, ?, ?)`

		_, execErr := tx.Exec(itemQuery, invoiceID, item.ProductID, item.Quantity, item.UnitPrice, item.VATRate, item.VATAmount, item.TotalAmount)
		if execErr != nil {
			return execErr
		}
	}

	return tx.Commit()
}

func (d *Database) GetPurchaseInvoices() ([]PurchaseInvoice, error) {
	query := `
		SELECT 
			pi.id, pi.invoice_number, pi.supplier_id, pi.issue_date, pi.due_date, 
			pi.sub_total, pi.vat_amount, pi.vat_rate, pi.vat_inclusive, pi.total_amount, pi.status, pi.notes, pi.notes_arabic, 
			pi.created_at, pi.updated_at, pi.created_by, pi.updated_by,
			s.id, s.company_name, s.contact_person, s.email, s.phone, s.address, s.vat_number
		FROM purchase_invoices pi
		LEFT JOIN suppliers s ON pi.supplier_id = s.id
		ORDER BY pi.created_at DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []PurchaseInvoice
	for rows.Next() {
		var inv PurchaseInvoice
		var supplier Supplier
		var supplierID interface{}
		var issueDate, dueDate time.Time
		var vatAmount float64
		
		scanErr := rows.Scan(
			&inv.ID, &inv.InvoiceNumber, &inv.SupplierID, 
			&issueDate, &dueDate, &inv.SubTotal, &vatAmount, &inv.VATRate, &inv.VATInclusive, &inv.TotalAmount, 
			&inv.Status, &inv.Notes, &inv.NotesArabic, &inv.CreatedAt, &inv.UpdatedAt, &inv.CreatedBy, &inv.UpdatedBy,
			&supplierID, &supplier.CompanyName, &supplier.ContactPerson, &supplier.Email, &supplier.Phone, 
			&supplier.Address, &supplier.VATNumber)
		if scanErr != nil {
			return nil, scanErr
		}
		
		// Convert time.Time to custom Date type
		inv.IssueDate = Date{Time: issueDate}
		inv.DueDate = Date{Time: dueDate}
		
		// Assign float64 directly
		inv.VATAmount = vatAmount
		
		// Only assign supplier if it exists
		if supplierID != nil {
			supplier.ID = int(supplierID.(int64))
			inv.Supplier = &supplier
		}
		
		invoices = append(invoices, inv)
	}
	return invoices, nil
}

func (d *Database) GetPurchaseInvoiceByID(id int) (*PurchaseInvoice, error) {
	query := `SELECT id, invoice_number, supplier_id, issue_date, due_date, sub_total, vat_amount, vat_rate, vat_inclusive, total_amount, status, notes, notes_arabic, created_at, updated_at, created_by, updated_by FROM purchase_invoices WHERE id = ?`

	var inv PurchaseInvoice
	var issueDate, dueDate time.Time
	var vatAmount float64
	err := d.db.QueryRow(query, id).Scan(&inv.ID, &inv.InvoiceNumber, &inv.SupplierID, &issueDate, &dueDate,
		&inv.SubTotal, &vatAmount, &inv.VATRate, &inv.VATInclusive, &inv.TotalAmount, &inv.Status, &inv.Notes, &inv.NotesArabic,
		&inv.CreatedAt, &inv.UpdatedAt, &inv.CreatedBy, &inv.UpdatedBy)
	if err != nil {
		return nil, err
	}

	// Convert time.Time to custom Date type
	inv.IssueDate = Date{Time: issueDate}
	inv.DueDate = Date{Time: dueDate}
	
	// Assign float64 directly
	inv.VATAmount = vatAmount

	// Get supplier only if SupplierID is not 0
	if inv.SupplierID > 0 {
		supplier, supplierErr := d.GetSupplierByID(inv.SupplierID)
		if supplierErr == nil {
			inv.Supplier = supplier
		}
		// If supplier retrieval fails, inv.Supplier remains nil
	}

	// Get purchase invoice items
	items, itemsErr := d.GetPurchaseInvoiceItems(inv.ID)
	if itemsErr == nil {
		inv.Items = items
	}

	return &inv, nil
}

func (d *Database) UpdatePurchaseInvoice(invoice *PurchaseInvoice) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update purchase invoice
	query := `
		UPDATE purchase_invoices 
		SET invoice_number = ?, supplier_id = ?, issue_date = ?, due_date = ?, 
		    sub_total = ?, vat_amount = ?, vat_rate = ?, vat_inclusive = ?, total_amount = ?, status = ?, notes = ?, notes_arabic = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err = tx.Exec(query, invoice.InvoiceNumber, invoice.SupplierID, invoice.IssueDate.Time, invoice.DueDate.Time,
		invoice.SubTotal, invoice.VATAmount, invoice.VATRate, invoice.VATInclusive, invoice.TotalAmount, invoice.Status, invoice.Notes, invoice.NotesArabic, invoice.ID)
	if err != nil {
		return err
	}

	// Delete existing items
	_, err = tx.Exec("DELETE FROM purchase_invoice_items WHERE invoice_id = ?", invoice.ID)
	if err != nil {
		return err
	}

	// Insert updated items
	for _, item := range invoice.Items {
		itemQuery := `
			INSERT INTO purchase_invoice_items (invoice_id, product_id, quantity, unit_price, vat_rate, vat_amount, total_amount)
			VALUES (?, ?, ?, ?, ?, ?, ?)`

		_, execErr := tx.Exec(itemQuery, invoice.ID, item.ProductID, item.Quantity, item.UnitPrice, item.VATRate, item.VATAmount, item.TotalAmount)
		if execErr != nil {
			return execErr
		}
	}

	return tx.Commit()
}

func (d *Database) DeletePurchaseInvoice(id int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete purchase invoice items first (due to foreign key constraint)
	_, err = tx.Exec("DELETE FROM purchase_invoice_items WHERE invoice_id = ?", id)
	if err != nil {
		return err
	}

	// Delete purchase invoice
	_, err = tx.Exec("DELETE FROM purchase_invoices WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *Database) GetPurchaseInvoiceItems(invoiceID int) ([]PurchaseInvoiceItem, error) {
	query := `SELECT id, invoice_id, product_id, quantity, unit_price, vat_rate, vat_amount, total_amount, created_at FROM purchase_invoice_items WHERE invoice_id = ?`

	rows, err := d.db.Query(query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []PurchaseInvoiceItem
	for rows.Next() {
		var item PurchaseInvoiceItem
		scanErr := rows.Scan(&item.ID, &item.InvoiceID, &item.ProductID, &item.Quantity, &item.UnitPrice,
			&item.VATRate, &item.VATAmount, &item.TotalAmount, &item.CreatedAt)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	return items, nil
}

func (d *Database) generatePurchaseInvoiceNumber() string {
	var count int
	d.db.QueryRow("SELECT COUNT(*) FROM purchase_invoices").Scan(&count)
	return fmt.Sprintf("PI-%06d", count+1)
}