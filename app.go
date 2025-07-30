package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"dijibill/database"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx                context.Context
	db                 *database.Database
	htmlInvoiceService *HTMLInvoiceService
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

	dbPath := filepath.Join(homeDir, "dijibill.db")
	a.db, err = database.NewDatabase(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize HTML invoice service
	a.htmlInvoiceService = NewHTMLInvoiceService(a.ctx, a.db)

	// Create sample data for testing
	if err := a.CreateSampleData(); err != nil {
		log.Printf("Warning: Could not create sample data: %v", err)
	}

	log.Println("Application started successfully")
}

// Customer Management Methods

func (a *App) CreateCustomer(customer database.Customer) error {
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	return a.db.CreateCustomer(&customer)
}

func (a *App) GetCustomers() ([]database.Customer, error) {
	return a.db.GetCustomers()
}

func (a *App) GetCustomerByID(id int) (*database.Customer, error) {
	return a.db.GetCustomerByID(id)
}

func (a *App) UpdateCustomer(customer database.Customer) error {
	customer.UpdatedAt = time.Now()
	return a.db.UpdateCustomer(&customer)
}

func (a *App) DeleteCustomer(id int) error {
	return a.db.DeleteCustomer(id)
}

// Product Category Management Methods

func (a *App) CreateProductCategory(category database.ProductCategory) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	return a.db.CreateProductCategory(&category)
}

func (a *App) GetProductCategories() ([]database.ProductCategory, error) {
	return a.db.GetProductCategories()
}

// Product Management Methods

func (a *App) CreateProduct(product database.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	return a.db.CreateProduct(&product)
}

func (a *App) GetProducts() ([]database.Product, error) {
	return a.db.GetProducts()
}

func (a *App) UpdateProduct(product database.Product) error {
	product.UpdatedAt = time.Now()
	return a.db.UpdateProduct(&product)
}

// Invoice Management Methods

func (a *App) CreateInvoice(invoice database.Invoice) error {
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

	// Create the invoice (QR codes will be generated on-demand when needed)
	return a.db.CreateInvoice(&invoice)
}

func (a *App) GetInvoices() ([]database.Invoice, error) {
	return a.db.GetInvoices()
}

func (a *App) GetInvoiceByID(id int) (*database.Invoice, error) {
	return a.db.GetInvoiceByID(id)
}

// Company Management Methods

func (a *App) GetCompany() (*database.Company, error) {
	return a.db.GetCompany()
}

func (a *App) UpdateCompany(company database.Company) error {
	return a.db.UpdateCompany(&company)
}

// Payment Type Management Methods

func (a *App) CreatePaymentType(paymentType database.PaymentType) error {
	paymentType.CreatedAt = time.Now()
	paymentType.UpdatedAt = time.Now()
	return a.db.CreatePaymentType(&paymentType)
}

func (a *App) GetPaymentTypes() ([]database.PaymentType, error) {
	return a.db.GetPaymentTypes()
}

func (a *App) UpdatePaymentType(paymentType database.PaymentType) error {
	paymentType.UpdatedAt = time.Now()
	return a.db.UpdatePaymentType(&paymentType)
}

func (a *App) DeletePaymentType(id int) error {
	return a.db.DeletePaymentType(id)
}

// Payment Management Methods

func (a *App) CreatePayment(payment database.Payment) error {
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	if payment.Status == "" {
		payment.Status = "completed"
	}
	createdPayment, err := a.db.CreatePayment(payment)
	if err != nil {
		return err
	}
	// Update the original payment with the created payment data
	payment = createdPayment
	return nil
}

func (a *App) GetPayments() ([]database.Payment, error) {
	return a.db.GetPayments()
}

func (a *App) GetPaymentsByInvoiceID(invoiceID int) ([]database.Payment, error) {
	return a.db.GetPaymentsByInvoiceID(invoiceID)
}

func (a *App) GetPaymentByID(id int) (database.Payment, error) {
	return a.db.GetPaymentByID(id)
}

func (a *App) UpdatePayment(payment database.Payment) error {
	payment.UpdatedAt = time.Now()
	return a.db.UpdatePayment(payment)
}

func (a *App) DeletePayment(id int) error {
	return a.db.DeletePayment(id)
}

// Sales Category Management Methods

func (a *App) CreateSalesCategory(salesCategory database.SalesCategory) error {
	salesCategory.CreatedAt = time.Now()
	salesCategory.UpdatedAt = time.Now()
	return a.db.CreateSalesCategory(&salesCategory)
}

