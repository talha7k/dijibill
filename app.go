package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	db         *Database
	pdfService *PDFService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get user home directory:", err)
	}

	dbPath := filepath.Join(homeDir, "dijinvoice.db")
	a.db, err = NewDatabase(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize PDF service
	a.pdfService = NewPDFService(a.db)

	log.Println("Application started successfully")
}

// Customer Management Methods

func (a *App) CreateCustomer(customer Customer) error {
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	return a.db.CreateCustomer(&customer)
}

func (a *App) GetCustomers() ([]Customer, error) {
	return a.db.GetCustomers()
}

func (a *App) GetCustomerByID(id int) (*Customer, error) {
	return a.db.GetCustomerByID(id)
}

func (a *App) UpdateCustomer(customer Customer) error {
	customer.UpdatedAt = time.Now()
	return a.db.UpdateCustomer(&customer)
}

func (a *App) DeleteCustomer(id int) error {
	return a.db.DeleteCustomer(id)
}

// Product Category Management Methods

func (a *App) CreateProductCategory(category ProductCategory) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	return a.db.CreateProductCategory(&category)
}

func (a *App) GetProductCategories() ([]ProductCategory, error) {
	return a.db.GetProductCategories()
}

// Product Management Methods

func (a *App) CreateProduct(product Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	return a.db.CreateProduct(&product)
}

func (a *App) GetProducts() ([]Product, error) {
	return a.db.GetProducts()
}

// Invoice Management Methods

func (a *App) CreateInvoice(invoice Invoice) error {
	// Calculate totals
	var subTotal, vatAmount, totalAmount float64

	for i := range invoice.Items {
		item := &invoice.Items[i]
		item.VATAmount = (item.UnitPrice * item.Quantity * item.VATRate) / 100
		item.TotalAmount = (item.UnitPrice * item.Quantity) + item.VATAmount

		subTotal += item.UnitPrice * item.Quantity
		vatAmount += item.VATAmount
		totalAmount += item.TotalAmount
	}

	invoice.SubTotal = subTotal
	invoice.VATAmount = vatAmount
	invoice.TotalAmount = totalAmount
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	if invoice.Status == "" {
		invoice.Status = "draft"
	}

	return a.db.CreateInvoice(&invoice)
}

func (a *App) GetInvoices() ([]Invoice, error) {
	return a.db.GetInvoices()
}

func (a *App) GetInvoiceByID(id int) (*Invoice, error) {
	return a.db.GetInvoiceByID(id)
}

// Company Management Methods

func (a *App) GetCompany() (*Company, error) {
	return a.db.GetCompany()
}

func (a *App) UpdateCompany(company Company) error {
	return a.db.UpdateCompany(&company)
}

// Payment Type Management Methods

func (a *App) CreatePaymentType(paymentType PaymentType) error {
	paymentType.CreatedAt = time.Now()
	paymentType.UpdatedAt = time.Now()
	return a.db.CreatePaymentType(&paymentType)
}

func (a *App) GetPaymentTypes() ([]PaymentType, error) {
	return a.db.GetPaymentTypes()
}

func (a *App) UpdatePaymentType(paymentType PaymentType) error {
	paymentType.UpdatedAt = time.Now()
	return a.db.UpdatePaymentType(&paymentType)
}

func (a *App) DeletePaymentType(id int) error {
	return a.db.DeletePaymentType(id)
}

// Sales Category Management Methods

func (a *App) CreateSalesCategory(salesCategory SalesCategory) error {
	salesCategory.CreatedAt = time.Now()
	salesCategory.UpdatedAt = time.Now()
	return a.db.CreateSalesCategory(&salesCategory)
}

func (a *App) GetSalesCategories() ([]SalesCategory, error) {
	return a.db.GetSalesCategories()
}

func (a *App) UpdateSalesCategory(salesCategory SalesCategory) error {
	salesCategory.UpdatedAt = time.Now()
	return a.db.UpdateSalesCategory(&salesCategory)
}

func (a *App) DeleteSalesCategory(id int) error {
	return a.db.DeleteSalesCategory(id)
}

// Tax Rate Management Methods

func (a *App) CreateTaxRate(taxRate TaxRate) error {
	taxRate.CreatedAt = time.Now()
	taxRate.UpdatedAt = time.Now()
	return a.db.CreateTaxRate(&taxRate)
}

func (a *App) GetTaxRates() ([]TaxRate, error) {
	return a.db.GetTaxRates()
}

func (a *App) UpdateTaxRate(taxRate TaxRate) error {
	taxRate.UpdatedAt = time.Now()
	return a.db.UpdateTaxRate(&taxRate)
}

func (a *App) DeleteTaxRate(id int) error {
	return a.db.DeleteTaxRate(id)
}

// Unit of Measurement Management Methods

