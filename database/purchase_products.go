package database

import (
	"database/sql"
	"fmt"
	"time"
)

// GetPurchaseProductCategories retrieves all purchase product categories
func (d *Database) GetPurchaseProductCategories() ([]PurchaseProductCategory, error) {
	query := `
		SELECT id, name, name_arabic, description, description_arabic, is_active, created_at, updated_at
		FROM purchase_product_categories
		WHERE is_active = 1
		ORDER BY name ASC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying purchase product categories: %v", err)
	}
	defer rows.Close()

	var categories []PurchaseProductCategory
	for rows.Next() {
		var category PurchaseProductCategory
		var nameArabic, description, descriptionArabic sql.NullString

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&nameArabic,
			&description,
			&descriptionArabic,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning purchase product category: %v", err)
		}

		category.NameArabic = nameArabic.String
		category.Description = description.String
		category.DescriptionArabic = descriptionArabic.String

		categories = append(categories, category)
	}

	return categories, nil
}

// CreatePurchaseProductCategory creates a new purchase product category
func (d *Database) CreatePurchaseProductCategory(category PurchaseProductCategory) (*PurchaseProductCategory, error) {
	query := `
		INSERT INTO purchase_product_categories (name, name_arabic, description, description_arabic, is_active)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(query,
		category.Name,
		nullString(category.NameArabic),
		nullString(category.Description),
		nullString(category.DescriptionArabic),
		category.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating purchase product category: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %v", err)
	}

	category.ID = int(id)
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	return &category, nil
}

// UpdatePurchaseProductCategory updates an existing purchase product category
func (d *Database) UpdatePurchaseProductCategory(category PurchaseProductCategory) error {
	query := `
		UPDATE purchase_product_categories 
		SET name = ?, name_arabic = ?, description = ?, description_arabic = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := d.db.Exec(query,
		category.Name,
		nullString(category.NameArabic),
		nullString(category.Description),
		nullString(category.DescriptionArabic),
		category.IsActive,
		category.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating purchase product category: %v", err)
	}

	return nil
}

// DeletePurchaseProductCategory soft deletes a purchase product category
func (d *Database) DeletePurchaseProductCategory(id int) error {
	query := `UPDATE purchase_product_categories SET is_active = 0, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting purchase product category: %v", err)
	}
	return nil
}

// GetPurchaseProducts retrieves all purchase products with their categories
func (d *Database) GetPurchaseProducts() ([]PurchaseProduct, error) {
	query := `
		SELECT 
			pp.id, pp.name, pp.name_arabic, pp.description, pp.description_arabic,
			pp.category_id, pp.unit_price, pp.vat_rate, pp.unit, pp.unit_arabic,
			pp.sku, pp.barcode, pp.is_active, pp.notes, pp.notes_arabic,
			pp.created_at, pp.updated_at,
			ppc.name as category_name
		FROM purchase_products pp
		LEFT JOIN purchase_product_categories ppc ON pp.category_id = ppc.id
		WHERE pp.is_active = 1
		ORDER BY pp.name ASC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying purchase products: %v", err)
	}
	defer rows.Close()

	var products []PurchaseProduct
	for rows.Next() {
		var product PurchaseProduct
		var nameArabic, description, descriptionArabic sql.NullString
		var categoryID sql.NullInt64
		var unitArabic, sku, barcode, notes, notesArabic sql.NullString
		var categoryName sql.NullString

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&nameArabic,
			&description,
			&descriptionArabic,
			&categoryID,
			&product.UnitPrice,
			&product.VATRate,
			&product.Unit,
			&unitArabic,
			&sku,
			&barcode,
			&product.IsActive,
			&notes,
			&notesArabic,
			&product.CreatedAt,
			&product.UpdatedAt,
			&categoryName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning purchase product: %v", err)
		}

		product.NameArabic = nameArabic.String
		product.Description = description.String
		product.DescriptionArabic = descriptionArabic.String
		if categoryID.Valid {
			product.CategoryID = int(categoryID.Int64)
		}
		product.UnitArabic = unitArabic.String
		product.SKU = sku.String
		product.Barcode = barcode.String
		product.Notes = notes.String
		product.NotesArabic = notesArabic.String
		product.CategoryName = categoryName.String

		products = append(products, product)
	}

	return products, nil
}

// CreatePurchaseProduct creates a new purchase product
func (d *Database) CreatePurchaseProduct(product PurchaseProduct) (*PurchaseProduct, error) {
	query := `
		INSERT INTO purchase_products (
			name, name_arabic, description, description_arabic, category_id,
			unit_price, vat_rate, unit, unit_arabic, sku, barcode,
			is_active, notes, notes_arabic
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(query,
		product.Name,
		nullString(product.NameArabic),
		nullString(product.Description),
		nullString(product.DescriptionArabic),
		nullInt(product.CategoryID),
		product.UnitPrice,
		product.VATRate,
		product.Unit,
		nullString(product.UnitArabic),
		nullString(product.SKU),
		nullString(product.Barcode),
		product.IsActive,
		nullString(product.Notes),
		nullString(product.NotesArabic),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating purchase product: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %v", err)
	}

	product.ID = int(id)
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	return &product, nil
}

// UpdatePurchaseProduct updates an existing purchase product
func (d *Database) UpdatePurchaseProduct(product PurchaseProduct) error {
	query := `
		UPDATE purchase_products 
		SET name = ?, name_arabic = ?, description = ?, description_arabic = ?, 
			category_id = ?, unit_price = ?, vat_rate = ?, unit = ?, unit_arabic = ?,
			sku = ?, barcode = ?, is_active = ?, notes = ?, notes_arabic = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := d.db.Exec(query,
		product.Name,
		nullString(product.NameArabic),
		nullString(product.Description),
		nullString(product.DescriptionArabic),
		nullInt(product.CategoryID),
		product.UnitPrice,
		product.VATRate,
		product.Unit,
		nullString(product.UnitArabic),
		nullString(product.SKU),
		nullString(product.Barcode),
		product.IsActive,
		nullString(product.Notes),
		nullString(product.NotesArabic),
		product.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating purchase product: %v", err)
	}

	return nil
}

// DeletePurchaseProduct soft deletes a purchase product
func (d *Database) DeletePurchaseProduct(id int) error {
	query := `UPDATE purchase_products SET is_active = 0, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting purchase product: %v", err)
	}
	return nil
}

// Helper functions for null handling
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullInt(i int) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(i), Valid: true}
}
