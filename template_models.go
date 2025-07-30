package main

import (
	"dijibill/database"
)

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
