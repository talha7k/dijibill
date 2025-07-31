<script>
  import { createEventDispatcher, onMount } from 'svelte'
  import { GetTaxRates, GetPurchaseProducts } from '../../wailsjs/go/main/App.js'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'
  
  /** @type {boolean} */
  export let show = false
  /** @type {any} */
  export let editingInvoice = null
  /** @type {boolean} */
  export let loading = false
  /** @type {Array<{id: number, company_name: string}>} */
  export let suppliers = []

  const dispatch = createEventDispatcher()

  /** @type {Array<any>} */
  let taxRates = []
  /** @type {Array<any>} */
  let purchaseProducts = []

  /** @type {{invoice_number: string, supplier_id: string, issue_date: string, due_date: string, amount: string, vat_amount: string, total_amount: string, currency: string, status: string, notes: string, payment_terms: string, reference_number: string, vat_rate: string, vat_inclusive: boolean}} */
  let invoiceForm = {
    invoice_number: '',
    supplier_id: '',
    issue_date: '',
    due_date: '',
    amount: '',
    vat_amount: '',
    total_amount: '',
    currency: 'SAR',
    status: 'pending',
    notes: '',
    payment_terms: 'net_30',
    reference_number: '',
    vat_rate: '15', // Default VAT rate for Saudi Arabia
    vat_inclusive: false
  }

  /** @type {Record<string, string>} */
  let errors = {}

  // Invoice items for line-by-line product selection
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
  let vatTotal = 0
  let totalPrice = 0

  // Form validation errors for items
  let formErrors = {
    invoice_number: '',
    supplier_id: '',
    issue_date: '',
    items: []
  }

  onMount(async () => {
    try {
      const [taxRatesResult, purchaseProductsResult] = await Promise.all([
        GetTaxRates(),
        GetPurchaseProducts()
      ])
      taxRates = taxRatesResult || []
      purchaseProducts = purchaseProductsResult || []
    } catch (error) {
      console.error('Error loading data:', error)
      taxRates = []
      purchaseProducts = []
    }
  })

  // Calculate totals when items change
  $: if (invoiceItems) {
    calculateTotals()
  }

  // Reset form when modal opens/closes or when editing invoice changes
  $: if (show) {
    if (editingInvoice && typeof editingInvoice === 'object') {
      const invoice = editingInvoice
      invoiceForm = {
        invoice_number: invoice?.invoice_number || '',
        supplier_id: invoice?.supplier_id?.toString() || '',
        issue_date: invoice?.issue_date ? new Date(invoice.issue_date).toISOString().split('T')[0] : '',
        due_date: invoice?.due_date ? new Date(invoice.due_date).toISOString().split('T')[0] : '',
        amount: invoice?.amount?.toString() || '',
        vat_amount: invoice?.vat_amount?.toString() || '',
        total_amount: invoice?.total_amount?.toString() || '',
        currency: invoice?.currency || 'SAR',
        status: invoice?.status || 'pending',
        notes: invoice?.notes || '',
        payment_terms: invoice?.payment_terms || 'net_30',
        reference_number: invoice?.reference_number || '',
        vat_rate: invoice?.vat_rate?.toString() || '15',
        vat_inclusive: invoice?.vat_inclusive || false
      }
      // TODO: Load invoice items from backend when editing
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
    } else {
      invoiceForm = {
        invoice_number: '',
        supplier_id: '',
        issue_date: new Date().toISOString().split('T')[0],
        due_date: '',
        amount: '',
        vat_amount: '',
        total_amount: '',
        currency: 'SAR',
        status: 'pending',
        notes: '',
        payment_terms: 'net_30',
        reference_number: '',
        vat_rate: '15',
        vat_inclusive: false
      }
      // Reset invoice items for new invoice
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
    errors = {}
    formErrors = {
      invoice_number: '',
      supplier_id: '',
      issue_date: '',
      items: []
    }
  }

  // Calculate due date based on payment terms
  $: if (invoiceForm.issue_date && invoiceForm.payment_terms) {
    const invoiceDate = new Date(invoiceForm.issue_date)
    let daysToAdd = 30
    
    switch (invoiceForm.payment_terms) {
      case 'net_15': daysToAdd = 15; break
      case 'net_30': daysToAdd = 30; break
      case 'net_45': daysToAdd = 45; break
      case 'net_60': daysToAdd = 60; break
      case 'cash_on_delivery': daysToAdd = 0; break
      case 'advance_payment': daysToAdd = -1; break
      default: daysToAdd = 30
    }
    
    if (daysToAdd >= 0) {
      const dueDate = new Date(invoiceDate)
      dueDate.setDate(dueDate.getDate() + daysToAdd)
      invoiceForm.due_date = dueDate.toISOString().split('T')[0]
    }
  }

  // Transform tax rates for FormField options
  $: taxRateOptions = (taxRates || []).map(rate => ({ 
    value: rate.rate, 
    label: `${rate.name} (${rate.rate}%)` 
  }))

  // Transform purchase products for select options
  $: purchaseProductOptions = (purchaseProducts || []).map(product => ({
    value: product.id,
    label: product.name
  }))

  // Item management functions
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
    
    totalPrice = subtotal + vatTotal
    
    // Update the form amounts to reflect calculated totals
    invoiceForm.amount = subtotal.toFixed(2)
    invoiceForm.vat_amount = vatTotal.toFixed(2)
    invoiceForm.total_amount = totalPrice.toFixed(2)
  }

  function onProductChange(index, productId) {
    if (!purchaseProducts || purchaseProducts.length === 0) return
    
    const product = purchaseProducts.find(p => p.id === parseInt(productId))
    if (product) {
      invoiceItems[index].unit_price = product.unit_price || 0
      invoiceItems[index].vat_rate = product.vat_rate || 15
      invoiceItems = [...invoiceItems] // Trigger reactivity
    }
  }

  function handleProductChange(index, event) {
    onProductChange(index, event.target.value)
  }

  function validateForm() {
    errors = {}
    formErrors = {
      invoice_number: '',
      supplier_id: '',
      issue_date: '',
      items: []
    }

    // Basic form validation
    if (!invoiceForm.invoice_number?.trim()) {
      errors.invoice_number = 'Invoice number is required'
      formErrors.invoice_number = 'Invoice number is required'
    }

    if (!invoiceForm.supplier_id) {
      errors.supplier_id = 'Supplier is required'
      formErrors.supplier_id = 'Supplier is required'
    }

    if (!invoiceForm.issue_date) {
      errors.issue_date = 'Invoice date is required'
      formErrors.issue_date = 'Invoice date is required'
    }

    // Validate items
    let hasValidItems = false
    invoiceItems.forEach((item, index) => {
      const itemErrors = {}
      
      if (!item.product_id || item.product_id === 0) {
        itemErrors.product_id = 'Product is required'
      }
      
      if (!item.quantity || item.quantity <= 0) {
        itemErrors.quantity = 'Quantity must be greater than 0'
      }
      
      if (!item.unit_price || item.unit_price <= 0) {
        itemErrors.unit_price = 'Unit price must be greater than 0'
      }
      
      formErrors.items[index] = itemErrors
      
      if (item.product_id && item.quantity > 0 && item.unit_price > 0) {
        hasValidItems = true
      }
    })

    if (!hasValidItems) {
      errors.items = 'At least one valid item is required'
    }

    return Object.keys(errors).length === 0
  }

  function handleSave() {
    if (validateForm()) {
      dispatch('save', invoiceForm)
    }
  }

  function handleClose() {
    dispatch('close')
  }

  function handleCancel() {
    dispatch('close')
  }

  $: modalTitle = editingInvoice ? 'Edit Purchase Invoice' : 'Add New Purchase Invoice'
  $: primaryButtonText = editingInvoice ? 'Update Invoice' : 'Create Invoice'
  $: isFormValid = invoiceForm.invoice_number.trim().length > 0 && 
                   invoiceForm.supplier_id && 
                   invoiceForm.issue_date && 
                   parseFloat(invoiceForm.amount) > 0

  const statusOptions = [
    { value: 'pending', label: 'Pending' },
    { value: 'approved', label: 'Approved' },
    { value: 'paid', label: 'Paid' },
    { value: 'overdue', label: 'Overdue' },
    { value: 'cancelled', label: 'Cancelled' }
  ]

  const paymentTermsOptions = [
    { value: 'net_15', label: 'Net 15 days' },
    { value: 'net_30', label: 'Net 30 days' },
    { value: 'net_45', label: 'Net 45 days' },
    { value: 'net_60', label: 'Net 60 days' },
    { value: 'cash_on_delivery', label: 'Cash on Delivery' },
    { value: 'advance_payment', label: 'Advance Payment' }
  ]

  const currencyOptions = [
    { value: 'SAR', label: 'SAR - Saudi Riyal' },
    { value: 'USD', label: 'USD - US Dollar' },
    { value: 'EUR', label: 'EUR - Euro' },
    { value: 'AED', label: 'AED - UAE Dirham' }
  ]

  $: supplierOptions = suppliers.map(supplier => ({
    value: supplier.id,
    label: supplier.company_name
  }))
