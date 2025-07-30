<script>
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'
  
  export let show = false
  export let editingInvoice = null
  export let loading = false
  export let suppliers = []

  const dispatch = createEventDispatcher()

  let invoiceForm = {
    invoice_number: '',
    supplier_id: '',
    invoice_date: '',
    due_date: '',
    amount: '',
    tax_amount: '',
    total_amount: '',
    currency: 'SAR',
    status: 'pending',
    description: '',
    description_arabic: '',
    payment_terms: 'net_30',
    reference_number: '',
    notes: ''
  }

  let errors = {}

  // Reset form when modal opens/closes or when editing invoice changes
  $: if (show && editingInvoice) {
    const invoice = editingInvoice
    invoiceForm = {
      invoice_number: invoice.invoice_number || '',
      supplier_id: invoice.supplier_id || '',
      invoice_date: invoice.invoice_date || '',
      due_date: invoice.due_date || '',
      amount: invoice.amount || '',
      tax_amount: invoice.tax_amount || '',
      total_amount: invoice.total_amount || '',
      currency: invoice.currency || 'SAR',
      status: invoice.status || 'pending',
      description: invoice.description || '',
      description_arabic: invoice.description_arabic || '',
      payment_terms: invoice.payment_terms || 'net_30',
      reference_number: invoice.reference_number || '',
      notes: invoice.notes || ''
    }
    errors = {}
  } else if (show && !editingInvoice) {
    invoiceForm = {
      invoice_number: '',
      supplier_id: '',
      invoice_date: new Date().toISOString().split('T')[0],
      due_date: '',
      amount: '',
      tax_amount: '',
      total_amount: '',
      currency: 'SAR',
      status: 'pending',
      description: '',
      description_arabic: '',
      payment_terms: 'net_30',
      reference_number: '',
      notes: ''
    }
    errors = {}
  }

  // Calculate due date based on payment terms
  $: if (invoiceForm.invoice_date && invoiceForm.payment_terms) {
    const invoiceDate = new Date(invoiceForm.invoice_date)
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

  // Calculate total amount when amount or tax changes
  $: if (invoiceForm.amount || invoiceForm.tax_amount) {
    const amount = parseFloat(invoiceForm.amount) || 0
    const taxAmount = parseFloat(invoiceForm.tax_amount) || 0
    invoiceForm.total_amount = (amount + taxAmount).toFixed(2)
  }

  function validateForm() {
    errors = {}
    
    if (!invoiceForm.invoice_number.trim()) {
      errors.invoice_number = 'Invoice number is required'
    }
    
    if (!invoiceForm.supplier_id) {
      errors.supplier_id = 'Supplier is required'
    }
    
    if (!invoiceForm.invoice_date) {
      errors.invoice_date = 'Invoice date is required'
    }
    
    if (!invoiceForm.amount || parseFloat(invoiceForm.amount) <= 0) {
      errors.amount = 'Amount must be greater than 0'
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
                   invoiceForm.invoice_date && 
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
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
    <!-- Invoice Information -->
    <div class="glass-card p-6">
      <h3 class="heading-4 mb-4">Invoice Information</h3>
      
      <div class="space-y-4">
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
      
      <div class="space-y-4">
        <FormField
          label="Invoice Date"
          labelArabic="تاريخ الفاتورة"
          type="date"
          bind:value={invoiceForm.invoice_date}
          required={true}
          error={errors.invoice_date}
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

    <!-- Amount Information -->
    <div class="glass-card p-6">
      <h3 class="heading-4 mb-4">Amount Information</h3>
      
      <div class="space-y-4">
        <FormField
          label="Amount (Excluding Tax)"
          labelArabic="المبلغ (بدون ضريبة)"
          type="number"
          bind:value={invoiceForm.amount}
          placeholder="0.00"
          step="0.01"
          min="0"
          required={true}
          error={errors.amount}
        />

        <FormField
          label="Tax Amount"
          labelArabic="مبلغ الضريبة"
          type="number"
          bind:value={invoiceForm.tax_amount}
          placeholder="0.00"
          step="0.01"
          min="0"
          error={errors.tax_amount}
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

    <!-- Description & Notes -->
    <div class="glass-card p-6">
      <h3 class="heading-4 mb-4">Description & Notes</h3>
      
      <div class="space-y-4">
        <FormField
          label="Description"
          labelArabic="الوصف"
          type="textarea"
          bind:value={invoiceForm.description}
          placeholder="Enter description"
          rows={3}
          error={errors.description}
        />

        <FormField
          label="Description (Arabic)"
          labelArabic="الوصف بالعربية"
          type="textarea"
          bind:value={invoiceForm.description_arabic}
          placeholder="أدخل الوصف بالعربية"
          dir="rtl"
          rows={3}
          error={errors.description_arabic}
        />

        <FormField
          label="Notes"
          labelArabic="ملاحظات"
          type="textarea"
          bind:value={invoiceForm.notes}
          placeholder="Enter additional notes"
          rows={2}
          error={errors.notes}
        />
      </div>
    </div>
  </div>
</Modal>