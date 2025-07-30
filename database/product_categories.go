package database

import (
	"database/sql"
	"time"
)

// CreateProductCategory creates a new product category
func (d *Database) CreateProductCategory(category *ProductCategory) error {
	query := `
		INSERT INTO product_categories (name, name_arabic, description, description_arabic, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	result, err := d.db.Exec(query, category.Name, category.NameArabic, category.Description, category.DescriptionArabic, now, now)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	category.ID = int(id)
	category.CreatedAt = now
	category.UpdatedAt = now
	
	return nil
}

// GetProductCategories retrieves all product categories
func (d *Database) GetProductCategories() ([]ProductCategory, error) {
	query := `
		SELECT id, name, name_arabic, description, description_arabic, created_at, updated_at
		FROM product_categories
		ORDER BY name
	`
	
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var categories []ProductCategory
	for rows.Next() {
		var category ProductCategory
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.NameArabic,
			&category.Description,
			&category.DescriptionArabic,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	
	return categories, rows.Err()
}

// GetProductCategoryByID retrieves a product category by ID
func (d *Database) GetProductCategoryByID(id int) (*ProductCategory, error) {
	query := `
		SELECT id, name, name_arabic, description, description_arabic, created_at, updated_at
		FROM product_categories
		WHERE id = ?
	`
	
	var category ProductCategory
	err := d.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.NameArabic,
		&category.Description,
		&category.DescriptionArabic,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &category, nil
}

// UpdateProductCategory updates an existing product category
func (d *Database) UpdateProductCategory(category *ProductCategory) error {
	query := `
		UPDATE product_categories
		SET name = ?, name_arabic = ?, description = ?, description_arabic = ?, updated_at = ?
		WHERE id = ?
	`
	
	now := time.Now()
	_, err := d.db.Exec(query, category.Name, category.NameArabic, category.Description, category.DescriptionArabic, now, category.ID)
	if err != nil {
		return err
	}
	
	category.UpdatedAt = now
	return nil
}

// DeleteProductCategory deletes a product category
func (d *Database) DeleteProductCategory(id int) error {
	// First check if there are any products using this category
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM products WHERE category_id = ?", id).Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return sql.ErrConnDone // Return an error indicating category is in use
	}
	
	query := "DELETE FROM product_categories WHERE id = ?"
	_, err = d.db.Exec(query, id)
	return err
}

// GetProductCategoriesWithProductCount retrieves all product categories with their product counts
func (d *Database) GetProductCategoriesWithProductCount() ([]ProductCategoryWithCount, error) {
	query := `
		SELECT 
			pc.id, pc.name, pc.name_arabic, pc.description, pc.description_arabic, 
			pc.created_at, pc.updated_at,
			COALESCE(COUNT(p.id), 0) as product_count
		FROM product_categories pc
		LEFT JOIN products p ON pc.id = p.category_id
		GROUP BY pc.id, pc.name, pc.name_arabic, pc.description, pc.description_arabic, pc.created_at, pc.updated_at
		ORDER BY pc.name
	`
	
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var categories []ProductCategoryWithCount
	for rows.Next() {
		var category ProductCategoryWithCount
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.NameArabic,
			&category.Description,
			&category.DescriptionArabic,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.ProductCount,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	
	return categories, rows.Err()
}