package database

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateProductCategory creates a new product category
func (db *DB) CreateProductCategory(category *ProductCategory) error {
	query := `
		INSERT INTO product_categories (name, name_arabic, description, description_arabic, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now
	
	result, err := db.conn.Exec(query, 
		category.Name, 
		category.NameArabic, 
		category.Description, 
		category.DescriptionArabic, 
		category.CreatedAt, 
		category.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create product category: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	
	category.ID = int(id)
	return nil
}

// GetProductCategories retrieves all product categories
func (db *DB) GetProductCategories() ([]ProductCategory, error) {
	query := `
		SELECT id, name, name_arabic, description, description_arabic, created_at, updated_at
		FROM product_categories
		ORDER BY name ASC
	`
	
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query product categories: %w", err)
	}
	defer rows.Close()
	
	var categories []ProductCategory
	for rows.Next() {
		var category ProductCategory
		var nameArabic, description, descriptionArabic sql.NullString
		
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&nameArabic,
			&description,
			&descriptionArabic,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product category: %w", err)
		}
		
		category.NameArabic = nameArabic.String
		category.Description = description.String
		category.DescriptionArabic = descriptionArabic.String
		
		categories = append(categories, category)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product categories: %w", err)
	}
	
	return categories, nil
}

// GetProductCategoryByID retrieves a product category by its ID
func (db *DB) GetProductCategoryByID(id int) (*ProductCategory, error) {
	query := `
		SELECT id, name, name_arabic, description, description_arabic, created_at, updated_at
		FROM product_categories
		WHERE id = ?
	`
	
	var category ProductCategory
	var nameArabic, description, descriptionArabic sql.NullString
	
	err := db.conn.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&nameArabic,
		&description,
		&descriptionArabic,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product category with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get product category: %w", err)
	}
	
	category.NameArabic = nameArabic.String
	category.Description = description.String
	category.DescriptionArabic = descriptionArabic.String
	
	return &category, nil
}

// UpdateProductCategory updates an existing product category
func (db *DB) UpdateProductCategory(category *ProductCategory) error {
	query := `
		UPDATE product_categories 
		SET name = ?, name_arabic = ?, description = ?, description_arabic = ?, updated_at = ?
		WHERE id = ?
	`
	
	category.UpdatedAt = time.Now()
	
	result, err := db.conn.Exec(query,
		category.Name,
		category.NameArabic,
		category.Description,
		category.DescriptionArabic,
		category.UpdatedAt,
		category.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update product category: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("product category with id %d not found", category.ID)
	}
	
	return nil
}

// DeleteProductCategory deletes a product category by ID
func (db *DB) DeleteProductCategory(id int) error {
	// First check if any products are using this category
	var count int
	checkQuery := "SELECT COUNT(*) FROM products WHERE category_id = ?"
	err := db.conn.QueryRow(checkQuery, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check product category usage: %w", err)
	}
	
	if count > 0 {
		return fmt.Errorf("cannot delete product category: %d products are using this category", count)
	}
	
	query := "DELETE FROM product_categories WHERE id = ?"
	result, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product category: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("product category with id %d not found", id)
	}
	
	return nil
}

// GetProductCategoriesWithProductCount retrieves all product categories with their product counts
func (db *DB) GetProductCategoriesWithProductCount() ([]ProductCategoryWithCount, error) {
	query := `
		SELECT 
			pc.id, 
			pc.name, 
			pc.name_arabic, 
			pc.description, 
			pc.description_arabic, 
			pc.created_at, 
			pc.updated_at,
			COUNT(p.id) as product_count
		FROM product_categories pc
		LEFT JOIN products p ON pc.id = p.category_id
		GROUP BY pc.id, pc.name, pc.name_arabic, pc.description, pc.description_arabic, pc.created_at, pc.updated_at
		ORDER BY pc.name ASC
	`
	
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query product categories with count: %w", err)
	}
	defer rows.Close()
	
	var categories []ProductCategoryWithCount
	for rows.Next() {
		var category ProductCategoryWithCount
		var nameArabic, description, descriptionArabic sql.NullString
		
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&nameArabic,
			&description,
			&descriptionArabic,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.ProductCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product category with count: %w", err)
		}
		
		category.NameArabic = nameArabic.String
		category.Description = description.String
		category.DescriptionArabic = descriptionArabic.String
		
		categories = append(categories, category)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product categories with count: %w", err)
	}
	
	return categories, nil
}