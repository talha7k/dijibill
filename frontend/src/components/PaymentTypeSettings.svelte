<script>
  import { GetPaymentTypes, CreatePaymentType, UpdatePaymentType, DeletePaymentType } from '../../wailsjs/go/main/App.js'
  import { database } from '../../wailsjs/go/models'
  import { createEventDispatcher } from 'svelte'
  import DataTable from './DataTable.svelte'
  import FormField from './FormField.svelte'
  import ConfirmationModal from './ConfirmationModal.svelte'

  const dispatch = createEventDispatcher()

  export let paymentTypes = []
  export let isLoading = false

  let newPaymentType = { name: '', name_arabic: '', code: '', description: '', is_default: false, is_active: true }
  let showPaymentTypeForm = false
  let editingPaymentType = null
  let showDeleteConfirm = false
  let paymentTypeToDelete = null
  let loading = false

  async function addPaymentType() {
    if (newPaymentType.name && newPaymentType.code) {
      try {
        const paymentTypeData = new database.PaymentType({
          id: 0,
          name: newPaymentType.name,
          name_arabic: newPaymentType.name_arabic || '',
          code: newPaymentType.code,
          description: newPaymentType.description || '',
          is_default: newPaymentType.is_default || false,
          is_active: true,
          created_at: null,
          updated_at: null
        });
        await CreatePaymentType(paymentTypeData)
        dispatch('reload')
        resetPaymentTypeForm()
      } catch (error) {
        console.error('Error adding payment type:', error)
        dispatch('error', { message: 'Error adding payment type: ' + error.message })
      }
    }
  }

  function editPaymentType(paymentType) {
    editingPaymentType = paymentType.id
    newPaymentType = { ...paymentType }
    showPaymentTypeForm = true
  }

  async function updatePaymentType() {
    try {
      const paymentTypeData = new database.PaymentType({
        id: editingPaymentType,
        name: newPaymentType.name,
        name_arabic: newPaymentType.name_arabic || '',
        code: newPaymentType.code,
        description: newPaymentType.description || '',
        is_default: newPaymentType.is_default || false,
        is_active: newPaymentType.is_active !== undefined ? newPaymentType.is_active : true,
        created_at: newPaymentType.created_at || null,
        updated_at: newPaymentType.updated_at || null
      });
      await UpdatePaymentType(paymentTypeData)
      dispatch('reload')
      resetPaymentTypeForm()
    } catch (error) {
      console.error('Error updating payment type:', error)
      dispatch('error', { message: 'Error updating payment type: ' + error.message })
    }
  }

  function showDeleteConfirmation(id) {
    paymentTypeToDelete = id
    showDeleteConfirm = true
  }

  function cancelDelete() {
    showDeleteConfirm = false
    paymentTypeToDelete = null
  }

  async function confirmDelete() {
    if (paymentTypeToDelete) {
      loading = true
      try {
        await DeletePaymentType(paymentTypeToDelete)
        dispatch('reload')
        showDeleteConfirm = false
        paymentTypeToDelete = null
      } catch (error) {
        console.error('Error deleting payment type:', error)
        dispatch('error', { message: 'Error deleting payment type: ' + error.message })
      } finally {
        loading = false
      }
    }
  }

  function resetPaymentTypeForm() {
    newPaymentType = { name: '', name_arabic: '', code: '', description: '', is_default: false, is_active: true }
    showPaymentTypeForm = false
    editingPaymentType = null
  }

  async function setDefaultPaymentType(id) {
    try {
      // First, update all payment types to not be default
      for (const paymentType of paymentTypes) {
        if (paymentType.is_default) {
          const updatedPaymentType = new database.PaymentType({ ...paymentType, is_default: false });
          await UpdatePaymentType(updatedPaymentType);
        }
      }
      // Then set the selected one as default
      const selectedPaymentType = paymentTypes.find(pt => pt.id === id);
      if (selectedPaymentType) {
        const updatedPaymentType = new database.PaymentType({ ...selectedPaymentType, is_default: true });
        await UpdatePaymentType(updatedPaymentType);
        dispatch('reload')
      }
    } catch (error) {
      console.error('Error setting default payment type:', error)
      dispatch('error', { message: 'Error setting default payment type: ' + error.message })
    }
  }

  // DataTable configuration
  let searchTerm = ''
  
  $: filteredPaymentTypes = paymentTypes.filter(paymentType => 
    paymentType.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (paymentType.name_arabic && paymentType.name_arabic.includes(searchTerm)) ||
    paymentType.code.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (paymentType.description && paymentType.description.toLowerCase().includes(searchTerm.toLowerCase()))
  )

  const columns = [
    {
      key: 'name',
      label: 'Name',
      sortable: true,
      render: (paymentType) => `
        <div>
          <div class="font-medium text-gray-100">${paymentType.name}</div>
          ${paymentType.name_arabic ? `<div class="text-sm text-gray-300">${paymentType.name_arabic}</div>` : ''}
        </div>
      `
    },
    {
      key: 'code',
      label: 'Code',
      sortable: true,
      render: (paymentType) => `<span class="font-mono text-sm">${paymentType.code}</span>`
    },
    {
      key: 'description',
      label: 'Description',
      render: (paymentType) => paymentType.description || '-'
    },
    {
      key: 'status',
      label: 'Status',
      render: (paymentType) => `
        <div class="flex items-center space-x-2">
          ${paymentType.is_default ? '<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">Default</span>' : ''}
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${paymentType.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
            ${paymentType.is_active ? 'Active' : 'Inactive'}
          </span>
        </div>
      `
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (paymentType) => `
        <div class="flex items-center space-x-2">
          ${!paymentType.is_default ? `<button class="action-btn" data-action="setDefault" data-id="${paymentType.id}">Set Default</button>` : ''}
          <button class="action-btn" data-action="edit" data-id="${paymentType.id}">Edit</button>
          <button class="action-btn text-red-600" data-action="delete" data-id="${paymentType.id}">Delete</button>
        </div>
      `
    }
  ]

  const primaryAction = {
    label: 'Add Payment Type',
    icon: 'fa-plus'
  }

  // Event handlers for DataTable
  function handlePrimaryAction() {
    showPaymentTypeForm = true
  }

  function handleSearch(event) {
    searchTerm = event.detail.searchTerm
  }

  function handleRowAction(event) {
    const { action, id } = event.detail
    const paymentType = paymentTypes.find(pt => pt.id === parseInt(id))
    
    switch (action) {
      case 'edit':
        editPaymentType(paymentType)
        break
      case 'delete':
        showDeleteConfirmation(parseInt(id))
        break
      case 'setDefault':
        setDefaultPaymentType(parseInt(id))
        break
    }
  }
