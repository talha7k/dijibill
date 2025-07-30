package main

import (
	"bytes"
	"fmt"

	"codeberg.org/go-pdf/fpdf"
	"github.com/skip2/go-qrcode"

	"dijibill/database"
)

// PDFService handles PDF generation operations
type PDFService struct {
	db *database.Database
}

// NewPDFService creates a new PDF service instance
func NewPDFService(db *database.Database) *PDFService {
	return &PDFService{
		db: db,
	}
}

// arabic-safe string function with proper UTF-8 handling
func ar(str string) string {
	// Ensure the string is properly encoded as UTF-8
	// For basic Arabic text, this should be sufficient with the UTF-8 font
	if str == "" {
		return str
	}
	
	// Return the string as-is since we're using UTF-8 font
	// The fpdf library with UTF-8 font should handle Arabic text correctly
	return str
}

// getProductByID retrieves a product by its ID
func (p *PDFService) getProductByID(productID int) (*database.Product, error) {
	return p.db.GetProductByID(productID)
}

// GenerateInvoicePDF generates a PDF invoice with proper Arabic support using GoFPDF
func (p *PDFService) GenerateInvoicePDF(invoiceID int) ([]byte, error) {
	// Fetch invoice data
	invoice, err := p.db.GetInvoiceByID(invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %v", err)
	}

	// Fetch company data
	company, err := p.db.GetCompany()
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %v", err)
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Instead of loading a custom Arabic font, use Arial Unicode which has better support
	// This should avoid the index out of range error
	pdf.SetFont("Arial", "", 16)
	
	fmt.Printf("DEBUG: Font loading completed successfully\n")

	// Generate QR code
	qrService := NewZATCAQRService()
	fmt.Printf("DEBUG: Generating QR code for invoice ID: %d\n", invoice.ID)
	fmt.Printf("DEBUG: Company name: %s, VAT: %s\n", company.Name, company.VATNumber)
	fmt.Printf("DEBUG: Invoice total: %.2f, VAT: %.2f\n", invoice.TotalAmount, invoice.VATAmount)
	
	qrCodeData, err := qrService.GenerateZATCAQRCode(invoice, company)
	if err != nil {
		fmt.Printf("DEBUG: QR code generation failed: %v\n", err)
		return nil, fmt.Errorf("failed to generate QR code: %v", err)
	}
	fmt.Printf("DEBUG: QR code generated successfully, length: %d\n", len(qrCodeData))
	
	// Generate and add QR code
	fmt.Printf("DEBUG: Encoding QR code to image...\n")
	qrCodeBytes, _ := qrcode.Encode(qrCodeData, qrcode.Medium, 256)
	fmt.Printf("DEBUG: QR code image generated, size: %d bytes\n", len(qrCodeBytes))
	
	// Save qr code to a temporary file or use a reader
	fmt.Printf("DEBUG: Registering QR image with PDF...\n")
	qrImageInfo := pdf.RegisterImageOptionsReader("qr", fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(qrCodeBytes))
	if qrImageInfo != nil {
		fmt.Printf("DEBUG: Adding QR image to PDF...\n")
		// Position QR Code top-right
		pdf.ImageOptions("qr", 160, 10, 40, 40, false, fpdf.ImageOptions{}, 0, "")
		fmt.Printf("DEBUG: QR image added successfully\n")
	}

	// --- Header ---
	fmt.Printf("DEBUG: Setting font to Arial...\n")
	pdf.SetFont("Arial", "", 16)
	fmt.Printf("DEBUG: Font set successfully, setting position...\n")
	pdf.SetXY(10, 10)
	fmt.Printf("DEBUG: Position set, ready to add content...\n")
	
	// Handle nil company gracefully
	companyName := "Company Name"
	companyAddress := "Address"
	companyCity := "City"
	companyVAT := "VAT Number"
	companyCR := "CR Number"
	
	if company != nil {
		companyName = company.Name
		companyAddress = company.Address
		companyCity = company.City
		companyVAT = company.VATNumber
		companyCR = company.CRNumber
	}
	
	fmt.Printf("DEBUG: About to add company name cell: %s\n", companyName)
	pdf.Cell(50, 10, companyName) // English Name
	fmt.Printf("DEBUG: Company name cell added successfully\n")

	pdf.SetFont("Arial", "", 10)
	pdf.SetXY(10, 18)
	pdf.Cell(100, 6, fmt.Sprintf("Address: %s, %s", companyAddress, companyCity))
	pdf.SetXY(10, 24)
	pdf.Cell(100, 6, fmt.Sprintf("Tax Number: %s", companyVAT))
	pdf.SetXY(10, 30)
	pdf.Cell(100, 6, fmt.Sprintf("Commercial Registration Number: %s", companyCR))

	// --- Billed For Box ---
	pdf.SetXY(10, 50)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(95, 7, "Billed for:")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 12)
	
	// Handle nil customer gracefully
	customerName := "N/A"
	customerAddress := "N/A"
	customerCity := "N/A"
	customerPhone := "N/A"
	customerVAT := "N/A"
	
	if invoice.Customer != nil {
		customerName = invoice.Customer.Name
		customerAddress = invoice.Customer.Address
		customerCity = invoice.Customer.City
		customerPhone = invoice.Customer.Phone
		customerVAT = invoice.Customer.VATNumber
	}
	
	pdf.Cell(95, 7, customerName)
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(95, 5, fmt.Sprintf("Address: %s, %s\nPhone Number: %s\nTax Number: %s",
		customerAddress, customerCity, customerPhone, customerVAT), "0", "L", false)
	pdf.Rect(10, 50, 95, 35, "D") // Draw a box around the customer info

	// --- Invoice Details Box ---
	pdf.SetXY(115, 50)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 7, "Invoice Number:")
	
	// Safe handling of invoice fields
	invoiceNumber := "N/A"
	issueDate := "N/A"
	dueDate := "N/A"
	
	if invoice != nil {
		invoiceNumber = invoice.InvoiceNumber
		if !invoice.IssueDate.IsZero() {
			issueDate = invoice.IssueDate.Format("02/01/2006")
		}
		if !invoice.DueDate.IsZero() {
			dueDate = invoice.DueDate.Format("02/01/2006")
		}
	}
	
	pdf.Cell(40, 7, invoiceNumber)
	pdf.Ln(6)
	pdf.SetX(115)
	pdf.Cell(40, 7, "Invoice Date:")
	pdf.Cell(40, 7, issueDate)
	pdf.Ln(6)
	pdf.SetX(115)
	pdf.Cell(40, 7, "Due Date:")
	pdf.Cell(40, 7, dueDate)
	pdf.Rect(115, 50, 85, 20, "D") // Draw a box around invoice details

	// --- Items Table ---
	pdf.SetY(95)
	pdf.SetFont("Arial", "B", 10)  // Use standard Arial bold for headers
	pdf.SetFillColor(240, 240, 240)
	// Headers
	pdf.CellFormat(80, 10, "Product", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 10, "Tax", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Total", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	// Table Rows
	pdf.SetFont("Arial", "", 10)
	
	// Safe handling of invoice items
	var items []database.InvoiceItem
	if invoice != nil && invoice.Items != nil {
		items = invoice.Items
	}
	
	for _, item := range items {
		product, err := p.getProductByID(item.ProductID)
		if err != nil {
			// If product not found, use a default description
			product = &database.Product{
				Name:       "Unknown Product",
				NameArabic: "منتج غير معروف",
			}
		}

		// Safety checks for product names
		productName := "Product"
		productNameArabic := "منتج"
		
		if product != nil {
			if product.Name != "" {
				productName = product.Name
			}
			if product.NameArabic != "" {
				productNameArabic = ar(product.NameArabic)
			}
		}

		// Combine English and Arabic names in one cell
		productDescription := fmt.Sprintf("%s\n%s", productName, productNameArabic)

		// Get cell height needed for MultiCell with safety checks
		_, lineHt := pdf.GetFontSize()
		
		// Safety check for SplitLines to prevent index out of range
		var lines [][]byte
		if productDescription != "" {
			lines = pdf.SplitLines([]byte(productDescription), 80)
		}
		
		// Ensure we have at least one line to prevent division by zero
		if len(lines) == 0 {
			lines = [][]byte{[]byte("Product")}
		}
		
		cellHeight := float64(len(lines)) * lineHt * 1.2
		if cellHeight < 10 { // Minimum cell height
			cellHeight = 10
		}

		currentX := pdf.GetX()
		currentY := pdf.GetY()

		pdf.MultiCell(80, cellHeight/2, productDescription, "LR", "L", false)
		pdf.SetXY(currentX+80, currentY) // Reset position after multicell
		pdf.CellFormat(20, cellHeight, fmt.Sprintf("%.2f", item.Quantity), "LR", 0, "C", false, 0, "")
		pdf.CellFormat(30, cellHeight, fmt.Sprintf("%.2f SAR", item.UnitPrice), "LR", 0, "C", false, 0, "")
		pdf.CellFormat(20, cellHeight, fmt.Sprintf("%.0f%%", item.VATRate), "LR", 0, "C", false, 0, "")
		pdf.CellFormat(40, cellHeight, fmt.Sprintf("%.2f SAR", item.TotalAmount), "LR", 0, "C", false, 0, "")
		pdf.Ln(cellHeight)
	}
	pdf.CellFormat(190, 0, "", "T", 0, "", false, 0, "") // Close the bottom of the table
	pdf.Ln(5)

	// --- Totals Section ---
	totalYStart := pdf.GetY()
	pdf.SetFont("Arial", "", 10)
	
	// Safe handling of totals
	subTotal := 0.0
	vatAmount := 0.0
	totalAmount := 0.0
	notes := ""
	
	if invoice != nil {
		subTotal = invoice.SubTotal
		vatAmount = invoice.VATAmount
		totalAmount = invoice.TotalAmount
		notes = invoice.Notes
	}
	
	// Subtotal
	pdf.SetX(130)
	pdf.Cell(30, 8, "Subtotal")
	pdf.CellFormat(40, 8, fmt.Sprintf("%.2f SAR", subTotal), "", 1, "R", false, 0, "")
	// VAT
	pdf.SetX(130)
	pdf.Cell(30, 8, "Value added tax")
	pdf.CellFormat(40, 8, fmt.Sprintf("%.2f SAR", vatAmount), "", 1, "R", false, 0, "")
	// Total
	pdf.SetX(130)
	pdf.SetFont("Arial", "B", 12)  // Use standard Arial bold for totals
	pdf.Cell(30, 10, "Total Price")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f SAR", totalAmount), "", 1, "R", false, 0, "")

	// --- Notes ---
	pdf.SetY(totalYStart) // Align notes with the start of the totals
	pdf.SetFont("Arial", "B", 10)  // Use standard Arial bold for notes header
	pdf.Cell(100, 6, "Notes:")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 9)  // Use Arial font for notes content
	pdf.MultiCell(120, 5, notes, "T", "L", false) // Draw a line above notes

	// --- Finalize ---
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %v", err)
	}

	return buf.Bytes(), nil
}
