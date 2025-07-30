<script>
  import { onMount } from 'svelte'
  import { GetPurchaseInvoices, GetSuppliers, DeletePurchaseInvoice, CreatePurchaseInvoice, UpdatePurchaseInvoice } from '../wailsjs/go/main/App.js'
  import PageLayout from './components/PageLayout.svelte'
  import DataTable from './components/DataTable.svelte'
  import PurchaseInvoiceModal from './components/PurchaseInvoiceModal.svelte'
  import StatusBadge from './components/StatusBadge.svelte'

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

  function handleEditInvoice(invoice) {
    editingInvoice = invoice
    showInvoiceModal = true
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
          tax_amount: parseFloat(invoiceData.tax_amount),
          total_amount: parseFloat(invoiceData.total_amount)
        })
      } else {
        // Create new invoice
        await CreatePurchaseInvoice({
          ...invoiceData,
          supplier_id: parseInt(invoiceData.supplier_id),
          amount: parseFloat(invoiceData.amount),
          tax_amount: parseFloat(invoiceData.tax_amount),
          total_amount: parseFloat(invoiceData.total_amount)
        })
      }
      
      showInvoiceModal = false
      await loadData() // Refresh the list
    } catch (error) {
      console.error('Error saving invoice:', error)
      alert('Error saving invoice: ' + error.message)
    }
  }

  function handleInvoiceClose() {
    showInvoiceModal = false
    editingInvoice = null
  }

  async function handleDeleteInvoice(invoice) {
    if (confirm(`Are you sure you want to delete invoice ${invoice.invoice_number}?`)) {
      try {
        await DeletePurchaseInvoice(invoice.id)
        await loadData() // Refresh the list
      } catch (error) {
        console.error('Error deleting invoice:', error)
        alert('Error deleting invoice: ' + error.message)
      }
    }
  }

  function handleInvoiceModalClose() {
    showInvoiceModal = false
    editingInvoice = null
    loading = false
  }

  function handleImport() {
    // TODO: Implement import functionality
    alert('Import functionality will be implemented soon')
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
      sortable: true,
      render: (invoice) => getStatusBadge(invoice.status)
    },
    {
      label: 'Actions',
      actions: [
        { key: 'edit', text: 'Edit', icon: 'fa-edit', class: 'btn-primary' },
        { key: 'delete', text: 'Delete', icon: 'fa-trash', class: 'btn-error' }
      ]
    }
  ]
</script>

<PageLayout
  title="Purchase Invoices"
  icon="fa-file-invoice"
>
  <svelte:fragment slot="actions">
    <button class="btn btn-secondary-outline" on:click={handleImport}>
      <i class="fas fa-upload"></i>
      Import
    </button>
    <button class="btn btn-primary" on:click={handleAddInvoice}>
      <i class="fas fa-plus"></i>
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
    emptyStateIcon="fa-file-invoice"
    primaryAction={{ text: 'Add Invoice', icon: 'fa-plus' }}
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
  />
</PageLayout>

<PurchaseInvoiceModal
  show={showInvoiceModal}
  {editingInvoice}
  {suppliers}
  {loading}
  on:save={handleInvoiceSave}
  on:close={handleInvoiceModalClose}
/>