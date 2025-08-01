import { 
  GetProducts, 
  GetProductCategories, 
  GetSalesCategories, 
  GetCustomers, 
  CreateSalesInvoice,
  UpdateSalesInvoice,
  GetOpenSalesInvoices,
  GetSalesInvoices,
  GetPaymentTypes,
  CreatePayment,
  GetPaymentsByInvoiceID,
  UpdatePayment,
  GetSalesInvoiceByID
} from '../../wailsjs/go/main/App.js'
import { showDbSuccess, showDbError } from '../stores/notificationStore.js'
import { database } from '../../wailsjs/go/models'
import { 
  categories,
  salesCategories,
  customers,
  products,
  selectedSalesCategory,
  currentView,
  selectedCategory,
  loading,
  currentSale,
  resetSale,
  showRefundModal,
  showTransferModal,
  showPaymentModal,
  transferToInvoiceId,
  selectedTransferInvoice,
  refundReason,
  openInvoices,
  pastInvoices,
  paymentTypes,
  currentInvoiceForPayment,
  paymentItems,
  remainingAmount,
  selectedInvoice,
  invoiceDateFilter
} from '../stores/posStore.js'
import { get } from 'svelte/store'

export async function loadInitialData() {
  try {
    loading.set(true)
    const [categoriesResult, salesCategoriesResult, customersResult, paymentTypesResult] = await Promise.all([
      GetProductCategories(),
      GetSalesCategories(),
      GetCustomers(),
      GetPaymentTypes()
    ])
    
    categories.set(categoriesResult || [])
    salesCategories.set(salesCategoriesResult || [])
    customers.set(customersResult || [])
    paymentTypes.set(paymentTypesResult || [])
    
    // Set default sales category if available
    const salesCats = salesCategoriesResult || []
    const defaultSalesCategory = salesCats.find(cat => cat.name === 'POS' || cat.name === 'Default')
    if (defaultSalesCategory) {
      selectedSalesCategory.set(defaultSalesCategory)
    } else if (salesCats.length > 0) {
      selectedSalesCategory.set(salesCats[0])
    }
  } catch (error) {
    console.error('Error loading data:', error)
  } finally {
    loading.set(false)
  }
}

export async function loadOpenInvoices() {
  try {
    loading.set(true)
    const result = await GetOpenSalesInvoices()
    openInvoices.set(result || [])
  } catch (error) {
    console.error('Error loading open invoices:', error)
    openInvoices.set([])
  } finally {
    loading.set(false)
  }
}

export async function selectCategory(category) {
  try {
    selectedCategory.set(category)
    currentView.set('products')
    loading.set(true)
    
    const productsResult = await GetProducts()
    products.set(productsResult || [])
  } catch (error) {
    console.error('Error loading products:', error)
    products.set([])
  } finally {
    loading.set(false)
  }
}

export async function saveSale() {
  try {
    loading.set(true)
    
    const sale = get(currentSale)
    const selectedCategory = get(selectedSalesCategory)
    
    if (!selectedCategory) {
      showDbError('Please select a sales category before saving the sale.')
      return
    }
    
    if (sale.items.length === 0) {
      showDbError('Cannot save an empty sale. Please add items first.')
      return
    }

    const invoiceData = {
      id: 0,
      invoice_number: '',
      customer_id: sale.customer_id || 1,
      sales_category_id: selectedCategory.id,
      table_number: sale.table_number || null,
      issue_date: new Date().toISOString().split('T')[0],
      due_date: '',
      sub_total: sale.subtotal,
      vat_amount: sale.vat_total,
      total_amount: sale.total,
      status: 'paid',
      notes: '',
      notes_arabic: '',
      qr_code: '',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      items: sale.items.map(item => ({
        product_id: item.product_id,
        quantity: item.quantity,
        unit_price: item.unit_price,
        vat_rate: item.vat_rate,
        vat_amount: (item.quantity * item.unit_price * item.vat_rate) / 100,
        total_amount: item.total
      }))
    }

    const invoiceObj = new database.SalesInvoice(invoiceData)
    const createdInvoice = await CreateSalesInvoice(invoiceObj)
    
    showDbSuccess('Sale saved successfully!')
    resetSale()
  } catch (error) {
    console.error('Error saving sale:', error)
    showDbError('Error saving sale: ' + error.message)
  } finally {
    loading.set(false)
  }
}

