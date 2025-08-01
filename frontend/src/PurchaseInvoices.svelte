<script>
  import { onMount } from 'svelte'
  import { GetPurchaseInvoices, GetSuppliers, DeletePurchaseInvoice, CreatePurchaseInvoice, UpdatePurchaseInvoice, GetPurchaseInvoiceByID, MarkPurchaseInvoiceReceived } from '../wailsjs/go/main/App.js'
  import PageLayout from './components/PageLayout.svelte'
  import DataTable from './components/DataTable.svelte'
  import PurchaseInvoiceModal from './components/PurchaseInvoiceModal.svelte'
  import ClickableStatusBadge from './components/ClickableStatusBadge.svelte'
  import ConfirmationModal from './components/ConfirmationModal.svelte'
  import InventoryUpdateModal from './components/InventoryUpdateModal.svelte'
  import { showDbSuccess, showDbError } from './stores/notificationStore.js'

  /** @type {Array<any>} */
  let purchaseInvoices = []

  /** @type {Array<any>} */
  let suppliers = []

  /** @type {string} */
  let searchTerm = ''
  /** @type {boolean} */
  let showInvoiceModal = false
  /** @type {any} */
  let editingInvoice = null
  /** @type {boolean} */
  let loading = false
  /** @type {boolean} */
  let isLoading = false
  
  // Confirmation modal state
  let showDeleteConfirm = false
  let invoiceToDelete = null
  let confirmLoading = false

  // Inventory update modal state
  let showInventoryUpdateModal = false
  let inventoryUpdateMessage = ''
  let updatedProducts = []
  let markingReceived = false

  onMount(async () => {
    await loadData()
  })

  async function loadData() {
    isLoading = true
    try {
      const [invoicesResult, suppliersResult] = await Promise.all([
        GetPurchaseInvoices(),
        GetSuppliers()
      ])
      purchaseInvoices = invoicesResult || []
      suppliers = suppliersResult || []
    } catch (error) {
      console.error('Error loading data:', error)
      purchaseInvoices = []
      suppliers = []
    } finally {
      isLoading = false
    }
  }

  // Filter invoices based on search term
  $: filteredInvoices = purchaseInvoices.filter(invoice => 
    invoice.invoice_number?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (invoice.supplier && invoice.supplier.company_name?.toLowerCase().includes(searchTerm.toLowerCase())) ||
    invoice.notes?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    invoice.reference_number?.toLowerCase().includes(searchTerm.toLowerCase())
  )

  function handleAddInvoice() {
    editingInvoice = null
    showInvoiceModal = true
  }

  async function handleEditInvoice(invoice) {
    try {
      loading = true
      // Fetch complete invoice data with items and supplier details
      const completeInvoice = await GetPurchaseInvoiceByID(invoice.id)
      editingInvoice = completeInvoice
      showInvoiceModal = true
    } catch (error) {
      console.error('Error fetching invoice details:', error)
      showDbError('load', 'Purchase Invoice', error)
    } finally {
      loading = false
    }
  }

  async function handleInvoiceSave(event) {
    try {
      const invoiceData = event.detail
      
      if (editingInvoice) {
        // Update existing invoice
        await UpdatePurchaseInvoice({
          ...invoiceData,
          id: editingInvoice.id,
          supplier_id: parseInt(invoiceData.supplier_id),
          amount: parseFloat(invoiceData.amount),
          vat_amount: parseFloat(invoiceData.vat_amount) || 0,
          total_amount: parseFloat(invoiceData.total_amount),
          vat_rate: parseFloat(invoiceData.vat_rate) || 0,
          vat_inclusive: invoiceData.vat_inclusive || false
        })
        showDbSuccess('update', 'Purchase Invoice');
      } else {
        // Create new invoice
        await CreatePurchaseInvoice({
          ...invoiceData,
          supplier_id: parseInt(invoiceData.supplier_id),
          amount: parseFloat(invoiceData.amount),
          vat_amount: parseFloat(invoiceData.vat_amount) || 0,
          total_amount: parseFloat(invoiceData.total_amount),
          vat_rate: parseFloat(invoiceData.vat_rate) || 0,
          vat_inclusive: invoiceData.vat_inclusive || false
        })
        showDbSuccess('create', 'Purchase Invoice');
      }
      
      showInvoiceModal = false
      await loadData() // Refresh the list
    } catch (error) {
      console.error('Error saving invoice:', error)
      const operation = editingInvoice ? 'update' : 'create';
      showDbError(operation, 'Purchase Invoice', error)
    }
  }

  function handleInvoiceClose() {
    showInvoiceModal = false
    editingInvoice = null
  }

  function showDeleteConfirmation(invoice) {
    invoiceToDelete = invoice;
    showDeleteConfirm = true;
  }

  function cancelDelete() {
    showDeleteConfirm = false;
    invoiceToDelete = null;
  }

  async function confirmDelete() {
    if (!invoiceToDelete) return;
    
    try {
      confirmLoading = true;
      await DeletePurchaseInvoice(invoiceToDelete.id);
      await loadData(); // Refresh the list
      showDbSuccess('delete', 'Purchase Invoice');
      showDeleteConfirm = false;
      invoiceToDelete = null;
    } catch (error) {
      console.error('Error deleting invoice:', error);
      showDbError('delete', 'Purchase Invoice', error);
    } finally {
      confirmLoading = false;
    }
  }

  async function handleDeleteInvoice(invoice) {
    showDeleteConfirmation(invoice);
  }

  async function handleMarkAsReceived(invoice) {
    if (invoice.status === 'received') {
      showDbError('update', 'Purchase Invoice', new Error('Invoice is already marked as received'))
      return
    }

    try {
      markingReceived = true
      
      // Get complete invoice data with items
      const completeInvoice = await GetPurchaseInvoiceByID(invoice.id)
      
      // Mark invoice as received and update inventory
      await MarkPurchaseInvoiceReceived(invoice.id)
      
      // Prepare updated products list for display
      updatedProducts = completeInvoice.items?.map(item => ({
        name: item.product?.name || `Product ID: ${item.product_id}`,
        quantity: item.quantity,
        newStock: (item.product?.stock || 0) + Number(item.quantity)
      })) || []
      
      inventoryUpdateMessage = `Invoice ${invoice.invoice_number} has been marked as received and inventory has been updated.`
      showInventoryUpdateModal = true
      
      // Refresh the invoice list
      await loadData()
    } catch (error) {
      console.error('Error marking invoice as received:', error)
      showDbError('update', 'Purchase Invoice status', error)
    } finally {
      markingReceived = false
    }
  }

  function handleInventoryUpdateClose() {
    showInventoryUpdateModal = false
    inventoryUpdateMessage = ''
    updatedProducts = []
  }

  function handleInvoiceModalClose() {
    showInvoiceModal = false
    editingInvoice = null
    loading = false
  }

  function handleImport() {
    // TODO: Implement import functionality
    showDbError('Import functionality will be implemented soon')
  }

  function handleSearch(event) {
    searchTerm = event.detail.searchTerm
  }

  function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  }

  function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'SAR',
      minimumFractionDigits: 2
    }).format(amount)
  }

  function getStatusBadge(status) {
    const statusMap = {
      'pending': 'warning',
      'approved': 'info',
      'paid': 'success',
      'overdue': 'error',
      'cancelled': 'neutral'
    }
    const statusLabels = {
      'pending': 'Pending',
      'approved': 'Approved',
      'paid': 'Paid',
      'overdue': 'Overdue',
      'cancelled': 'Cancelled'
    }
    const badgeType = statusMap[status] || 'neutral'
    const label = statusLabels[status] || status
    return `<span class="status-badge status-badge-${badgeType}">${label}</span>`
  }

  /** @type {Array<{label: string, key?: string, labelArabic?: string, sortable?: boolean, class?: string, render?: Function, actions?: Array<{key: string, text: string, icon?: string, class?: string, title?: string}>}>} */
  const columns = [
    {
      key: 'invoice_number',
      label: 'Invoice Number',
      labelArabic: 'رقم الفاتورة',
      sortable: true,
      render: (invoice) => `
        <div>
          <div class="font-medium text-white">${invoice.invoice_number}</div>
          <div class="text-sm text-white/60">${invoice.reference_number || 'No reference'}</div>
        </div>
      `
    },
    {
      key: 'supplier',
      label: 'Supplier',
      labelArabic: 'المورد',
      sortable: true,
      render: (invoice) => `
        <div class="font-medium text-white">${invoice.supplier ? invoice.supplier.company_name : 'No supplier'}</div>
      `
    },
    {
      key: 'issue_date',
      label: 'Date',
      labelArabic: 'التاريخ',
      sortable: true,
      render: (invoice) => `
        <div>
          <div class="text-white">${formatDate(invoice.issue_date)}</div>
          <div class="text-sm text-white/60">Due: ${formatDate(invoice.due_date)}</div>
        </div>
      `
    },
    {
      key: 'total_amount',
      label: 'Amount',
      labelArabic: 'المبلغ',
      sortable: true,
      render: (invoice) => `
        <div>
          <div class="font-medium text-white">${formatCurrency(invoice.total_amount)}</div>
          <div class="text-sm text-white/60">VAT: ${formatCurrency(invoice.vat_amount)}</div>
        </div>
      `
    },
    {
      key: 'status',
      label: 'Status',
      labelArabic: 'الحالة',
      sortable: true
    },
    {
      label: 'Actions',
      actions: [
        { key: 'edit', text: 'Edit', icon: 'edit', class: 'btn-primary' },
        { key: 'delete', text: 'Delete', icon: 'delete', class: 'btn-error' }
      ]
    }
  ]
