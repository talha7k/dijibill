import { writable, derived } from 'svelte/store'

// Current sale state
export const currentSale = writable({
  items: [],
  customer_id: null,
  table_number: null,
  subtotal: 0,
  vat_total: 0,
  total: 0,
  status: 'pending'
})

// Data stores
export const categories = writable([])
export const salesCategories = writable([])
export const selectedSalesCategory = writable(null)
export const products = writable([])
export const customers = writable([])
export const openInvoices = writable([])
export const paymentTypes = writable([])

// UI state
export const currentView = writable('categories') // categories, products
export const selectedCategory = writable(null)
export const loading = writable(false)
export const searchQuery = writable('')

// Modal states
export const showCustomerModal = writable(false)
export const showTableModal = writable(false)
export const showSalesCategoryModal = writable(false)
export const showTransferModal = writable(false)
export const showRefundModal = writable(false)
export const showPaymentModal = writable(false)
export const showInvoiceSelectionModal = writable(false)

// Transfer and refund states
export const transferToInvoiceId = writable(null)
export const selectedTransferInvoice = writable(null)
export const refundReason = writable('')

// Invoice selection states
export const pastInvoices = writable([])
export const selectedInvoice = writable(null)
export const invoiceSelectionTab = writable('open') // 'open' or 'past'
export const invoiceDateFilter = writable({
  start_date: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0], // 30 days ago
  end_date: new Date().toISOString().split('T')[0] // today
})

// Payment states
export const currentInvoiceForPayment = writable(null)
export const paymentItems = writable([])
export const remainingAmount = writable(0)

// Derived stores
export const filteredProducts = derived(
  [products, selectedCategory, searchQuery, currentView],
  ([$products, $selectedCategory, $searchQuery, $currentView]) => {
    if ($currentView !== 'products' || !$selectedCategory) return []
    
    let filtered = $products.filter(product => 
      product.category_id === $selectedCategory.id
    )
    
    if ($searchQuery) {
      filtered = filtered.filter(product =>
        product.name.toLowerCase().includes($searchQuery.toLowerCase())
      )
    }
    
    return filtered
  }
)

// Helper functions
export function calculateTotals() {
  currentSale.update(sale => {
    sale.subtotal = 0
    sale.vat_total = 0
    
    sale.items.forEach(item => {
      const lineTotal = item.quantity * item.unit_price
      sale.subtotal += lineTotal
      
      const vatAmount = (lineTotal * item.vat_rate) / 100
      sale.vat_total += vatAmount
      
      item.total = lineTotal + vatAmount
    })
    
    sale.total = sale.subtotal + sale.vat_total
    return sale
  })
}

export function addProductToSale(product) {
  currentSale.update(sale => {
    const existingItem = sale.items.find(item => item.product_id === product.id)
    
    if (existingItem) {
      existingItem.quantity += 1
    } else {
      sale.items.push({
        product_id: product.id,
        product_name: product.name,
        quantity: 1,
        unit_price: product.unit_price || 0,
        vat_rate: product.vat_rate || 15,
        total: product.unit_price || 0
      })
    }
    
    return sale
  })
  calculateTotals()
}

export function removeItem(index) {
  currentSale.update(sale => {
    sale.items.splice(index, 1)
    return sale
  })
  calculateTotals()
}

export function updateQuantity(index, quantity) {
  if (quantity <= 0) {
    removeItem(index)
    return
  }
  
  currentSale.update(sale => {
    sale.items[index].quantity = quantity
    return sale
  })
  calculateTotals()
}

export function resetSale() {
  currentSale.set({
    items: [],
    customer_id: null,
    table_number: null,
    subtotal: 0,
    vat_total: 0,
    total: 0,
    status: 'pending'
  })
}

export function getCustomerName(customerId, customersList) {
  const customer = customersList.find(c => c.id === customerId)
  return customer ? customer.name : 'Walk-in Customer'
}