export async function refundSale() {
  try {
    loading.set(true)
    
    const sale = get(currentSale)
    const selectedCategory = get(selectedSalesCategory)
    const reason = get(refundReason)
    
    if (!selectedCategory) {
      showDbError('Please select a sales category before processing refund.')
      return
    }
    
    if (sale.items.length === 0) {
      showDbError('Cannot process refund for an empty sale.')
      return
    }

    const refundInvoiceData = {
      id: 0,
      invoice_number: '',
      customer_id: sale.customer_id || 1,
      sales_category_id: selectedCategory.id,
      table_number: sale.table_number || null,
      issue_date: new Date().toISOString().split('T')[0],
      due_date: '',
      sub_total: -sale.subtotal,
      vat_amount: -sale.vat_total,
      total_amount: -sale.total,
      status: 'refunded',
      notes: `Refund - ${reason}`,
      notes_arabic: '',
      qr_code: '',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      items: sale.items.map(item => ({
        product_id: item.product_id,
        quantity: -item.quantity,
        unit_price: item.unit_price,
        vat_rate: item.vat_rate,
        vat_amount: -(item.quantity * item.unit_price * item.vat_rate) / 100,
        total_amount: -item.total
      }))
    }

    const refundInvoiceObj = new database.SalesInvoice(refundInvoiceData)
    const createdRefundInvoice = await CreateSalesInvoice(refundInvoiceObj)
    
    showDbSuccess('Refund processed successfully!')
    resetSale()
    showRefundModal.set(false)
    refundReason.set('')
  } catch (error) {
    console.error('Error processing refund:', error)
    showDbError('Error processing refund: ' + error.message)
  } finally {
    loading.set(false)
  }
}

export async function transferItems() {
  const sale = get(currentSale)
  const selectedSalesCat = get(selectedSalesCategory)
  const selectedInvoice = get(selectedTransferInvoice)
  
  if (sale.items.length === 0) {
    showDbError('No items to transfer')
    return
  }

  if (!selectedInvoice) {
    showDbError('Please select an invoice to transfer to')
    return
  }

  try {
    loading.set(true)
    
    const transferInvoiceObj = {
      id: 0,
      invoice_number: '',
      customer_id: sale.customer_id || 1,
      sales_category_id: selectedSalesCat?.id || 1,
      table_number: sale.table_number || null,
      issue_date: new Date().toISOString().split('T')[0],
      due_date: new Date().toISOString().split('T')[0],
      sub_total: sale.subtotal,
      vat_amount: sale.vat_total,
      total_amount: sale.total,
      status: 'transferred',
      notes: `Transferred from POS${sale.table_number ? ` | Table: ${sale.table_number}` : ''} to Invoice #${selectedInvoice.invoice_number}`,
      notes_arabic: '',
      qr_code: '',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      items: sale.items.map(item => ({
        id: 0,
        invoice_id: 0,
        product_id: item.product_id,
        quantity: item.quantity,
        unit_price: item.unit_price,
        vat_rate: item.vat_rate,
        vat_amount: 0,
        total_amount: 0
      }))
    }

    const createdTransferInvoice = await CreateSalesInvoice(new database.SalesInvoice(transferInvoiceObj))
    resetSale()
    showTransferModal.set(false)
    selectedTransferInvoice.set(null)
    showDbSuccess(`Items transferred successfully to Invoice #${selectedInvoice.invoice_number}!`)
  } catch (error) {
    console.error('Error transferring items:', error)
    showDbError('Error transferring items')
  } finally {
    loading.set(false)
  }
}

