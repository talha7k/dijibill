<script>
  import { createEventDispatcher } from 'svelte'
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

  /** @type {{invoice_number: string, supplier_id: string, issue_date: string, due_date: string, amount: string, vat_amount: string, total_amount: string, currency: string, status: string, notes: string, payment_terms: string, reference_number: string}} */
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
    reference_number: ''
  }

  /** @type {Record<string, string>} */
  let errors = {}

  // Reset form when modal opens/closes or when editing invoice changes
  $: if (show) {
    if (editingInvoice && typeof editingInvoice === 'object') {
      const invoice = editingInvoice
      invoiceForm = {
        invoice_number: invoice?.invoice_number || '',
        supplier_id: invoice?.supplier_id || '',
        issue_date: invoice?.issue_date || '',
        due_date: invoice?.due_date || '',
        amount: invoice?.amount || '',
        vat_amount: invoice?.vat_amount || '',
        total_amount: invoice?.total_amount || '',
        currency: invoice?.currency || 'SAR',
        status: invoice?.status || 'pending',
        notes: invoice?.notes || '',
        payment_terms: invoice?.payment_terms || 'net_30',
        reference_number: invoice?.reference_number || ''
      }
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
        reference_number: ''
      }
    }
    errors = {}
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

  // Calculate total amount when amount or tax changes
  $: if (invoiceForm.amount || invoiceForm.vat_amount) {
    const amount = parseFloat(invoiceForm.amount) || 0
    const vatAmount = parseFloat(invoiceForm.vat_amount) || 0
    invoiceForm.total_amount = (amount + vatAmount).toFixed(2)
  }

  function validateForm() {
    errors = {}
    
    if (!invoiceForm.invoice_number.trim()) {
      errors.invoice_number = 'Invoice number is required'
    }
    
    if (!invoiceForm.supplier_id) {
      errors.supplier_id = 'Supplier is required'
    }
    
    if (!invoiceForm.issue_date) {
      errors.issue_date = 'Invoice date is required'
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
          label="VAT Amount"
          labelArabic="مبلغ ضريبة القيمة المضافة"
          type="number"
          bind:value={invoiceForm.vat_amount}
          placeholder="0.00"
          step="0.01"
          min="0"
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