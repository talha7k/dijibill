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
  $: if (show && editingCustomer && typeof editingCustomer === 'object') {
    console.log('🔄 Resetting form for editing customer:', editingCustomer);
    customerForm = {
      name: editingCustomer.name || '',
      name_arabic: editingCustomer.name_arabic || '',
      vat_number: editingCustomer.vat_number || '',
      email: editingCustomer.email || '',
      phone: editingCustomer.phone || '',
      address: editingCustomer.address || '',
      address_arabic: editingCustomer.address_arabic || '',
      city: editingCustomer.city || '',
      city_arabic: editingCustomer.city_arabic || '',
      country: editingCustomer.country || 'SA',
      country_arabic: editingCustomer.country_arabic || 'المملكة العربية السعودية'
    }
    errors = {}
  } else if (show && !editingCustomer) {
    console.log('🔄 Resetting form for new customer');
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
    console.log('💾 CustomerModal handleSave called');
    console.log('💾 Current customerForm:', customerForm);
    console.log('💾 editingCustomer:', editingCustomer);
    
    if (validateForm()) {
      console.log('✅ Form validation passed');
      const eventData = {
        customer: customerForm,
        isEditing: !!editingCustomer
      };
      console.log('💾 Dispatching save event with data:', eventData);
      dispatch('save', eventData);
    } else {
      console.log('❌ Form validation failed, errors:', errors);
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