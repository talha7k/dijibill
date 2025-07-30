<script>
  import { createEventDispatcher } from 'svelte'

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
          {editingCategory ? 'Edit Category' : 'New Category'}
        </h3>
        <button class="btn btn-sm btn-circle btn-ghost text-white" on:click={closeModal}>
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>
      
      <div class="space-y-4">
        <!-- Category Name -->
        <div class="form-control">
          <label class="label">
            <span class="label-text text-white">Category Name</span>
          </label>
          <input 
            type="text" 
            bind:value={categoryForm.name}
            placeholder="Enter category name" 
            class="input input-bordered bg-white/10 text-white placeholder-white/50"
            required
          />
        </div>

        <!-- Category Name Arabic -->
        <div class="form-control">
          <label class="label">
            <span class="label-text text-white">Category Name (Arabic)</span>
          </label>
          <input 
            type="text" 
            bind:value={categoryForm.name_arabic}
            placeholder="أدخل اسم الفئة" 
            class="input input-bordered bg-white/10 text-white placeholder-white/50"
            dir="rtl"
          />
        </div>

        <!-- Description -->
        <div class="form-control">
          <label class="label">
            <span class="label-text text-white">Description</span>
          </label>
          <textarea 
            bind:value={categoryForm.description}
            placeholder="Enter category description" 
            class="textarea textarea-bordered bg-white/10 text-white placeholder-white/50"
          ></textarea>
        </div>

        <!-- Description Arabic -->
        <div class="form-control">
          <label class="label">
            <span class="label-text text-white">Description (Arabic)</span>
          </label>
          <textarea 
            bind:value={categoryForm.description_arabic}
            placeholder="أدخل وصف الفئة" 
            class="textarea textarea-bordered bg-white/10 text-white placeholder-white/50"
            dir="rtl"
          ></textarea>
        </div>
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