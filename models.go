package main

import (
	"time"
)

// Customer represents a customer in the system
type Customer struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	NameArabic    string    `json:"name_arabic"`
	VATNumber     string    `json:"vat_number"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	AddressArabic string    `json:"address_arabic"`
	City          string    `json:"city"`
	CityArabic    string    `json:"city_arabic"`
	Country       string    `json:"country"`
	CountryArabic string    `json:"country_arabic"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ProductCategory represents a category for products
type ProductCategory struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	NameArabic        string    `json:"name_arabic"`
	Description       string    `json:"description"`
	DescriptionArabic string    `json:"description_arabic"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Product represents a product that can be sold
type Product struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	NameArabic        string    `json:"name_arabic"`
	Description       string    `json:"description"`
	DescriptionArabic string    `json:"description_arabic"`
	CategoryID        int       `json:"category_id"`
	Category          *ProductCategory `json:"category,omitempty"`
	UnitPrice         float64   `json:"unit_price"`
	VATRate           float64   `json:"vat_rate"`
	Unit              string    `json:"unit"`
	UnitArabic        string    `json:"unit_arabic"`
	SKU               string    `json:"sku"`
	Barcode           string    `json:"barcode"`
	Stock             int       `json:"stock"`
	MinStock          int       `json:"min_stock"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Invoice represents an invoice
type Invoice struct {
	ID            int           `json:"id"`
	InvoiceNumber string        `json:"invoice_number"`
	CustomerID    int           `json:"customer_id"`
	Customer      *Customer     `json:"customer,omitempty"`
	IssueDate     time.Time     `json:"issue_date"`
	DueDate       time.Time     `json:"due_date"`
	SubTotal      float64       `json:"sub_total"`
	VATAmount     float64       `json:"vat_amount"`
	TotalAmount   float64       `json:"total_amount"`
	Status        string        `json:"status"` // draft, sent, paid, cancelled
	Notes         string        `json:"notes"`
	NotesArabic   string        `json:"notes_arabic"`
	QRCode        string        `json:"qr_code"`
	Items         []InvoiceItem `json:"items,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// InvoiceItem represents an item in an invoice
type InvoiceItem struct {
	ID          int      `json:"id"`
	InvoiceID   int      `json:"invoice_id"`
	ProductID   int      `json:"product_id"`
	Product     *Product `json:"product,omitempty"`
	Quantity    float64  `json:"quantity"`
	UnitPrice   float64  `json:"unit_price"`
	VATRate     float64  `json:"vat_rate"`
	VATAmount   float64  `json:"vat_amount"`
	TotalAmount float64  `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

// Company represents the seller company information
type Company struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	NameArabic    string `json:"name_arabic"`
	VATNumber     string `json:"vat_number"`
	CRNumber      string `json:"cr_number"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	AddressArabic string `json:"address_arabic"`
	City          string `json:"city"`
	CityArabic    string `json:"city_arabic"`
	Country       string `json:"country"`
	CountryArabic string `json:"country_arabic"`
	Logo          string `json:"logo"`
}

// PaymentType represents different payment methods
type PaymentType struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	NameArabic  string    `json:"name_arabic"`
	Code        string    `json:"code"` // cash, card, bank_transfer, etc.
	IsDefault   bool      `json:"is_default"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SalesCategory represents categories for sales/invoices
type SalesCategory struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	NameArabic        string    `json:"name_arabic"`
	Code              string    `json:"code"` // retail, wholesale, service, etc.
	Description       string    `json:"description"`
	DescriptionArabic string    `json:"description_arabic"`
	IsDefault         bool      `json:"is_default"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TaxRate represents different tax rates
type TaxRate struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	NameArabic  string    `json:"name_arabic"`
	Rate        float64   `json:"rate"`
	IsDefault   bool      `json:"is_default"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UnitOfMeasurement represents units for products
type UnitOfMeasurement struct {
	ID         int       `json:"id"`
	Value      string    `json:"value"`      // kg, pcs, m, etc.
	Label      string    `json:"label"`      // Kilograms, Pieces, Meters
	Arabic     string    `json:"arabic"`     // Arabic translation
	IsDefault  bool      `json:"is_default"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DefaultProductSettings represents default settings for creating new products
type DefaultProductSettings struct {
	ID                    int     `json:"id"`
	DefaultStock          int     `json:"default_stock"`           // Default stock quantity
	DefaultTaxRateID      int     `json:"default_tax_rate_id"`     // Default tax rate ID
	DefaultUnitID         int     `json:"default_unit_id"`         // Default unit of measurement ID
	DefaultPaymentTypeID  int     `json:"default_payment_type_id"` // Default payment type ID
	DefaultSalesCategoryID int    `json:"default_sales_category_id"` // Default sales category ID
	DefaultProductType    string  `json:"default_product_type"`    // "product" or "service"
	DefaultProductStatus  bool    `json:"default_product_status"`  // true = active, false = inactive
	DefaultMarkup         float64 `json:"default_markup"`          // Default markup percentage
	DefaultPriceIncludesTax bool  `json:"default_price_includes_tax"` // Whether default price includes tax
	DefaultPriceChangeAllowed bool `json:"default_price_change_allowed"` // Whether price changes are allowed
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
