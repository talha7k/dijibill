package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"dijibill/database"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// HTMLInvoiceService handles HTML invoice generation
type HTMLInvoiceService struct {
	ctx context.Context
	db  *database.Database
}

// NewHTMLInvoiceService creates a new HTML invoice service
func NewHTMLInvoiceService(ctx context.Context, db *database.Database) *HTMLInvoiceService {
	return &HTMLInvoiceService{
		ctx: ctx,
		db:  db,
	}
}

// InvoiceData represents the data structure for invoice template
type InvoiceData struct {
	Invoice *database.Invoice
	Company *database.Company
	Items   []InvoiceItemData
}

// InvoiceItemData represents invoice item with product details
type InvoiceItemData struct {
	*database.InvoiceItem
	Product *database.Product
}

// GenerateInvoiceHTML generates HTML content for an invoice
func (h *HTMLInvoiceService) GenerateInvoiceHTML(invoiceID int) (string, error) {
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
			Name:        "Your Company Name",
			NameArabic:  "اسم شركتك",
			Address:     "Company Address",
			AddressArabic: "عنوان الشركة",
			City:        "City",
			CityArabic:  "المدينة",
			Country:     "Country",
			CountryArabic: "البلد",
			VATNumber:   "VAT Number",
			CRNumber:    "CR Number",
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
				Name:    "Walk-in Customer",
				Address: "N/A",
				City:    "N/A",
				Country: "N/A",
			}
		}
	} else {
		// No customer ID provided - create a placeholder
		customer = &database.Customer{
			Name:    "Walk-in Customer",
			Address: "N/A",
			City:    "N/A",
			Country: "N/A",
		}
	}
	invoice.Customer = customer

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

	// Prepare template data
	data := InvoiceData{
		Invoice: invoice,
		Company: company,
		Items:   items,
	}

	// Parse and execute template
	tmpl, err := template.New("invoice").Parse(invoiceHTMLTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

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

// SaveInvoiceHTML saves the invoice as an HTML file
func (h *HTMLInvoiceService) SaveInvoiceHTML(invoiceID int) error {
	// Get invoice details for default filename
	invoice, err := h.db.GetInvoiceByID(invoiceID)
	if err != nil {
		return err
	}

	defaultFilename := fmt.Sprintf("Invoice_%s_%s.html", invoice.InvoiceNumber, time.Now().Format("20060102_150405"))

	// Show save dialog
	savePath, err := runtime.SaveFileDialog(h.ctx, runtime.SaveDialogOptions{
		DefaultFilename: defaultFilename,
		Title:           "Save Invoice HTML",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "HTML Files (*.html)",
				Pattern:     "*.html",
			},
		},
	})
	if err != nil {
		return err
	}

	// If user cancelled the dialog, savePath will be empty
	if savePath == "" {
		return nil // User cancelled, not an error
	}

	// Generate HTML
	htmlContent, err := h.GenerateInvoiceHTML(invoiceID)
	if err != nil {
		return err
	}

	// Save to chosen location
	if err := os.WriteFile(savePath, []byte(htmlContent), 0644); err != nil {
		return err
	}

	// Show success message
	runtime.MessageDialog(h.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "Success",
		Message: fmt.Sprintf("Invoice HTML saved successfully to:\n%s", savePath),
	})
	return nil
}

