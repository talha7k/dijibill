package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func (a *App) DeleteProduct(id int) error {
	return a.db.DeleteProduct(id)
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

// Language-specific HTML Invoice Generation Methods

func (a *App) GenerateInvoiceHTMLEnglish(invoiceID int) (string, error) {
	return a.htmlInvoiceService.GenerateInvoiceHTMLWithLanguage(invoiceID, "english")
}

func (a *App) GenerateInvoiceHTMLArabic(invoiceID int) (string, error) {
	return a.htmlInvoiceService.GenerateInvoiceHTMLWithLanguage(invoiceID, "arabic")
}

func (a *App) GenerateInvoiceHTMLBilingual(invoiceID int) (string, error) {
	return a.htmlInvoiceService.GenerateInvoiceHTMLWithLanguage(invoiceID, "bilingual")
}

func (a *App) ViewInvoiceHTMLEnglish(invoiceID int) error {
	htmlContent, err := a.htmlInvoiceService.GenerateInvoiceHTMLWithLanguage(invoiceID, "english")
	if err != nil {
		return err
	}
	return a.htmlInvoiceService.openHTMLInBrowser(htmlContent)
}

func (a *App) ViewInvoiceHTMLArabic(invoiceID int) error {
	htmlContent, err := a.htmlInvoiceService.GenerateInvoiceHTMLWithLanguage(invoiceID, "arabic")
	if err != nil {
		return err
	}
	return a.htmlInvoiceService.openHTMLInBrowser(htmlContent)
}

func (a *App) ViewInvoiceHTMLBilingual(invoiceID int) error {
	htmlContent, err := a.htmlInvoiceService.GenerateInvoiceHTMLWithLanguage(invoiceID, "bilingual")
	if err != nil {
		return err
	}
	return a.htmlInvoiceService.openHTMLInBrowser(htmlContent)
}

func (a *App) SaveInvoiceHTMLEnglish(invoiceID int) error {
	htmlContent, err := a.htmlInvoiceService.GenerateInvoiceHTMLWithLanguage(invoiceID, "english")
	if err != nil {
		return err
	}
	return a.htmlInvoiceService.saveHTMLFile(htmlContent, invoiceID, "english")
}

func (a *App) SaveInvoiceHTMLArabic(invoiceID int) error {
	htmlContent, err := a.htmlInvoiceService.GenerateInvoiceHTMLWithLanguage(invoiceID, "arabic")
	if err != nil {
		return err
	}
	return a.htmlInvoiceService.saveHTMLFile(htmlContent, invoiceID, "arabic")
}

func (a *App) SaveInvoiceHTMLBilingual(invoiceID int) error {
	htmlContent, err := a.htmlInvoiceService.GenerateInvoiceHTMLWithLanguage(invoiceID, "bilingual")
	if err != nil {
		return err
	}
	return a.htmlInvoiceService.saveHTMLFile(htmlContent, invoiceID, "bilingual")
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

	// Generate the QR code on-demand and return the TLV Base64 content (not the PNG image)
	// This returns the actual QR code content that should be scanned
	qrData := ZATCAQRData{
		SellerName:  sampleCompany.Name,
		VATNumber:   sampleCompany.VATNumber,
		Timestamp:   sampleInvoice.IssueDate,
		TotalAmount: sampleInvoice.TotalAmount,
		VATAmount:   sampleInvoice.VATAmount,
	}

	// Encode to TLV format
	tlvBytes, err := qrService.encodeTLV(qrData)
	if err != nil {
		return "", fmt.Errorf("failed to encode TLV: %v", err)
	}

	// Return the TLV Base64 content (this is what should be in the QR code)
	return base64.StdEncoding.EncodeToString(tlvBytes), nil
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

// PopulateSampleData creates comprehensive sample data for testing (5 items each)
func (a *App) PopulateSampleData() error {
	// Create 5 sample customers
	customers := []database.Customer{
		{
			Name:          "Ahmed Al-Rashid",
			NameArabic:    "أحمد الراشد",
			VATNumber:     "300000000000001",
			Email:         "ahmed@example.com",
			Phone:         "+966501234567",
			Address:       "123 King Fahd Road",
			AddressArabic: "طريق الملك فهد ١٢٣",
			City:          "Riyadh",
			CityArabic:    "الرياض",
			Country:       "Saudi Arabia",
			CountryArabic: "المملكة العربية السعودية",
		},
		{
			Name:          "Fatima Al-Zahra",
			NameArabic:    "فاطمة الزهراء",
			VATNumber:     "300000000000002",
			Email:         "fatima@example.com",
			Phone:         "+966502345678",
			Address:       "456 Prince Sultan Street",
			AddressArabic: "شارع الأمير سلطان ٤٥٦",
			City:          "Jeddah",
			CityArabic:    "جدة",
			Country:       "Saudi Arabia",
			CountryArabic: "المملكة العربية السعودية",
		},
		{
			Name:          "Mohammed Al-Otaibi",
			NameArabic:    "محمد العتيبي",
			VATNumber:     "300000000000003",
			Email:         "mohammed@example.com",
			Phone:         "+966503456789",
			Address:       "789 Al-Madinah Avenue",
			AddressArabic: "شارع المدينة ٧٨٩",
			City:          "Dammam",
			CityArabic:    "الدمام",
			Country:       "Saudi Arabia",
			CountryArabic: "المملكة العربية السعودية",
		},
		{
			Name:          "Sarah Al-Mansouri",
			NameArabic:    "سارة المنصوري",
			VATNumber:     "300000000000004",
			Email:         "sarah@example.com",
			Phone:         "+966504567890",
			Address:       "321 Corniche Road",
			AddressArabic: "طريق الكورنيش ٣٢١",
			City:          "Khobar",
			CityArabic:    "الخبر",
			Country:       "Saudi Arabia",
			CountryArabic: "المملكة العربية السعودية",
		},
		{
			Name:          "Omar Al-Harbi",
			NameArabic:    "عمر الحربي",
			VATNumber:     "300000000000005",
			Email:         "omar@example.com",
			Phone:         "+966505678901",
			Address:       "654 Tahlia Street",
			AddressArabic: "شارع التحلية ٦٥٤",
			City:          "Riyadh",
			CityArabic:    "الرياض",
			Country:       "Saudi Arabia",
			CountryArabic: "المملكة العربية السعودية",
		},
	}

	for _, customer := range customers {
		if err := a.CreateCustomer(customer); err != nil {
			return fmt.Errorf("failed to create customer %s: %v", customer.Name, err)
		}
	}

	// Create 5 sample products
	products := []database.Product{
		{
			Name:        "Laptop Computer",
			NameArabic:  "جهاز كمبيوتر محمول",
			Description: "High-performance laptop for business use",
			UnitPrice:   2500.0,
			VATRate:     15.0,
			Unit:        "pcs",
			UnitArabic:  "قطعة",
			SKU:         "LAPTOP-001",
			Stock:       25,
			IsActive:    true,
		},
		{
			Name:        "Office Chair",
			NameArabic:  "كرسي مكتب",
			Description: "Ergonomic office chair with lumbar support",
			UnitPrice:   450.0,
			VATRate:     15.0,
			Unit:        "pcs",
			UnitArabic:  "قطعة",
			SKU:         "CHAIR-001",
			Stock:       50,
			IsActive:    true,
		},
		{
			Name:        "Wireless Mouse",
			NameArabic:  "فأرة لاسلكية",
			Description: "Bluetooth wireless mouse with precision tracking",
			UnitPrice:   85.0,
			VATRate:     15.0,
			Unit:        "pcs",
			UnitArabic:  "قطعة",
			SKU:         "MOUSE-001",
			Stock:       100,
			IsActive:    true,
		},
		{
			Name:        "Monitor 24 inch",
			NameArabic:  "شاشة ٢٤ بوصة",
			Description: "Full HD LED monitor with HDMI connectivity",
			UnitPrice:   650.0,
			VATRate:     15.0,
			Unit:        "pcs",
			UnitArabic:  "قطعة",
			SKU:         "MONITOR-001",
			Stock:       30,
			IsActive:    true,
		},
		{
			Name:        "Desk Lamp",
			NameArabic:  "مصباح مكتب",
			Description: "LED desk lamp with adjustable brightness",
			UnitPrice:   120.0,
			VATRate:     15.0,
			Unit:        "pcs",
			UnitArabic:  "قطعة",
			SKU:         "LAMP-001",
			Stock:       75,
			IsActive:    true,
		},
	}

	for _, product := range products {
		if err := a.CreateProduct(product); err != nil {
			return fmt.Errorf("failed to create product %s: %v", product.Name, err)
		}
	}

	// Create 5 sample payment types
	paymentTypes := []database.PaymentType{
		{
			Name:        "Cash",
			NameArabic:  "نقدي",
			Description: "Cash payment",
			IsActive:    true,
		},
		{
			Name:        "Credit Card",
			NameArabic:  "بطاقة ائتمان",
			Description: "Credit card payment",
			IsActive:    true,
		},
		{
			Name:        "Bank Transfer",
			NameArabic:  "تحويل بنكي",
			Description: "Bank wire transfer",
			IsActive:    true,
		},
		{
			Name:        "Check",
			NameArabic:  "شيك",
			Description: "Check payment",
			IsActive:    true,
		},
		{
			Name:        "Digital Wallet",
			NameArabic:  "محفظة رقمية",
			Description: "Digital wallet payment (Apple Pay, Google Pay, etc.)",
			IsActive:    true,
		},
	}

	for _, paymentType := range paymentTypes {
		if err := a.CreatePaymentType(paymentType); err != nil {
			return fmt.Errorf("failed to create payment type %s: %v", paymentType.Name, err)
		}
	}

	// Get created data for invoice creation
	createdCustomers, err := a.db.GetCustomers()
	if err != nil {
		return err
	}

	createdProducts, err := a.db.GetProducts()
	if err != nil {
		return err
	}

	salesCategories, err := a.db.GetSalesCategories()
	if err != nil {
		return err
	}

	if len(salesCategories) == 0 {
		return fmt.Errorf("no sales categories found - please ensure default data is created")
	}

	// Create 5 sample invoices
	invoices := []database.Invoice{
		{
			InvoiceNumber:   "INV-2024-001",
			CustomerID:      createdCustomers[0].ID,
			SalesCategoryID: salesCategories[0].ID,
			IssueDate:       time.Now().AddDate(0, 0, -10),
			DueDate:         time.Now().AddDate(0, 0, 20),
			Status:          "paid",
			Notes:           "Laptop purchase for office setup",
			NotesArabic:     "شراء جهاز كمبيوتر محمول لإعداد المكتب",
			Items: []database.InvoiceItem{
				{
					ProductID: createdProducts[0].ID, // Laptop
					Quantity:  1,
					UnitPrice: 2500.0,
					VATRate:   15.0,
				},
				{
					ProductID: createdProducts[2].ID, // Mouse
					Quantity:  1,
					UnitPrice: 85.0,
					VATRate:   15.0,
				},
			},
		},
		{
			InvoiceNumber:   "INV-2024-002",
			CustomerID:      createdCustomers[1].ID,
			SalesCategoryID: salesCategories[0].ID,
			IssueDate:       time.Now().AddDate(0, 0, -5),
			DueDate:         time.Now().AddDate(0, 0, 25),
			Status:          "pending",
			Notes:           "Office furniture order",
			NotesArabic:     "طلب أثاث مكتبي",
			Items: []database.InvoiceItem{
				{
					ProductID: createdProducts[1].ID, // Chair
					Quantity:  3,
					UnitPrice: 450.0,
					VATRate:   15.0,
				},
				{
					ProductID: createdProducts[4].ID, // Lamp
					Quantity:  2,
					UnitPrice: 120.0,
					VATRate:   15.0,
				},
			},
		},
		{
			InvoiceNumber:   "INV-2024-003",
			CustomerID:      createdCustomers[2].ID,
			SalesCategoryID: salesCategories[0].ID,
			IssueDate:       time.Now().AddDate(0, 0, -3),
			DueDate:         time.Now().AddDate(0, 0, 27),
			Status:          "draft",
			Notes:           "Monitor and accessories",
			NotesArabic:     "شاشة وملحقات",
			Items: []database.InvoiceItem{
				{
					ProductID: createdProducts[3].ID, // Monitor
					Quantity:  2,
					UnitPrice: 650.0,
					VATRate:   15.0,
				},
			},
		},
		{
			InvoiceNumber:   "INV-2024-004",
			CustomerID:      createdCustomers[3].ID,
			SalesCategoryID: salesCategories[0].ID,
			IssueDate:       time.Now().AddDate(0, 0, -1),
			DueDate:         time.Now().AddDate(0, 0, 29),
			Status:          "sent",
			Notes:           "Complete workstation setup",
			NotesArabic:     "إعداد محطة عمل كاملة",
			Items: []database.InvoiceItem{
				{
					ProductID: createdProducts[0].ID, // Laptop
					Quantity:  1,
					UnitPrice: 2500.0,
					VATRate:   15.0,
				},
				{
					ProductID: createdProducts[1].ID, // Chair
					Quantity:  1,
					UnitPrice: 450.0,
					VATRate:   15.0,
				},
				{
					ProductID: createdProducts[3].ID, // Monitor
					Quantity:  1,
					UnitPrice: 650.0,
					VATRate:   15.0,
				},
			},
		},
		{
			InvoiceNumber:   "INV-2024-005",
			CustomerID:      createdCustomers[4].ID,
			SalesCategoryID: salesCategories[0].ID,
			IssueDate:       time.Now(),
			DueDate:         time.Now().AddDate(0, 0, 30),
			Status:          "draft",
			Notes:           "Accessories bundle",
			NotesArabic:     "حزمة ملحقات",
			Items: []database.InvoiceItem{
				{
					ProductID: createdProducts[2].ID, // Mouse
					Quantity:  5,
					UnitPrice: 85.0,
					VATRate:   15.0,
				},
				{
					ProductID: createdProducts[4].ID, // Lamp
					Quantity:  3,
					UnitPrice: 120.0,
					VATRate:   15.0,
				},
			},
		},
	}

	for _, invoice := range invoices {
		if err := a.CreateInvoice(invoice); err != nil {
			return fmt.Errorf("failed to create invoice %s: %v", invoice.InvoiceNumber, err)
		}
	}

	return nil
}

// TestCustomerNAHandling creates a test invoice with empty customer data to verify N/A handling across all templates
func (a *App) TestCustomerNAHandling() error {
	// Create a test customer with empty fields
	customer := database.Customer{
		Name:          "Test Customer",
		NameArabic:    "عميل اختبار",
		VATNumber:     "", // Empty field
		Email:         "", // Empty field
		Phone:         "", // Empty field
		Address:       "", // Empty field
		AddressArabic: "", // Empty field
		City:          "", // Empty field
		CityArabic:    "", // Empty field
		Country:       "", // Empty field
		CountryArabic: "", // Empty field
	}

	if err := a.CreateCustomer(customer); err != nil {
		return err
	}

	// Get the created customer
	customers, err := a.db.GetCustomers()
	if err != nil {
		return err
	}

	var testCustomer *database.Customer
	for _, c := range customers {
		if c.Name == "Test Customer" {
			testCustomer = &c
			break
		}
	}

	if testCustomer == nil {
		return fmt.Errorf("test customer not found")
	}

	// Get sales categories
	salesCategories, err := a.db.GetSalesCategories()
	if err != nil {
		return err
	}

	if len(salesCategories) == 0 {
		return fmt.Errorf("no sales categories found")
	}

	// Get products
	products, err := a.db.GetProducts()
	if err != nil {
		return err
	}

	if len(products) == 0 {
		return fmt.Errorf("no products found")
	}

	// Create test invoice
	invoice := database.Invoice{
		InvoiceNumber:   "TEST-NA-001",
		CustomerID:      testCustomer.ID,
		SalesCategoryID: salesCategories[0].ID,
		IssueDate:       time.Now(),
		DueDate:         time.Now().AddDate(0, 0, 30),
		Status:          "draft",
		Notes:           "Test invoice for N/A handling verification",
		NotesArabic:     "فاتورة اختبار للتحقق من معالجة غير متوفر",
		Items: []database.InvoiceItem{
			{
				ProductID: products[0].ID,
				Quantity:  1,
				UnitPrice: 100.0,
				VATRate:   15.0,
			},
		},
	}

	if err := a.CreateInvoice(invoice); err != nil {
		return err
	}

	// Get the created invoice
	invoices, err := a.db.GetInvoices()
	if err != nil {
		return err
	}

	var testInvoice *database.Invoice
	for _, inv := range invoices {
		if inv.InvoiceNumber == "TEST-NA-001" {
			testInvoice = &inv
			break
		}
	}

	if testInvoice == nil {
		return fmt.Errorf("test invoice not found")
	}

	// Test all three templates
	languages := []string{"english", "arabic", "bilingual"}
	for _, lang := range languages {
		htmlContent, err := a.htmlInvoiceService.GenerateInvoiceHTMLWithLanguage(testInvoice.ID, lang)
		if err != nil {
			return fmt.Errorf("failed to generate %s template: %v", lang, err)
		}

		// Check if the HTML contains "n/a" or "غير متوفر" for empty fields
		if lang == "arabic" {
			if !strings.Contains(htmlContent, "غير متوفر") {
				return fmt.Errorf("Arabic template does not contain 'غير متوفر' for empty fields")
			}
		} else {
			if !strings.Contains(htmlContent, "n/a") {
				return fmt.Errorf("%s template does not contain 'n/a' for empty fields", lang)
			}
		}
	}

	return nil
}

// Greet returns a greeting for the given name (keeping original method for compatibility)
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, Welcome to DijiBill!", name)
}
