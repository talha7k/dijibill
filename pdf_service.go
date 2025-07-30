package main

import (
	"bytes"
	"fmt"

	"dijibill/database"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

// Helper function to handle Arabic text display
func (p *PDFService) handleArabicText(text string) string {
	// For now, we'll transliterate common Arabic words to Latin characters
	// This is a temporary solution until proper Arabic font support is added
	arabicToLatin := map[string]string{
		"فاتورة":                    "FATURAH",
		"تاريخ الإصدار":             "TARIKH AL-ISDAR",
		"تاريخ الاستحقاق":           "TARIKH AL-ISTIHQAQ",
		"إرسال الفاتورة إلى":        "IRSAL AL-FATURAH ILA",
		"عميل مباشر":               "AMEEL MUBASHIR",
		"الرقم الضريبي":            "AL-RAQAM AL-DARIBI",
		"رقم السجل التجاري":         "RAQAM AL-SIJIL AL-TIJARI",
		"الهاتف":                   "AL-HATIF",
		"البريد الإلكتروني":        "AL-BAREED AL-ILIKTROONI",
		"الوصف":                    "AL-WASF",
		"المجموع الفرعي":           "AL-MAJMOO AL-FAR'EE",
		"ضريبة القيمة المضافة":      "DAREEBAT AL-QEEMAH AL-MUDAFAH",
		"المجموع":                  "AL-MAJMOO",
		"ملاحظات":                  "MULAHAZAT",
		"شكراً لتعاملكم معنا":        "SHUKRAN LI-TA'AMULUKUM MA'ANA",
		"هذه فاتورة مُنشأة بواسطة الحاسوب": "HADHIHI FATURAH MUNSHA'AH BI-WASIITAT AL-HASOOB",
	}
	
	if latin, exists := arabicToLatin[text]; exists {
		return latin
	}
	return text
}

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

	// Use Arial Unicode MS for better character support
	// For Arabic text, we'll use a workaround since gofpdf has limited Unicode support
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

	// Header - Company Information
	p.addCompanyHeader(pdf, company)

	// Invoice Title
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(95, 10, "INVOICE / "+p.handleArabicText("فاتورة"))
	pdf.Cell(95, 10, fmt.Sprintf("Invoice No: %s", invoice.InvoiceNumber))
	pdf.Ln(8)

	// Invoice Details
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(95, 6, fmt.Sprintf("Issue Date: %s", invoice.IssueDate.Format("2006-01-02")))
	pdf.Cell(95, 6, fmt.Sprintf("%s: %s", p.handleArabicText("تاريخ الإصدار"), invoice.IssueDate.Format("2006-01-02")))
	pdf.Ln(6)

	if !invoice.DueDate.IsZero() {
		pdf.Cell(95, 6, fmt.Sprintf("Due Date: %s", invoice.DueDate.Format("2006-01-02")))
		pdf.Cell(95, 6, fmt.Sprintf("%s: %s", p.handleArabicText("تاريخ الاستحقاق"), invoice.DueDate.Format("2006-01-02")))
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

	// Notes section
	if invoice.Notes != "" || invoice.NotesArabic != "" {
		pdf.Ln(10)
		p.addNotes(pdf, invoice)
	}

	// QR Code
	pdf.Ln(10)
	p.addQRCode(pdf, qrCodeBytes)

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
	pdf.Cell(95, 8, p.handleArabicText(company.NameArabic))
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(95, 6, fmt.Sprintf("VAT No: %s", company.VATNumber))
	pdf.Cell(95, 6, fmt.Sprintf("%s: %s", p.handleArabicText("الرقم الضريبي"), company.VATNumber))
	pdf.Ln(6)

	if company.CRNumber != "" {
		pdf.Cell(95, 6, fmt.Sprintf("CR No: %s", company.CRNumber))
		pdf.Cell(95, 6, fmt.Sprintf("%s: %s", p.handleArabicText("رقم السجل التجاري"), company.CRNumber))
		pdf.Ln(6)
	}

	pdf.Cell(95, 6, company.Address)
	pdf.Cell(95, 6, p.handleArabicText(company.AddressArabic))
	pdf.Ln(6)

	pdf.Cell(95, 6, fmt.Sprintf("%s, %s", company.City, company.Country))
	pdf.Cell(95, 6, fmt.Sprintf("%s، %s", p.handleArabicText(company.CityArabic), p.handleArabicText(company.CountryArabic)))
	pdf.Ln(6)

	if company.Phone != "" {
		pdf.Cell(95, 6, fmt.Sprintf("Phone: %s", company.Phone))
		pdf.Cell(95, 6, fmt.Sprintf("%s: %s", p.handleArabicText("الهاتف"), company.Phone))
		pdf.Ln(6)
	}

	if company.Email != "" {
		pdf.Cell(95, 6, fmt.Sprintf("Email: %s", company.Email))
		pdf.Cell(95, 6, fmt.Sprintf("%s: %s", p.handleArabicText("البريد الإلكتروني"), company.Email))
		pdf.Ln(6)
	}
}

func (p *PDFService) addCustomerInfo(pdf *gofpdf.Fpdf, customer *database.Customer) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(95, 8, "Bill To:")
	pdf.Cell(95, 8, p.handleArabicText("إرسال الفاتورة إلى")+":")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	
	// Handle nil customer case
	if customer == nil {
		pdf.Cell(95, 6, "Walk-in Customer")
		pdf.Cell(95, 6, p.handleArabicText("عميل مباشر"))
		pdf.Ln(6)
		return
	}

	// Customer name
	customerName := customer.Name
	if customerName == "" {
		customerName = "Walk-in Customer"
	}
	customerNameArabic := customer.NameArabic
	if customerNameArabic == "" {
		customerNameArabic = p.handleArabicText("عميل مباشر")
	}
	
	pdf.Cell(95, 6, customerName)
	pdf.Cell(95, 6, p.handleArabicText(customerNameArabic))
	pdf.Ln(6)

	if customer.VATNumber != "" {
		pdf.Cell(95, 6, fmt.Sprintf("VAT No: %s", customer.VATNumber))
		pdf.Cell(95, 6, fmt.Sprintf("%s: %s", p.handleArabicText("الرقم الضريبي"), customer.VATNumber))
		pdf.Ln(6)
	}

	if customer.Address != "" {
		pdf.Cell(95, 6, customer.Address)
		pdf.Cell(95, 6, p.handleArabicText(customer.AddressArabic))
		pdf.Ln(6)
	}

	if customer.City != "" {
		pdf.Cell(95, 6, fmt.Sprintf("%s, %s", customer.City, customer.Country))
		pdf.Cell(95, 6, fmt.Sprintf("%s، %s", p.handleArabicText(customer.CityArabic), p.handleArabicText(customer.CountryArabic)))
		pdf.Ln(6)
	}

	if customer.Phone != "" {
		pdf.Cell(95, 6, fmt.Sprintf("Phone: %s", customer.Phone))
		pdf.Cell(95, 6, fmt.Sprintf("%s: %s", p.handleArabicText("الهاتف"), customer.Phone))
		pdf.Ln(6)
	}

	if customer.Email != "" {
		pdf.Cell(95, 6, fmt.Sprintf("Email: %s", customer.Email))
		pdf.Cell(95, 6, fmt.Sprintf("%s: %s", p.handleArabicText("البريد الإلكتروني"), customer.Email))
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
	pdf.CellFormat(45, 8, p.handleArabicText("الوصف"), "1", 1, "C", true, 0, "")

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
	pdf.Cell(35, 8, "Subtotal / "+p.handleArabicText("المجموع الفرعي")+":")
	pdf.Cell(35, 8, fmt.Sprintf("%.2f SAR", invoice.SubTotal))
	pdf.Ln(8)

	// VAT
	pdf.Cell(120, 8, "")
	pdf.Cell(35, 8, "VAT / "+p.handleArabicText("ضريبة القيمة المضافة")+":")
	pdf.Cell(35, 8, fmt.Sprintf("%.2f SAR", invoice.VATAmount))
	pdf.Ln(8)

	// Total
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(120, 10, "")
	pdf.Cell(35, 10, "Total / "+p.handleArabicText("المجموع")+":")
	pdf.Cell(35, 10, fmt.Sprintf("%.2f SAR", invoice.TotalAmount))
	pdf.Ln(10)
}

func (p *PDFService) addNotes(pdf *gofpdf.Fpdf, invoice *database.Invoice) {
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 8, "Notes / "+p.handleArabicText("ملاحظات")+":")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 9)
	if invoice.Notes != "" {
		pdf.Cell(95, 6, invoice.Notes)
	}
	if invoice.NotesArabic != "" {
		pdf.Cell(95, 6, p.handleArabicText(invoice.NotesArabic))
	}
	pdf.Ln(6)
}

func (p *PDFService) addQRCode(pdf *gofpdf.Fpdf, qrCodeBytes []byte) {
	// Create a temporary file for the QR code image
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 8, "ZATCA QR Code:")
	pdf.Ln(10)
	
	// Register the QR code image from bytes
	reader := bytes.NewReader(qrCodeBytes)
	pdf.RegisterImageReader("qrcode", "PNG", reader)
	
	// Add the QR code image to PDF
	pdf.Image("qrcode", 10, pdf.GetY(), 30, 30, false, "", 0, "")
	
	pdf.Ln(35)
}

func (p *PDFService) addFooter(pdf *gofpdf.Fpdf, company *database.Company) {
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 6, "Thank you for your business! / "+p.handleArabicText("شكراً لتعاملكم معنا")+"!")
	pdf.Ln(4)
	pdf.Cell(0, 6, "This is a computer generated invoice. / "+p.handleArabicText("هذه فاتورة مُنشأة بواسطة الحاسوب")+".")
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
