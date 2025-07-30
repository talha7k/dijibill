package database

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