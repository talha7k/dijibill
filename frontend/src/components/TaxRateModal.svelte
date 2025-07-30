<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'

  const dispatch = createEventDispatcher()

  export let show = false
  export let taxRate: TaxRate | null = null
  export let loading = false

  interface TaxRate {
    id?: number
    name: string
    name_arabic?: string
    rate: number
    description?: string
    is_default: boolean
    is_active: boolean
    created_at?: string | null
    updated_at?: string | null
  }

  let formData = {
    name: '',
    name_arabic: '',
    rate: '0',
    description: '',
    is_default: false,
    is_active: true
  }

  $: isEditing = taxRate && taxRate.id
  $: modalTitle = isEditing ? 'Edit Tax Rate' : 'Add New Tax Rate'
  $: primaryButtonText = isEditing ? 'Update' : 'Add'

  // Reset form when modal opens/closes or taxRate changes
  $: if (show && taxRate) {
    formData = {
      name: taxRate.name || '',
      name_arabic: taxRate.name_arabic || '',
      rate: (taxRate.rate || 0).toString(),
      description: taxRate.description || '',
      is_default: taxRate.is_default || false,
      is_active: taxRate.is_active !== undefined ? taxRate.is_active : true
    }
  } else if (show && !taxRate) {
    formData = {
      name: '',
      name_arabic: '',
      rate: '0',
      description: '',
      is_default: false,
      is_active: true
    }
  }

  function handleClose() {
    dispatch('close')
  }

  function handleSave() {
    if (!formData.name || parseFloat(formData.rate) < 0) {
      return
    }

    const taxRateData = {
      ...formData,
      rate: parseFloat(formData.rate),
      id: isEditing ? taxRate.id : 0,
      created_at: isEditing ? taxRate.created_at : null,
      updated_at: isEditing ? taxRate.updated_at : null
    }

    dispatch('save', taxRateData)
  }

  function handleCancel() {
    dispatch('close')
  }

  $: isFormValid = formData.name.trim() && parseFloat(formData.rate) >= 0
</script>

<Modal
  {show}
  title={modalTitle}
  size="md"
  {loading}
  primaryButtonText={primaryButtonText}
  secondaryButtonText="Cancel"
  primaryButtonDisabled={!isFormValid}
  on:close={handleClose}
  on:primary={handleSave}
  on:secondary={handleCancel}
>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <FormField
      label="Name"
      labelArabic="الاسم"
      type="text"
      bind:value={formData.name}
      placeholder="Enter tax rate name"
      required
    />

    <FormField
      label="Name (Arabic)"
      labelArabic="الاسم بالعربية"
      type="text"
      bind:value={formData.name_arabic}
      placeholder="اسم معدل الضريبة"
      dir="rtl"
    />

    <FormField
      label="Rate (%)"
      labelArabic="المعدل (%)"
      type="number"
      bind:value={formData.rate}
      placeholder="Enter tax rate percentage"
      min="0"
      max="100"
      step="0.01"
      required
    />

    <FormField
      label="Description"
      labelArabic="الوصف"
      type="text"
      bind:value={formData.description}
      placeholder="Enter description"
    />

    <div class="flex items-center space-x-6 col-span-2">
      <FormField
        type="checkbox"
        bind:checked={formData.is_default}
        placeholder="Set as default"
      />
      <FormField
        type="checkbox"
        bind:checked={formData.is_active}
        placeholder="Active"
      />
    </div>
  </div>
</Modal>