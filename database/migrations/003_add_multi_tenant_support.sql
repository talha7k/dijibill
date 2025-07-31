-- Migration 003: Add multi-tenant support
-- This migration adds the users table and company_id columns to all relevant tables

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user', -- admin, manager, user
    is_active BOOLEAN NOT NULL DEFAULT 1,
    company_id INTEGER NOT NULL,
    last_login DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);

-- Add company_id to existing tables
ALTER TABLE customers ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE suppliers ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE product_categories ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE products ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE sales_invoices ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE purchase_invoices ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE payment_types ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE payments ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE sales_categories ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE tax_rates ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE units_of_measurement ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE default_product_settings ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE purchase_product_categories ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE purchase_products ADD COLUMN company_id INTEGER DEFAULT 1;
ALTER TABLE system_settings ADD COLUMN company_id INTEGER DEFAULT 1;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE INDEX IF NOT EXISTS idx_customers_company_id ON customers(company_id);
CREATE INDEX IF NOT EXISTS idx_suppliers_company_id ON suppliers(company_id);
CREATE INDEX IF NOT EXISTS idx_product_categories_company_id ON product_categories(company_id);
CREATE INDEX IF NOT EXISTS idx_products_company_id ON products(company_id);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_company_id ON sales_invoices(company_id);
CREATE INDEX IF NOT EXISTS idx_purchase_invoices_company_id ON purchase_invoices(company_id);
CREATE INDEX IF NOT EXISTS idx_payment_types_company_id ON payment_types(company_id);
CREATE INDEX IF NOT EXISTS idx_payments_company_id ON payments(company_id);
CREATE INDEX IF NOT EXISTS idx_sales_categories_company_id ON sales_categories(company_id);
CREATE INDEX IF NOT EXISTS idx_tax_rates_company_id ON tax_rates(company_id);
CREATE INDEX IF NOT EXISTS idx_units_of_measurement_company_id ON units_of_measurement(company_id);
CREATE INDEX IF NOT EXISTS idx_default_product_settings_company_id ON default_product_settings(company_id);
CREATE INDEX IF NOT EXISTS idx_purchase_product_categories_company_id ON purchase_product_categories(company_id);
CREATE INDEX IF NOT EXISTS idx_purchase_products_company_id ON purchase_products(company_id);
CREATE INDEX IF NOT EXISTS idx_system_settings_company_id ON system_settings(company_id);

-- Create a default admin user for the first company
INSERT OR IGNORE INTO users (
    username, 
    email, 
    password, 
    first_name, 
    last_name, 
    role, 
    is_active, 
    company_id
) VALUES (
    'admin', 
    'admin@company.com', 
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: "password"
    'Admin', 
    'User', 
    'admin', 
    1, 
    1
);

-- Update existing data to belong to company 1
UPDATE customers SET company_id = 1 WHERE company_id IS NULL;
UPDATE suppliers SET company_id = 1 WHERE company_id IS NULL;
UPDATE product_categories SET company_id = 1 WHERE company_id IS NULL;
UPDATE products SET company_id = 1 WHERE company_id IS NULL;
UPDATE sales_invoices SET company_id = 1 WHERE company_id IS NULL;
UPDATE purchase_invoices SET company_id = 1 WHERE company_id IS NULL;
UPDATE payment_types SET company_id = 1 WHERE company_id IS NULL;
UPDATE payments SET company_id = 1 WHERE company_id IS NULL;
UPDATE sales_categories SET company_id = 1 WHERE company_id IS NULL;
UPDATE tax_rates SET company_id = 1 WHERE company_id IS NULL;
UPDATE units_of_measurement SET company_id = 1 WHERE company_id IS NULL;
UPDATE default_product_settings SET company_id = 1 WHERE company_id IS NULL;
UPDATE purchase_product_categories SET company_id = 1 WHERE company_id IS NULL;
UPDATE purchase_products SET company_id = 1 WHERE company_id IS NULL;
UPDATE system_settings SET company_id = 1 WHERE company_id IS NULL;