<script>
  import { GetPaymentTypes, CreatePaymentType, UpdatePaymentType, DeletePaymentType } from '../../wailsjs/go/main/App.js'
  import { database } from '../../wailsjs/go/models'
  import { createEventDispatcher } from 'svelte'
  import { showDbSuccess, showDbError } from '../stores/notificationStore.js'
  import DataTable from './DataTable.svelte'
  import PaymentTypeModal from './PaymentTypeModal.svelte'
  import ConfirmationModal from './ConfirmationModal.svelte'

  const dispatch = createEventDispatcher()

  export let paymentTypes = []
  export let isLoading = false

  let showPaymentTypeModal = false
  let editingPaymentType = null
  let showDeleteConfirm = false
  let paymentTypeToDelete = null
  let loading = false

  async function handlePaymentTypeSave(event) {
    const paymentTypeData = event.detail
    
    try {
      loading = true
      const dbPaymentType = new database.PaymentType({
        id: paymentTypeData.id || 0,
        name: paymentTypeData.name,
        name_arabic: paymentTypeData.name_arabic || '',
        code: paymentTypeData.code,
        description: paymentTypeData.description || '',
        is_default: paymentTypeData.is_default || false,
        is_active: paymentTypeData.is_active !== undefined ? paymentTypeData.is_active : true,
        created_at: paymentTypeData.created_at || null,
        updated_at: paymentTypeData.updated_at || null
      })

      if (editingPaymentType) {
        await UpdatePaymentType(dbPaymentType)
        showDbSuccess('update', 'Payment Type')
      } else {
        await CreatePaymentType(dbPaymentType)
        showDbSuccess('create', 'Payment Type')
      }
      
      dispatch('reload')
      closePaymentTypeModal()
    } catch (error) {
      console.error('Error saving payment type:', error)
      showDbError('save', 'Payment Type', error)
    } finally {
      loading = false
    }
  }

  function editPaymentType(paymentType) {
    editingPaymentType = paymentType
    showPaymentTypeModal = true
  }

  function closePaymentTypeModal() {
    showPaymentTypeModal = false
    editingPaymentType = null
    loading = false
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
        showDbSuccess('delete', 'Payment Type')
        dispatch('reload')
        showDeleteConfirm = false
        paymentTypeToDelete = null
      } catch (error) {
        console.error('Error deleting payment type:', error)
        showDbError('delete', 'Payment Type', error)
      } finally {
        loading = false
      }
    }
  }

  function resetPaymentTypeForm() {
    closePaymentTypeModal()
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
      label: 'Actions',
      actions: [
        {
          key: 'edit',
          text: 'Edit',
          icon: 'edit',
          class: 'btn-secondary',
          title: 'Edit payment type'
        },
        {
          key: 'delete',
          text: 'Delete',
          icon: 'delete',
          class: 'btn-error',
          title: 'Delete payment type'
        }
      ]
    }
  ]

  const primaryAction = {
    label: 'Add Payment Type',
    icon: 'add'
  }

  // Event handlers for DataTable
  function handlePrimaryAction() {
    showPaymentTypeModal = true
  }

  function handleSearch(event) {
    searchTerm = event.detail.searchTerm
  }

  function handleRowAction(event) {
    const { action, item } = event.detail
    
    switch (action) {
      case 'edit':
        editPaymentType(item)
        break
      case 'delete':
        showDeleteConfirmation(item.id)
        break
    }
  }
</script>

<div class="space-y-6">
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
    emptyStateIcon="payment"
    on:primary-action={handlePrimaryAction}
    on:search={handleSearch}
    on:row-action={handleRowAction}
  />
</div>

<!-- Payment Type Modal -->
<PaymentTypeModal
  show={showPaymentTypeModal}
  paymentType={editingPaymentType}
  {loading}
  on:save={handlePaymentTypeSave}
  on:close={closePaymentTypeModal}
/>

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