export async function openPaymentModal() {
  const sale = get(currentSale)
  const selectedCategory = get(selectedSalesCategory)
  
  if (!selectedCategory) {
    showDbError('Please select a sales category before processing payment.')
    return
  }
  
  if (sale.items.length === 0) {
    showDbError('Cannot process payment for an empty sale.')
    return
  }

  try {
    loading.set(true)
    
    // Create a draft invoice first
    const invoiceData = {
      id: 0,
      invoice_number: '',
      customer_id: sale.customer_id || 1,
      sales_category_id: selectedCategory.id,
      table_number: sale.table_number || null,
      issue_date: new Date().toISOString().split('T')[0],
      due_date: '',
      sub_total: sale.subtotal,
      vat_amount: sale.vat_total,
      total_amount: sale.total,
      status: 'draft',
      notes: '',
      notes_arabic: '',
      qr_code: '',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      items: sale.items.map(item => ({
        product_id: item.product_id,
        quantity: item.quantity,
        unit_price: item.unit_price,
        vat_rate: item.vat_rate,
        vat_amount: (item.quantity * item.unit_price * item.vat_rate) / 100,
        total_amount: item.total
      }))
    }

    const invoiceObj = new database.SalesInvoice(invoiceData)
    const createdInvoice = await CreateSalesInvoice(invoiceObj)
    
    // Set up payment modal with the created invoice (which now has the correct ID)
    currentInvoiceForPayment.set(createdInvoice)
    remainingAmount.set(sale.total)
    paymentItems.set([])
    showPaymentModal.set(true)
  } catch (error) {
    console.error('Error creating invoice for payment:', error)
    showDbError('Error creating invoice: ' + error.message)
  } finally {
    loading.set(false)
  }
}

export function addPaymentItem(paymentTypeId, amount) {
  const currentItems = get(paymentItems)
  const currentRemaining = get(remainingAmount)
  const paymentTypesList = get(paymentTypes)
  
  const paymentType = paymentTypesList.find(pt => pt.id === paymentTypeId)
  if (!paymentType) {
    showDbError('Invalid payment type selected')
    return
  }
  
  if (amount <= 0) {
    showDbError('Payment amount must be greater than 0')
    return
  }
  
  if (amount > currentRemaining) {
    showDbError('Payment amount cannot exceed remaining amount')
    return
  }
  
  const newPaymentItem = {
    id: Date.now(), // temporary ID for UI
    payment_type_id: paymentTypeId,
    payment_type_name: paymentType.name,
    amount: amount
  }
  
  paymentItems.set([...currentItems, newPaymentItem])
  remainingAmount.set(currentRemaining - amount)
}

export function removePaymentItem(itemId) {
  const currentItems = get(paymentItems)
  const itemToRemove = currentItems.find(item => item.id === itemId)
  
  if (itemToRemove) {
    const newItems = currentItems.filter(item => item.id !== itemId)
    paymentItems.set(newItems)
    
    const currentRemaining = get(remainingAmount)
    remainingAmount.set(currentRemaining + itemToRemove.amount)
  }
}

export async function processPayments() {
  const invoice = get(currentInvoiceForPayment)
  const payments = get(paymentItems)
  const remaining = get(remainingAmount)
  
  if (!invoice) {
    showDbError('No invoice selected for payment')
    return
  }
  
  if (payments.length === 0) {
    showDbError('Please add at least one payment')
    return
  }
  
  try {
    loading.set(true)
    
    // Create all payments
    for (const payment of payments) {
      const paymentData = {
        invoice_id: invoice.id, // Now using the correct invoice ID
        payment_type_id: payment.payment_type_id,
        amount: payment.amount,
        payment_date: new Date().toISOString(),
        reference: '',
        notes: 'POS Payment',
        notes_arabic: '',
        status: 'completed'
      }
      
      await CreatePayment(new database.Payment(paymentData))
    }
    
    // Determine invoice status based on remaining amount
    let invoiceStatus = 'paid'
    if (remaining > 0.01) { // Allow for small rounding differences
      invoiceStatus = 'partially_paid'
    }
    
    // Update invoice status
    const updatedInvoice = {
      ...invoice,
      status: invoiceStatus,
      updated_at: new Date().toISOString()
    }
    
    await UpdateSalesInvoice(new database.SalesInvoice(updatedInvoice))
    
    const statusMessage = invoiceStatus === 'paid' ? 'Payment completed successfully!' : 'Partial payment processed successfully!'
    showDbSuccess(statusMessage)
    
    resetSale()
    showPaymentModal.set(false)
    currentInvoiceForPayment.set(null)
    paymentItems.set([])
    remainingAmount.set(0)
  } catch (error) {
    console.error('Error processing payments:', error)
    showDbError('Error processing payments: ' + error.message)
  } finally {
    loading.set(false)
  }
}

