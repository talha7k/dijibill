package database

import (
	"database/sql"
	"time"
)

// CreatePayment creates a new payment record
func (d *Database) CreatePayment(payment Payment) (Payment, error) {
	query := `
		INSERT INTO payments (invoice_id, payment_type_id, amount, payment_date, reference, notes, notes_arabic, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	result, err := d.db.Exec(query, payment.InvoiceID, payment.PaymentTypeID, payment.Amount, payment.PaymentDate, payment.Reference, payment.Notes, payment.NotesArabic, payment.Status, now, now)
	if err != nil {
		return Payment{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Payment{}, err
	}

	payment.ID = int(id)
	payment.CreatedAt = now
	payment.UpdatedAt = now

	return payment, nil
}

// GetPayments retrieves all payments with optional filtering
func (d *Database) GetPayments() ([]Payment, error) {
	query := `
		SELECT p.id, p.invoice_id, p.payment_type_id, p.amount, p.payment_date, p.reference, p.notes, p.notes_arabic, p.status, p.created_at, p.updated_at,
			   pt.name as payment_type_name, pt.name_arabic as payment_type_name_arabic, pt.code as payment_type_code
		FROM payments p
		LEFT JOIN payment_types pt ON p.payment_type_id = pt.id
		ORDER BY p.payment_date DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var payment Payment
		var paymentType PaymentType
		var paymentTypeName, paymentTypeNameArabic, paymentTypeCode sql.NullString

		err := rows.Scan(
			&payment.ID, &payment.InvoiceID, &payment.PaymentTypeID, &payment.Amount, &payment.PaymentDate,
			&payment.Reference, &payment.Notes, &payment.NotesArabic, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt,
			&paymentTypeName, &paymentTypeNameArabic, &paymentTypeCode,
		)
		if err != nil {
			return nil, err
		}

		if paymentTypeName.Valid {
			paymentType.ID = payment.PaymentTypeID
			paymentType.Name = paymentTypeName.String
			paymentType.NameArabic = paymentTypeNameArabic.String
			paymentType.Code = paymentTypeCode.String
			payment.PaymentType = &paymentType
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

// GetPaymentsByInvoiceID retrieves all payments for a specific invoice
func (d *Database) GetPaymentsByInvoiceID(invoiceID int) ([]Payment, error) {
	query := `
		SELECT p.id, p.invoice_id, p.payment_type_id, p.amount, p.payment_date, p.reference, p.notes, p.notes_arabic, p.status, p.created_at, p.updated_at,
			   pt.name as payment_type_name, pt.name_arabic as payment_type_name_arabic, pt.code as payment_type_code
		FROM payments p
		LEFT JOIN payment_types pt ON p.payment_type_id = pt.id
		WHERE p.invoice_id = ?
		ORDER BY p.payment_date DESC
	`

	rows, err := d.db.Query(query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var payment Payment
		var paymentType PaymentType
		var paymentTypeName, paymentTypeNameArabic, paymentTypeCode sql.NullString

		err := rows.Scan(
			&payment.ID, &payment.InvoiceID, &payment.PaymentTypeID, &payment.Amount, &payment.PaymentDate,
			&payment.Reference, &payment.Notes, &payment.NotesArabic, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt,
			&paymentTypeName, &paymentTypeNameArabic, &paymentTypeCode,
		)
		if err != nil {
			return nil, err
		}

		if paymentTypeName.Valid {
			paymentType.ID = payment.PaymentTypeID
			paymentType.Name = paymentTypeName.String
			paymentType.NameArabic = paymentTypeNameArabic.String
			paymentType.Code = paymentTypeCode.String
			payment.PaymentType = &paymentType
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

// GetPaymentByID retrieves a payment by its ID
func (d *Database) GetPaymentByID(id int) (Payment, error) {
	query := `
		SELECT p.id, p.invoice_id, p.payment_type_id, p.amount, p.payment_date, p.reference, p.notes, p.notes_arabic, p.status, p.created_at, p.updated_at,
			   pt.name as payment_type_name, pt.name_arabic as payment_type_name_arabic, pt.code as payment_type_code
		FROM payments p
		LEFT JOIN payment_types pt ON p.payment_type_id = pt.id
		WHERE p.id = ?
	`

	var payment Payment
	var paymentType PaymentType
	var paymentTypeName, paymentTypeNameArabic, paymentTypeCode sql.NullString

	err := d.db.QueryRow(query, id).Scan(
		&payment.ID, &payment.InvoiceID, &payment.PaymentTypeID, &payment.Amount, &payment.PaymentDate,
		&payment.Reference, &payment.Notes, &payment.NotesArabic, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt,
		&paymentTypeName, &paymentTypeNameArabic, &paymentTypeCode,
	)
	if err != nil {
		return Payment{}, err
	}

	if paymentTypeName.Valid {
		paymentType.ID = payment.PaymentTypeID
		paymentType.Name = paymentTypeName.String
		paymentType.NameArabic = paymentTypeNameArabic.String
		paymentType.Code = paymentTypeCode.String
		payment.PaymentType = &paymentType
	}

	return payment, nil
}

// UpdatePayment updates an existing payment
func (d *Database) UpdatePayment(payment Payment) error {
	query := `
		UPDATE payments 
		SET payment_type_id = ?, amount = ?, payment_date = ?, reference = ?, notes = ?, notes_arabic = ?, status = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := d.db.Exec(query, payment.PaymentTypeID, payment.Amount, payment.PaymentDate, payment.Reference, payment.Notes, payment.NotesArabic, payment.Status, time.Now(), payment.ID)
	return err
}

// DeletePayment deletes a payment by its ID
func (d *Database) DeletePayment(id int) error {
	query := `DELETE FROM payments WHERE id = ?`
	_, err := d.db.Exec(query, id)
	return err
}