</script>

<Modal
  {show}
  title={modalTitle}
  size="xl"
  {loading}
  primaryButtonText={primaryButtonText}
  primaryButtonDisabled={!isFormValid}
  on:primary={handleSave}
  on:secondary={handleCancel}
  on:close={handleClose}
>
  <div class="space-y-6">
    <!-- Invoice Information -->
    <div class="glass-card p-6">
      <h3 class="heading-4 mb-4">Invoice Information</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <FormField
          label="Invoice Number"
          labelArabic="رقم الفاتورة"
          type="text"
          bind:value={invoiceForm.invoice_number}
          placeholder="Enter invoice number"
          required={true}
          error={errors.invoice_number}
        />

        <FormField
          label="Supplier"
          labelArabic="المورد"
          type="select"
          bind:value={invoiceForm.supplier_id}
          options={supplierOptions}
          placeholder="Select supplier"
          required={true}
          error={errors.supplier_id}
        />

        <FormField
          label="Reference Number"
          labelArabic="الرقم المرجعي"
          type="text"
          bind:value={invoiceForm.reference_number}
          placeholder="Enter reference number"
          error={errors.reference_number}
        />

        <FormField
          label="Status"
          labelArabic="الحالة"
          type="select"
          bind:value={invoiceForm.status}
          options={statusOptions}
          error={errors.status}
        />
      </div>
    </div>

    <!-- Date & Payment Information -->
    <div class="glass-card p-6">
      <h3 class="heading-4 mb-4">Date & Payment</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <FormField
          label="Invoice Date"
          labelArabic="تاريخ الفاتورة"
          type="date"
          bind:value={invoiceForm.issue_date}
          required={true}
          error={errors.issue_date}
        />

        <FormField
          label="Payment Terms"
          labelArabic="شروط الدفع"
          type="select"
          bind:value={invoiceForm.payment_terms}
          options={paymentTermsOptions}
          error={errors.payment_terms}
        />

        <FormField
          label="Due Date"
          labelArabic="تاريخ الاستحقاق"
          type="date"
          bind:value={invoiceForm.due_date}
          error={errors.due_date}
        />

        <FormField
          label="Currency"
          labelArabic="العملة"
          type="select"
          bind:value={invoiceForm.currency}
          options={currencyOptions}
          error={errors.currency}
        />
      </div>
    </div>

    <!-- Purchase Items -->
    <div class="glass-card p-6">
      <div class="flex justify-between items-center mb-4">
        <h3 class="heading-4">Purchase Items</h3>
        <button
          type="button"
          class="btn btn-secondary btn-sm"
          on:click={addItem}
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
          </svg>
          Add Item
        </button>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-gray-200 dark:border-gray-700">
              <th class="text-left py-3 px-2 font-medium text-gray-700 dark:text-gray-300">Product</th>
              <th class="text-left py-3 px-2 font-medium text-gray-700 dark:text-gray-300">Quantity</th>
              <th class="text-left py-3 px-2 font-medium text-gray-700 dark:text-gray-300">Unit Price</th>
              <th class="text-left py-3 px-2 font-medium text-gray-700 dark:text-gray-300">VAT %</th>
              <th class="text-left py-3 px-2 font-medium text-gray-700 dark:text-gray-300">Total</th>
              <th class="text-left py-3 px-2 font-medium text-gray-700 dark:text-gray-300">Action</th>
            </tr>
          </thead>
          <tbody>
            {#each invoiceItems as item, index}
              <tr class="border-b border-gray-100 dark:border-gray-800">
                <td class="py-3 px-2">
                  <select
                    bind:value={item.product_id}
                    on:change={(e) => handleProductChange(index, e)}
                    class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value={0}>Select Product</option>
                    {#each purchaseProductOptions as option}
                      <option value={option.value}>{option.label}</option>
                    {/each}
                  </select>
                </td>
                <td class="py-3 px-2">
                  <input
                    type="number"
                    bind:value={item.quantity}
                    min="1"
                    step="1"
                    class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </td>
                <td class="py-3 px-2">
                  <input
                    type="number"
                    bind:value={item.unit_price}
                    min="0"
                    step="0.01"
                    class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </td>
                <td class="py-3 px-2">
                  <input
                    type="number"
                    bind:value={item.vat_rate}
                    min="0"
                    max="100"
                    step="0.01"
                    class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </td>
                <td class="py-3 px-2 text-right font-medium">
                  {calculateItemTotal(item).toFixed(2)}
                </td>
                <td class="py-3 px-2">
                  <button
                    type="button"
                    class="btn btn-danger btn-sm"
                    on:click={() => removeItem(index)}
                    disabled={invoiceItems.length === 1}
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

      <!-- Totals Summary -->
      <div class="mt-6 flex justify-end">
        <div class="w-80 space-y-2">
          <div class="flex justify-between py-2 border-b border-gray-200 dark:border-gray-700">
            <span class="text-gray-600 dark:text-gray-400">Subtotal:</span>
            <span class="font-medium">{subtotal.toFixed(2)}</span>
          </div>
          <div class="flex justify-between py-2 border-b border-gray-200 dark:border-gray-700">
            <span class="text-gray-600 dark:text-gray-400">VAT Total:</span>
            <span class="font-medium">{vatTotal.toFixed(2)}</span>
          </div>
          <div class="flex justify-between py-2 text-lg font-bold border-t-2 border-gray-300 dark:border-gray-600">
            <span>Total:</span>
            <span>{totalPrice.toFixed(2)}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Amount Information -->
    <div class="glass-card p-6">
      <h3 class="heading-4 mb-4">Amount Summary (Calculated from Items)</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <FormField
          label="Subtotal (Excluding VAT)"
          labelArabic="المبلغ الفرعي (بدون ضريبة)"
          type="number"
          bind:value={invoiceForm.amount}
          placeholder="0.00"
          step="0.01"
          min="0"
          readonly={true}
          error={errors.amount}
        />

        <FormField
          label="VAT Amount"
          labelArabic="مبلغ ضريبة القيمة المضافة"
          type="number"
          bind:value={invoiceForm.vat_amount}
          placeholder="0.00"
          step="0.01"
          min="0"
          readonly={true}
          error={errors.vat_amount}
        />

        <FormField
          label="Total Amount"
          labelArabic="المبلغ الإجمالي"
          type="number"
          bind:value={invoiceForm.total_amount}
          placeholder="0.00"
          step="0.01"
          min="0"
          readonly={true}
          error={errors.total_amount}
        />
      </div>
    </div>

    <!-- Notes -->
    <div class="glass-card p-6">
      <h3 class="heading-4 mb-4">Notes</h3>
      
      <div class="space-y-4">
        <FormField
          label="Notes"
          labelArabic="ملاحظات"
          type="textarea"
          bind:value={invoiceForm.notes}
          placeholder="Enter additional notes"
          rows={4}
          error={errors.notes}
        />
      </div>
    </div>
  </div>
</Modal>