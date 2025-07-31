import { 
  GetProducts, 
  GetProductCategories, 
  GetSalesCategories, 
  GetCustomers, 
  CreateSalesInvoice,
  GetOpenSalesInvoices,
  GetPaymentTypes,
  CreatePayment
} from '../../wailsjs/go/main/App.js'
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
  paymentTypes,
  currentInvoiceForPayment,
  paymentItems,
  remainingAmount
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
      alert('Please select a sales category before saving the sale.')
      return
    }
    
    if (sale.items.length === 0) {
      alert('Cannot save an empty sale. Please add items first.')
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
    await CreateSalesInvoice(invoiceObj)
    
    alert('Sale saved successfully!')
    resetSale()
  } catch (error) {
    console.error('Error saving sale:', error)
    alert('Error saving sale: ' + error.message)
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
      alert('Please select a sales category before processing refund.')
      return
    }
    
    if (sale.items.length === 0) {
      alert('Cannot process refund for an empty sale.')
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
    await CreateSalesInvoice(refundInvoiceObj)
    
    alert('Refund processed successfully!')
    resetSale()
    showRefundModal.set(false)
    refundReason.set('')
  } catch (error) {
    console.error('Error processing refund:', error)
    alert('Error processing refund: ' + error.message)
  } finally {
    loading.set(false)
  }
}

export async function transferItems() {
  const sale = get(currentSale)
  const selectedSalesCat = get(selectedSalesCategory)
  const selectedInvoice = get(selectedTransferInvoice)
  
  if (sale.items.length === 0) {
    alert('No items to transfer')
    return
  }

  if (!selectedInvoice) {
    alert('Please select an invoice to transfer to')
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

    await CreateSalesInvoice(new database.SalesInvoice(transferInvoiceObj))
    resetSale()
    showTransferModal.set(false)
    selectedTransferInvoice.set(null)
    alert(`Items transferred successfully to Invoice #${selectedInvoice.invoice_number}!`)
  } catch (error) {
    console.error('Error transferring items:', error)
    alert('Error transferring items')
  } finally {
    loading.set(false)
  }
}

export async function openPaymentModal() {
  const sale = get(currentSale)
  const selectedCategory = get(selectedSalesCategory)
  
  if (!selectedCategory) {
    alert('Please select a sales category before processing payment.')
    return
  }
  
  if (sale.items.length === 0) {
    alert('Cannot process payment for an empty sale.')
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
    await CreateSalesInvoice(invoiceObj)
    
    // Set up payment modal
    currentInvoiceForPayment.set(invoiceData)
    remainingAmount.set(sale.total)
    paymentItems.set([])
    showPaymentModal.set(true)
  } catch (error) {
    console.error('Error creating invoice for payment:', error)
    alert('Error creating invoice: ' + error.message)
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
    alert('Invalid payment type selected')
    return
  }
  
  if (amount <= 0) {
    alert('Payment amount must be greater than 0')
    return
  }
  
  if (amount > currentRemaining) {
    alert('Payment amount cannot exceed remaining amount')
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
    alert('No invoice selected for payment')
    return
  }
  
  if (payments.length === 0) {
    alert('Please add at least one payment')
    return
  }
  
  if (remaining > 0.01) { // Allow for small rounding differences
    alert('Please complete all payments before processing')
    return
  }
  
  try {
    loading.set(true)
    
    // Create all payments
    for (const payment of payments) {
      const paymentData = {
        invoice_id: invoice.id || 0, // This will need to be the actual invoice ID
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
    
    alert('Payment processed successfully!')
    resetSale()
    showPaymentModal.set(false)
    currentInvoiceForPayment.set(null)
    paymentItems.set([])
    remainingAmount.set(0)
  } catch (error) {
    console.error('Error processing payments:', error)
    alert('Error processing payments: ' + error.message)
  } finally {
    loading.set(false)
  }
}