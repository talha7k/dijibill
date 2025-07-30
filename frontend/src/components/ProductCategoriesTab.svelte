<script>
  import { createEventDispatcher } from 'svelte'
  import { CreateProductCategory } from '../../wailsjs/go/main/App.js'
  import { main } from '../../wailsjs/go/models'

  export let categories = []
  export let isLoading = false

  const dispatch = createEventDispatcher()

  let showCategoryModal = false
  let searchTerm = ''

  let categoryForm = {
    name: '',
    name_arabic: '',
    description: '',
    description_arabic: ''
  }

  // Filtered categories based on search
  $: filteredCategories = categories.filter(category => 
    category.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (category.name_arabic && category.name_arabic.includes(searchTerm))
  )

  function handleAddCategory() {
    resetCategoryForm()
    showCategoryModal = true
  }

  function handleCloseCategoryModal() {
    showCategoryModal = false
    resetCategoryForm()
  }

  async function handleSaveCategory() {
    try {
      isLoading = true
      const category = main.ProductCategory.createFrom(categoryForm)
      await CreateProductCategory(category)
      
      dispatch('categoriesUpdated')
      showCategoryModal = false
      resetCategoryForm()
    } catch (error) {
      console.error('Error saving category:', error)
      alert('Error saving category: ' + error.message)
    } finally {
      isLoading = false
    }
  }

  function resetCategoryForm() {
    categoryForm = {
      name: '',
      name_arabic: '',
      description: '',
      description_arabic: ''
    }
  }
</script>

<div class="space-y-6">
  <!-- Search and Add Section -->
  <div class="glass-card">
    <div class="p-6">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h2 class="heading-2 mb-2">Product Categories</h2>
          <p class="text-secondary">Organize your products into categories</p>
        </div>
        <button 
          class="btn-primary"
          on:click={handleAddCategory}
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
          Add Category
        </button>
      </div>

      <!-- Search -->
      <div class="form-control max-w-md">
        <label class="label-glass">Search Categories</label>
        <input 
          type="text" 
          class="input-glass"
          placeholder="Search by name..."
          bind:value={searchTerm}
        />
      </div>
    </div>
  </div>

  <!-- Categories List -->
  <div class="glass-card">
    <div class="p-6">
      {#if isLoading}
        <div class="loading-overlay">
          <div class="loading-spinner-lg"></div>
        </div>
      {:else if filteredCategories.length === 0}
        <div class="text-center py-12">
          <svg class="w-16 h-16 mx-auto text-white/40 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
          </svg>
          <h3 class="heading-3 mb-2">No categories found</h3>
          <p class="text-secondary mb-4">
            {searchTerm ? 'No categories match your search.' : 'Get started by creating your first category.'}
          </p>
          {#if !searchTerm}
            <button class="btn-primary" on:click={handleAddCategory}>
              Create First Category
            </button>
          {/if}
        </div>
      {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {#each filteredCategories as category}
            <div class="glass-card-subtle">
              <div class="p-4">
                <div class="flex items-start justify-between mb-3">
                  <div class="flex-1">
                    <h3 class="heading-4 mb-1">{category.name}</h3>
                    {#if category.name_arabic}
                      <p class="text-secondary text-sm" dir="rtl">{category.name_arabic}</p>
                    {/if}
                  </div>
                </div>
                
                {#if category.description}
                  <p class="text-secondary text-sm mb-2">{category.description}</p>
                {/if}
                
                {#if category.description_arabic}
                  <p class="text-secondary text-sm" dir="rtl">{category.description_arabic}</p>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>
</div>

<!-- Category Modal -->
{#if showCategoryModal}
  <div class="modal-glass">
    <div class="modal-backdrop" on:click={handleCloseCategoryModal}></div>
    <div class="modal-content max-w-2xl">
      <div class="modal-header">
        <h3 class="heading-3">New Category</h3>
        <button class="btn-icon" on:click={handleCloseCategoryModal}>
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>
      
      <div class="modal-body">
        <form on:submit|preventDefault={handleSaveCategory} class="space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Category Name -->
            <div class="form-control">
              <label class="label-glass">Category Name</label>
              <input 
                type="text" 
                class="input-glass"
                bind:value={categoryForm.name}
                placeholder="Enter category name"
                required
              />
            </div>

            <!-- Category Name Arabic -->
            <div class="form-control">
              <label class="label-glass">Category Name (Arabic)</label>
              <input 
                type="text" 
                class="input-glass"
                bind:value={categoryForm.name_arabic}
                placeholder="أدخل اسم الفئة"
                dir="rtl"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Description -->
            <div class="form-control">
              <label class="label-glass">Description</label>
              <textarea 
                class="textarea-glass"
                bind:value={categoryForm.description}
                placeholder="Enter category description"
                rows="3"
              ></textarea>
            </div>

            <!-- Description Arabic -->
            <div class="form-control">
              <label class="label-glass">Description (Arabic)</label>
              <textarea 
                class="textarea-glass"
                bind:value={categoryForm.description_arabic}
                placeholder="أدخل وصف الفئة"
                dir="rtl"
                rows="3"
              ></textarea>
            </div>
          </div>
        </form>
      </div>
      
      <div class="modal-footer">
        <button type="button" class="btn-glass" on:click={handleCloseCategoryModal}>
          Cancel
        </button>
        <button 
          type="button" 
          class="btn-primary" 
          on:click={handleSaveCategory}
          disabled={isLoading || !categoryForm.name.trim()}
        >
          {#if isLoading}
            <div class="loading-spinner-sm mr-2"></div>
          {/if}
          Create Category
        </button>
      </div>
    </div>
  </div>
{/if}