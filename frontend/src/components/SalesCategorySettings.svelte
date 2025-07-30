<script>
  import { GetSalesCategories, CreateSalesCategory, UpdateSalesCategory, DeleteSalesCategory } from '../../wailsjs/go/main/App.js'
  import { database } from '../../wailsjs/go/models'
  import { createEventDispatcher } from 'svelte'
  import DataTable from './DataTable.svelte'
  import FormField from './FormField.svelte'

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

  // DataTable configuration
  let searchTerm = ''
  
  $: filteredSalesCategories = salesCategories.filter(salesCategory => 
    salesCategory.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (salesCategory.name_arabic && salesCategory.name_arabic.includes(searchTerm)) ||
    salesCategory.code.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (salesCategory.description && salesCategory.description.toLowerCase().includes(searchTerm.toLowerCase())) ||
    (salesCategory.description_arabic && salesCategory.description_arabic.includes(searchTerm))
  )

  const columns = [
    {
      key: 'name',
      label: 'Name',
      sortable: true,
      render: (salesCategory) => `
        <div>
          <div class="font-medium text-gray-100">${salesCategory.name}</div>
          ${salesCategory.name_arabic ? `<div class="text-sm text-gray-300">${salesCategory.name_arabic}</div>` : ''}
        </div>
      `
    },
    {
      key: 'code',
      label: 'Code',
      sortable: true,
      render: (salesCategory) => `<span class="font-mono text-sm">${salesCategory.code}</span>`
    },
    {
      key: 'description',
      label: 'Description',
      render: (salesCategory) => `
        <div>
          ${salesCategory.description ? `<div>${salesCategory.description}</div>` : ''}
          ${salesCategory.description_arabic ? `<div class="text-sm text-gray-300" dir="rtl">${salesCategory.description_arabic}</div>` : ''}
        </div>
      `
    },
    {
      key: 'status',
      label: 'Status',
      render: (salesCategory) => `
        <div class="flex items-center space-x-2">
          ${salesCategory.is_default ? '<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">Default</span>' : ''}
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${salesCategory.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
            ${salesCategory.is_active ? 'Active' : 'Inactive'}
          </span>
        </div>
      `
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (salesCategory) => `
        <div class="flex items-center space-x-2">
          ${!salesCategory.is_default ? `<button class="action-btn" data-action="setDefault" data-id="${salesCategory.id}">Set Default</button>` : ''}
          <button class="action-btn" data-action="edit" data-id="${salesCategory.id}">Edit</button>
          <button class="action-btn text-red-600" data-action="delete" data-id="${salesCategory.id}">Delete</button>
        </div>
      `
    }
  ]

  const primaryAction = {
    label: 'Add Sales Category',
    icon: 'fa-plus'
  }

  // Event handlers for DataTable
  function handlePrimaryAction() {
    showSalesCategoryForm = true
  }

  function handleSearch(event) {
    searchTerm = event.detail.searchTerm
  }

  function handleRowAction(event) {
    const { action, id } = event.detail
    const salesCategory = salesCategories.find(sc => sc.id === parseInt(id))
    
    switch (action) {
      case 'edit':
        editSalesCategory(salesCategory)
        break
      case 'delete':
        deleteSalesCategory(parseInt(id))
        break
      case 'setDefault':
        setDefaultSalesCategory(parseInt(id))
        break
    }
  }
</script>

<div class="space-y-6">
  <!-- Sales Category Form -->
  {#if showSalesCategoryForm}
    <div class="bg-gray-50 p-6 rounded-lg border">
      <h4 class="text-lg font-medium text-gray-100 mb-6">
        {editingSalesCategory ? 'Edit Sales Category' : 'Add New Sales Category'}
      </h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <FormField
          label="Name"
          labelArabic="الاسم"
          type="text"
          bind:value={newSalesCategory.name}
          placeholder="Enter category name"
          required
        />

        <FormField
          label="Name (Arabic)"
          labelArabic="الاسم بالعربية"
          type="text"
          bind:value={newSalesCategory.name_arabic}
          placeholder="اسم الفئة"
          dir="rtl"
        />

        <FormField
          label="Code"
          labelArabic="الرمز"
          type="text"
          bind:value={newSalesCategory.code}
          placeholder="Enter category code"
          required
        />

        <FormField
          label="Description"
          labelArabic="الوصف"
          type="text"
          bind:value={newSalesCategory.description}
          placeholder="Enter description"
        />

        <div class="col-span-2">
          <FormField
            label="Description (Arabic)"
            labelArabic="الوصف بالعربية"
            type="text"
            bind:value={newSalesCategory.description_arabic}
            placeholder="الوصف"
            dir="rtl"
          />
        </div>

        <div class="flex items-center space-x-6 col-span-2">
          <FormField
            type="checkbox"
            bind:checked={newSalesCategory.is_default}
            placeholder="Set as default"
          />
          <FormField
            type="checkbox"
            bind:checked={newSalesCategory.is_active}
            placeholder="Active"
          />
        </div>
      </div>
      
      <div class="flex justify-end space-x-3 mt-6">
        <button
          on:click={resetSalesCategoryForm}
          class="px-4 py-2 border border-gray-300 text-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
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

  <!-- Sales Categories DataTable -->
  <DataTable
    data={filteredSalesCategories}
    {columns}
    {primaryAction}
    loading={isLoading}
    searchTerm={searchTerm}
    searchPlaceholder="Search sales categories..."
    emptyStateTitle="No sales categories found"
    emptyStateMessage="Start by adding your first sales category"
    emptyStateIcon="fa-tags"
    on:primary-action={handlePrimaryAction}
    on:search={handleSearch}
    on:row-action={handleRowAction}
  />
</div>