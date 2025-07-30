package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"dijibill/database"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// HTMLInvoiceService handles HTML invoice generation
type HTMLInvoiceService struct {
	ctx             context.Context
	db              *database.Database
	templateService *TemplateService
}

// NewHTMLInvoiceService creates a new HTML invoice service
func NewHTMLInvoiceService(ctx context.Context, db *database.Database) *HTMLInvoiceService {
	templateService, err := NewTemplateService()
	if err != nil {
		log.Printf("Warning: Failed to initialize template service: %v", err)
		templateService = nil
	}

	return &HTMLInvoiceService{
		ctx:             ctx,
		db:              db,
		templateService: templateService,
	}
}

// sanitizeCustomerData replaces empty string fields with "n/a"
func (h *HTMLInvoiceService) sanitizeCustomerData(customer *database.Customer) *database.Customer {
	if customer == nil {
		return &database.Customer{
			Name:          "n/a",
			NameArabic:    "غير متوفر",
			Address:       "n/a",
			AddressArabic: "غير متوفر",
			City:          "n/a",
			CityArabic:    "غير متوفر",
			Country:       "n/a",
			CountryArabic: "غير متوفر",
			VATNumber:     "n/a",
			Email:         "n/a",
			Phone:         "n/a",
		}
	}

	// Create a copy to avoid modifying the original
	sanitized := *customer

	// Replace empty strings with "n/a"
	if sanitized.Name == "" {
		sanitized.Name = "n/a"
	}
	if sanitized.NameArabic == "" {
		sanitized.NameArabic = "غير متوفر"
	}
	if sanitized.Address == "" {
		sanitized.Address = "n/a"
	}
	if sanitized.AddressArabic == "" {
		sanitized.AddressArabic = "غير متوفر"
	}
	if sanitized.City == "" {
		sanitized.City = "n/a"
	}
	if sanitized.CityArabic == "" {
		sanitized.CityArabic = "غير متوفر"
	}
	if sanitized.Country == "" {
		sanitized.Country = "n/a"
	}
	if sanitized.CountryArabic == "" {
		sanitized.CountryArabic = "غير متوفر"
	}
	if sanitized.VATNumber == "" {
		sanitized.VATNumber = "n/a"
	}
	if sanitized.Email == "" {
		sanitized.Email = "n/a"
	}
	if sanitized.Phone == "" {
		sanitized.Phone = "n/a"
	}

	return &sanitized
}

// GenerateInvoiceHTML generates HTML content for an invoice with on-demand QR code generation
func (h *HTMLInvoiceService) GenerateInvoiceHTML(invoiceID int) (string, error) {
	return h.GenerateInvoiceHTMLWithLanguage(invoiceID, "english")
}