</script>

<div class="space-y-6">
  <!-- Payment Type Form -->
  {#if showPaymentTypeForm}
    <div class="bg-gray-50 p-6 rounded-lg border">
      <h4 class="text-lg font-medium text-gray-100 mb-6">
        {editingPaymentType ? 'Edit Payment Type' : 'Add New Payment Type'}
      </h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <FormField
          label="Name"
          labelArabic="الاسم"
          type="text"
          bind:value={newPaymentType.name}
          placeholder="Enter payment type name"
          required
        />

        <FormField
          label="Name (Arabic)"
          labelArabic="الاسم بالعربية"
          type="text"
          bind:value={newPaymentType.name_arabic}
          placeholder="اسم نوع الدفع"
          dir="rtl"
        />

        <FormField
          label="Code"
          labelArabic="الرمز"
          type="text"
          bind:value={newPaymentType.code}
          placeholder="Enter payment type code"
          required
        />

        <FormField
          label="Description"
          labelArabic="الوصف"
          type="text"
          bind:value={newPaymentType.description}
          placeholder="Enter description"
        />

        <div class="flex items-center space-x-6 col-span-2">
          <FormField
            type="checkbox"
            bind:checked={newPaymentType.is_default}
            placeholder="Set as default"
          />
          <FormField
            type="checkbox"
            bind:checked={newPaymentType.is_active}
            placeholder="Active"
          />
        </div>
      </div>
      
      <div class="flex justify-end space-x-3 mt-6">
        <button
          on:click={resetPaymentTypeForm}
          class="px-4 py-2 border border-gray-300 text-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          Cancel
        </button>
        <button
          on:click={editingPaymentType ? updatePaymentType : addPaymentType}
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          {editingPaymentType ? 'Update' : 'Add'} Payment Type
        </button>
      </div>
    </div>
  {/if}

  <!-- Payment Types DataTable -->
  <DataTable
    data={filteredPaymentTypes}
    {columns}
    {primaryAction}
    loading={isLoading}
    searchTerm={searchTerm}
    searchPlaceholder="Search payment types..."
    emptyStateTitle="No payment types found"
    emptyStateMessage="Start by adding your first payment type"
    emptyStateIcon="fa-credit-card"
    on:primary-action={handlePrimaryAction}
    on:search={handleSearch}
    on:row-action={handleRowAction}
  />
</div>

{#if showDeleteConfirm}
  <ConfirmationModal
    show={showDeleteConfirm}
    title="Delete Payment Type"
    message="Are you sure you want to delete this payment type? This action cannot be undone."
    confirmText="Delete"
    cancelText="Cancel"
    confirmClass="btn-error"
    {loading}
    on:confirm={confirmDelete}
    on:cancel={cancelDelete}
  />
{/if}