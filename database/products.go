package database

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

func (d *Database) GetProductByID(id int) (*Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.name_arabic, p.description, p.description_arabic, 
			p.category_id, COALESCE(pc.name, '') as category_name,
			p.unit_price, p.vat_rate, p.unit, p.unit_arabic, 
			p.sku, p.barcode, p.stock, p.min_stock, p.is_active, 
			p.created_at, p.updated_at
		FROM products p
		LEFT JOIN product_categories pc ON p.category_id = pc.id
		WHERE p.id = ?`

	var product Product
	err := d.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.NameArabic, &product.Description, &product.DescriptionArabic,
		&product.CategoryID, &product.CategoryName, &product.UnitPrice, &product.VATRate, &product.Unit, &product.UnitArabic,
		&product.SKU, &product.Barcode, &product.Stock, &product.MinStock, &product.IsActive, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	
	// Set default values for new fields that don't exist in the database yet
	product.Color = ""
	product.ImageURL = ""
	product.ServiceNotUsingStock = false
	
	return &product, nil
}

func (d *Database) DeleteProduct(id int) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := d.db.Exec(query, id)
	return err
}