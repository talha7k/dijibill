package database

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