<script>
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'

  export let show = false
  export let editingCategory = null
  export let categoryForm = {}
  export let isLoading = false

  const dispatch = createEventDispatcher()

  function closeModal() {
    dispatch('close')
  }

  function saveCategory() {
    dispatch('save')
  }
</script>

<Modal
  {show}
  title={editingCategory ? 'Edit Category' : 'New Category'}
  size="md"
  loading={isLoading}
  primaryButtonText={editingCategory ? 'Update Category' : 'Create Category'}
  secondaryButtonText="Cancel"
  primaryButtonDisabled={isLoading}
  on:close={closeModal}
  on:primary={saveCategory}
  on:secondary={closeModal}
>
  <div class="space-y-4">
    <!-- Category Name -->
    <FormField
      label="Category Name"
      type="text"
      bind:value={categoryForm.name}
      placeholder="Enter category name"
      required={true}
    />

    <!-- Category Name Arabic -->
    <FormField
      label="Category Name (Arabic)"
      labelArabic="اسم الفئة"
      type="text"
      bind:value={categoryForm.name_arabic}
      placeholder="أدخل اسم الفئة"
      dir="rtl"
    />

    <!-- Description -->
    <FormField
      label="Description"
      type="textarea"
      bind:value={categoryForm.description}
      placeholder="Enter category description"
    />

    <!-- Description Arabic -->
    <FormField
      label="Description (Arabic)"
      labelArabic="الوصف"
      type="textarea"
      bind:value={categoryForm.description_arabic}
      placeholder="أدخل وصف الفئة"
      dir="rtl"
    />
  </div>
</Modal>