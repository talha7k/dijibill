<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'

  const dispatch = createEventDispatcher()

  export let show = false
  export let unit: Unit | null = null
  export let loading = false

  interface Unit {
    id?: number
    value: string
    label: string
    arabic?: string
    is_default: boolean
    is_active: boolean
    created_at?: string | null
    updated_at?: string | null
  }

  let formData = {
    value: '',
    label: '',
    arabic: '',
    is_default: false,
    is_active: true
  }

  $: isEditing = unit && unit.id
  $: modalTitle = isEditing ? 'Edit Unit' : 'Add New Unit'
  $: primaryButtonText = isEditing ? 'Update' : 'Add'

  // Reset form when modal opens/closes or unit changes
  $: if (show && unit) {
    formData = {
      value: unit.value || '',
      label: unit.label || '',
      arabic: unit.arabic || '',
      is_default: unit.is_default || false,
      is_active: unit.is_active !== undefined ? unit.is_active : true
    }
  } else if (show && !unit) {
    formData = {
      value: '',
      label: '',
      arabic: '',
      is_default: false,
      is_active: true
    }
  }

  function handleClose() {
    dispatch('close')
  }

  function handleSave() {
    if (!formData.value || !formData.label) {
      return
    }

    const unitData = {
      ...formData,
      id: isEditing ? unit.id : 0,
      created_at: isEditing ? unit.created_at : null,
      updated_at: isEditing ? unit.updated_at : null
    }

    dispatch('save', unitData)
  }

  function handleCancel() {
    dispatch('close')
  }

  $: isFormValid = formData.value.trim() && formData.label.trim()
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
      label="Value"
      labelArabic="القيمة"
      type="text"
      bind:value={formData.value}
      placeholder="Enter unit value (e.g., kg, m, pcs)"
      required
    />

    <FormField
      label="Label"
      labelArabic="التسمية"
      type="text"
      bind:value={formData.label}
      placeholder="Enter unit label"
      required
    />

    <FormField
      label="Arabic Label"
      labelArabic="التسمية بالعربية"
      type="text"
      bind:value={formData.arabic}
      placeholder="التسمية بالعربية"
      dir="rtl"
    />

    <div class="flex items-center space-x-6">
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