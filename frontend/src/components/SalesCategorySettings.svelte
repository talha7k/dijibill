<script>
  import { GetSalesCategories, CreateSalesCategory, UpdateSalesCategory, DeleteSalesCategory } from '../../wailsjs/go/main/App.js'
  import { database } from '../../wailsjs/go/models'
  import { createEventDispatcher } from 'svelte'
  import DataTable from './DataTable.svelte'
  import SalesCategoryModal from './SalesCategoryModal.svelte'
  import ConfirmationModal from './ConfirmationModal.svelte'
  import { showDbSuccess, showDbError } from '../stores/notificationStore.js'

  const dispatch = createEventDispatcher()

  export let salesCategories = []
  export let isLoading = false

  let showSalesCategoryModal = false
  let editingSalesCategory = null
  let showDeleteConfirm = false
  let salesCategoryToDelete = null
  let loading = false

  function editSalesCategory(salesCategory) {
    editingSalesCategory = salesCategory
    showSalesCategoryModal = true
  }

  function closeSalesCategoryModal() {
    showSalesCategoryModal = false
    editingSalesCategory = null
  }

  async function handleSalesCategorySave(event) {
    const salesCategoryData = event.detail
    loading = true
    
    try {
      if (editingSalesCategory) {
        // Update existing sales category
        const updatedSalesCategory = new database.SalesCategory({
          ...editingSalesCategory,
          ...salesCategoryData
        })
        await UpdateSalesCategory(updatedSalesCategory)
        showDbSuccess('update', 'Sales Category')
      } else {
        // Create new sales category
        const newSalesCategory = new database.SalesCategory({
          id: 0,
          ...salesCategoryData,
          created_at: null,
          updated_at: null
        })
        await CreateSalesCategory(newSalesCategory)
        showDbSuccess('create', 'Sales Category')
      }
      
      dispatch('reload')
      closeSalesCategoryModal()
    } catch (error) {
      console.error('Error saving sales category:', error)
      const operation = editingSalesCategory ? 'update' : 'create'
      showDbError(operation, 'Sales Category', error)
    } finally {
      loading = false
    }
  }

  function showDeleteConfirmation(id) {
    salesCategoryToDelete = id
    showDeleteConfirm = true
  }

  function cancelDelete() {
    showDeleteConfirm = false
    salesCategoryToDelete = null
  }

  async function confirmDelete() {
    if (salesCategoryToDelete) {
      loading = true
      try {
        await DeleteSalesCategory(salesCategoryToDelete)
        showDbSuccess('delete', 'Sales Category')
        dispatch('reload')
        showDeleteConfirm = false
        salesCategoryToDelete = null
      } catch (error) {
        console.error('Error deleting sales category:', error)
        showDbError('delete', 'Sales Category', error)
      } finally {
        loading = false
      }
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
      label: 'Actions',
      actions: [
        {
          key: 'edit',
          text: 'Edit',
          icon: 'fa-edit',
          class: 'btn-secondary',
          title: 'Edit sales category'
        },
        {
          key: 'delete',
          text: 'Delete',
          icon: 'fa-trash',
          class: 'btn-error',
          title: 'Delete sales category'
        }
      ]
    }
  ]

  const primaryAction = {
    label: 'Add Sales Category',
    icon: 'fa-plus'
  }

  // Event handlers for DataTable
  function handlePrimaryAction() {
    showSalesCategoryModal = true
  }

  function handleSearch(event) {
    searchTerm = event.detail.searchTerm
  }

  function handleRowAction(event) {
    const { action, item } = event.detail
    
    switch (action) {
      case 'edit':
        editSalesCategory(item)
        break
      case 'delete':
        showDeleteConfirmation(item.id)
        break
    }
  }
</script>

<div class="space-y-6">
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

<!-- Sales Category Modal -->
<SalesCategoryModal
  show={showSalesCategoryModal}
  salesCategory={editingSalesCategory}
  {loading}
  on:save={handleSalesCategorySave}
  on:close={closeSalesCategoryModal}
/>

{#if showDeleteConfirm}
  <ConfirmationModal
    show={showDeleteConfirm}
    title="Delete Sales Category"
    message="Are you sure you want to delete this sales category? This action cannot be undone."
    confirmText="Delete"
    cancelText="Cancel"
    confirmClass="btn-error"
    {loading}
    on:confirm={confirmDelete}
    on:cancel={cancelDelete}
  />
{/if}