</script>

<PageLayout
  title="Purchase Invoices"
  icon="invoice"
>
  <svelte:fragment slot="actions">
    <button class="btn btn-secondary-outline" on:click={handleImport}>
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
      </svg>
      Import
    </button>
    <button class="btn btn-primary" on:click={handleAddInvoice}>
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
      </svg>
      Add Invoice
    </button>
  </svelte:fragment>

  <DataTable
    data={filteredInvoices}
    {columns}
    {loading}
    searchPlaceholder="Search invoices..."
    emptyStateTitle="No purchase invoices found"
    emptyStateMessage="Start by adding your first purchase invoice"
    emptyStateIcon="invoice"
    primaryAction={{ text: 'Add Invoice', icon: 'add' }}
    on:primary-action={handleAddInvoice}
    on:search={handleSearch}
    on:row-action={(e) => {
      const { action, item } = e.detail
      if (action === 'edit') {
        handleEditInvoice(item)
      } else if (action === 'delete') {
        handleDeleteInvoice(item)
      }
    }}
  >
    <svelte:fragment slot="cell" let:item let:column>
      {#if column.key === 'status'}
        <ClickableStatusBadge
          status={item.status}
          clickable={item.status === 'pending'}
          loading={markingReceived}
          on:click={() => handleMarkAsReceived(item)}
        />
      {/if}
    </svelte:fragment>
  </DataTable>
</PageLayout>

<PurchaseInvoiceModal
  show={showInvoiceModal}
  {editingInvoice}
  {suppliers}
  {loading}
  on:save={handleInvoiceSave}
  on:close={handleInvoiceModalClose}
/>

<!-- Confirmation Modal -->
<ConfirmationModal
  show={showDeleteConfirm}
  title="Delete Purchase Invoice"
  message={invoiceToDelete ? `Are you sure you want to delete invoice ${invoiceToDelete.invoice_number}? This action cannot be undone.` : ''}
  confirmText="Delete"
  cancelText="Cancel"
  loading={confirmLoading}
  on:confirm={confirmDelete}
  on:cancel={cancelDelete}
/>

<!-- Inventory Update Modal -->
<InventoryUpdateModal
  show={showInventoryUpdateModal}
  message={inventoryUpdateMessage}
  {updatedProducts}
  on:close={handleInventoryUpdateClose}
  on:confirm={handleInventoryUpdateClose}
/>