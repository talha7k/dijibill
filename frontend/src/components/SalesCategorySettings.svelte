<script>
  import { GetSalesCategories, CreateSalesCategory, UpdateSalesCategory, DeleteSalesCategory } from '../../wailsjs/go/main/App.js'
  import { database } from '../../wailsjs/go/models'
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  export let salesCategories = []
  export let isLoading = false

  let newSalesCategory = { name: '', name_arabic: '', code: '', description: '', description_arabic: '', is_default: false, is_active: true }
  let showSalesCategoryForm = false
  let editingSalesCategory = null

  async function addSalesCategory() {
    if (newSalesCategory.name && newSalesCategory.code) {
      try {
        const salesCategoryData = new database.SalesCategory({
          id: 0,
          name: newSalesCategory.name,
          name_arabic: newSalesCategory.name_arabic || '',
          code: newSalesCategory.code,
          description: newSalesCategory.description || '',
          description_arabic: newSalesCategory.description_arabic || '',
          is_default: newSalesCategory.is_default || false,
          is_active: true,
          created_at: null,
          updated_at: null
        });
        await CreateSalesCategory(salesCategoryData)
        dispatch('reload')
        resetSalesCategoryForm()
      } catch (error) {
        console.error('Error adding sales category:', error)
        dispatch('error', { message: 'Error adding sales category: ' + error.message })
      }
    }
  }

  function editSalesCategory(salesCategory) {
    editingSalesCategory = salesCategory.id
    newSalesCategory = { ...salesCategory }
    showSalesCategoryForm = true
  }

  async function updateSalesCategory() {
    try {
      const salesCategoryData = new database.SalesCategory({
        id: editingSalesCategory,
        name: newSalesCategory.name,
        name_arabic: newSalesCategory.name_arabic || '',
        code: newSalesCategory.code,
        description: newSalesCategory.description || '',
        description_arabic: newSalesCategory.description_arabic || '',
        is_default: newSalesCategory.is_default || false,
        is_active: newSalesCategory.is_active !== undefined ? newSalesCategory.is_active : true,
        created_at: newSalesCategory.created_at || null,
        updated_at: newSalesCategory.updated_at || null
      });
      await UpdateSalesCategory(salesCategoryData)
      dispatch('reload')
      resetSalesCategoryForm()
    } catch (error) {
      console.error('Error updating sales category:', error)
      dispatch('error', { message: 'Error updating sales category: ' + error.message })
    }
  }

  async function deleteSalesCategory(id) {
    if (confirm('Are you sure you want to delete this sales category?')) {
      try {
        await DeleteSalesCategory(id)
        dispatch('reload')
      } catch (error) {
        console.error('Error deleting sales category:', error)
        dispatch('error', { message: 'Error deleting sales category: ' + error.message })
      }
    }
  }

  function resetSalesCategoryForm() {
    newSalesCategory = { name: '', name_arabic: '', code: '', description: '', description_arabic: '', is_default: false, is_active: true }
    showSalesCategoryForm = false
    editingSalesCategory = null
  }

  async function setDefaultSalesCategory(id) {
    try {
      // First, update all sales categories to not be default
      for (const salesCategory of salesCategories) {
        if (salesCategory.is_default) {
          const updatedSalesCategory = new database.SalesCategory({ ...salesCategory, is_default: false });
          await UpdateSalesCategory(updatedSalesCategory);
        }
      }
      // Then set the selected one as default
      const selectedSalesCategory = salesCategories.find(sc => sc.id === id);
      if (selectedSalesCategory) {
        const updatedSalesCategory = new database.SalesCategory({ ...selectedSalesCategory, is_default: true });
        await UpdateSalesCategory(updatedSalesCategory);
        dispatch('reload')
      }
    } catch (error) {
      console.error('Error setting default sales category:', error)
      dispatch('error', { message: 'Error setting default sales category: ' + error.message })
    }
  }
</script>

