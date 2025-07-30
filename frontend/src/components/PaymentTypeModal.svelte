<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'

  const dispatch = createEventDispatcher()

  export let show = false
  export let paymentType: PaymentType | null = null
  export let loading = false

  interface PaymentType {
    id?: number
    name: string
    name_arabic?: string
    code: string
    description?: string
    is_default: boolean
    is_active: boolean
    created_at?: string | null
    updated_at?: string | null
  }

  let formData = {
    name: '',
    name_arabic: '',
    code: '',
    description: '',
    is_default: false,
    is_active: true
  }

  $: isEditing = paymentType && paymentType.id
  $: modalTitle = isEditing ? 'Edit Payment Type' : 'Add New Payment Type'
  $: primaryButtonText = isEditing ? 'Update' : 'Add'

  // Reset form when modal opens/closes or paymentType changes
  $: if (show && paymentType) {
    formData = {
      name: paymentType.name || '',
      name_arabic: paymentType.name_arabic || '',
      code: paymentType.code || '',
      description: paymentType.description || '',
      is_default: paymentType.is_default || false,
      is_active: paymentType.is_active !== undefined ? paymentType.is_active : true
    }
  } else if (show && !paymentType) {
    formData = {
      name: '',
      name_arabic: '',
      code: '',
      description: '',
      is_default: false,
      is_active: true
    }
  }

  function handleClose() {
    dispatch('close')
  }

  function handleSave() {
    if (!formData.name || !formData.code) {
      return
    }

    const paymentTypeData = {
      ...formData,
      id: isEditing ? paymentType.id : 0,
      created_at: isEditing ? paymentType.created_at : null,
      updated_at: isEditing ? paymentType.updated_at : null
    }

    dispatch('save', paymentTypeData)
  }

  function handleCancel() {
    dispatch('close')
  }

  $: isFormValid = formData.name.trim() && formData.code.trim()
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
      placeholder="Enter payment type name"
      required
    />

    <FormField
      label="Name (Arabic)"
      labelArabic="الاسم بالعربية"
      type="text"
      bind:value={formData.name_arabic}
      placeholder="اسم نوع الدفع"
      dir="rtl"
    />

    <FormField
      label="Code"
      labelArabic="الرمز"
      type="text"
      bind:value={formData.code}
      placeholder="Enter payment type code"
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