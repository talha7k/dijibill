<script>
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'

  const dispatch = createEventDispatcher()

  export let show = false
  export let salesCategory = null
  export let loading = false

  let formData = {
    name: '',
    name_arabic: '',
    code: '',
    description: '',
    is_default: false,
    is_active: true
  }

  $: isEditing = salesCategory && salesCategory.id
  $: modalTitle = isEditing ? 'Edit Sales Category' : 'Add New Sales Category'
  $: primaryButtonText = isEditing ? 'Update' : 'Add'

  // Reset form when modal opens/closes or salesCategory changes
  $: if (show && salesCategory) {
    formData = {
      name: salesCategory.name || '',
      name_arabic: salesCategory.name_arabic || '',
      code: salesCategory.code || '',
      description: salesCategory.description || '',
      is_default: salesCategory.is_default || false,
      is_active: salesCategory.is_active !== undefined ? salesCategory.is_active : true
    }
  } else if (show && !salesCategory) {
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

    const salesCategoryData = {
      ...formData,
      id: isEditing ? salesCategory.id : 0,
      created_at: isEditing ? salesCategory.created_at : null,
      updated_at: isEditing ? salesCategory.updated_at : null
    }

    dispatch('save', salesCategoryData)
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
      placeholder="Enter sales category name"
      required
    />

    <FormField
      label="Name (Arabic)"
      labelArabic="الاسم بالعربية"
      type="text"
      bind:value={formData.name_arabic}
      placeholder="اسم فئة المبيعات"
      dir="rtl"
    />

    <FormField
      label="Code"
      labelArabic="الرمز"
      type="text"
      bind:value={formData.code}
      placeholder="Enter sales category code"
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