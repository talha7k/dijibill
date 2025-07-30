package database

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