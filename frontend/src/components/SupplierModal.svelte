<script>
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'
  
  export let show = false
  export let editingSupplier = null
  export let loading = false

  const dispatch = createEventDispatcher()

  let supplierForm = {
    company_name: '',
    company_name_arabic: '',
    contact_person: '',
    contact_person_arabic: '',
    vat_number: '',
    email: '',
    phone: '',
    address: '',
    address_arabic: '',
    city: '',
    city_arabic: '',
    country: 'SA',
    country_arabic: 'المملكة العربية السعودية',
    payment_terms: 'net_30',
    active: true
  }

  let errors = {}

  // Reset form when modal opens/closes or when editing supplier changes
  $: if (show && editingSupplier) {
    const supplier = editingSupplier
    supplierForm = {
      company_name: supplier.company_name || '',
      company_name_arabic: supplier.company_name_arabic || '',
      contact_person: supplier.contact_person || '',
      contact_person_arabic: supplier.contact_person_arabic || '',
      vat_number: supplier.vat_number || '',
      email: supplier.email || '',
      phone: supplier.phone || '',
      address: supplier.address || '',
      address_arabic: supplier.address_arabic || '',
      city: supplier.city || '',
      city_arabic: supplier.city_arabic || '',
      country: supplier.country || 'SA',
      country_arabic: supplier.country_arabic || 'المملكة العربية السعودية',
      payment_terms: supplier.payment_terms || 'net_30',
      active: supplier.active !== undefined ? supplier.active : true
    }
    errors = {}
  } else if (show && !editingSupplier) {
    supplierForm = {
      company_name: '',
      company_name_arabic: '',
      contact_person: '',
      contact_person_arabic: '',
      vat_number: '',
      email: '',
      phone: '',
      address: '',
      address_arabic: '',
      city: '',
      city_arabic: '',
      country: 'SA',
      country_arabic: 'المملكة العربية السعودية',
      payment_terms: 'net_30',
      active: true
    }
    errors = {}
  }

  function validateForm() {
    errors = {}
    
    if (!supplierForm.company_name.trim()) {
      errors.company_name = 'Company name is required'
    }
    
    if (!supplierForm.contact_person.trim()) {
      errors.contact_person = 'Contact person is required'
    }
    
    if (supplierForm.email && !isValidEmail(supplierForm.email)) {
      errors.email = 'Please enter a valid email address'
    }
    
    if (supplierForm.vat_number && supplierForm.vat_number.length > 0 && supplierForm.vat_number.length < 15) {
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
      dispatch('save', supplierForm)
    }
  }

  function handleClose() {
    dispatch('close')
  }

  function handleCancel() {
    dispatch('close')
  }

  $: modalTitle = editingSupplier ? 'Edit Supplier' : 'Add New Supplier'
  $: primaryButtonText = editingSupplier ? 'Update Supplier' : 'Create Supplier'
  $: isFormValid = supplierForm.company_name.trim().length > 0 && supplierForm.contact_person.trim().length > 0

  const paymentTermsOptions = [
    { value: 'net_15', label: 'Net 15 days' },
    { value: 'net_30', label: 'Net 30 days' },
    { value: 'net_45', label: 'Net 45 days' },
    { value: 'net_60', label: 'Net 60 days' },
    { value: 'cash_on_delivery', label: 'Cash on Delivery' },
    { value: 'advance_payment', label: 'Advance Payment' }
  ]

  const countryOptions = [
    { value: 'SA', label: 'Saudi Arabia' },
    { value: 'AE', label: 'United Arab Emirates' },
    { value: 'KW', label: 'Kuwait' },
    { value: 'QA', label: 'Qatar' },
    { value: 'BH', label: 'Bahrain' },
    { value: 'OM', label: 'Oman' }
  ]
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
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
    <!-- Company Information -->
    <div class="glass-card p-6">
      <h3 class="heading-4 mb-4">Company Information</h3>
      
      <div class="space-y-4">
        <FormField
          label="Company Name"
          labelArabic="اسم الشركة"
          type="text"
          bind:value={supplierForm.company_name}
          placeholder="Enter company name"
          required={true}
          error={errors.company_name}
        />

        <FormField
          label="Company Name (Arabic)"
          labelArabic="اسم الشركة بالعربية"
          type="text"
          bind:value={supplierForm.company_name_arabic}
          placeholder="أدخل اسم الشركة بالعربية"
          dir="rtl"
          error={errors.company_name_arabic}
        />

        <FormField
          label="VAT Number"
          labelArabic="الرقم الضريبي"
          type="text"
          bind:value={supplierForm.vat_number}
          placeholder="Enter VAT number"
          error={errors.vat_number}
        />

        <FormField
          label="Payment Terms"
          labelArabic="شروط الدفع"
          type="select"
          bind:value={supplierForm.payment_terms}
          options={paymentTermsOptions}
          error={errors.payment_terms}
        />
      </div>
    </div>

    <!-- Contact Information -->
    <div class="glass-card p-6">
      <h3 class="heading-4 mb-4">Contact Information</h3>
      
      <div class="space-y-4">
        <FormField
          label="Contact Person"
          labelArabic="الشخص المسؤول"
          type="text"
          bind:value={supplierForm.contact_person}
          placeholder="Enter contact person name"
          required={true}
          error={errors.contact_person}
        />

        <FormField
          label="Contact Person (Arabic)"
          labelArabic="الشخص المسؤول بالعربية"
          type="text"
          bind:value={supplierForm.contact_person_arabic}
          placeholder="أدخل اسم الشخص المسؤول بالعربية"
          dir="rtl"
          error={errors.contact_person_arabic}
        />

        <FormField
          label="Email"
          labelArabic="البريد الإلكتروني"
          type="email"
          bind:value={supplierForm.email}
          placeholder="Enter email address"
          error={errors.email}
        />

        <FormField
          label="Phone"
          labelArabic="رقم الهاتف"
          type="tel"
          bind:value={supplierForm.phone}
          placeholder="Enter phone number"
          error={errors.phone}
        />
      </div>
    </div>

    <!-- Address Information -->
    <div class="glass-card p-6 lg:col-span-2">
      <h3 class="heading-4 mb-4">Address Information</h3>
      
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <FormField
          label="Address"
          labelArabic="العنوان"
          type="textarea"
          bind:value={supplierForm.address}
          placeholder="Enter address"
          rows={3}
          error={errors.address}
        />

        <FormField
          label="Address (Arabic)"
          labelArabic="العنوان بالعربية"
          type="textarea"
          bind:value={supplierForm.address_arabic}
          placeholder="أدخل العنوان بالعربية"
          dir="rtl"
          rows={3}
          error={errors.address_arabic}
        />

        <FormField
          label="City"
          labelArabic="المدينة"
          type="text"
          bind:value={supplierForm.city}
          placeholder="Enter city"
          error={errors.city}
        />

        <FormField
          label="City (Arabic)"
          labelArabic="المدينة بالعربية"
          type="text"
          bind:value={supplierForm.city_arabic}
          placeholder="أدخل المدينة بالعربية"
          dir="rtl"
          error={errors.city_arabic}
        />

        <FormField
          label="Country"
          labelArabic="البلد"
          type="select"
          bind:value={supplierForm.country}
          options={countryOptions}
          error={errors.country}
        />

        <div class="flex items-center pt-6">
          <FormField
            type="checkbox"
            bind:checked={supplierForm.active}
            placeholder="Active supplier"
          />
        </div>
      </div>
    </div>
  </div>
</Modal>