func (a *App) CreateUnitOfMeasurement(unit UnitOfMeasurement) error {
	unit.CreatedAt = time.Now()
	unit.UpdatedAt = time.Now()
	return a.db.CreateUnitOfMeasurement(&unit)
}

func (a *App) GetUnitsOfMeasurement() ([]UnitOfMeasurement, error) {
	return a.db.GetUnitsOfMeasurement()
}

func (a *App) UpdateUnitOfMeasurement(unit UnitOfMeasurement) error {
	unit.UpdatedAt = time.Now()
	return a.db.UpdateUnitOfMeasurement(&unit)
}

func (a *App) DeleteUnitOfMeasurement(id int) error {
	return a.db.DeleteUnitOfMeasurement(id)
}

// Default Product Settings Management Methods

func (a *App) GetDefaultProductSettings() (*DefaultProductSettings, error) {
	return a.db.GetDefaultProductSettings()
}

func (a *App) UpdateDefaultProductSettings(settings DefaultProductSettings) error {
	settings.UpdatedAt = time.Now()
	return a.db.UpdateDefaultProductSettings(&settings)
}

// PDF Generation Methods

func (a *App) GenerateInvoicePDF(invoiceID int) (string, error) {
	pdfBytes, err := a.pdfService.GenerateInvoicePDF(invoiceID)
	if err != nil {
		return "", err
	}

	// Save PDF to user's Documents folder
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	documentsDir := filepath.Join(homeDir, "Documents", "DijInvoice")
	if err := os.MkdirAll(documentsDir, 0755); err != nil {
		return "", err
	}

	invoice, err := a.db.GetInvoiceByID(invoiceID)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("Invoice_%s_%s.pdf", invoice.InvoiceNumber, time.Now().Format("20060102_150405"))
	filepath := filepath.Join(documentsDir, filename)

	if err := os.WriteFile(filepath, pdfBytes, 0644); err != nil {
		return "", err
	}

	return filepath, nil
}

// OpenPDFInViewer opens a PDF file in the default system viewer
func (a *App) OpenPDFInViewer(filepath string) error {
	runtime.BrowserOpenURL(a.ctx, "file://"+filepath)
	return nil
}

// QR Code Methods

func (a *App) ValidateZATCAQRCode(qrCodeBase64 string) error {
	qrService := NewZATCAQRService()
	return qrService.ValidateZATCAQRCode(qrCodeBase64)
}

func (a *App) GetQRCodeInfo(qrCodeBase64 string) (map[string]interface{}, error) {
	qrService := NewZATCAQRService()
	return qrService.GetQRCodeInfo(qrCodeBase64)
}

func (a *App) GenerateSampleQRCode() (string, error) {
	qrService := NewZATCAQRService()
	
	// Create sample invoice data
	sampleInvoice := &Invoice{
		InvoiceNumber: "INV-2024-001",
		IssueDate:     time.Now(),
		TotalAmount:   115.0,
		VATAmount:     15.0,
	}
	
	// Create sample company data
	sampleCompany := &Company{
		Name:      "Sample Company Ltd",
		VATNumber: "300000000000003",
	}
	
	return qrService.GenerateZATCAQRCode(sampleInvoice, sampleCompany)
}

func (a *App) RegenerateInvoiceQRCode(invoiceID int) (string, error) {
	// Get invoice and company data
	invoice, err := a.db.GetInvoiceByID(invoiceID)
	if err != nil {
		return "", fmt.Errorf("failed to get invoice: %v", err)
	}

	company, err := a.db.GetCompany()
	if err != nil {
		return "", fmt.Errorf("failed to get company: %v", err)
	}

	// Generate new ZATCA-compliant QR code
	qrService := NewZATCAQRService()
	qrCodeBase64, err := qrService.GenerateZATCAQRCode(invoice, company)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	// Update invoice with new QR code
	err = a.db.UpdateInvoiceQRCode(invoiceID, qrCodeBase64)
	if err != nil {
		return "", fmt.Errorf("failed to update invoice QR code: %v", err)
	}

	return qrCodeBase64, nil
}

// Utility Methods

func (a *App) ShowSaveDialog(defaultFilename string) (string, error) {
	options := runtime.SaveDialogOptions{
		DefaultFilename: defaultFilename,
		Title:           "Save Invoice PDF",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PDF Files (*.pdf)",
				Pattern:     "*.pdf",
			},
		},
	}

	return runtime.SaveFileDialog(a.ctx, options)
}

func (a *App) ShowMessage(title, message string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   title,
		Message: message,
	})
}

func (a *App) ShowError(title, message string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   title,
		Message: message,
	})
}

// Greet returns a greeting for the given name (keeping original method for compatibility)
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, Welcome to DijInvoice!", name)
}
