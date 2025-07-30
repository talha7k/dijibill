package database

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