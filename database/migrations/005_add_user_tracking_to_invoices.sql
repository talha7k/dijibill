-- Add created_by and updated_by fields to sales_invoices table
ALTER TABLE sales_invoices ADD COLUMN created_by INTEGER;
ALTER TABLE sales_invoices ADD COLUMN updated_by INTEGER;

-- Add created_by and updated_by fields to purchase_invoices table  
ALTER TABLE purchase_invoices ADD COLUMN created_by INTEGER;
ALTER TABLE purchase_invoices ADD COLUMN updated_by INTEGER;

-- Add foreign key constraints (SQLite doesn't support adding foreign keys to existing tables,
-- but we'll handle this in the application logic)

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_sales_invoices_created_by ON sales_invoices(created_by);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_updated_by ON sales_invoices(updated_by);
CREATE INDEX IF NOT EXISTS idx_purchase_invoices_created_by ON purchase_invoices(created_by);
CREATE INDEX IF NOT EXISTS idx_purchase_invoices_updated_by ON purchase_invoices(updated_by);