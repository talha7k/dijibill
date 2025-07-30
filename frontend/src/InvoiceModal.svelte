<script>
  import { createEventDispatcher, onMount } from 'svelte'
  import { GetCustomers, GetProducts, GetPaymentTypes, GetSalesCategories, CreateInvoice } from '../wailsjs/go/main/App.js'
  import {database} from '../wailsjs/go/models'
  
  const dispatch = createEventDispatcher()
  
  export let isOpen = false
  
  // Form data
  let customers = []
  let products = []
  let paymentTypes = []
  let salesCategories = []
  let isLoading = false
  let isSaving = false
  
  // Invoice form data
  let invoiceData = {
    customer_id: 0,
    sales_category_id: 0,
    issue_date: new Date().toISOString().split('T')[0],
    due_date: '',
    payment_terms: '',
    notes: ''
  }
  
  // Form validation errors
  let formErrors = {
    customer_id: '',
    sales_category_id: '',
    issue_date: '',
    items: []
  }
  
  // Add initial empty item
  let invoiceItems = [
    {
      product_id: 0,
      quantity: 1,
      unit_price: 0,
      vat_rate: 15,
      vat_amount: 0,
      total_amount: 0
    }
  ]
  
  // Calculated totals
  let subtotal = 0
  let discount = 0
  let deliveryFees = 0
  let vatTotal = 0
  let totalPrice = 0
  
  onMount(async () => {
    if (isOpen) {
      await loadData()
    }
  })
  
  $: if (isOpen) {
    loadData()
  }
  
  $: if (invoiceItems) {
    calculateTotals()
  }
  
  async function loadData() {
    isLoading = true
    try {
      const [customersResult, productsResult, paymentTypesResult, salesCategoriesResult] = await Promise.all([
        GetCustomers(),
        GetProducts(),
        GetPaymentTypes(),
        GetSalesCategories()
      ])
      customers = customersResult || []
      products = productsResult || []
      paymentTypes = paymentTypesResult || []
      salesCategories = salesCategoriesResult || []
    } catch (error) {
      console.error('Error loading data:', error)
      customers = []
      products = []
      paymentTypes = []
      salesCategories = []
    } finally {
      isLoading = false
    }
  }
  
  function addItem() {
    invoiceItems = [...invoiceItems, {
      product_id: 0,
      quantity: 1,
      unit_price: 0,
      vat_rate: 15,
      vat_amount: 0,
      total_amount: 0
    }]
  }
  
  function removeItem(index) {
    if (invoiceItems.length > 1) {
      invoiceItems = invoiceItems.filter((_, i) => i !== index)
    }
  }
  
  function calculateItemTotal(item) {
    const lineTotal = item.quantity * item.unit_price
    item.vat_amount = (lineTotal * item.vat_rate) / 100
    item.total_amount = lineTotal + item.vat_amount
    return item.total_amount
  }
  
  function calculateTotals() {
    subtotal = 0
    vatTotal = 0
    
    invoiceItems.forEach(item => {
      const lineTotal = item.quantity * item.unit_price
      subtotal += lineTotal
      item.vat_amount = (lineTotal * item.vat_rate) / 100
      vatTotal += item.vat_amount
      item.total_amount = lineTotal + item.vat_amount
    })
    
    totalPrice = subtotal - discount + deliveryFees + vatTotal
  }
  
  function onProductChange(index, productId) {
    if (!products || products.length === 0) return
    
    const product = products.find(p => p.id === parseInt(productId))
    if (product) {
      invoiceItems[index].unit_price = product.unit_price || 0
      invoiceItems = [...invoiceItems] // Trigger reactivity
    }
  }
  
  function handleProductChange(index, event) {
    onProductChange(index, event.target.value)
  }
  
  function validateForm() {
    console.log('validateForm called')
    // Reset errors
    formErrors = {
      customer_id: '',
      sales_category_id: '',
      issue_date: '',
      items: []
    }
    
    let isValid = true
    
    // Validate issue date
    if (!invoiceData.issue_date) {
      formErrors.issue_date = 'Please select an invoice date'
      isValid = false
      console.log('Issue date validation failed')
    }
    
    // Validate invoice items
    formErrors.items = []
    invoiceItems.forEach((item, index) => {
      let itemErrors = {
        product_id: '',
        quantity: '',
        unit_price: ''
      }
      
      if (!item.product_id || item.product_id === 0) {
        itemErrors.product_id = 'Please select a product'
        isValid = false
        console.log(`Item ${index + 1} product validation failed`)
      }
      
      if (!item.quantity || item.quantity <= 0) {
        itemErrors.quantity = 'Quantity must be greater than 0'
        isValid = false
        console.log(`Item ${index + 1} quantity validation failed`)
      }
      
      if (!item.unit_price || item.unit_price < 0) {
        itemErrors.unit_price = 'Unit price must be 0 or greater'
        isValid = false
        console.log(`Item ${index + 1} unit price validation failed`)
      }
      
      formErrors.items[index] = itemErrors
    })
    
    console.log('Final formErrors:', formErrors)
    console.log('Validation result:', isValid)
    return isValid
  }
  
  function clearErrors() {
    formErrors = {
      customer_id: '',
      sales_category_id: '',
      issue_date: '',
      items: []
    }
  }

  async function saveInvoice() {
    console.log('=== SAVE INVOICE START ===')
    console.log('saveInvoice called')
    console.log('Current invoiceData:', JSON.stringify(invoiceData, null, 2))
    console.log('Current invoiceItems:', JSON.stringify(invoiceItems, null, 2))
    
    // Validate form first
    if (!validateForm()) {
      console.log('Form validation failed:', formErrors)
      return
    }
    
    console.log('Form validation passed')
    
    // Check if any items have valid product_id
    const validItems = invoiceItems.filter(item => item.product_id > 0)
    console.log('Valid items:', validItems)
    if (validItems.length === 0) {
      console.log('No valid items found')
      alert('Please select products for your invoice items')
      return
    }
    
    isSaving = true
    try {
      console.log('Creating invoice object...')
      const invoiceObj = {
        id: 0,
        invoice_number: '',
        customer_id: invoiceData.customer_id ? parseInt(invoiceData.customer_id.toString()) : 0,
        sales_category_id: invoiceData.sales_category_id ? parseInt(invoiceData.sales_category_id.toString()) : 0,
        issue_date: new Date(invoiceData.issue_date),
        due_date: invoiceData.due_date ? new Date(invoiceData.due_date) : new Date(),
        sub_total: 0,
        vat_amount: 0,
        total_amount: 0,
        status: invoiceData.status || 'draft',
        notes: invoiceData.notes || '',
        notes_arabic: '',
        qr_code: '',
        created_at: new Date(),
        updated_at: new Date(),
        items: validItems.map(item => ({
          id: 0,
          invoice_id: 0,
          product_id: parseInt(item.product_id.toString()),
          quantity: parseFloat(item.quantity.toString()),
          unit_price: parseFloat(item.unit_price.toString()),
          vat_rate: parseFloat(item.vat_rate.toString()),
          vat_amount: 0,
          total_amount: 0,
          created_at: new Date()
        }))
      }
      
      console.log('invoiceObj created:', JSON.stringify(invoiceObj, null, 2))
      
      console.log('Calling CreateInvoice backend function...')
      const result = await CreateInvoice(new database.Invoice(invoiceObj))
      console.log('CreateInvoice result:', result)
      console.log('CreateInvoice completed successfully')
      
      alert('Invoice saved successfully!')
      dispatch('saved')
      closeModal()
    } catch (error) {
      console.error('=== ERROR SAVING INVOICE ===')
      console.error('Error type:', typeof error)
      console.error('Error message:', error.message)
      console.error('Error stack:', error.stack)
      console.error('Full error object:', error)
      alert('Error saving invoice: ' + (error.message || error.toString()))
    } finally {
      isSaving = false
      console.log('=== SAVE INVOICE END ===')
    }
  }
  
  function closeModal() {
    isOpen = false
    dispatch('close')
    resetForm()
  }
  
  function resetForm() {
    invoiceData = {
      customer_id: 0,
      sales_category_id: 0,
      issue_date: new Date().toISOString().split('T')[0],
      due_date: '',
      payment_terms: '',
      notes: '',
      status: 'draft'
    }
    invoiceItems = [{
      product_id: 0,
      quantity: 1,
      unit_price: 0,
      vat_rate: 15,
      vat_amount: 0,
      total_amount: 0
    }]
    discount = 0
    deliveryFees = 0
  }
