package database

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Date is a custom type that can handle both "YYYY-MM-DD" and RFC3339 formats
type Date struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler interface to handle date parsing
func (d *Date) UnmarshalJSON(data []byte) error {
	// Remove quotes from JSON string
	str := strings.Trim(string(data), `"`)
	
	if str == "null" || str == "" {
		d.Time = time.Time{}
		return nil
	}
	
	// Try parsing as "YYYY-MM-DD" format first (from frontend)
	if t, err := time.Parse("2006-01-02", str); err == nil {
		d.Time = t
		return nil
	}
	
	// Try parsing as RFC3339 format (standard Go format)
	if t, err := time.Parse(time.RFC3339, str); err == nil {
		d.Time = t
		return nil
	}
	
	// Try parsing as "YYYY-MM-DD HH:MM:SS" format
	if t, err := time.Parse("2006-01-02 15:04:05", str); err == nil {
		d.Time = t
		return nil
	}
	
	return fmt.Errorf("cannot parse date: %s", str)
}

// MarshalJSON implements json.Marshaler interface
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Time.Format("2006-01-02"))
}

// Customer represents a customer in the system
type Customer struct {
	ID            int       `json:"id"`
	CompanyID     int       `json:"company_id"`  // Customer belongs to a company
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

// Supplier represents a supplier in the system
type Supplier struct {
	ID                  int       `json:"id"`
	CompanyID           int       `json:"company_id"`  // Supplier belongs to a company
	CompanyName         string    `json:"company_name"`
	CompanyNameArabic   string    `json:"company_name_arabic"`
	ContactPerson       string    `json:"contact_person"`
	ContactPersonArabic string    `json:"contact_person_arabic"`
	VATNumber           string    `json:"vat_number"`
	Email               string    `json:"email"`
	Phone               string    `json:"phone"`
	Address             string    `json:"address"`
	AddressArabic       string    `json:"address_arabic"`
	City                string    `json:"city"`
	CityArabic          string    `json:"city_arabic"`
	Country             string    `json:"country"`
	CountryArabic       string    `json:"country_arabic"`
	PaymentTerms        string    `json:"payment_terms"` // net_15, net_30, net_45, net_60, cash_on_delivery, advance_payment
	Active              bool      `json:"active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ProductCategory represents a category for products
type ProductCategory struct {
	ID                int       `json:"id"`
	CompanyID         int       `json:"company_id"`  // Category belongs to a company
	Name              string    `json:"name"`
	NameArabic        string    `json:"name_arabic"`
	Description       string    `json:"description"`
	DescriptionArabic string    `json:"description_arabic"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ProductCategoryWithCount represents a product category with product count
type ProductCategoryWithCount struct {
	ProductCategory
	ProductCount int `json:"product_count"`
}

// Product represents a product that can be sold
type Product struct {
	ID                     int       `json:"id"`
	CompanyID              int       `json:"company_id"`              // Product belongs to a company
	Name                   string    `json:"name"`
	NameArabic             string    `json:"name_arabic"`
	Description            string    `json:"description"`
	DescriptionArabic      string    `json:"description_arabic"`
	CategoryID             int       `json:"category_id"`
	CategoryName           string    `json:"category_name"`           // Joined from product_categories table
	Category               *ProductCategory `json:"category,omitempty"`
	UnitPrice              float64   `json:"unit_price"`
	VATRate                float64   `json:"vat_rate"`
	Unit                   string    `json:"unit"`
	UnitArabic             string    `json:"unit_arabic"`
	SKU                    string    `json:"sku"`
	Barcode                string    `json:"barcode"`
	Stock                  int       `json:"stock"`
	MinStock               int       `json:"min_stock"`
	IsActive               bool      `json:"is_active"`
	Color                  string    `json:"color"`                   // Product color (optional)
	ImageURL               string    `json:"image_url"`               // Product image URL (optional)
	ServiceNotUsingStock   bool      `json:"service_not_using_stock"` // Whether this is a service that doesn't use stock
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// SalesInvoice represents a sales invoice (to customers)
type SalesInvoice struct {
	ID               int                `json:"id"`
	CompanyID        int                `json:"company_id"`  // Invoice belongs to a company
	InvoiceNumber    string             `json:"invoice_number"`
	CustomerID       int                `json:"customer_id"`
	Customer         *Customer          `json:"customer,omitempty"`
	SalesCategoryID  int                `json:"sales_category_id"`
	SalesCategory    *SalesCategory     `json:"sales_category,omitempty"`
	TableNumber      *string            `json:"table_number,omitempty"` // Optional table number for restaurant/cafe POS
	IssueDate        Date               `json:"issue_date"`
	DueDate          Date               `json:"due_date"`
	SubTotal         float64            `json:"sub_total"`
	VATAmount        float64            `json:"vat_amount"`
	TotalAmount      float64            `json:"total_amount"`
	Status           string             `json:"status"` // draft, sent, paid, cancelled
	Notes            string             `json:"notes"`
	NotesArabic      string             `json:"notes_arabic"`
	QRCode           string             `json:"qr_code"`
	Items            []SalesInvoiceItem `json:"items,omitempty"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

// SalesInvoiceItem represents an item in a sales invoice
type SalesInvoiceItem struct {
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

// PurchaseInvoice represents a purchase invoice (from suppliers)
type PurchaseInvoice struct {
	ID               int                   `json:"id"`
	CompanyID        int                   `json:"company_id"`  // Invoice belongs to a company
	InvoiceNumber    string                `json:"invoice_number"`
	SupplierID       int                   `json:"supplier_id"`
	Supplier         *Supplier             `json:"supplier,omitempty"`
	IssueDate        Date                  `json:"issue_date"`
	DueDate          Date                  `json:"due_date"`
	SubTotal         float64               `json:"sub_total"`
	VATAmount        float64               `json:"vat_amount"`
	VATRate          float64               `json:"vat_rate"`
	VATInclusive     bool                  `json:"vat_inclusive"`
	TotalAmount      float64               `json:"total_amount"`
	Status           string                `json:"status"` // draft, received, paid, cancelled
	Notes            string                `json:"notes"`
	NotesArabic      string                `json:"notes_arabic"`
	Items            []PurchaseInvoiceItem `json:"items,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}

// PurchaseInvoiceItem represents an item in a purchase invoice
type PurchaseInvoiceItem struct {
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

// Legacy Invoice type for backward compatibility (maps to SalesInvoice)
type Invoice = SalesInvoice
type InvoiceItem = SalesInvoiceItem

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
	CompanyID   int       `json:"company_id"`  // Payment type belongs to a company
	Name        string    `json:"name"`
	NameArabic  string    `json:"name_arabic"`
	Code        string    `json:"code"` // cash, card, bank_transfer, etc.
	IsDefault   bool      `json:"is_default"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Payment represents a payment made against an invoice
type Payment struct {
	ID            int          `json:"id"`
	CompanyID     int          `json:"company_id"`  // Payment belongs to a company
	InvoiceID     int          `json:"invoice_id"`
	Invoice       *Invoice     `json:"invoice,omitempty"`
	PaymentTypeID int          `json:"payment_type_id"`
	PaymentType   *PaymentType `json:"payment_type,omitempty"`
	Amount        float64      `json:"amount"`
	PaymentDate   time.Time    `json:"payment_date"`
	Reference     string       `json:"reference"`     // Check number, transaction ID, etc.
	Notes         string       `json:"notes"`
	NotesArabic   string       `json:"notes_arabic"`
	Status        string       `json:"status"`        // pending, completed, failed, cancelled
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// SalesCategory represents categories for sales/invoices
type SalesCategory struct {
	ID                int       `json:"id"`
	CompanyID         int       `json:"company_id"`  // Sales category belongs to a company
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
	CompanyID   int       `json:"company_id"`  // Tax rate belongs to a company
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
	CompanyID  int       `json:"company_id"`  // Unit belongs to a company
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
	ID                        int     `json:"id"`
	CompanyID                 int     `json:"company_id"`  // Settings belong to a company
	DefaultStock              int     `json:"default_stock"`               // Default stock quantity
	DefaultTaxRateID          int     `json:"default_tax_rate_id"`         // Default tax rate ID
	DefaultUnitID             int     `json:"default_unit_id"`             // Default unit of measurement ID
	DefaultProductType        string  `json:"default_product_type"`        // "product" or "service"
	DefaultProductStatus      bool    `json:"default_product_status"`      // true = active, false = inactive
	DefaultMarkup             float64 `json:"default_markup"`              // Default markup percentage
	DefaultPriceIncludesTax   bool    `json:"default_price_includes_tax"`  // Whether default price includes tax
	DefaultPriceChangeAllowed bool    `json:"default_price_change_allowed"` // Whether price changes are allowed
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

// PurchaseProductCategory represents categories for purchase products
type PurchaseProductCategory struct {
	ID                int       `json:"id"`
	CompanyID         int       `json:"company_id"`  // Category belongs to a company
	Name              string    `json:"name"`
	NameArabic        string    `json:"name_arabic"`
	Description       string    `json:"description"`
	DescriptionArabic string    `json:"description_arabic"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// PurchaseProduct represents products that can be purchased from suppliers
type PurchaseProduct struct {
	ID                int                       `json:"id"`
	CompanyID         int                       `json:"company_id"`  // Product belongs to a company
	Name              string                    `json:"name"`
	NameArabic        string                    `json:"name_arabic"`
	Description       string                    `json:"description"`
	DescriptionArabic string                    `json:"description_arabic"`
	CategoryID        int                       `json:"category_id"`
	CategoryName      string                    `json:"category_name"`           // Joined from purchase_product_categories table
	Category          *PurchaseProductCategory  `json:"category,omitempty"`
	UnitPrice         float64                   `json:"unit_price"`              // Expected purchase price
	VATRate           float64                   `json:"vat_rate"`
	Unit              string                    `json:"unit"`
	UnitArabic        string                    `json:"unit_arabic"`
	SKU               string                    `json:"sku"`
	Barcode           string                    `json:"barcode"`
	IsActive          bool                      `json:"is_active"`
	Notes             string                    `json:"notes"`
	NotesArabic       string                    `json:"notes_arabic"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
}

// User represents a user in the system
type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"` // This should be hashed
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Role        string    `json:"role"`        // admin, manager, user
	IsActive    bool      `json:"is_active"`
	CompanyID   int       `json:"company_id"`  // User belongs to a company
	Company     *Company  `json:"company,omitempty"`
	LastLogin   *time.Time `json:"last_login,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SystemSettings represents system-wide configuration settings
type SystemSettings struct {
	ID               int        `json:"id"`
	CompanyID        int        `json:"company_id"`  // Settings per company
	Currency         string     `json:"currency"`
	Language         string     `json:"language"`
	Timezone         string     `json:"timezone"`
	DateFormat       string     `json:"date_format"`
	InvoiceLanguage  string     `json:"invoice_language"`
	ZatcaEnabled     bool       `json:"zatca_enabled"`
	AutoBackup       bool       `json:"auto_backup"`
	BackupFrequency  string     `json:"backup_frequency"`
	LastBackupTime   *time.Time `json:"last_backup_time,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}