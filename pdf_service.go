package main

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"dijibill/database"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type PDFService struct {
	db        *database.Database
	qrService *ZATCAQRService
}

func NewPDFService(db *database.Database) *PDFService {
	return &PDFService{
		db:        db,
		qrService: NewZATCAQRService(),
	}
}

// GenerateInvoicePDF generates a PDF invoice with Arabic/English support
func (p *PDFService) GenerateInvoicePDF(invoiceID int) ([]byte, error) {
	// Get invoice data
	invoice, err := p.db.GetInvoiceByID(invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %v", err)
	}

	// Get company data
	company, err := p.db.GetCompany()
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %v", err)
	}

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add Arabic font support (you may need to add actual Arabic font files)
	// For now, we'll use the default font and handle Arabic text as best as possible
	pdf.SetFont("Arial", "", 12)

	// Generate ZATCA-compliant QR Code
	qrCodeBase64, err := p.generateQRCodeData(invoice, company)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ZATCA QR code: %v", err)
	}

	// Generate QR code image for PDF embedding
	qrCodeBytes, err := qrcode.Encode(qrCodeBase64, qrcode.Medium, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code image: %v", err)
	}

	// Convert QR code image to base64 for embedding
	qrCodeImageBase64 := base64.StdEncoding.EncodeToString(qrCodeBytes)

	// Header - Company Information
	p.addCompanyHeader(pdf, company)

	// Invoice Title
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(95, 10, "INVOICE / فاتورة")
	pdf.Cell(95, 10, fmt.Sprintf("Invoice No: %s", invoice.InvoiceNumber))
	pdf.Ln(8)

	// Invoice Details
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(95, 6, fmt.Sprintf("Issue Date: %s", invoice.IssueDate.Format("2006-01-02")))
	pdf.Cell(95, 6, fmt.Sprintf("تاريخ الإصدار: %s", invoice.IssueDate.Format("2006-01-02")))
	pdf.Ln(6)

	if !invoice.DueDate.IsZero() {
		pdf.Cell(95, 6, fmt.Sprintf("Due Date: %s", invoice.DueDate.Format("2006-01-02")))
		pdf.Cell(95, 6, fmt.Sprintf("تاريخ الاستحقاق: %s", invoice.DueDate.Format("2006-01-02")))
		pdf.Ln(6)
	}

	pdf.Ln(5)

	// Customer Information
	p.addCustomerInfo(pdf, invoice.Customer)

	pdf.Ln(10)

	// Invoice Items Table
	p.addInvoiceItemsTable(pdf, invoice)

	pdf.Ln(10)

	// Totals
	p.addInvoiceTotals(pdf, invoice)

	// QR Code
	pdf.Ln(10)
	p.addQRCode(pdf, qrCodeImageBase64)

	// Store the ZATCA-compliant QR code in the invoice record
	invoice.QRCode = qrCodeBase64
	err = p.db.UpdateInvoiceQRCode(invoice.ID, qrCodeBase64)
	if err != nil {
		// Log error but don't fail PDF generation
		fmt.Printf("Warning: failed to update invoice QR code: %v\n", err)
	}

	// Footer
	p.addFooter(pdf, company)

	// Output PDF to buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %v", err)
	}

	return buf.Bytes(), nil
}

func (p *PDFService) addCompanyHeader(pdf *gofpdf.Fpdf, company *database.Company) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(95, 8, company.Name)
	pdf.Cell(95, 8, company.NameArabic)
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(95, 6, fmt.Sprintf("VAT No: %s", company.VATNumber))
	pdf.Cell(95, 6, fmt.Sprintf("الرقم الضريبي: %s", company.VATNumber))
	pdf.Ln(6)

	if company.CRNumber != "" {
		pdf.Cell(95, 6, fmt.Sprintf("CR No: %s", company.CRNumber))
		pdf.Cell(95, 6, fmt.Sprintf("رقم السجل التجاري: %s", company.CRNumber))
		pdf.Ln(6)
	}

	pdf.Cell(95, 6, company.Address)
	pdf.Cell(95, 6, company.AddressArabic)
	pdf.Ln(6)

	pdf.Cell(95, 6, fmt.Sprintf("%s, %s", company.City, company.Country))
	pdf.Cell(95, 6, fmt.Sprintf("%s، %s", company.CityArabic, company.CountryArabic))
	pdf.Ln(6)

	if company.Phone != "" {
		pdf.Cell(95, 6, fmt.Sprintf("Phone: %s", company.Phone))
		pdf.Cell(95, 6, fmt.Sprintf("الهاتف: %s", company.Phone))
		pdf.Ln(6)
	}

	if company.Email != "" {
		pdf.Cell(95, 6, fmt.Sprintf("Email: %s", company.Email))
		pdf.Cell(95, 6, fmt.Sprintf("البريد الإلكتروني: %s", company.Email))
		pdf.Ln(6)
	}
}