<div class="space-y-6">
  <!-- Add Sales Category Button -->
  <div class="flex justify-between items-center">
    <h3 class="text-lg font-medium text-gray-900">Sales Categories</h3>
    <button
      on:click={() => showSalesCategoryForm = true}
      class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
    >
      Add Sales Category
    </button>
  </div>

  <!-- Sales Category Form -->
  {#if showSalesCategoryForm}
    <div class="bg-gray-50 p-4 rounded-lg">
      <h4 class="text-md font-medium text-gray-900 mb-4">
        {editingSalesCategory ? 'Edit Sales Category' : 'Add New Sales Category'}
      </h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Name</label>
          <input
            type="text"
            bind:value={newSalesCategory.name}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter category name"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Name (Arabic)</label>
          <input
            type="text"
            bind:value={newSalesCategory.name_arabic}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="اسم الفئة"
            dir="rtl"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Code</label>
          <input
            type="text"
            bind:value={newSalesCategory.code}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter category code"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
          <input
            type="text"
            bind:value={newSalesCategory.description}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter description"
          />
        </div>
        <div class="col-span-2">
          <label class="block text-sm font-medium text-gray-700 mb-2">Description (Arabic)</label>
          <input
            type="text"
            bind:value={newSalesCategory.description_arabic}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="الوصف"
            dir="rtl"
          />
        </div>
        <div class="flex items-center space-x-4 col-span-2">
          <label class="flex items-center">
            <input
              type="checkbox"
              bind:checked={newSalesCategory.is_default}
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span class="ml-2 text-sm text-gray-700">Default</span>
          </label>
          <label class="flex items-center">
            <input
              type="checkbox"
              bind:checked={newSalesCategory.is_active}
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span class="ml-2 text-sm text-gray-700">Active</span>
          </label>
        </div>
      </div>
      <div class="flex justify-end space-x-2 mt-4">
        <button
          on:click={resetSalesCategoryForm}
          class="px-4 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          Cancel
        </button>
        <button
          on:click={editingSalesCategory ? updateSalesCategory : addSalesCategory}
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          {editingSalesCategory ? 'Update' : 'Add'} Sales Category
        </button>
      </div>
    </div>
  {/if}

  <!-- Sales Categories Table -->
  <div class="bg-white shadow overflow-hidden sm:rounded-md">
    <ul class="divide-y divide-gray-200">
      {#each salesCategories as salesCategory (salesCategory.id)}
        <li class="px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex-1">
              <div class="flex items-center">
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900">
                    {salesCategory.name}
                    {#if salesCategory.name_arabic}
                      <span class="text-gray-500">({salesCategory.name_arabic})</span>
                    {/if}
                  </p>
                  <p class="text-sm text-gray-500">Code: {salesCategory.code}</p>
                  {#if salesCategory.description}
                    <p class="text-sm text-gray-500">{salesCategory.description}</p>
                  {/if}
                </div>
                <div class="flex items-center space-x-2">
                  {#if salesCategory.is_default}
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      Default
                    </span>
                  {/if}
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {salesCategory.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
                    {salesCategory.is_active ? 'Active' : 'Inactive'}
                  </span>
                </div>
              </div>
            </div>
            <div class="flex items-center space-x-2">
              {#if !salesCategory.is_default}
                <button
                  on:click={() => setDefaultSalesCategory(salesCategory.id)}
                  class="text-sm text-blue-600 hover:text-blue-900"
                >
                  Set Default
                </button>
              {/if}
              <button
                on:click={() => editSalesCategory(salesCategory)}
                class="text-sm text-indigo-600 hover:text-indigo-900"
              >
                Edit
              </button>
              <button
                on:click={() => deleteSalesCategory(salesCategory.id)}
                class="text-sm text-red-600 hover:text-red-900"
              >
                Delete
              </button>
            </div>
          </div>
        </li>
      {/each}
    </ul>
  </div>
</div>