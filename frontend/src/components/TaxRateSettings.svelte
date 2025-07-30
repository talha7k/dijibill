<script>
  import { GetTaxRates, CreateTaxRate, UpdateTaxRate, DeleteTaxRate } from '../../wailsjs/go/main/App.js'
   import { database } from '../../wailsjs/go/models'
  import { createEventDispatcher } from 'svelte'
  import DataTable from './DataTable.svelte'
  import FormField from './FormField.svelte'
  import ConfirmationModal from './ConfirmationModal.svelte'

  const dispatch = createEventDispatcher()

  /** @type {Array<any>} */
  export let taxRates = []
  /** @type {boolean} */
  export let isLoading = false

  /** @type {{name: string, name_arabic: string, rate: number, is_default: boolean, is_active: boolean, description?: string, created_at?: any, updated_at?: any}} */
  let newTaxRate = { name: '', name_arabic: '', rate: 0, is_default: false, is_active: true }
  /** @type {boolean} */
  let showTaxRateForm = false
  /** @type {any} */
  let editingTaxRate = null
  /** @type {string} */
  let searchTerm = ''
  
  // Confirmation modal state
  let showDeleteConfirm = false
  let taxRateToDelete = null
  let loading = false

  // DataTable configuration
  /** @type {Array<{label: string, key?: string, class?: string, render?: Function, actions?: Array<{key: string, text: string, icon?: string, class?: string, title?: string}>}>} */
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
          key: 'setDefault',
          text: 'Set Default',
          icon: 'fa-star',
          class: 'btn-secondary',
          title: 'Set as default tax rate'
        },
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

  /** @type {{text: string, icon: string}} */
  const primaryAction = {
    text: 'Add Tax Rate',
    icon: 'fa-plus'
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

  async function addTaxRate() {
    if (newTaxRate.name && newTaxRate.rate >= 0) {
      try {
        const taxRateData = new database.TaxRate({
          id: 0,
          name: newTaxRate.name,
          name_arabic: newTaxRate.name_arabic || '',
          rate: newTaxRate.rate,
          is_default: newTaxRate.is_default || false,
          is_active: true,
          description: '',
          created_at: null,
          updated_at: null
        });
        await CreateTaxRate(taxRateData)
        dispatch('reload')
        resetTaxRateForm()
      } catch (error) {
        console.error('Error adding tax rate:', error)
        dispatch('error', { message: 'Error adding tax rate: ' + error.message })
      }
    }
  }

  function editTaxRate(taxRate) {
    editingTaxRate = taxRate.id
    newTaxRate = { ...taxRate }
    showTaxRateForm = true
  }

  async function updateTaxRate() {
    try {
      const taxRateData = new database.TaxRate({
        id: editingTaxRate,
        name: newTaxRate.name,
        name_arabic: newTaxRate.name_arabic || '',
        rate: newTaxRate.rate,
        is_default: newTaxRate.is_default || false,
        is_active: newTaxRate.is_active !== undefined ? newTaxRate.is_active : true,
        description: newTaxRate.description || '',
        created_at: newTaxRate.created_at || null,
        updated_at: newTaxRate.updated_at || null
      });
      await UpdateTaxRate(taxRateData)
      dispatch('reload')
      resetTaxRateForm()
    } catch (error) {
      console.error('Error updating tax rate:', error)
      dispatch('error', { message: 'Error updating tax rate: ' + error.message })
    }
  }

  function showDeleteConfirmation(taxRate) {
    taxRateToDelete = taxRate;
    showDeleteConfirm = true;
  }

  function cancelDelete() {
    showDeleteConfirm = false;
    taxRateToDelete = null;
  }

  async function confirmDelete() {
    if (!taxRateToDelete) return;
    
    try {
      loading = true;
      await DeleteTaxRate(taxRateToDelete.id);
      dispatch('reload');
      showDeleteConfirm = false;
      taxRateToDelete = null;
    } catch (error) {
      console.error('Error deleting tax rate:', error);
      dispatch('error', { message: 'Error deleting tax rate: ' + error.message });
    } finally {
      loading = false;
    }
  }

  async function deleteTaxRate(id) {
    if (confirm('Are you sure you want to delete this tax rate?')) {
      try {
        await DeleteTaxRate(id)
        dispatch('reload')
      } catch (error) {
        console.error('Error deleting tax rate:', error)
        dispatch('error', { message: 'Error deleting tax rate: ' + error.message })
      }
    }
  }

  function resetTaxRateForm() {
    newTaxRate = { name: '', name_arabic: '', rate: 0, is_default: false, is_active: true }
    showTaxRateForm = false
    editingTaxRate = null
  }

  async function setDefaultTaxRate(id) {
    try {
      // First, update all tax rates to not be default
      for (const taxRate of taxRates) {
        if (taxRate.is_default) {
          const updatedTaxRate = new database.TaxRate({ ...taxRate, is_default: false });
          await UpdateTaxRate(updatedTaxRate);
        }
      }
      // Then set the selected one as default
      const selectedTaxRate = taxRates.find(tr => tr.id === id);
      if (selectedTaxRate) {
        const updatedTaxRate = new database.TaxRate({ ...selectedTaxRate, is_default: true });
        await UpdateTaxRate(updatedTaxRate);
        dispatch('reload')
      }
    } catch (error) {
      console.error('Error setting default tax rate:', error)
      dispatch('error', { message: 'Error setting default tax rate: ' + error.message })
    }
  }

  // Computed properties for FormField string conversion
  $: rateValue = newTaxRate.rate?.toString() || ''
  
  function handleRateChange(event) {
    const value = parseFloat(event.target.value)
    newTaxRate.rate = isNaN(value) ? 0 : value
  }

  // DataTable event handlers
  function handlePrimaryAction() {
    showTaxRateForm = true
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
      case 'setDefault':
        if (!item.is_default) {
          setDefaultTaxRate(item.id)
        }
        break
    }
  }
</script>

<div class="space-y-6">
  <!-- Tax Rate Form -->
  {#if showTaxRateForm}
    <div class=" p-6 rounded-lg border">
      <h4 class="text-lg font-medium text-gray-100 mb-6">
        {editingTaxRate ? 'Edit Tax Rate' : 'Add New Tax Rate'}
      </h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <FormField
          label="Name"
          labelArabic="الاسم"
          type="text"
          bind:value={newTaxRate.name}
          placeholder="Enter tax rate name"
          required
        />

        <FormField
          label="Name (Arabic)"
          labelArabic="الاسم بالعربية"
          type="text"
          bind:value={newTaxRate.name_arabic}
          placeholder="اسم معدل الضريبة"
          dir="rtl"
        />

        <FormField
           label="Rate (%)"
           labelArabic="المعدل (%)"
           type="number"
           value={rateValue}
           placeholder="Enter tax rate percentage"
           min="0"
           max="100"
           step="0.01"
           required
           on:input={handleRateChange}
         />

        <div class="flex items-center space-x-6">
          <FormField
            type="checkbox"
            bind:checked={newTaxRate.is_default}
            placeholder="Set as default"
          />
          <FormField
            type="checkbox"
            bind:checked={newTaxRate.is_active}
            placeholder="Active"
          />
        </div>
      </div>
      
      <div class="flex justify-end space-x-3 mt-6">
        <button
          on:click={resetTaxRateForm}
          class="px-4 py-2 border border-gray-300 text-gray-300 rounded-md hover: focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          Cancel
        </button>
        <button
          on:click={editingTaxRate ? updateTaxRate : addTaxRate}
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          {editingTaxRate ? 'Update' : 'Add'} Tax Rate
        </button>
      </div>
    </div>
  {/if}

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