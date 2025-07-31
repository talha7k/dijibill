<script>
  import { createEventDispatcher, onMount } from 'svelte'
  import { GetCustomers, GetProducts, GetPaymentTypes, GetSalesCategories, CreateSalesInvoice, UpdateSalesInvoice } from '../wailsjs/go/main/App.js'
  import {database} from '../wailsjs/go/models'
  import Modal from './components/Modal.svelte'
  import FormField from './components/FormField.svelte'
  
  const dispatch = createEventDispatcher()
  
  export let isOpen = false
  export let editingInvoice = null
  
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

  // Reactive title based on edit mode
  $: modalTitle = editingInvoice ? 'Edit Sales Invoice' : 'New Sales Invoice'

  // Watch for changes to editingInvoice and reinitialize form
  $: if (editingInvoice !== null && isOpen) {
    console.log('editingInvoice changed, reinitializing form:', editingInvoice)
    initializeFormData()
  }

  // Transform data for FormField options
  $: customerOptions = (customers || []).map(customer => ({ 
    value: customer.id.toString(), 
    label: customer.name 
  }))
  $: salesCategoryOptions = (salesCategories || []).map(category => ({ 
    value: category.id.toString(), 
    label: category.name 
  }))
  $: productOptions = (products || []).map(product => ({ 
    value: product.id.toString(), 
    label: product.name 
  }))
  $: statusOptions = [
    { value: 'draft', label: 'Draft' },
    { value: 'sent', label: 'Sent' },
    { value: 'partially_paid', label: 'Partially Paid' },
    { value: 'paid', label: 'Paid' },
    { value: 'refunded', label: 'Refunded' },
    { value: 'overdue', label: 'Overdue' },
    { value: 'cancelled', label: 'Cancelled' }
  ]

  // String versions for FormField binding
  let customerIdString = ''
  let salesCategoryIdString = ''

  // Convert string values back to numbers for the form data
  $: if (customerIdString && customerIdString !== invoiceData.customer_id?.toString()) {
    invoiceData.customer_id = parseInt(customerIdString) || 0
  }
  $: if (salesCategoryIdString && salesCategoryIdString !== invoiceData.sales_category_id?.toString()) {
    invoiceData.sales_category_id = parseInt(salesCategoryIdString) || 0
  }
  
  onMount(async () => {
    if (isOpen) {
      await loadData()
    }
  })
  
  $: if (isOpen) {
    loadData().then(() => {
      initializeFormData()
    })
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
      
      // Initialize string values from existing data
      if (invoiceData.customer_id) {
        customerIdString = invoiceData.customer_id.toString()
      }
      if (invoiceData.sales_category_id) {
        salesCategoryIdString = invoiceData.sales_category_id.toString()
      }
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

  function initializeFormData() {
    console.log('initializeFormData called with editingInvoice:', editingInvoice)
    
    if (editingInvoice) {
      // Populate form with existing invoice data
      invoiceData = {
        customer_id: editingInvoice.customer_id || 0,
        sales_category_id: editingInvoice.sales_category_id || 0,
        issue_date: editingInvoice.issue_date ? new Date(editingInvoice.issue_date).toISOString().split('T')[0] : new Date().toISOString().split('T')[0],
        due_date: editingInvoice.due_date ? new Date(editingInvoice.due_date).toISOString().split('T')[0] : '',
        payment_terms: editingInvoice.payment_terms || '',
        notes: editingInvoice.notes || '',
        status: editingInvoice.status || 'draft'
      }
      
      // Set string values for FormField binding
      customerIdString = invoiceData.customer_id.toString()
      salesCategoryIdString = invoiceData.sales_category_id.toString()
      
      // Populate invoice items if they exist
      if (editingInvoice.items && editingInvoice.items.length > 0) {
        console.log('Populating items:', editingInvoice.items)
        invoiceItems = editingInvoice.items.map(item => ({
          product_id: item.product_id || 0,
          quantity: item.quantity || 1,
          unit_price: item.unit_price || 0,
          vat_rate: item.vat_rate || 15,
          vat_amount: item.vat_amount || 0,
          total_amount: item.total_amount || 0
        }))
      } else {
        console.log('No items found, using default item')
        invoiceItems = [
          {
            product_id: 0,
            quantity: 1,
            unit_price: 0,
            vat_rate: 15,
            vat_amount: 0,
            total_amount: 0
          }
        ]
      }
      
      // Set additional charges if they exist
      discount = editingInvoice.discount || 0
      deliveryFees = editingInvoice.delivery_fees || 0
      
    } else {
      // Reset form for new invoice
      invoiceData = {
        customer_id: 0,
        sales_category_id: 0,
        issue_date: new Date().toISOString().split('T')[0],
        due_date: '',
        payment_terms: '',
        notes: '',
        status: 'draft'
      }
      customerIdString = ''
      salesCategoryIdString = ''
      invoiceItems = [
        {
          product_id: 0,
          quantity: 1,
          unit_price: 0,
          vat_rate: 15,
          vat_amount: 0,
          total_amount: 0
        }
      ]
      discount = 0
      deliveryFees = 0
    }
    
    // Force reactivity update
    invoiceItems = [...invoiceItems]
    calculateTotals()
    console.log('Form initialized with items:', invoiceItems)
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
      
      if (!item.product_id || parseInt(item.product_id.toString()) === 0) {
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
    const validItems = invoiceItems.filter(item => item.product_id && parseInt(item.product_id.toString()) > 0)
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
        id: editingInvoice ? editingInvoice.id : 0,
        invoice_number: editingInvoice ? editingInvoice.invoice_number : '',
        customer_id: invoiceData.customer_id ? parseInt(invoiceData.customer_id.toString()) : 0,
        sales_category_id: invoiceData.sales_category_id ? parseInt(invoiceData.sales_category_id.toString()) : 0,
        issue_date: invoiceData.issue_date,
        due_date: invoiceData.due_date || new Date().toISOString().split('T')[0],
        sub_total: 0,
        vat_amount: 0,
        total_amount: 0,
        status: invoiceData.status || 'draft',
        notes: invoiceData.notes || '',
        notes_arabic: '',
        qr_code: editingInvoice ? editingInvoice.qr_code || '' : '',
        items: validItems.map(item => ({
          id: 0,
          invoice_id: editingInvoice ? editingInvoice.id : 0,
          product_id: parseInt(item.product_id.toString()),
          quantity: parseFloat(item.quantity.toString()),
          unit_price: parseFloat(item.unit_price.toString()),
          vat_rate: parseFloat(item.vat_rate.toString()),
          vat_amount: 0,
          total_amount: 0
        }))
      }
      
      console.log('invoiceObj created:', JSON.stringify(invoiceObj, null, 2))
      
      let result
      if (editingInvoice) {
        console.log('Calling UpdateSalesInvoice backend function...')
        result = await UpdateSalesInvoice(new database.SalesInvoice(invoiceObj))
        console.log('UpdateSalesInvoice result:', result)
        console.log('UpdateSalesInvoice completed successfully')
        alert('Invoice updated successfully!')
      } else {
        console.log('Calling CreateSalesInvoice backend function...')
        result = await CreateSalesInvoice(new database.SalesInvoice(invoiceObj))
        console.log('CreateSalesInvoice result:', result)
        console.log('CreateSalesInvoice completed successfully')
        alert('Invoice saved successfully!')
      }
      
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
    customerIdString = ''
    salesCategoryIdString = ''
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
    
    // Clear form errors
    formErrors = {
      customer_id: '',
      sales_category_id: '',
      issue_date: '',
      items: []
    }
  }
</script>

<Modal
  show={isOpen}
  title={modalTitle}
  size="xl"
  loading={isLoading}
  primaryButtonText={isSaving ? (editingInvoice ? 'Updating...' : 'Saving...') : (editingInvoice ? 'Update Invoice' : 'Save Invoice')}
  secondaryButtonText="Cancel"
  primaryButtonDisabled={isSaving || isLoading}
  on:close={closeModal}
  on:primary={saveInvoice}
  on:secondary={closeModal}
>
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
    <!-- Invoice Header Information -->
    <div class="glass-card p-6">
      <h3 class="text-lg font-semibold text-white mb-4">Invoice Information</h3>
      
      <!-- Customer and Invoice Details -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Left Column -->
        <div class="space-y-6">
          <!-- Customer Selection -->
          <FormField
            label="Customer"
            labelArabic="العميل"
            type="select"
            bind:value={customerIdString}
            options={customerOptions}
            placeholder="Select Customer (Optional)"
            error={formErrors.customer_id}
          />
          
          <!-- Sales Category Selection -->
          <FormField
            label="Sales Category"
            labelArabic="فئة المبيعات"
            type="select"
            bind:value={salesCategoryIdString}
            options={salesCategoryOptions}
            placeholder="Select Sales Category (Optional)"
            error={formErrors.sales_category_id}
          />

          <!-- Invoice Date -->
          <FormField
            label="Invoice Date"
            labelArabic="تاريخ الفاتورة"
            type="date"
            bind:value={invoiceData.issue_date}
            required={true}
            error={formErrors.issue_date}
          />
        </div>

        <!-- Right Column -->
        <div class="space-y-6">
          <!-- Due Date -->
          <FormField
            label="Due Date"
            labelArabic="تاريخ الاستحقاق"
            type="date"
            bind:value={invoiceData.due_date}
            placeholder="Optional"
          />
          
          <!-- Status -->
          <FormField
            label="Status"
            labelArabic="الحالة"
            type="select"
            bind:value={invoiceData.status}
            options={statusOptions}
            placeholder="Select Status"
          />

          <!-- Payment Terms -->
          <FormField
            label="Payment Terms"
            labelArabic="شروط الدفع"
            type="text"
            bind:value={invoiceData.payment_terms}
            placeholder="e.g., Net 30, Due on receipt (Optional)"
          />
        </div>
      </div>
    </div>
    
    <!-- Item Details Section -->
    <div class="glass-card p-6">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-semibold text-white">Item Details</h3>
        <button 
          type="button"
          class="btn-primary btn-sm"
          on:click={addItem}
        >
          + Add Item
        </button>
      </div>
      
      <!-- Items Table -->
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-white/20">
              <th class="text-left py-3 px-2 text-sm font-medium text-white/70 min-w-[200px]">Product</th>
              <th class="text-left py-3 px-2 text-sm font-medium text-white/70 w-24">Quantity</th>
              <th class="text-left py-3 px-2 text-sm font-medium text-white/70 w-28">Unit Price</th>
              <th class="text-left py-3 px-2 text-sm font-medium text-white/70 w-20">VAT %</th>
              <th class="text-left py-3 px-2 text-sm font-medium text-white/70 w-28">Total</th>
              <th class="text-left py-3 px-2 text-sm font-medium text-white/70 w-16">Action</th>
            </tr>
          </thead>
          <tbody>
            {#each invoiceItems as item, index}
              <tr class="border-b border-white/10">
                <!-- Product Selection -->
                <td class="py-3 px-2">
                  <select 
                    id="product-{index}"
                    bind:value={item.product_id}
                    on:change={(e) => handleProductChange(index, e)}
                    class="select-glass w-full text-sm"
                    required
                  >
                    <option value={0}>Select Product</option>
                    {#each products as product}
                      <option value={product.id}>{product.name}</option>
                    {/each}
                  </select>
                  {#if formErrors.items[index]?.product_id}
                    <p class="text-error text-xs mt-1">{formErrors.items[index].product_id}</p>
                  {/if}
                </td>
                
                <!-- Quantity -->
                <td class="py-3 px-2">
                  <input 
                    id="quantity-{index}"
                    type="number"
                    bind:value={item.quantity}
                    on:input={calculateTotals}
                    min="1"
                    step="0.01"
                    class="input-glass w-full text-sm"
                    required
                  />
                  {#if formErrors.items[index]?.quantity}
                    <p class="text-error text-xs mt-1">{formErrors.items[index].quantity}</p>
                  {/if}
                </td>
                
                <!-- Unit Price -->
                <td class="py-3 px-2">
                  <input 
                    id="unit-price-{index}"
                    type="number"
                    bind:value={item.unit_price}
                    on:input={calculateTotals}
                    min="0"
                    step="0.01"
                    class="input-glass w-full text-sm"
                    required
                  />
                  {#if formErrors.items[index]?.unit_price}
                    <p class="text-error text-xs mt-1">{formErrors.items[index].unit_price}</p>
                  {/if}
                </td>
                
                <!-- VAT Rate -->
                <td class="py-3 px-2">
                  <input 
                    id="vat-rate-{index}"
                    type="number"
                    bind:value={item.vat_rate}
                    on:input={calculateTotals}
                    min="0"
                    max="100"
                    step="0.01"
                    class="input-glass w-full text-sm"
                  />
                </td>
                
                <!-- Total -->
                <td class="py-3 px-2">
                  <div class="input-glass w-full bg-white/5 text-white/70 text-sm px-3 py-2">
                    {calculateItemTotal(item).toFixed(2)}
                  </div>
                </td>
                
                <!-- Remove Button -->
                <td class="py-3 px-2">
                  <button 
                    type="button"
                    class="btn-icon text-error hover:bg-error/20 p-2"
                    on:click={() => removeItem(index)}
                    disabled={invoiceItems.length === 1}
                    aria-label="Remove item {index + 1}"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                    </svg>
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
    
    <!-- Totals Section -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Left Column: Notes and Additional Charges -->
      <div class="glass-card p-6">
        <h3 class="text-lg font-semibold text-white mb-4">Additional Information</h3>
        
        <!-- Notes -->
        <div class="form-group mb-6">
          <label for="invoice-notes" class="label-standard">Notes <span class="text-white/50">(Optional)</span></label>
          <textarea 
            id="invoice-notes"
            bind:value={invoiceData.notes}
            placeholder="Add any additional notes or comments about this invoice"
            class="textarea-glass w-full h-24 resize-y"
          ></textarea>
        </div>
        
        <!-- Additional Charges -->
        <div class="space-y-4">
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
      <div class="glass-card p-6">
        <h3 class="text-lg font-semibold text-white mb-4">Invoice Summary</h3>
        
        <div class="space-y-3">
          <div class="flex justify-between text-white/70">
            <span>Subtotal</span>
            <span class="font-medium">{subtotal.toFixed(2)}</span>
          </div>
          <div class="flex justify-between text-white/70">
            <span>Discount</span>
            <span class="font-medium text-green-400">-{discount.toFixed(2)}</span>
          </div>
          <div class="flex justify-between text-white/70">
            <span>VAT</span>
            <span class="font-medium">{vatTotal.toFixed(2)}</span>
          </div>
          <div class="flex justify-between text-white/70">
            <span>Delivery Fees</span>
            <span class="font-medium">{deliveryFees.toFixed(2)}</span>
          </div>
          <div class="border-t border-white/20 pt-3 mt-4">
            <div class="flex justify-between text-xl font-bold text-white">
              <span>Total Amount</span>
              <span class="text-primary">{totalPrice.toFixed(2)}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </form>
</Modal>