/**
 * Refunds an existing invoice by marking the invoice and all its payments as refunded
 * @param {number} invoiceId - The ID of the invoice to refund
 * @param {string} refundReason - The reason for the refund
 */
export async function refundInvoice(invoiceId, refundReason = '') {
  try {
    loading.set(true)
    
    // Get the full invoice details
    const invoice = await GetSalesInvoiceByID(invoiceId)
    if (!invoice) {
      showDbError('Invoice not found')
      return
    }
    
    // Get the invoice payments
    const payments = await GetPaymentsByInvoiceID(invoiceId)
    
    // Update all payments to refunded status
    for (const payment of payments) {
      const updatedPayment = {
        ...payment,
        status: 'refunded',
        notes: payment.notes ? `${payment.notes} | REFUNDED: ${refundReason}` : `REFUNDED: ${refundReason}`,
        updated_at: new Date().toISOString()
      }
      
      await UpdatePayment(new database.Payment(updatedPayment))
    }
    
    // Update the invoice with full details and refunded status
    const updatedInvoice = {
      ...invoice,
      status: 'refunded',
      notes: invoice.notes ? `${invoice.notes} | REFUNDED: ${refundReason}` : `REFUNDED: ${refundReason}`,
      updated_at: new Date().toISOString()
    }
    
    await UpdateSalesInvoice(new database.SalesInvoice(updatedInvoice))
    
    showDbSuccess('Invoice and all associated payments have been marked as refunded successfully!')
    
    // Reload open invoices if needed
    await loadOpenInvoices()
    
  } catch (error) {
    console.error('Error refunding invoice:', error)
    showDbError('Error refunding invoice: ' + error.message)
  } finally {
    loading.set(false)
  }
}

/**
 * Loads all past invoices for selection
 */
export async function loadPastInvoices() {
  try {
    loading.set(true)
    const invoices = await GetSalesInvoices()
    
    // Filter out open invoices and sort by date (newest first)
    const pastInvoicesList = invoices
      .filter(invoice => invoice.status !== 'pending' && invoice.status !== 'draft')
      .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    
    pastInvoices.set(pastInvoicesList)
  } catch (error) {
    console.error('Error loading past invoices:', error)
    showDbError('Error loading past invoices: ' + error.message)
  } finally {
    loading.set(false)
  }
}

/**
 * Filters past invoices by date range
 * @param {string} startDate - Start date in YYYY-MM-DD format
 * @param {string} endDate - End date in YYYY-MM-DD format
 */
export function filterPastInvoicesByDate(startDate, endDate) {
  const allPastInvoices = get(pastInvoices)
  
  if (!startDate && !endDate) {
    return allPastInvoices
  }
  
  return allPastInvoices.filter(invoice => {
    const invoiceDate = new Date(invoice.created_at).toISOString().split('T')[0]
    
    if (startDate && invoiceDate < startDate) return false
    if (endDate && invoiceDate > endDate) return false
    
    return true
  })
}

/**
 * Selects an invoice for actions (refund, reprint, etc.)
 * @param {Object} invoice - The invoice to select
 */
export function selectInvoice(invoice) {
  selectedInvoice.set(invoice)
}

/**
 * Reprints a receipt for the selected invoice
 * @param {number} invoiceId - The ID of the invoice to reprint
 */
export async function reprintReceipt(invoiceId) {
  try {
    loading.set(true)
    
    // Get the full invoice details
    const invoice = await GetSalesInvoiceByID(invoiceId)
    if (!invoice) {
      showDbError('Invoice not found')
      return
    }
    
    // For now, we'll just show a notification. In a real implementation,
    // this would trigger the receipt printing functionality
    showDbSuccess(`Reprinting receipt for Invoice #${invoice.invoice_number || invoice.id}`)
    
    // TODO: Implement actual receipt printing logic here
    // This could involve:
    // 1. Generating a receipt template
    // 2. Sending to a printer
    // 3. Opening a print dialog
    // 4. Generating a PDF
    
  } catch (error) {
    console.error('Error reprinting receipt:', error)
    showDbError('Error reprinting receipt: ' + error.message)
  } finally {
    loading.set(false)
  }
}