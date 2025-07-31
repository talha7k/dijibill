<script>
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'

  /**
   * Defines the structure of a Payment object for type checking.
   * @typedef {object} Payment
   * @property {number} id - The unique ID of the payment.
   * @property {number} invoice_id - The ID of the associated invoice.
   * @property {number} payment_type_id - The ID of the payment type.
   * @property {number | string} amount - The payment amount.
   * @property {string} payment_date - The date of the payment (YYYY-MM-DD).
   * @property {string | null} [reference] - Optional payment reference.
   * @property {string | null} [notes] - Optional notes in English.
   * @property {string | null} [notes_arabic] - Optional notes in Arabic.
   * @property {string} status - The payment status (e.g., 'completed').
   */

  // Props
  export let show = false
  /** @type {Payment | null} */
  export let editingPayment = null
  export let loading = false
  export let invoices = []
  export let paymentTypes = []

  const dispatch = createEventDispatcher()

  // Form state
  let paymentForm = {
    invoice_id: '',
    payment_type_id: '',
    amount: '',
    payment_date: new Date().toISOString().split('T')[0],
    reference: '',
    notes: '',
    notes_arabic: '',
    status: 'completed'
  }
  
  let errors = {}

  // Reset form when modal opens/closes or when editing payment changes
  $: if (show && editingPayment) {
    paymentForm = {
      invoice_id: editingPayment.invoice_id?.toString() ?? '',
      payment_type_id: editingPayment.payment_type_id?.toString() ?? '',
      amount: editingPayment.amount?.toString() ?? '',
      payment_date: editingPayment.payment_date ? new Date(editingPayment.payment_date).toISOString().split('T')[0] : new Date().toISOString().split('T')[0],
      reference: editingPayment.reference ?? '',
      notes: editingPayment.notes ?? '',
      notes_arabic: editingPayment.notes_arabic ?? '',
      status: editingPayment.status ?? 'completed'
    }
    errors = {}
  } else if (show && !editingPayment) {
    paymentForm = {
      invoice_id: '',
      payment_type_id: '',
      amount: '',
      payment_date: new Date().toISOString().split('T')[0],
      reference: '',
      notes: '',
      notes_arabic: '',
      status: 'completed'
    }
    errors = {}
  }

  // Create options for select fields
  $: invoiceOptions = invoices.map(invoice => ({
    value: invoice.id.toString(),
    label: invoice.invoice_number || `Invoice #${invoice.id}`
  }))

  $: paymentTypeOptions = paymentTypes.map(type => ({
    value: type.id.toString(),
    label: type.name
  }))

  $: statusOptions = [
    { value: 'completed', label: 'Completed' },
    { value: 'pending', label: 'Pending' },
    { value: 'refunded', label: 'Refunded' },
    { value: 'failed', label: 'Failed' },
    { value: 'cancelled', label: 'Cancelled' }
  ]

  function validateForm() {
    errors = {}
    
    if (!paymentForm.invoice_id || paymentForm.invoice_id === '') {
      errors.invoice_id = 'Please select an invoice'
    }
    
    if (!paymentForm.payment_type_id || paymentForm.payment_type_id === '') {
      errors.payment_type_id = 'Please select a payment method'
    }
    
    if (!paymentForm.amount || parseFloat(paymentForm.amount) <= 0) {
      errors.amount = 'Amount must be greater than 0'
    }
    
    if (!paymentForm.payment_date) {
      errors.payment_date = 'Payment date is required'
    }

    return Object.keys(errors).length === 0
  }

  function handleSave() {
    if (validateForm()) {
      // Convert string values back to numbers for the API
      const paymentData = {
        ...paymentForm,
        invoice_id: parseInt(paymentForm.invoice_id),
        payment_type_id: parseInt(paymentForm.payment_type_id),
        amount: parseFloat(paymentForm.amount)
      }
      
      dispatch('save', {
        payment: paymentData,
        isEditing: !!editingPayment
      })
    }
  }

  function handleClose() {
    dispatch('close')
  }

  function handleCancel() {
    dispatch('close')
  }

  $: modalTitle = editingPayment ? 'Edit Payment' : 'Record New Payment'
  $: primaryButtonText = editingPayment ? 'Update Payment' : 'Record Payment'
  $: isFormValid = paymentForm.invoice_id !== '' && paymentForm.payment_type_id !== '' && parseFloat(paymentForm.amount || '0') > 0
</script>

<Modal
  {show}
  title={modalTitle}
  size="lg"
  {loading}
  primaryButtonText={primaryButtonText}
  primaryButtonDisabled={!isFormValid}
  on:primary={handleSave}
  on:secondary={handleCancel}
  on:close={handleClose}
>
  <div class="space-y-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <FormField
        label="Invoice"
        labelArabic="الفاتورة"
        type="select"
        bind:value={paymentForm.invoice_id}
        options={invoiceOptions}
        placeholder="Select an invoice"
        required={true}
        error={errors.invoice_id}
      />

      <FormField
        label="Payment Method"
        labelArabic="طريقة الدفع"
        type="select"
        bind:value={paymentForm.payment_type_id}
        options={paymentTypeOptions}
        placeholder="Select payment method"
        required={true}
        error={errors.payment_type_id}
      />

      <FormField
        label="Amount"
        labelArabic="المبلغ"
        type="number"
        bind:value={paymentForm.amount}
        placeholder="0.00"
        min="0"
        step="0.01"
        required={true}
        error={errors.amount}
      />

      <FormField
        label="Payment Date"
        labelArabic="تاريخ الدفع"
        type="date"
        bind:value={paymentForm.payment_date}
        required={true}
        error={errors.payment_date}
      />

      <FormField
        label="Reference"
        labelArabic="المرجع"
        type="text"
        bind:value={paymentForm.reference}
        placeholder="Payment reference (optional)"
        error={errors.reference}
      />

      <FormField
        label="Status"
        labelArabic="الحالة"
        type="select"
        bind:value={paymentForm.status}
        options={statusOptions}
        required={true}
        error={errors.status}
      />
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <FormField
        label="Notes (English)"
        labelArabic="ملاحظات (إنجليزي)"
        type="textarea"
        bind:value={paymentForm.notes}
        placeholder="Additional notes about this payment"
        rows={3}
        error={errors.notes}
      />

      <FormField
        label="Notes (Arabic)"
        placeholder="ملاحظات إضافية حول هذه الدفعة"
        labelArabic="ملاحظات (عربي)"
        type="textarea"
        bind:value={paymentForm.notes_arabic}
        dir="rtl"
        rows={3}
        error={errors.notes_arabic}
      />
    </div>
  </div>
</Modal>