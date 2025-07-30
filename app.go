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

// Item Category Management Methods

func (a *App) CreateItemCategory(category ItemCategory) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	return a.db.CreateItemCategory(&category)
}

func (a *App) GetItemCategories() ([]ItemCategory, error) {
	return a.db.GetItemCategories()
}

// Sale Item Management Methods

func (a *App) CreateSaleItem(item SaleItem) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	return a.db.CreateSaleItem(&item)
}

func (a *App) GetSaleItems() ([]SaleItem, error) {
	return a.db.GetSaleItems()
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