func (a *App) GetSalesCategories() ([]database.SalesCategory, error) {
	return a.db.GetSalesCategories()
}

func (a *App) UpdateSalesCategory(salesCategory database.SalesCategory) error {
	salesCategory.UpdatedAt = time.Now()
	return a.db.UpdateSalesCategory(&salesCategory)
}

func (a *App) DeleteSalesCategory(id int) error {
	return a.db.DeleteSalesCategory(id)
}

// Tax Rate Management Methods

func (a *App) CreateTaxRate(taxRate database.TaxRate) error {
	taxRate.CreatedAt = time.Now()
	taxRate.UpdatedAt = time.Now()
	return a.db.CreateTaxRate(&taxRate)
}

func (a *App) GetTaxRates() ([]database.TaxRate, error) {
	return a.db.GetTaxRates()
}

func (a *App) UpdateTaxRate(taxRate database.TaxRate) error {
	taxRate.UpdatedAt = time.Now()
	return a.db.UpdateTaxRate(&taxRate)
}

func (a *App) DeleteTaxRate(id int) error {
	return a.db.DeleteTaxRate(id)
}

// Unit of Measurement Management Methods

func (a *App) CreateUnitOfMeasurement(unit database.UnitOfMeasurement) error {
	unit.CreatedAt = time.Now()
	unit.UpdatedAt = time.Now()
	return a.db.CreateUnitOfMeasurement(&unit)
}

func (a *App) GetUnitsOfMeasurement() ([]database.UnitOfMeasurement, error) {
	return a.db.GetUnitsOfMeasurement()
}

func (a *App) UpdateUnitOfMeasurement(unit database.UnitOfMeasurement) error {
	unit.UpdatedAt = time.Now()
	return a.db.UpdateUnitOfMeasurement(&unit)
}

func (a *App) DeleteUnitOfMeasurement(id int) error {
	return a.db.DeleteUnitOfMeasurement(id)
}

// Default Product Settings Management Methods

func (a *App) GetDefaultProductSettings() (*database.DefaultProductSettings, error) {
	return a.db.GetDefaultProductSettings()
}

func (a *App) UpdateDefaultProductSettings(settings database.DefaultProductSettings) error {
	settings.UpdatedAt = time.Now()
	return a.db.UpdateDefaultProductSettings(&settings)
}

// HTML Invoice Generation Methods

func (a *App) GenerateInvoiceHTML(invoiceID int) (string, error) {
	return a.htmlInvoiceService.GenerateInvoiceHTML(invoiceID)
}

func (a *App) ViewInvoiceHTML(invoiceID int) error {
	return a.htmlInvoiceService.ViewInvoiceHTML(invoiceID)
}

func (a *App) PrintInvoiceHTML(invoiceID int) error {
	return a.htmlInvoiceService.PrintInvoiceHTML(invoiceID)
}

func (a *App) SaveInvoiceHTML(invoiceID int) error {
	return a.htmlInvoiceService.SaveInvoiceHTML(invoiceID)
}

// Legacy PDF methods - kept for backward compatibility but now generate HTML
func (a *App) GenerateInvoicePDF(invoiceID int) (string, error) {
	// Generate HTML and save as .html file instead
	htmlContent, err := a.htmlInvoiceService.GenerateInvoiceHTML(invoiceID)
	if err != nil {
		return "", err
	}

	// Save HTML to user's Documents folder
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	documentsDir := filepath.Join(homeDir, "Documents", "dijibill")
	if err := os.MkdirAll(documentsDir, 0755); err != nil {
		return "", err
	}

	invoice, err := a.db.GetInvoiceByID(invoiceID)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("Invoice_%s_%s.html", invoice.InvoiceNumber, time.Now().Format("20060102_150405"))
	filepath := filepath.Join(documentsDir, filename)

	if err := os.WriteFile(filepath, []byte(htmlContent), 0644); err != nil {
		return "", err
	}

	return filepath, nil
}

// ViewInvoicePDF generates HTML and opens it in browser for preview (legacy method name)
func (a *App) ViewInvoicePDF(invoiceID int) error {
	return a.htmlInvoiceService.ViewInvoiceHTML(invoiceID)
}