// invoiceHTMLTemplate is the HTML template for invoices
const invoiceHTMLTemplate = `
<!DOCTYPE html>
<html lang="en" dir="ltr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Invoice {{.Invoice.InvoiceNumber}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Arial', sans-serif;
            line-height: 1.6;
            color: #333;
            background: white;
        }

        .invoice-container {
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background: white;
        }

        .header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 30px;
            border-bottom: 2px solid #007bff;
            padding-bottom: 20px;
        }

        .company-info {
            flex: 1;
        }

        .company-name {
            font-size: 24px;
            font-weight: bold;
            color: #007bff;
            margin-bottom: 5px;
        }

        .company-name-arabic {
            font-size: 20px;
            font-weight: bold;
            color: #007bff;
            direction: rtl;
            text-align: right;
        }

        .company-details {
            font-size: 14px;
            color: #666;
        }

        .company-details-arabic {
            font-size: 14px;
            color: #666;
            direction: rtl;
            text-align: right;
            margin-top: 5px;
        }

        .qr-code {
            width: 100px;
            height: 100px;
            border: 1px solid #ddd;
            display: flex;
            align-items: center;
            justify-content: center;
            background: #f8f9fa;
        }

        .invoice-details {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 30px;
            margin-bottom: 30px;
        }

        .bill-to, .invoice-info {
            padding: 20px;
            border: 1px solid #ddd;
            border-radius: 5px;
            background: #f8f9fa;
        }

        .section-title {
            font-size: 16px;
            font-weight: bold;
            color: #007bff;
            margin-bottom: 10px;
            border-bottom: 1px solid #ddd;
            padding-bottom: 5px;
        }

        .section-title-arabic {
            font-size: 16px;
            font-weight: bold;
            color: #007bff;
            direction: rtl;
            text-align: right;
        }

        .customer-info {
            font-size: 14px;
        }

        .customer-name {
            font-weight: bold;
            margin-bottom: 5px;
        }

        .customer-name-arabic {
            font-weight: bold;
            direction: rtl;
            text-align: right;
            margin-bottom: 5px;
        }

        .invoice-meta {
            display: flex;
            justify-content: space-between;
            margin-bottom: 5px;
        }

        .items-table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 30px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }

        .items-table th {
            background: #007bff;
            color: white;
            padding: 12px;
            text-align: center;
            font-weight: bold;
        }

        .items-table td {
            padding: 12px;
            border-bottom: 1px solid #ddd;
            text-align: center;
        }

        .items-table tr:nth-child(even) {
            background: #f8f9fa;
        }

        .item-description {
            text-align: left;
            padding-left: 15px;
        }

        .item-description-arabic {
            text-align: right;
            direction: rtl;
            font-size: 12px;
            color: #666;
            margin-top: 3px;
        }

        .totals {
            float: right;
            width: 300px;
            margin-top: 20px;
        }

        .total-row {
            display: flex;
            justify-content: space-between;
            padding: 8px 0;
            border-bottom: 1px solid #eee;
        }

        .total-row.final {
            border-bottom: 2px solid #007bff;
            font-weight: bold;
            font-size: 18px;
            color: #007bff;
        }

        .notes {
            clear: both;
            margin-top: 40px;
            padding: 20px;
            background: #f8f9fa;
            border-left: 4px solid #007bff;
        }

        .notes-title {
            font-weight: bold;
            margin-bottom: 10px;
            color: #007bff;
        }

        .notes-arabic {
            direction: rtl;
            text-align: right;
            margin-top: 10px;
            color: #666;
        }

        .footer {
            margin-top: 40px;
            text-align: center;
            font-size: 12px;
            color: #666;
            border-top: 1px solid #ddd;
            padding-top: 20px;
        }

        @media print {
            body {
                margin: 0;
                padding: 0;
            }
            
            .invoice-container {
                max-width: none;
                margin: 0;
                padding: 15px;
            }
            
            .header {
                margin-bottom: 20px;
                padding-bottom: 15px;
            }
            
            .invoice-details {
                margin-bottom: 20px;
            }
            
            .items-table {
                margin-bottom: 20px;
            }
            
            .notes {
                margin-top: 20px;
            }
            
            .footer {
                margin-top: 20px;
            }
        }

        @media (max-width: 768px) {
            .header {
                flex-direction: column;
                text-align: center;
            }
            
            .invoice-details {
                grid-template-columns: 1fr;
                gap: 20px;
            }
            
            .totals {
                float: none;
                width: 100%;
            }
        }
    </style>
</head>
<body>
    <div class="invoice-container">
        <!-- Header -->
        <div class="header">
            <div class="company-info">
                <div class="company-name">{{if .Company.Name}}{{.Company.Name}}{{else}}Your Company Name{{end}}</div>
                {{if .Company.NameArabic}}<div class="company-name-arabic">{{.Company.NameArabic}}</div>{{end}}
                <div class="company-details">
                    {{if or .Company.Address .Company.City}}{{.Company.Address}}{{if and .Company.Address .Company.City}}, {{end}}{{.Company.City}}{{else}}Company Address{{end}}<br>
                    {{if .Company.VATNumber}}VAT: {{.Company.VATNumber}}{{else}}VAT: Not Set{{end}}{{if .Company.CRNumber}} | CR: {{.Company.CRNumber}}{{end}}
                </div>
                {{if or .Company.AddressArabic .Company.CityArabic}}
                <div class="company-details-arabic">
                    {{.Company.AddressArabic}}{{if and .Company.AddressArabic .Company.CityArabic}}، {{end}}{{.Company.CityArabic}}<br>
                    {{if .Company.VATNumber}}الرقم الضريبي: {{.Company.VATNumber}}{{end}}{{if .Company.CRNumber}} | س.ت: {{.Company.CRNumber}}{{end}}
                </div>
                {{end}}
            </div>
            <div class="qr-code">
                {{if .Invoice.QRCode}}
                    <img src="data:image/png;base64,{{.Invoice.QRCode}}" alt="QR Code" style="max-width: 100%; max-height: 100%;">
                {{else}}
                    QR Code
                {{end}}
            </div>
        </div>

        <!-- Invoice Details -->
        <div class="invoice-details">
            <div class="bill-to">
                <div class="section-title">BILL TO</div>
                <div class="section-title-arabic">الفاتورة إلى</div>
                <div class="customer-info">
                    <div class="customer-name">{{.Invoice.Customer.Name}}</div>
                    {{if .Invoice.Customer.NameArabic}}<div class="customer-name-arabic">{{.Invoice.Customer.NameArabic}}</div>{{end}}
                    {{if .Invoice.Customer.Address}}<div>{{.Invoice.Customer.Address}}</div>{{end}}
                    {{if .Invoice.Customer.AddressArabic}}<div style="direction: rtl; text-align: right;">{{.Invoice.Customer.AddressArabic}}</div>{{end}}
                    {{if or .Invoice.Customer.City .Invoice.Customer.Country}}<div>{{.Invoice.Customer.City}}{{if and .Invoice.Customer.City .Invoice.Customer.Country}}, {{end}}{{.Invoice.Customer.Country}}</div>{{end}}
                    {{if .Invoice.Customer.VATNumber}}<div>VAT: {{.Invoice.Customer.VATNumber}}</div>{{end}}
                    {{if .Invoice.Customer.Phone}}<div>Phone: {{.Invoice.Customer.Phone}}</div>{{end}}
                    {{if .Invoice.Customer.Email}}<div>Email: {{.Invoice.Customer.Email}}</div>{{end}}
                </div>
            </div>
            
            <div class="invoice-info">
                <div class="section-title">INVOICE DETAILS</div>
                <div class="invoice-meta">
                    <span><strong>Invoice #:</strong></span>
                    <span>{{.Invoice.InvoiceNumber}}</span>
                </div>
                <div class="invoice-meta">
                    <span><strong>Date:</strong></span>
                    <span>{{.Invoice.IssueDate.Format "2006-01-02"}}</span>
                </div>
                <div class="invoice-meta">
                    <span><strong>Due Date:</strong></span>
                    <span>{{.Invoice.DueDate.Format "2006-01-02"}}</span>
                </div>
                <div class="invoice-meta">
                    <span><strong>Status:</strong></span>
                    <span style="text-transform: capitalize;">{{.Invoice.Status}}</span>
                </div>
            </div>
        </div>

        <!-- Items Table -->
        <table class="items-table">
            <thead>
                <tr>
                    <th style="width: 40%;">Item / البند</th>
                    <th style="width: 10%;">Qty / الكمية</th>
                    <th style="width: 15%;">Unit Price / السعر</th>
                    <th style="width: 10%;">VAT % / ضريبة</th>
                    <th style="width: 15%;">VAT Amount / مبلغ الضريبة</th>
                    <th style="width: 15%;">Total / الإجمالي</th>
                </tr>
            </thead>
            <tbody>
                {{range .Items}}
                <tr>
                    <td class="item-description">
                        <div>{{if .Product.Name}}{{.Product.Name}}{{else}}Product{{end}}</div>
                        {{if .Product.NameArabic}}<div class="item-description-arabic">{{.Product.NameArabic}}</div>{{end}}
                    </td>
                    <td>{{printf "%.2f" .Quantity}}</td>
                    <td>{{printf "%.2f" .UnitPrice}}</td>
                    <td>{{printf "%.1f%%" .VATRate}}</td>
                    <td>{{printf "%.2f" .VATAmount}}</td>
                    <td>{{printf "%.2f" .TotalAmount}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>

        <!-- Totals -->
        <div class="totals">
            <div class="total-row">
                <span>Subtotal:</span>
                <span>{{printf "%.2f" .Invoice.SubTotal}}</span>
            </div>
            <div class="total-row">
                <span>VAT Amount:</span>
                <span>{{printf "%.2f" .Invoice.VATAmount}}</span>
            </div>
            <div class="total-row final">
                <span>Total Amount:</span>
                <span>{{printf "%.2f" .Invoice.TotalAmount}}</span>
            </div>
        </div>

        <!-- Notes -->
        {{if or .Invoice.Notes .Invoice.NotesArabic}}
        <div class="notes">
            <div class="notes-title">Notes / ملاحظات</div>
            {{if .Invoice.Notes}}
            <div>{{.Invoice.Notes}}</div>
            {{end}}
            {{if .Invoice.NotesArabic}}
            <div class="notes-arabic">{{.Invoice.NotesArabic}}</div>
            {{end}}
        </div>
        {{end}}

        <!-- Footer -->
        <div class="footer">
            <p>Thank you for your business! | شكراً لتعاملكم معنا!</p>
            <p>Generated on {{.Invoice.CreatedAt.Format "2006-01-02 15:04:05"}}</p>
        </div>
    </div>
</body>
</html>
`