</script>

{#if isOpen}
  <div class="modal-backdrop fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
    <div class="modal-content glass-card w-full max-w-6xl max-h-[90vh] overflow-y-auto">
      <div class="modal-header border-b border-white/20 p-6">
        <div class="flex items-center justify-between">
          <h2 class="text-2xl font-semibold text-white">New Invoice</h2>
          <button 
            class="btn-icon text-white/60 hover:text-white hover:bg-white/10"
            on:click={closeModal}
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
      </div>
      
      <div class="modal-body p-6">
        {#if isLoading}
          <div class="flex items-center justify-center py-8">
            <div class="loading loading-spinner loading-lg text-primary"></div>
          </div>
        {:else}
          <!-- Error Display -->
          {#if formErrors.issue_date || formErrors.items.some(item => item.product_id || item.quantity || item.unit_price)}
            <div class="bg-error/10 border border-error/20 rounded-lg p-4 mb-6">
              <div class="flex items-start gap-3">
                <svg class="w-5 h-5 text-error mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                <div>
                  <h4 class="text-error font-medium mb-2">Please fix the following errors:</h4>
                  <ul class="text-error text-sm space-y-1">
                    {#if formErrors.issue_date}
                      <li>• {formErrors.issue_date}</li>
                    {/if}
                    {#each formErrors.items as itemError, index}
                      {#if itemError && (itemError.product_id || itemError.quantity || itemError.unit_price)}
                        {#if itemError.product_id}
                          <li>• Item {index + 1}: {itemError.product_id}</li>
                        {/if}
                        {#if itemError.quantity}
                          <li>• Item {index + 1}: {itemError.quantity}</li>
                        {/if}
                        {#if itemError.unit_price}
                          <li>• Item {index + 1}: {itemError.unit_price}</li>
                        {/if}
                      {/if}
                    {/each}
                  </ul>
                </div>
              </div>
            </div>
          {/if}
          
          <form on:submit|preventDefault={saveInvoice} class="space-y-6">
            <!-- Customer and Invoice Details -->
            <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
              <!-- Customer Selection -->
              <div class="form-group">
                <label for="customer-select" class="label-standard">Customer <span class="text-white/50">(Optional)</span></label>
                <select 
                  id="customer-select"
                  bind:value={invoiceData.customer_id}
                  class="select-glass w-full"
                >
                  <option value={0}>Select Customer</option>
                  {#each customers as customer}
                    <option value={customer.id}>{customer.name}</option>
                  {/each}
                </select>
                {#if formErrors.customer_id}
                  <p class="text-error text-sm mt-1">{formErrors.customer_id}</p>
                {/if}
              </div>
              
              <!-- Sales Category Selection -->
              <div class="form-group">
                <label for="sales-category-select" class="label-standard">Sales Category <span class="text-white/50">(Optional)</span></label>
                <select 
                  id="sales-category-select"
                  bind:value={invoiceData.sales_category_id}
                  class="select-glass w-full"
                >
                  <option value={0}>Select Sales Category</option>
                  {#each salesCategories as category}
                    <option value={category.id}>{category.name}</option>
                  {/each}
                </select>
                {#if formErrors.sales_category_id}
                  <p class="text-error text-sm mt-1">{formErrors.sales_category_id}</p>
                {/if}
              </div>
              
              <!-- Invoice Date -->
              <div class="form-group">
                <label for="invoice-date" class="label-standard">Invoice Date</label>
                <input 
                  id="invoice-date"
                  type="date"
                  bind:value={invoiceData.issue_date}
                  class="input-glass w-full"
                  required
                />
                {#if formErrors.issue_date}
                  <p class="text-error text-sm mt-1">{formErrors.issue_date}</p>
                {/if}
              </div>
              
              <!-- Due Date (Optional) -->
              <div class="form-group">
                <label for="due-date" class="label-standard">Due Date <span class="text-white/50">(Optional)</span></label>
                <input 
                  id="due-date"
                  type="date"
                  bind:value={invoiceData.due_date}
                  class="input-glass w-full"
                />
              </div>
              
              <!-- Status -->
              <div class="form-group">
                <label for="status-select" class="label-standard">Status</label>
                <select 
                  id="status-select"
                  bind:value={invoiceData.status}
                  class="select-glass w-full"
                >
                  <option value="draft">Draft</option>
                  <option value="sent">Sent</option>
                  <option value="paid">Paid</option>
                  <option value="overdue">Overdue</option>
                  <option value="cancelled">Cancelled</option>
                </select>
              </div>
            </div>
            
            <!-- Payment Terms (Optional) -->
            <div class="form-group">
              <label for="payment-terms" class="label-standard">Payment Terms <span class="text-white/50">(Optional)</span></label>
              <input 
                id="payment-terms"
                type="text"
                bind:value={invoiceData.payment_terms}
                placeholder="e.g., Net 30, Due on receipt"
                class="input-glass w-full"
              />
            </div>
            
            <!-- Invoice Items -->
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold text-white">Item Details</h3>
                <button 
                  type="button"
                  class="btn-primary btn-sm"
                  on:click={addItem}
                >
                  + Add Item
                </button>
              </div>
              
              <!-- Items Table Header -->
              <div class="hidden lg:grid lg:grid-cols-12 gap-4 p-4 bg-white/5 rounded-lg border border-white/10">
                <div class="col-span-4 text-sm font-medium text-white/70">Item Details</div>
                <div class="col-span-2 text-sm font-medium text-white/70">Quantity</div>
                <div class="col-span-2 text-sm font-medium text-white/70">Unit Price</div>
                <div class="col-span-1 text-sm font-medium text-white/70">VAT %</div>
                <div class="col-span-2 text-sm font-medium text-white/70">Total</div>
                <div class="col-span-1 text-sm font-medium text-white/70">Action</div>
              </div>
              
              <!-- Invoice Items -->
              {#each invoiceItems as item, index}
                <div class="glass-card p-4">
                  <div class="grid grid-cols-1 lg:grid-cols-12 gap-4">
                    <!-- Product Selection -->
                    <div class="col-span-1 lg:col-span-4">
                      <label for="product-{index}" class="label-standard lg:hidden">Product</label>
                      <select 
                        id="product-{index}"
                        bind:value={item.product_id}
                        on:change={(e) => handleProductChange(index, e)}
                        class="select-glass w-full"
                        required
                      >
                        <option value={0}>Select Product</option>
                        {#each products as product}
                          <option value={product.id}>{product.name}</option>
                        {/each}
                      </select>
                      {#if formErrors.items[index]?.product_id}
                        <p class="text-error text-sm mt-1">{formErrors.items[index].product_id}</p>
                      {/if}
                    </div>
                    
                    <!-- Quantity -->
                    <div class="col-span-1 lg:col-span-2">
                      <label for="quantity-{index}" class="label-standard lg:hidden">Quantity</label>
                      <input 
                        id="quantity-{index}"
                        type="number"
                        bind:value={item.quantity}
                        on:input={calculateTotals}
                        min="1"
                        step="0.01"
                        class="input-glass w-full"
                        required
                      />
                      {#if formErrors.items[index]?.quantity}
                        <p class="text-error text-sm mt-1">{formErrors.items[index].quantity}</p>
                      {/if}
                    </div>
                    
                    <!-- Unit Price -->
                    <div class="col-span-1 lg:col-span-2">
                      <label for="unit-price-{index}" class="label-standard lg:hidden">Unit Price</label>
                      <input 
                        id="unit-price-{index}"
                        type="number"
                        bind:value={item.unit_price}
                        on:input={calculateTotals}
                        min="0"
                        step="0.01"
                        class="input-glass w-full"
                        required
                      />
                      {#if formErrors.items[index]?.unit_price}
                        <p class="text-error text-sm mt-1">{formErrors.items[index].unit_price}</p>
                      {/if}
                    </div>
                    
                    <!-- VAT Rate -->
                    <div class="col-span-1 lg:col-span-1">
                      <label for="vat-rate-{index}" class="label-standard lg:hidden">VAT %</label>
                      <input 
                        id="vat-rate-{index}"
                        type="number"
                        bind:value={item.vat_rate}
                        on:input={calculateTotals}
                        min="0"
                        max="100"
                        step="0.01"
                        class="input-glass w-full"
                      />
                    </div>
                    
                    <!-- Total -->
                    <div class="col-span-1 lg:col-span-2">
                      <div class="text-sm font-medium text-white/70 mb-2 lg:hidden">Total</div>
                      <div class="input-glass w-full bg-white/5 text-white/70 flex items-center">
                        {calculateItemTotal(item).toFixed(2)}
                      </div>
                    </div>
                    
                    <!-- Remove Button -->
                    <div class="col-span-1 lg:col-span-1 flex items-end">
                      <button 
                        type="button"
                        class="btn-icon text-error hover:bg-error/20 w-full lg:w-auto"
                        on:click={() => removeItem(index)}
                        disabled={invoiceItems.length === 1}
                        aria-label="Remove item {index + 1}"
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
              {/each}
            </div>
            
            <!-- Totals Section -->
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <!-- Left Column: Notes and Additional Charges -->
              <div class="space-y-4">
                <!-- Notes -->
                <div class="form-group">
                  <label for="invoice-notes" class="label-standard">Notes <span class="text-white/50">(Optional)</span></label>
                  <textarea 
                    id="invoice-notes"
                    bind:value={invoiceData.notes}
                    placeholder="Add any additional notes or comments about this invoice"
                    class="textarea-glass w-full h-24 resize-y"
                  ></textarea>
                </div>
                
                <!-- Additional Charges -->
                <div class="glass-card p-4 space-y-3">
                  <h4 class="text-sm font-medium text-white/70 mb-3">Additional Charges</h4>
                  
                  <!-- Discount -->
                  <div class="form-group">
                    <label for="invoice-discount" class="label-standard text-sm">Discount</label>
                    <input 
                      id="invoice-discount"
                      type="number"
                      bind:value={discount}
                      on:input={calculateTotals}
                      min="0"
                      step="0.01"
                      class="input-glass w-full"
                      placeholder="0.00"
                    />
                  </div>
                  
                  <!-- Delivery Fees -->
                  <div class="form-group">
                    <label for="delivery-fees" class="label-standard text-sm">Delivery Fees</label>
                    <input 
                      id="delivery-fees"
                      type="number"
                      bind:value={deliveryFees}
                      on:input={calculateTotals}
                      min="0"
                      step="0.01"
                      class="input-glass w-full"
                      placeholder="0.00"
                    />
                  </div>
                </div>
              </div>
              
              <!-- Right Column: Invoice Summary -->
              <div class="glass-card p-4 space-y-3">
                <h4 class="text-sm font-medium text-white/70 mb-3">Invoice Summary</h4>
                <div class="flex justify-between text-white/70">
                  <span>Subtotal</span>
                  <span>{subtotal.toFixed(2)}</span>
                </div>
                <div class="flex justify-between text-white/70">
                  <span>Discount</span>
                  <span>-{discount.toFixed(2)}</span>
                </div>
                <div class="flex justify-between text-white/70">
                  <span>VAT</span>
                  <span>{vatTotal.toFixed(2)}</span>
                </div>
                <div class="flex justify-between text-white/70">
                  <span>Delivery Fees</span>
                  <span>{deliveryFees.toFixed(2)}</span>
                </div>
                <div class="border-t border-white/20 pt-3">
                  <div class="flex justify-between text-lg font-semibold text-white">
                    <span>Total Price</span>
                    <span>{totalPrice.toFixed(2)}</span>
                  </div>
                </div>
              </div>
            </div>
          </form>
        {/if}
      </div>
      
      <div class="modal-footer border-t border-white/20 p-6">
        <div class="flex justify-end gap-3">
          <button 
            type="button"
            class="btn-secondary"
            on:click={closeModal}
            disabled={isSaving}
          >
            Cancel
          </button>
          <button 
            type="button"
            class="btn-primary"
            on:click={saveInvoice}
            disabled={isSaving || isLoading}
          >
            {#if isSaving}
              <span class="loading loading-spinner loading-sm"></span>
              Saving...
            {:else}
              Save Invoice
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    animation: fadeIn 0.2s ease-out;
  }
  
  .modal-content {
    animation: slideIn 0.3s ease-out;
  }
  
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  
  @keyframes slideIn {
    from { 
      opacity: 0;
      transform: translateY(-20px) scale(0.95);
    }
    to { 
      opacity: 1;
      transform: translateY(0) scale(1);
    }
  }
</style>