func (p *PDFService) addCustomerInfo(pdf *gofpdf.Fpdf, customer *database.Customer) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(95, 8, "Bill To:")
	pdf.Cell(95, 8, "إرسال الفاتورة إلى:")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(95, 6, customer.Name)
	pdf.Cell(95, 6, customer.NameArabic)
	pdf.Ln(6)

	if customer.VATNumber != "" {
		pdf.Cell(95, 6, fmt.Sprintf("VAT No: %s", customer.VATNumber))
		pdf.Cell(95, 6, fmt.Sprintf("الرقم الضريبي: %s", customer.VATNumber))
		pdf.Ln(6)
	}

	if customer.Address != "" {
		pdf.Cell(95, 6, customer.Address)
		pdf.Cell(95, 6, customer.AddressArabic)
		pdf.Ln(6)
	}

	if customer.City != "" {
		pdf.Cell(95, 6, fmt.Sprintf("%s, %s", customer.City, customer.Country))
		pdf.Cell(95, 6, fmt.Sprintf("%s، %s", customer.CityArabic, customer.CountryArabic))
		pdf.Ln(6)
	}

	if customer.Phone != "" {
		pdf.Cell(95, 6, fmt.Sprintf("Phone: %s", customer.Phone))
		pdf.Cell(95, 6, fmt.Sprintf("الهاتف: %s", customer.Phone))
		pdf.Ln(6)
	}

	if customer.Email != "" {
		pdf.Cell(95, 6, fmt.Sprintf("Email: %s", customer.Email))
		pdf.Cell(95, 6, fmt.Sprintf("البريد الإلكتروني: %s", customer.Email))
		pdf.Ln(6)
	}
}

func (p *PDFService) addInvoiceItemsTable(pdf *gofpdf.Fpdf, invoice *database.Invoice) {
	// Table headers
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240)

	// English headers
	pdf.CellFormat(30, 8, "Description", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, "Qty", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, "Unit Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, "VAT %", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, "VAT Amount", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, "Total", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 8, "الوصف", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 8)
	pdf.SetFillColor(255, 255, 255)

	// Get products for the invoice items
	for _, item := range invoice.Items {
		product, err := p.getProductByID(item.ProductID)
		if err != nil {
			continue
		}

		pdf.CellFormat(30, 6, p.truncateText(product.Name, 25), "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.2f", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("%.2f", item.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.1f%%", item.VATRate), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("%.2f", item.VATAmount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("%.2f", item.TotalAmount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(45, 6, p.truncateText(product.NameArabic, 35), "1", 1, "R", false, 0, "")
	}
}

func (p *PDFService) addInvoiceTotals(pdf *gofpdf.Fpdf, invoice *database.Invoice) {
	pdf.SetFont("Arial", "B", 10)

	// Subtotal
	pdf.Cell(120, 8, "")
	pdf.Cell(35, 8, "Subtotal / المجموع الفرعي:")
	pdf.Cell(35, 8, fmt.Sprintf("%.2f SAR", invoice.SubTotal))
	pdf.Ln(8)

	// VAT
	pdf.Cell(120, 8, "")
	pdf.Cell(35, 8, "VAT / ضريبة القيمة المضافة:")
	pdf.Cell(35, 8, fmt.Sprintf("%.2f SAR", invoice.VATAmount))
	pdf.Ln(8)

	// Total
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(120, 10, "")
	pdf.Cell(35, 10, "Total / المجموع:")
	pdf.Cell(35, 10, fmt.Sprintf("%.2f SAR", invoice.TotalAmount))
	pdf.Ln(10)
}

func (p *PDFService) addQRCode(pdf *gofpdf.Fpdf, qrCodeBase64 string) {
	// Decode base64 and check if valid
	_, err := base64.StdEncoding.DecodeString(qrCodeBase64)
	if err == nil {
		// For now, we'll just add text indicating QR code position
		pdf.SetFont("Arial", "", 8)
		pdf.Cell(0, 6, "QR Code for KSA VAT compliance would be placed here")
		pdf.Ln(6)
		pdf.Cell(0, 6, "رمز الاستجابة السريعة للامتثال لضريبة القيمة المضافة السعودية")
		pdf.Ln(6)
	}
}

func (p *PDFService) addFooter(pdf *gofpdf.Fpdf, company *database.Company) {
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 6, "Thank you for your business! / شكراً لتعاملكم معنا!")
	pdf.Ln(4)
	pdf.Cell(0, 6, "This is a computer generated invoice. / هذه فاتورة مُنشأة بواسطة الحاسوب.")
}

func (p *PDFService) generateQRCodeData(invoice *database.Invoice, company *database.Company) (string, error) {
	// Generate ZATCA-compliant QR code
	qrCodeBase64, err := p.qrService.GenerateZATCAQRCode(invoice, company)
	if err != nil {
		return "", fmt.Errorf("failed to generate ZATCA QR code: %v", err)
	}

	return qrCodeBase64, nil
}

func (p *PDFService) getProductByID(id int) (*database.Product, error) {
	// This is a simplified version - in a real app you'd have a proper method in the database
	products, err := p.db.GetProducts()
	if err != nil {
		return nil, err
	}

	for _, product := range products {
		if product.ID == id {
			return &product, nil
		}
	}

	return nil, fmt.Errorf("product not found")
}

func (p *PDFService) truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}
