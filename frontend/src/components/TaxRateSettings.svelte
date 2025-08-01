<script>
  import { GetTaxRates, CreateTaxRate, UpdateTaxRate, DeleteTaxRate } from '../../wailsjs/go/main/App.js'
  import { database } from '../../wailsjs/go/models'
  import { createEventDispatcher } from 'svelte'
  import DataTable from './DataTable.svelte'
  import TaxRateModal from './TaxRateModal.svelte'
  import ConfirmationModal from './ConfirmationModal.svelte'
  import { showDbSuccess, showDbError } from '../stores/notificationStore.js'

  const dispatch = createEventDispatcher()

  export let taxRates = []
  export let isLoading = false

  let showTaxRateModal = false
  let editingTaxRate = null
  let showDeleteConfirm = false
  let taxRateToDelete = null
  let loading = false
  let searchTerm = ''

  function editTaxRate(taxRate) {
    editingTaxRate = taxRate
    showTaxRateModal = true
  }

  function closeTaxRateModal() {
    showTaxRateModal = false
    editingTaxRate = null
  }

  async function handleTaxRateSave(event) {
    const taxRateData = event.detail
    loading = true
    
    try {
      if (editingTaxRate) {
        // Update existing tax rate
        const updatedTaxRate = new database.TaxRate({
          ...editingTaxRate,
          ...taxRateData
        })
        await UpdateTaxRate(updatedTaxRate)
        showDbSuccess('update', 'Tax Rate')
      } else {
        // Create new tax rate
        const newTaxRate = new database.TaxRate({
          id: 0,
          ...taxRateData,
          created_at: null,
          updated_at: null
        })
        await CreateTaxRate(newTaxRate)
        showDbSuccess('create', 'Tax Rate')
      }
      
      dispatch('reload')
      closeTaxRateModal()
    } catch (error) {
      console.error('Error saving tax rate:', error)
      const operation = editingTaxRate ? 'update' : 'create'
      showDbError(operation, 'Tax Rate', error)
    } finally {
      loading = false
    }
  }

  function showDeleteConfirmation(taxRate) {
    taxRateToDelete = taxRate
    showDeleteConfirm = true
  }

  function cancelDelete() {
    showDeleteConfirm = false
    taxRateToDelete = null
  }

  async function confirmDelete() {
    if (!taxRateToDelete) return
    
    try {
      loading = true
      await DeleteTaxRate(taxRateToDelete.id)
      showDbSuccess('delete', 'Tax Rate')
      dispatch('reload')
      showDeleteConfirm = false
      taxRateToDelete = null
    } catch (error) {
      console.error('Error deleting tax rate:', error)
      showDbError('delete', 'Tax Rate', error)
    } finally {
      loading = false
    }
  }

  // Filter tax rates based on search term
  $: filteredTaxRates = taxRates.filter(taxRate => {
    if (!searchTerm) return true
    const searchLower = searchTerm.toLowerCase()
    return (
      taxRate.name.toLowerCase().includes(searchLower) ||
      (taxRate.name_arabic && taxRate.name_arabic.includes(searchTerm)) ||
      taxRate.rate.toString().includes(searchTerm)
    )
  })

  // DataTable configuration
  const columns = [
    {
      key: 'name',
      label: 'Name',
      render: (item) => {
        let html = `<div class="font-medium">${item.name}</div>`
        if (item.name_arabic) {
          html += `<div class="text-sm text-gray-300">${item.name_arabic}</div>`
        }
        return html
      }
    },
    {
      key: 'rate',
      label: 'Rate',
      render: (item) => `${item.rate}%`
    },
    {
      key: 'status',
      label: 'Status',
      render: (item) => {
        let badges = ''
        if (item.is_default) {
          badges += '<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 mr-2">Default</span>'
        }
        badges += `<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${item.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">${item.is_active ? 'Active' : 'Inactive'}</span>`
        return badges
      }
    },
    {
      label: 'Actions',
      actions: [
        {
          key: 'edit',
          text: 'Edit',
          icon: 'fa-edit',
          class: 'btn-secondary',
          title: 'Edit tax rate'
        },
        {
          key: 'delete',
          text: 'Delete',
          icon: 'fa-trash',
          class: 'btn-error',
          title: 'Delete tax rate'
        }
      ]
    }
  ]

  const primaryAction = {
    label: 'Add Tax Rate',
    icon: 'fa-plus'
  }

  // Event handlers for DataTable
  function handlePrimaryAction() {
    showTaxRateModal = true
  }

  function handleSearch(event) {
    searchTerm = event.detail.searchTerm
  }

  function handleRowAction(event) {
    const { action, item } = event.detail
    
    switch (action) {
      case 'edit':
        editTaxRate(item)
        break
      case 'delete':
        showDeleteConfirmation(item)
        break
    }
  }
</script>
<div class="space-y-6">
  <!-- Tax Rates DataTable -->
  <DataTable
    data={filteredTaxRates}
    {columns}
    {primaryAction}
    loading={isLoading}
    searchTerm={searchTerm}
    searchPlaceholder="Search tax rates..."
    emptyStateTitle="No tax rates found"
    emptyStateMessage="Start by adding your first tax rate"
    emptyStateIcon="fa-percentage"
    on:primary-action={handlePrimaryAction}
    on:search={handleSearch}
    on:row-action={handleRowAction}
  />
</div>

<!-- Tax Rate Modal -->
<TaxRateModal
  show={showTaxRateModal}
  taxRate={editingTaxRate}
  {loading}
  on:save={handleTaxRateSave}
  on:close={closeTaxRateModal}
/>

<!-- Confirmation Modal -->
<ConfirmationModal
  show={showDeleteConfirm}
  title="Delete Tax Rate"
  message={taxRateToDelete ? `Are you sure you want to delete "${taxRateToDelete.name}"? This action cannot be undone.` : ''}
  confirmText="Delete"
  cancelText="Cancel"
  loading={loading}
  on:confirm={confirmDelete}
  on:cancel={cancelDelete}
/>