// DownloadInvoicePDF generates HTML and allows user to save it where they want (legacy method name)
func (a *App) DownloadInvoicePDF(invoiceID int) error {
	return a.htmlInvoiceService.SaveInvoiceHTML(invoiceID)
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
	sampleInvoice := &database.Invoice{
		InvoiceNumber: "INV-2024-001",
		IssueDate:     time.Now(),
		TotalAmount:   115.0,
		VATAmount:     15.0,
	}

	// Create sample company data
	sampleCompany := &database.Company{
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

// RegenerateAllMissingQRCodes generates QR codes for all invoices that don't have them
func (a *App) RegenerateAllMissingQRCodes() error {
	// Get all invoices
	invoices, err := a.db.GetInvoices()
	if err != nil {
		return fmt.Errorf("failed to get invoices: %v", err)
	}

	// Get company data
	company, err := a.db.GetCompany()
	if err != nil {
		return fmt.Errorf("failed to get company: %v", err)
	}

	qrService := NewZATCAQRService()
	regeneratedCount := 0

	for _, invoice := range invoices {
		// Check if invoice already has a QR code
		if invoice.QRCode != "" {
			continue
		}

		// Generate QR code
		qrCodeBase64, err := qrService.GenerateZATCAQRCode(&invoice, company)
		if err != nil {
			log.Printf("Warning: Could not generate QR code for invoice %d: %v", invoice.ID, err)
			continue
		}

		// Update invoice with QR code
		err = a.db.UpdateInvoiceQRCode(invoice.ID, qrCodeBase64)
		if err != nil {
			log.Printf("Warning: Could not update invoice %d with QR code: %v", invoice.ID, err)
			continue
		}

		regeneratedCount++
		log.Printf("Generated QR code for invoice %d (%s)", invoice.ID, invoice.InvoiceNumber)
	}

	log.Printf("Regenerated QR codes for %d invoices", regeneratedCount)
	return nil
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

// CreateSampleData creates sample data for testing
func (a *App) CreateSampleData() error {
	// Check if we already have invoices
	invoices, err := a.db.GetInvoices()
	if err != nil {
		return err
	}

	if len(invoices) > 0 {
		return nil // Already have data
	}

	// Create a sample customer first
	customer := database.Customer{
		Name:          "Sample Customer",
		NameArabic:    "عميل تجريبي",
		VATNumber:     "123456789012345",
		Email:         "customer@example.com",
		Phone:         "+966501234567",
		Address:       "123 Customer Street",
		AddressArabic: "شارع العميل ١٢٣",
		City:          "Riyadh",
		CityArabic:    "الرياض",
		Country:       "Saudi Arabia",
		CountryArabic: "المملكة العربية السعودية",
	}

	if err := a.CreateCustomer(customer); err != nil {
		return err
	}

	// Get the created customer
	customers, err := a.db.GetCustomers()
	if err != nil {
		return err
	}

	if len(customers) == 0 {
		return fmt.Errorf("no customers found")
	}

	// Get sales categories
	salesCategories, err := a.db.GetSalesCategories()
	if err != nil {
		return err
	}

	if len(salesCategories) == 0 {
		return fmt.Errorf("no sales categories found")
	}

	// Create a sample product
	product := database.Product{
		Name:        "Sample Product",
		NameArabic:  "منتج تجريبي",
		Description: "A sample product for testing",
		UnitPrice:   100.0,
		VATRate:     15.0,
		Unit:        "pcs",
		UnitArabic:  "قطعة",
		SKU:         "SAMPLE-001",
		Stock:       10,
		IsActive:    true,
	}

	if err := a.CreateProduct(product); err != nil {
		return err
	}

	// Get the created product
	products, err := a.db.GetProducts()
	if err != nil {
		return err
	}

	if len(products) == 0 {
		return fmt.Errorf("no products found")
	}

	// Create sample invoice
	invoice := database.Invoice{
		InvoiceNumber:   "INV-2024-001",
		CustomerID:      customers[0].ID,
		SalesCategoryID: salesCategories[0].ID,
		IssueDate:       time.Now(),
		DueDate:         time.Now().AddDate(0, 0, 30), // 30 days from now
		Status:          "draft",
		Notes:           "This is a sample invoice for testing",
		NotesArabic:     "هذه فاتورة تجريبية للاختبار",
		Items: []database.InvoiceItem{
			{
				ProductID: products[0].ID,
				Quantity:  2,
				UnitPrice: 100.0,
				VATRate:   15.0,
			},
		},
	}

	return a.CreateInvoice(invoice)
}

// Greet returns a greeting for the given name (keeping original method for compatibility)
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, Welcome to DijiBill!", name)
}
