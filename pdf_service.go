package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"time"

	"codeberg.org/go-pdf/fpdf"
	"github.com/abdullahdiaa/garabic" // For letter shaping
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

// ar handles Arabic text. It shapes the letters, then reverses the string
// for correct right-to-left rendering in FPDF.
// IMPORTANT: Only use this function on strings that are 100% Arabic.
func ar(text string) string {
	shapedText := garabic.Shape(text)
	runes := []rune(shapedText)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// getProductByID retrieves a product by its ID
func (p *PDFService) getProductByID(productID int) (*database.Product, error) {
	// This should be replaced with your actual database call
	// Returning dummy data for the example to work
	if productID == 123 {
		return &database.Product{
			ID:         123,
			Name:       "Professional Service",
			NameArabic: "خدمة احترافية",
		}, nil
	}
	return nil, fmt.Errorf("product not found")
}

// GenerateInvoicePDF generates a PDF invoice with proper Arabic support.
func (p *PDFService) GenerateInvoicePDF(invoiceID int) ([]byte, error) {
	// FIX: Declare the 'err' variable to resolve the "undefined: err" compiler error.
	var err error

	// In a real app, these would be fetched from the DB
	invoice := &database.Invoice{
		ID:            invoiceID,
		InvoiceNumber: "INV-000001",
		IssueDate:     time.Now(),
		DueDate:       time.Now().Add(30 * 24 * time.Hour),
		SubTotal:      125.00,
		VATAmount:     18.75,
		TotalAmount:   143.75,
		Items: []database.InvoiceItem{
			{ProductID: 123, Quantity: 1.00, UnitPrice: 125.00, VATRate: 15, TotalAmount: 143.75},
		},
		Customer: &database.Customer{
			Name:          "Valued Customer",
			NameArabic:    "العميل الكريم",
			Address:       "456 Client Avenue",
			AddressArabic: "٤٥٦ شارع العميل",
			City:          "Riyadh",
			CityArabic:    "الرياض",
			Phone:         "0501234567",
			VATNumber:     "310987654300003",
		},
		Notes:       "Payment is due within 30 days.\nThank you for your business.",
		NotesArabic: "الدفع مستحق خلال ٣٠ يوماً.\nشكراً لتعاملكم معنا.",
	}
	company := &database.Company{
		Name:          "Your Company Name",
		NameArabic:    "اسم شركتك",
		Address:       "123 Business Street",
		AddressArabic: "١٢٣ شارع الأعمال",
		City:          "Riyadh",
		CityArabic:    "الرياض",
		VATNumber:     "300123456700003",
		CRNumber:      "1010000000",
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Register fonts
	fontPath := filepath.Join("fonts", "NotoSansArabic-Regular.ttf")
	pdf.AddUTF8Font("arabic", "", fontPath)
	if pdf.Error() != nil {
		return nil, fmt.Errorf("failed to load arabic font: %w", pdf.Error())
	}

	// --- QR Code ---
	qrData := fmt.Sprintf("Seller:%s\nVAT:%s\nTime:%s\nTotal:%.2f\nVATTax:%.2f",
		company.Name, company.VATNumber, invoice.IssueDate.Format(time.RFC3339), invoice.TotalAmount, invoice.VATAmount)
	qrBytes, _ := qrcode.Encode(qrData, qrcode.Medium, 256)
	pdf.RegisterImageOptionsReader("qr", fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(qrBytes))
	pdf.ImageOptions("qr", 160, 10, 40, 40, false, fpdf.ImageOptions{}, 0, "")

	// --- Header ---
	yPos := 10.0
	// English (Left Side)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetXY(10, yPos)
	pdf.Cell(95, 8, company.Name)
	// Arabic (Right Side)
	pdf.SetFont("arabic", "", 14)
	pdf.SetXY(105, yPos)
	pdf.CellFormat(95, 8, ar(company.NameArabic), "", 0, "R", false, 0, "")
	yPos += 8

	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(10, yPos)
	pdf.Cell(95, 6, fmt.Sprintf("Address: %s, %s", company.Address, company.City))
	pdf.SetFont("arabic", "", 9)
	pdf.SetXY(105, yPos)
	pdf.CellFormat(95, 6, ar(company.AddressArabic+", "+company.CityArabic), "", 0, "R", false, 0, "")
	yPos += 6

	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(10, yPos)
	pdf.Cell(95, 6, "VAT: "+company.VATNumber+" | CR: "+company.CRNumber)
	pdf.SetFont("arabic", "", 9)
	pdf.SetXY(105, yPos)
	pdf.CellFormat(95, 6, ar("الرقم الضريبي: "+company.VATNumber+" | س.ت: "+company.CRNumber), "", 0, "R", false, 0, "")

	// --- Invoice Details ---
	yPos = 55.0
	// Left Box (Billed To)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(10, yPos)
	pdf.Cell(47.5, 7, "BILL TO")
	pdf.SetFont("arabic", "B", 10)
	pdf.SetX(57.5)
	pdf.CellFormat(47.5, 7, ar("الفاتورة إلى"), "", 0, "R", false, 0, "")
	yPos += 8

	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(10, yPos)
	pdf.Cell(95, 6, invoice.Customer.Name)
	pdf.SetFont("arabic", "", 9)
	pdf.SetXY(10, yPos+6)
	pdf.CellFormat(95, 6, ar(invoice.Customer.NameArabic), "", 0, "L", false, 0, "")
	yPos += 12

	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(10, yPos)
	pdf.MultiCell(95, 5, fmt.Sprintf("VAT #: %s\nPhone: %s", invoice.Customer.VATNumber, invoice.Customer.Phone), "", "L", false)
	pdf.Rect(10, 55, 95, 35, "S")

	// Right Box (Invoice Details)
	yPos = 55.0
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(115, yPos)
	pdf.Cell(42.5, 7, "INVOICE #")
	pdf.SetFont("Arial", "", 10)
	pdf.SetX(157.5)
	pdf.CellFormat(42.5, 7, invoice.InvoiceNumber, "", 0, "R", false, 0, "")
	yPos += 7

	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(115, yPos)
	pdf.Cell(42.5, 7, "DATE")
	pdf.SetFont("Arial", "", 10)
	pdf.SetX(157.5)
	pdf.CellFormat(42.5, 7, invoice.IssueDate.Format("2006-01-02"), "", 0, "R", false, 0, "")
	yPos += 7

	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(115, yPos)
	pdf.Cell(42.5, 7, "DUE DATE")
	pdf.SetFont("Arial", "", 10)
	pdf.SetX(157.5)
	pdf.CellFormat(42.5, 7, invoice.DueDate.Format("2006-01-02"), "", 0, "R", false, 0, "")
	pdf.Rect(115, 55, 85, 21, "S")

	// --- Items Table ---
	yPos = 100.0
	pdf.SetXY(10, yPos)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(100, 10, "Item / "+ar("البند"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 10, "Qty / "+ar("الكمية"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 10, "Unit Price / "+ar("السعر"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 10, "Total / "+ar("الإجمالي"), "1", 1, "C", true, 0, "")

	for _, item := range invoice.Items {
		product, _ := p.getProductByID(item.ProductID)

		// Use the Arabic font for the item description
		pdf.SetFont("arabic", "", 10)

		// Calculate height for MultiCell
		enDesc := product.Name
		arDesc := ar(product.NameArabic)
		_, lineHt := pdf.GetFontSize()
		cellHeight := lineHt * 2.5 // Height for two lines

		x := pdf.GetX()
		y := pdf.GetY()

		pdf.MultiCell(100, cellHeight/2, enDesc+"\n"+arDesc, "LR", "L", false)
		pdf.SetXY(x+100, y)

		// Use Arial for numbers
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(20, cellHeight, fmt.Sprintf("%.2f", item.Quantity), "LR", 0, "C", false, 0, "")
		pdf.CellFormat(35, cellHeight, fmt.Sprintf("%.2f", item.UnitPrice), "LR", 0, "C", false, 0, "")
		pdf.CellFormat(35, cellHeight, fmt.Sprintf("%.2f", item.TotalAmount), "LR", 1, "C", false, 0, "")
	}
	pdf.CellFormat(190, 0, "", "T", 1, "", false, 0, "") // Close table

	// --- Totals Section ---
	yPos = pdf.GetY() + 5
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(120, yPos)
	pdf.Cell(35, 8, "Subtotal:")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(45, 8, fmt.Sprintf("%.2f SAR", invoice.SubTotal), "", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(120, pdf.GetY())
	pdf.Cell(35, 8, "VAT (15%):")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(45, 8, fmt.Sprintf("%.2f SAR", invoice.VATAmount), "", 1, "R", false, 0, "")

	pdf.SetFillColor(230, 230, 230)
	pdf.SetXY(120, pdf.GetY())
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(35, 10, "TOTAL", "1", 0, "L", true, 0, "")
	pdf.CellFormat(45, 10, fmt.Sprintf("%.2f SAR", invoice.TotalAmount), "1", 1, "R", true, 0, "")

	// --- Notes Section ---
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(10, yPos)
	pdf.Cell(100, 6, "Notes:")
	pdf.SetFont("arabic", "", 9)
	pdf.SetXY(10, yPos+6)
	pdf.MultiCell(100, 5, invoice.Notes+"\n"+ar(invoice.NotesArabic), "T", "L", false)

	// --- Finalize ---
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF output: %v", err)
	}
	return buf.Bytes(), nil
}