// GenerateInvoiceHTMLWithLanguage generates HTML content for an invoice in the specified language
func (h *HTMLInvoiceService) GenerateInvoiceHTMLWithLanguage(invoiceID int, language string) (string, error) {
	// Get invoice data
	invoice, err := h.db.GetInvoiceByID(invoiceID)
	if err != nil {
		return "", fmt.Errorf("failed to get invoice: %v", err)
	}

	// Get company data (handle case where company data is missing)
	company, err := h.db.GetCompany()
	if err != nil {
		// Create a placeholder company if company data is not found
		company = &database.Company{
			Name:          "Your Company Name",
			NameArabic:    "اسم شركتك",
			Address:       "Company Address",
			AddressArabic: "عنوان الشركة",
			City:          "City",
			CityArabic:    "المدينة",
			Country:       "Country",
			CountryArabic: "البلد",
			VATNumber:     "VAT Number",
			CRNumber:      "CR Number",
		}
	}

	// Get customer (handle case where customer ID is 0 or customer doesn't exist)
	var customer *database.Customer
	if invoice.CustomerID > 0 {
		var err error
		customer, err = h.db.GetCustomerByID(invoice.CustomerID)
		if err != nil {
			// Log the error but don't fail - create a placeholder customer
			customer = &database.Customer{
				Name:          "Walk-in Customer",
				NameArabic:    "عميل عادي",
				Address:       "",
				AddressArabic: "",
				City:          "",
				CityArabic:    "",
				Country:       "",
				CountryArabic: "",
				VATNumber:     "",
				Email:         "",
				Phone:         "",
			}
		}
	} else {
		// No customer ID provided - create a placeholder
		customer = &database.Customer{
			Name:          "Walk-in Customer",
			NameArabic:    "عميل عادي",
			Address:       "",
			AddressArabic: "",
			City:          "",
			CityArabic:    "",
			Country:       "",
			CountryArabic: "",
			VATNumber:     "",
			Email:         "",
			Phone:         "",
		}
	}

	// Sanitize customer data to replace empty strings with "n/a"
	customer = h.sanitizeCustomerData(customer)
	invoice.Customer = customer

	// Generate QR code on-demand
	qrService := NewZATCAQRService()
	qrCodeBase64, err := qrService.GenerateZATCAQRCodeOnDemand(invoice, company)
	if err != nil {
		// Log the error but continue without QR code
		log.Printf("Warning: Could not generate QR code for invoice %d: %v", invoice.ID, err)
		qrCodeBase64 = ""
	}

	// Create a copy of the invoice with the generated QR code for template rendering
	invoiceWithQR := *invoice
	invoiceWithQR.QRCode = qrCodeBase64

	// Get invoice items with product details
	items := make([]InvoiceItemData, len(invoice.Items))
	for i, item := range invoice.Items {
		product, err := h.db.GetProductByID(item.ProductID)
		if err != nil {
			// Create a placeholder product if the product is not found
			product = &database.Product{
				Name:        fmt.Sprintf("Product ID: %d", item.ProductID),
				NameArabic:  fmt.Sprintf("منتج رقم: %d", item.ProductID),
				Description: "Product not found",
			}
		}
		items[i] = InvoiceItemData{
			InvoiceItem: &item,
			Product:     product,
		}
	}

	// Prepare template data with the QR code included
	data := InvoiceData{
		Invoice: &invoiceWithQR, // Use the copy with QR code
		Company: company,
		Items:   items,
	}

	// Check if template service is available
	if h.templateService == nil {
		return "", fmt.Errorf("template service not available")
	}

	// Get the appropriate template based on language
	tmpl := h.templateService.GetTemplate(language)
	if tmpl == nil {
		return "", fmt.Errorf("template not found for language: %s", language)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return buf.String(), nil
}

// ViewInvoiceHTML generates HTML and opens it in browser for preview
func (h *HTMLInvoiceService) ViewInvoiceHTML(invoiceID int) error {
	htmlContent, err := h.GenerateInvoiceHTML(invoiceID)
	if err != nil {
		return err
	}

	// Create temporary HTML file
	tempDir := os.TempDir()
	invoice, err := h.db.GetInvoiceByID(invoiceID)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("Invoice_%s_preview.html", invoice.InvoiceNumber)
	tempFilePath := filepath.Join(tempDir, filename)

	if err := os.WriteFile(tempFilePath, []byte(htmlContent), 0644); err != nil {
		return err
	}

	// Open in browser
	runtime.BrowserOpenURL(h.ctx, "file://"+tempFilePath)
	return nil
}

// Helper method to open HTML content in browser
func (h *HTMLInvoiceService) openHTMLInBrowser(htmlContent string) error {
	// Create temporary HTML file
	tempDir := os.TempDir()
	filename := fmt.Sprintf("Invoice_preview_%d.html", time.Now().Unix())
	tempFilePath := filepath.Join(tempDir, filename)

	if err := os.WriteFile(tempFilePath, []byte(htmlContent), 0644); err != nil {
		return err
	}

	// Open in browser
	runtime.BrowserOpenURL(h.ctx, "file://"+tempFilePath)
	return nil
}

// Helper method to save HTML file with language suffix
func (h *HTMLInvoiceService) saveHTMLFile(htmlContent string, invoiceID int, language string) error {
	invoice, err := h.db.GetInvoiceByID(invoiceID)
	if err != nil {
		return err
	}

	// Create temporary HTML file for printing with language indicator
	tempDir := os.TempDir()
	filename := fmt.Sprintf("Invoice_%s_%s_save.html", invoice.InvoiceNumber, language)
	tempFilePath := filepath.Join(tempDir, filename)

	// Create HTML with print-optimized styling
	printHTML := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Invoice %s (%s)</title>
    <style>
        @media print {
            body { margin: 0; }
            .no-print { display: none; }
        }
        @page {
            margin: 0.5in;
            size: A4;
        }
    </style>
</head>
<body>
    <div class="no-print" style="padding: 20px; background: #f0f0f0; text-align: center; margin-bottom: 20px;">
        <h3>Invoice %s (%s) - Ready to Save/Print</h3>
        <p>Use Ctrl+P (Cmd+P on Mac) or the button below to open print dialog.</p>
        <p>In the print dialog, you can choose "Save as PDF" or select a printer.</p>
        <button onclick="window.print()" style="padding: 10px 20px; font-size: 16px; background: #007bff; color: white; border: none; border-radius: 5px; cursor: pointer; margin-right: 10px;">Print / Save as PDF</button>
        <button onclick="window.close()" style="padding: 10px 20px; font-size: 16px; background: #6c757d; color: white; border: none; border-radius: 5px; cursor: pointer;">Close</button>
    </div>
    %s
    <script>
        // Auto-open print dialog after a short delay
        setTimeout(function() {
            window.print();
        }, 1000);
    </script>
</body>
</html>`, invoice.InvoiceNumber, language, invoice.InvoiceNumber, language, htmlContent)

	if err := os.WriteFile(tempFilePath, []byte(printHTML), 0644); err != nil {
		return err
	}

	// Open in browser which will automatically show print dialog
	runtime.BrowserOpenURL(h.ctx, "file://"+tempFilePath)
	return nil
}

// PrintInvoiceHTML generates HTML and opens print dialog
func (h *HTMLInvoiceService) PrintInvoiceHTML(invoiceID int) error {
	htmlContent, err := h.GenerateInvoiceHTML(invoiceID)
	if err != nil {
		return err
	}

	// Create temporary HTML file with print-specific styling
	printHTML := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Print Invoice</title>
    <style>
        @media print {
            body { margin: 0; }
            .no-print { display: none; }
        }
    </style>
</head>
<body>
    <div class="no-print" style="padding: 20px; background: #f0f0f0; text-align: center;">
        <button onclick="window.print()" style="padding: 10px 20px; font-size: 16px; background: #007bff; color: white; border: none; border-radius: 5px; cursor: pointer;">Print Invoice</button>
        <button onclick="window.close()" style="padding: 10px 20px; font-size: 16px; background: #6c757d; color: white; border: none; border-radius: 5px; cursor: pointer; margin-left: 10px;">Close</button>
    </div>
    %s
</body>
</html>`, htmlContent)

	tempDir := os.TempDir()
	invoice, err := h.db.GetInvoiceByID(invoiceID)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("Invoice_%s_print.html", invoice.InvoiceNumber)
	tempFilePath := filepath.Join(tempDir, filename)

	if err := os.WriteFile(tempFilePath, []byte(printHTML), 0644); err != nil {
		return err
	}

	// Open in browser for printing
	runtime.BrowserOpenURL(h.ctx, "file://"+tempFilePath)
	return nil
}

