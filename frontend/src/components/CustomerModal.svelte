<script>
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import CustomerForm from './CustomerForm.svelte'
  
  export let show = false
  export let editingCustomer = null
  export let loading = false

  const dispatch = createEventDispatcher()

  let customerForm = {
    name: '',
    name_arabic: '',
    vat_number: '',
    email: '',
    phone: '',
    address: '',
    address_arabic: '',
    city: '',
    city_arabic: '',
    country: 'SA',
    country_arabic: 'المملكة العربية السعودية'
  }

  let errors = {}

  // Reset form when modal opens/closes or when editing customer changes
  $: if (show && editingCustomer) {
    const customer = editingCustomer
    customerForm = {
      name: customer.name || '',
      name_arabic: customer.name_arabic || '',
      vat_number: customer.vat_number || '',
      email: customer.email || '',
      phone: customer.phone || '',
      address: customer.address || '',
      address_arabic: customer.address_arabic || '',
      city: customer.city || '',
      city_arabic: customer.city_arabic || '',
      country: customer.country || 'SA',
      country_arabic: customer.country_arabic || 'المملكة العربية السعودية'
    }
    errors = {}
  } else if (show && !editingCustomer) {
    customerForm = {
      name: '',
      name_arabic: '',
      vat_number: '',
      email: '',
      phone: '',
      address: '',
      address_arabic: '',
      city: '',
      city_arabic: '',
      country: 'SA',
      country_arabic: 'المملكة العربية السعودية'
    }
    errors = {}
  }

  function validateForm() {
    errors = {}
    
    if (!customerForm.name.trim()) {
      errors.name = 'Customer name is required'
    }
    
    if (customerForm.email && !isValidEmail(customerForm.email)) {
      errors.email = 'Please enter a valid email address'
    }
    
    if (customerForm.vat_number && customerForm.vat_number.length > 0 && customerForm.vat_number.length < 15) {
      errors.vat_number = 'VAT number should be at least 15 characters'
    }

    return Object.keys(errors).length === 0
  }

  function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(email)
  }

  function handleSave() {
    if (validateForm()) {
      dispatch('save', {
        customer: customerForm,
        isEditing: !!editingCustomer
      })
    }
  }

  function handleClose() {
    dispatch('close')
  }

  function handleCancel() {
    dispatch('close')
  }

  $: modalTitle = editingCustomer ? 'Edit Customer' : 'Add New Customer'
  $: primaryButtonText = editingCustomer ? 'Update Customer' : 'Create Customer'
  $: isFormValid = customerForm.name.trim().length > 0
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
  <CustomerForm bind:customerForm {errors} />
</Modal>