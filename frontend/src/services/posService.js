import { 
  GetProducts, 
  GetProductCategories, 
  GetSalesCategories, 
  GetCustomers, 
  CreateSalesInvoice 
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
  refundReason,
  transferToInvoiceId
} from '../stores/posStore.js'
import { get } from 'svelte/store'

export async function loadInitialData() {
  try {
    loading.set(true)
    const [categoriesResult, salesCategoriesResult, customersResult] = await Promise.all([
      GetProductCategories(),
      GetSalesCategories(),
      GetCustomers()
    ])
    
    categories.set(categoriesResult || [])
    salesCategories.set(salesCategoriesResult || [])
    customers.set(customersResult || [])
    
    // Set default sales category if available
    const salesCats = salesCategoriesResult || []
    const defaultSalesCategory = salesCats.find(cat => cat.is_default)
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
  const transferId = get(transferToInvoiceId)
  
  if (sale.items.length === 0) {
    alert('No items to transfer')
    return
  }

  if (!transferId) {
    alert('Please specify the invoice ID to transfer to')
    return
  }

  try {
    loading.set(true)
    
    const transferInvoiceObj = {
      id: 0,
      invoice_number: '',
      customer_id: sale.customer_id || 1,
      sales_category_id: selectedSalesCat?.id || 1,
      issue_date: new Date().toISOString().split('T')[0],
      due_date: new Date().toISOString().split('T')[0],
      sub_total: sale.subtotal,
      vat_amount: sale.vat_total,
      total_amount: sale.total,
      status: 'transferred',
      notes: `Transferred from POS${sale.table_number ? ` | Table: ${sale.table_number}` : ''}`,
      notes_arabic: '',
      qr_code: '',
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
    transferToInvoiceId.set(null)
    alert('Items transferred successfully!')
  } catch (error) {
    console.error('Error transferring items:', error)
    alert('Error transferring items')
  } finally {
    loading.set(false)
  }
}