// SaveInvoiceHTML opens the native print dialog for saving/printing the invoice
func (h *HTMLInvoiceService) SaveInvoiceHTML(invoiceID int) error {
	htmlContent, err := h.GenerateInvoiceHTML(invoiceID)
	if err != nil {
		return err
	}

	// Create temporary HTML file for printing
	tempDir := os.TempDir()
	invoice, err := h.db.GetInvoiceByID(invoiceID)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("Invoice_%s_save.html", invoice.InvoiceNumber)
	tempFilePath := filepath.Join(tempDir, filename)

	// Create HTML with print-optimized styling
	printHTML := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Invoice %s</title>
    <style>
        @media print {
            body { margin: 0; }
            .no-print { display: none; }
        }
        @page {
            margin: 0.5in;
            size: A4;
        }
    </style>
</head>
<body>
    <div class="no-print" style="padding: 20px; background: #f0f0f0; text-align: center; margin-bottom: 20px;">
        <h3>Invoice %s - Ready to Save/Print</h3>
        <p>Use Ctrl+P (Cmd+P on Mac) or the button below to open print dialog.</p>
        <p>In the print dialog, you can choose "Save as PDF" or select a printer.</p>
        <button onclick="window.print()" style="padding: 10px 20px; font-size: 16px; background: #007bff; color: white; border: none; border-radius: 5px; cursor: pointer; margin-right: 10px;">Print / Save as PDF</button>
        <button onclick="window.close()" style="padding: 10px 20px; font-size: 16px; background: #6c757d; color: white; border: none; border-radius: 5px; cursor: pointer;">Close</button>
    </div>
    %s
    <script>
        // Auto-open print dialog after a short delay
        setTimeout(function() {
            window.print();
        }, 1000);
    </script>
</body>
</html>`, invoice.InvoiceNumber, invoice.InvoiceNumber, htmlContent)

	if err := os.WriteFile(tempFilePath, []byte(printHTML), 0644); err != nil {
		return err
	}

	// Open in browser which will automatically show print dialog
	runtime.BrowserOpenURL(h.ctx, "file://"+tempFilePath)
	return nil
}
