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

// ItemCategory represents a category for items
type ItemCategory struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	NameArabic        string    `json:"name_arabic"`
	Description       string    `json:"description"`
	DescriptionArabic string    `json:"description_arabic"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// SaleItem represents an item that can be sold
type SaleItem struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	NameArabic        string    `json:"name_arabic"`
	Description       string    `json:"description"`
	DescriptionArabic string    `json:"description_arabic"`
	CategoryID        int       `json:"category_id"`
	UnitPrice         float64   `json:"unit_price"`
	VATRate           float64   `json:"vat_rate"`
	Unit              string    `json:"unit"`
	UnitArabic        string    `json:"unit_arabic"`
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
	ID          int       `json:"id"`
	InvoiceID   int       `json:"invoice_id"`
	SaleItemID  int       `json:"sale_item_id"`
	SaleItem    *SaleItem `json:"sale_item,omitempty"`
	Quantity    float64   `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	VATRate     float64   `json:"vat_rate"`
	VATAmount   float64   `json:"vat_amount"`
	TotalAmount float64   `json:"total_amount"`
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
