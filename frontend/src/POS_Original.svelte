<script>
  import { onMount } from 'svelte'
  import { GetProducts, GetProductCategories, GetSalesCategories, GetCustomers, CreateSalesInvoice, UpdateSalesInvoice, GetSalesInvoiceByID } from '../wailsjs/go/main/App.js'
  import { database } from '../wailsjs/go/models'
  import Modal from './components/Modal.svelte'
  import FormField from './components/FormField.svelte'

  // State management
  let currentView = 'categories' // 'categories' or 'products'
  let selectedCategory = null
  let categories = [] // Product categories for display
  let salesCategories = [] // Sales categories for invoice classification
  let selectedSalesCategory = null // Persistent sales category selection
  let products = []
  let filteredProducts = []
  let customers = []
  
  // Current sale state
  let currentSale = {
    items: [],
    customer_id: null,
    table_number: null,
    subtotal: 0,
    vat_total: 0,
    total: 0,
    status: 'pending'
  }

  // UI state
  let showCustomerModal = false
  let showTableModal = false
  let showSalesCategoryModal = false
  let showTransferModal = false
  let showRefundModal = false
  let loading = false
  let searchQuery = ''
  let transferToInvoiceId = null
  let refundReason = ''

  // Load initial data
  onMount(async () => {
    try {
      loading = true
      const [categoriesResult, salesCategoriesResult, customersResult] = await Promise.all([
        GetProductCategories(),
        GetSalesCategories(),
        GetCustomers()
      ])
      categories = categoriesResult || []
      salesCategories = salesCategoriesResult || []
      customers = customersResult || []
      
      // Set default sales category if available
      const defaultSalesCategory = salesCategories.find(cat => cat.is_default)
      if (defaultSalesCategory) {
        selectedSalesCategory = defaultSalesCategory
      } else if (salesCategories.length > 0) {
        selectedSalesCategory = salesCategories[0]
      }
    } catch (error) {
      console.error('Error loading data:', error)
    } finally {
      loading = false
    }
  })

  // Load products for selected category
  async function selectCategory(category) {
    try {
      selectedCategory = category
      currentView = 'products'
      loading = true
      
      const productsResult = await GetProducts()
      products = productsResult || []
      
      // Filter products by category
      filteredProducts = products.filter(product => 
        product.category_id === category.id
      )
    } catch (error) {
      console.error('Error loading products:', error)
      filteredProducts = []
    } finally {
      loading = false
    }
  }

  // Go back to categories view
  function backToCategories() {
    currentView = 'categories'
    selectedCategory = null
    filteredProducts = []
    searchQuery = ''
  }

  // Add product to current sale
  function addProductToSale(product) {
    const existingItem = currentSale.items.find(item => item.product_id === product.id)
    
    if (existingItem) {
      existingItem.quantity += 1
    } else {
      currentSale.items.push({
        product_id: product.id,
        product_name: product.name,
        quantity: 1,
        unit_price: product.price || 0,
        vat_rate: product.vat_rate || 15,
        total: product.price || 0
      })
    }
    
    calculateTotals()
  }

  // Remove item from sale
  function removeItem(index) {
    currentSale.items.splice(index, 1)
    currentSale.items = [...currentSale.items]
    calculateTotals()
  }

  // Update item quantity
  function updateQuantity(index, quantity) {
    if (quantity <= 0) {
      removeItem(index)
      return
    }
    
    currentSale.items[index].quantity = quantity
    calculateTotals()
  }

  // Calculate totals
  function calculateTotals() {
    currentSale.subtotal = 0
    currentSale.vat_total = 0
    
    currentSale.items.forEach(item => {
      const lineTotal = item.quantity * item.unit_price
      currentSale.subtotal += lineTotal
      
      const vatAmount = (lineTotal * item.vat_rate) / 100
      currentSale.vat_total += vatAmount
      
      item.total = lineTotal + vatAmount
    })
    
    currentSale.total = currentSale.subtotal + currentSale.vat_total
  }

  // Save current sale
  async function saveSale() {
    if (currentSale.items.length === 0) {
      alert('Please add items to the sale')
      return
    }

    if (!selectedSalesCategory) {
      alert('Please select a sales category')
      return
    }

    try {
      loading = true
      
      const invoiceObj = {
        id: 0,
        invoice_number: '',
        customer_id: currentSale.customer_id || 1,
        sales_category_id: selectedSalesCategory.id,
        issue_date: new Date().toISOString().split('T')[0],
        due_date: new Date().toISOString().split('T')[0],
        sub_total: currentSale.subtotal,
        vat_amount: currentSale.vat_total,
        total_amount: currentSale.total,
        status: 'completed',
        notes: currentSale.table_number ? `Table: ${currentSale.table_number}` : '',
        notes_arabic: '',
        qr_code: '',
        items: currentSale.items.map(item => ({
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

      await CreateSalesInvoice(new database.SalesInvoice(invoiceObj))
      
      // Reset current sale
      currentSale = {
        items: [],
        customer_id: null,
        table_number: null,
        subtotal: 0,
        vat_total: 0,
        total: 0,
        status: 'pending'
      }
      
      alert('Sale saved successfully!')
    } catch (error) {
      console.error('Error saving sale:', error)
      alert('Error saving sale')
    } finally {
      loading = false
    }
  }

  // Search products
  $: if (searchQuery && currentView === 'products') {
    filteredProducts = products.filter(product => 
      product.category_id === selectedCategory?.id &&
      product.name.toLowerCase().includes(searchQuery.toLowerCase())
    )
  } else if (currentView === 'products' && selectedCategory) {
    filteredProducts = products.filter(product => 
      product.category_id === selectedCategory.id
    )
  }

  // Get customer name
  function getCustomerName(customerId) {
    const customer = customers.find(c => c.id === customerId)
    return customer ? customer.name : 'Walk-in Customer'
  }

  // Refund current sale
  async function refundSale() {
    if (currentSale.items.length === 0) {
      alert('No items to refund')
      return
    }

    try {
      loading = true
      
      // Create a refund invoice
      const refundInvoiceObj = {
        id: 0,
        invoice_number: '',
        customer_id: currentSale.customer_id || 1,
        sales_category_id: selectedSalesCategory?.id || 1, // Use current sales category or default
        issue_date: new Date().toISOString().split('T')[0],
        due_date: new Date().toISOString().split('T')[0],
        sub_total: -currentSale.subtotal, // Negative amounts for refund
        vat_amount: -currentSale.vat_total,
        total_amount: -currentSale.total,
        status: 'refunded',
        notes: `Refund - ${refundReason}${currentSale.table_number ? ` | Table: ${currentSale.table_number}` : ''}`,
        notes_arabic: '',
        qr_code: '',
        items: currentSale.items.map(item => ({
          id: 0,
          invoice_id: 0,
          product_id: item.product_id,
          quantity: -item.quantity, // Negative quantity for refund
          unit_price: item.unit_price,
          vat_rate: item.vat_rate,
          vat_amount: 0,
          total_amount: 0
        }))
      }

      await CreateSalesInvoice(new database.SalesInvoice(refundInvoiceObj))
      
      // Reset current sale
      currentSale = {
        items: [],
        customer_id: null,
        table_number: null,
        subtotal: 0,
        vat_total: 0,
        total: 0,
        status: 'pending'
      }
      
      showRefundModal = false
      refundReason = ''
      alert('Refund processed successfully!')
    } catch (error) {
      console.error('Error processing refund:', error)
      alert('Error processing refund')
    } finally {
      loading = false
    }
  }

  // Transfer items to another invoice
  async function transferItems() {
    if (currentSale.items.length === 0) {
      alert('No items to transfer')
      return
    }

    if (!transferToInvoiceId) {
      alert('Please specify the invoice ID to transfer to')
      return
    }

    try {
      loading = true
      
      // For now, we'll create a new invoice with the transferred items
      // In a real implementation, you'd update the existing invoice
      const transferInvoiceObj = {
        id: 0,
        invoice_number: '',
        customer_id: currentSale.customer_id || 1,
        sales_category_id: selectedSalesCategory?.id || 1, // Use current sales category or default
        issue_date: new Date().toISOString().split('T')[0],
        due_date: new Date().toISOString().split('T')[0],
        sub_total: currentSale.subtotal,
        vat_amount: currentSale.vat_total,
        total_amount: currentSale.total,
        status: 'transferred',
        notes: `Transferred from POS${currentSale.table_number ? ` | Table: ${currentSale.table_number}` : ''}`,
        notes_arabic: '',
        qr_code: '',
        items: currentSale.items.map(item => ({
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
      
      // Reset current sale
      currentSale = {
        items: [],
        customer_id: null,
        table_number: null,
        subtotal: 0,
        vat_total: 0,
        total: 0,
        status: 'pending'
      }
      
      showTransferModal = false
      transferToInvoiceId = null
      alert('Items transferred successfully!')
    } catch (error) {
      console.error('Error transferring items:', error)
      alert('Error transferring items')
    } finally {
      loading = false
    }
  }
</script>

<div class="h-screen bg-gray-900 text-white flex">
  <!-- Left Panel - Current Sale -->
  <div class="w-1/3 bg-gray-800 flex flex-col">
    <!-- Header with action buttons -->
    <div class="p-4 border-b border-gray-700">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-bold">Current Sale</h2>
        <div class="flex space-x-2">
          <!-- Table Selection -->
          <button
            class="p-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
            on:click={() => showTableModal = true}
            title="Select Table"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M3 14h18m-9-4v8m-7 0h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
            </svg>
          </button>
          
          <!-- Sales Category Selection -->
          <button
            class="p-2 bg-purple-600 hover:bg-purple-700 rounded-lg transition-colors"
            on:click={() => showSalesCategoryModal = true}
            title="Select Sales Category"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"></path>
            </svg>
          </button>
          
          <!-- Customer Selection -->
          <button
            class="p-2 bg-green-600 hover:bg-green-700 rounded-lg transition-colors"
            on:click={() => showCustomerModal = true}
            title="Select Customer"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
            </svg>
          </button>
        </div>
      </div>

      <!-- Customer, Table, and Sales Category Info -->
      <div class="text-sm text-gray-300">
        <div>Customer: {getCustomerName(currentSale.customer_id)}</div>
        {#if currentSale.table_number}
          <div>Table: {currentSale.table_number}</div>
        {/if}
        <div>Category: {selectedSalesCategory?.name || 'None'}</div>
      </div>
    </div>

    <!-- Items List -->
    <div class="flex-1 overflow-y-auto p-4">
      {#if currentSale.items.length === 0}
        <div class="text-center text-gray-400 mt-8">
          <p>No Items</p>
          <p class="text-sm">Select products to add to sale</p>
        </div>
      {:else}
        <div class="space-y-2">
          {#each currentSale.items as item, index}
            <div class="bg-gray-700 p-3 rounded-lg">
              <div class="flex justify-between items-start mb-2">
                <div class="flex-1">
                  <h4 class="font-medium">{item.product_name}</h4>
                  <p class="text-sm text-gray-300">{item.unit_price.toFixed(2)} SAR</p>
                </div>
                <button
                  class="text-red-400 hover:text-red-300"
                  on:click={() => removeItem(index)}
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
              </div>
              <div class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                  <button
                    class="w-8 h-8 bg-gray-600 hover:bg-gray-500 rounded flex items-center justify-center"
                    on:click={() => updateQuantity(index, item.quantity - 1)}
                  >
                    -
                  </button>
                  <span class="w-8 text-center">{item.quantity}</span>
                  <button
                    class="w-8 h-8 bg-gray-600 hover:bg-gray-500 rounded flex items-center justify-center"
                    on:click={() => updateQuantity(index, item.quantity + 1)}
                  >
                    +
                  </button>
                </div>
                <div class="text-right">
                  <div class="font-medium">{item.total.toFixed(2)} SAR</div>
                </div>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Totals -->
    <div class="p-4 border-t border-gray-700">
      <div class="space-y-2 mb-4">
        <div class="flex justify-between">
          <span>Subtotal:</span>
          <span>{currentSale.subtotal.toFixed(2)} SAR</span>
        </div>
        <div class="flex justify-between">
          <span>VAT:</span>
          <span>{currentSale.vat_total.toFixed(2)} SAR</span>
        </div>
        <div class="flex justify-between text-xl font-bold border-t border-gray-600 pt-2">
          <span>TOTAL:</span>
          <span>{currentSale.total.toFixed(2)} SAR</span>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="grid grid-cols-2 gap-2">
        <button
          class="bg-red-600 hover:bg-red-700 text-white py-3 px-4 rounded-lg font-medium transition-colors"
          disabled={currentSale.items.length === 0}
          on:click={() => showRefundModal = true}
          title="Refund"
        >
          Refund
        </button>
        <button
          class="bg-blue-600 hover:bg-blue-700 text-white py-3 px-4 rounded-lg font-medium transition-colors"
          disabled={currentSale.items.length === 0}
          on:click={() => showTransferModal = true}
          title="Transfer"
        >
          Transfer
        </button>
        <button
          class="bg-green-600 hover:bg-green-700 text-white py-3 px-4 rounded-lg font-medium transition-colors col-span-2"
          disabled={currentSale.items.length === 0 || loading}
          on:click={saveSale}
        >
          {loading ? 'Saving...' : 'Save Sale'}
        </button>
      </div>
    </div>
  </div>

  <!-- Right Panel - Products -->
  <div class="flex-1 bg-gray-900 flex flex-col">
    <!-- Header -->
    <div class="p-4 border-b border-gray-700">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          {#if currentView === 'products'}
            <button
              class="p-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
              on:click={backToCategories}
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
              </svg>
            </button>
            <h2 class="text-xl font-bold">{selectedCategory?.name || 'Products'}</h2>
          {:else}
            <h2 class="text-xl font-bold">Product Categories</h2>
          {/if}
        </div>

        <!-- Search (only show in products view) -->
        {#if currentView === 'products'}
          <div class="relative">
            <input
              type="text"
              bind:value={searchQuery}
              placeholder="Search products..."
              class="w-64 pl-4 pr-10 py-2 bg-gray-800 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
              <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
              </svg>
            </div>
          </div>
        {/if}

        <!-- Placeholder for language change -->
        <button
          class="p-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors opacity-50 cursor-not-allowed"
          disabled
          title="Change Language (Coming Soon)"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129"></path>
          </svg>
        </button>
      </div>
    </div>

    <!-- Content Area -->
    <div class="flex-1 p-6 overflow-y-auto">
      {#if loading}
        <div class="flex items-center justify-center h-64">
          <div class="text-center">
            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
            <p>Loading...</p>
          </div>
        </div>
      {:else if currentView === 'categories'}
        <!-- Categories Grid -->
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {#each categories as category}
            <button
              class="aspect-square bg-gradient-to-br from-blue-500 to-blue-700 hover:from-blue-600 hover:to-blue-800 rounded-lg p-4 text-white font-medium text-lg transition-all transform hover:scale-105 shadow-lg"
              on:click={() => selectCategory(category)}
            >
              <div class="h-full flex flex-col items-center justify-center">
                <div class="text-3xl mb-2">📦</div>
                <div class="text-center">{category.name}</div>
              </div>
            </button>
          {/each}
        </div>
      {:else}
        <!-- Products Grid -->
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {#each filteredProducts as product}
            <button
              class="aspect-square bg-gradient-to-br from-green-500 to-green-700 hover:from-green-600 hover:to-green-800 rounded-lg p-4 text-white font-medium transition-all transform hover:scale-105 shadow-lg"
              on:click={() => addProductToSale(product)}
            >
              <div class="h-full flex flex-col items-center justify-center">
                <div class="text-2xl mb-2">🛍️</div>
                <div class="text-center text-sm mb-1">{product.name}</div>
                <div class="text-lg font-bold">{(product.price || 0).toFixed(2)} SAR</div>
              </div>
            </button>
          {/each}
        </div>
      {/if}
    </div>
  </div>
</div>

<!-- Customer Selection Modal -->
<Modal
  show={showCustomerModal}
  title="Select Customer"
  size="md"
  on:close={() => showCustomerModal = false}
  on:secondary={() => showCustomerModal = false}
>
  <div class="space-y-4">
    <FormField
      label="Customer"
      type="select"
      bind:value={currentSale.customer_id}
      options={customers.map(c => ({ value: c.id, label: c.name }))}
      placeholder="Select customer"
    />
    <div class="flex justify-end space-x-2">
      <button
        class="btn btn-secondary"
        on:click={() => showCustomerModal = false}
      >
        Cancel
      </button>
      <button
        class="btn btn-primary"
        on:click={() => showCustomerModal = false}
      >
        Select
      </button>
    </div>
  </div>
</Modal>

<!-- Table Selection Modal -->
<Modal
  show={showTableModal}
  title="Select Table"
  size="md"
  on:close={() => showTableModal = false}
  on:secondary={() => showTableModal = false}
>
  <div class="space-y-4">
    <FormField
      label="Table Number"
      type="text"
      bind:value={currentSale.table_number}
      placeholder="Enter table number"
    />
    <div class="flex justify-end space-x-2">
      <button
        class="btn btn-secondary"
        on:click={() => showTableModal = false}
      >
        Cancel
      </button>
      <button
        class="btn btn-primary"
        on:click={() => showTableModal = false}
      >
        Select
      </button>
    </div>
  </div>
</Modal>

<!-- Sales Category Selection Modal -->
<Modal
  show={showSalesCategoryModal}
  title="Select Sales Category"
  size="md"
  on:close={() => showSalesCategoryModal = false}
  on:secondary={() => showSalesCategoryModal = false}
>
  <div class="space-y-4">
    <FormField
      label="Sales Category"
      type="select"
      bind:value={selectedSalesCategory}
      options={salesCategories.map(c => ({ value: c, label: c.name }))}
      placeholder="Select sales category"
    />
    <div class="flex justify-end space-x-2">
      <button
        class="btn btn-secondary"
        on:click={() => showSalesCategoryModal = false}
      >
        Cancel
      </button>
      <button
        class="btn btn-primary"
        on:click={() => showSalesCategoryModal = false}
      >
        Select
      </button>
    </div>
  </div>
</Modal>

<!-- Transfer Modal -->
<Modal
  show={showTransferModal}
  title="Transfer Items"
  size="md"
  on:close={() => showTransferModal = false}
  on:secondary={() => showTransferModal = false}
>
  <div class="space-y-4">
    <FormField
      label="Transfer to Invoice ID"
      type="text"
      bind:value={transferToInvoiceId}
      placeholder="Enter invoice ID"
    />
    <div class="flex justify-end space-x-2">
      <button
        class="btn btn-secondary"
        on:click={() => showTransferModal = false}
      >
        Cancel
      </button>
      <button
        class="btn btn-primary"
        on:click={transferItems}
        disabled={!transferToInvoiceId}
      >
        Transfer
      </button>
    </div>
  </div>
</Modal>

<!-- Refund Modal -->
<Modal
  show={showRefundModal}
  title="Refund Sale"
  size="md"
  on:close={() => showRefundModal = false}
  on:secondary={() => showRefundModal = false}
>
  <div class="space-y-4">
    <FormField
      label="Refund Reason"
      type="text"
      bind:value={refundReason}
      placeholder="Enter refund reason"
    />
    <div class="flex justify-end space-x-2">
      <button
        class="btn btn-secondary"
        on:click={() => showRefundModal = false}
      >
        Cancel
      </button>
      <button
        class="btn btn-danger"
        on:click={refundSale}
        disabled={!refundReason}
      >
        Process Refund
      </button>
    </div>
  </div>
</Modal>