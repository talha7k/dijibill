<script>
  import { createEventDispatcher } from 'svelte'
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

{#if show}
  <div class="modal modal-open">
    <div class="modal-box bg-base-100/90 backdrop-blur-lg border border-white/20">
      <div class="flex items-center justify-between mb-6">
        <h3 class="font-bold text-lg text-white">
          {editingCategory ? 'Edit Purchase Product Category' : 'New Purchase Product Category'}
        </h3>
        <button class="btn btn-sm btn-circle btn-ghost text-white" on:click={closeModal}>
          <i class="fas fa-times w-4 h-4"></i>
        </button>
      </div>
      
      <div class="space-y-4">
        <!-- Category Name -->
        <FormField
          label="Category Name"
          type="text"
          bind:value={categoryForm.name}
          placeholder="Enter purchase product category name"
          required={true}
        />

        <!-- Category Name Arabic -->
        <FormField
          label="Category Name (Arabic)"
          labelArabic="اسم فئة منتجات الشراء"
          type="text"
          bind:value={categoryForm.name_arabic}
          placeholder="أدخل اسم فئة منتجات الشراء"
          dir="rtl"
        />

        <!-- Description -->
        <FormField
          label="Description"
          type="textarea"
          bind:value={categoryForm.description}
          placeholder="Enter purchase product category description"
        />

        <!-- Description Arabic -->
        <FormField
          label="Description (Arabic)"
          labelArabic="الوصف"
          type="textarea"
          bind:value={categoryForm.description_arabic}
          placeholder="أدخل وصف فئة منتجات الشراء"
          dir="rtl"
        />
      </div>

      <div class="modal-action">
        <button class="btn btn-ghost" on:click={closeModal}>Cancel</button>
        <button class="btn btn-primary" on:click={saveCategory} disabled={isLoading}>
          {#if isLoading}
            <span class="loading loading-spinner loading-sm"></span>
          {/if}
          {editingCategory ? 'Update' : 'Create'} Category
        </button>
      </div>
    </div>
  </div>
{/if}