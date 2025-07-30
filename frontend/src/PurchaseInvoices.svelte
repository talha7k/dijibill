<script>
  import PageLayout from './components/PageLayout.svelte'
  import DataTable from './components/DataTable.svelte'
  import PurchaseInvoiceModal from './components/PurchaseInvoiceModal.svelte'
  import StatusBadge from './components/StatusBadge.svelte'

  let purchaseInvoices = [
    {
      id: 1,
      invoice_number: 'PI-2024-001',
      supplier_id: 1,
      supplier_name: 'Tech Solutions Ltd',
      invoice_date: '2024-01-15',
      due_date: '2024-02-14',
      amount: 5000.00,
      tax_amount: 750.00,
      total_amount: 5750.00,
      currency: 'SAR',
      status: 'pending',
      description: 'Office equipment and supplies',
      description_arabic: 'معدات ولوازم مكتبية',
      payment_terms: 'net_30',
      reference_number: 'REF-001',
      notes: 'Urgent delivery required'
    },
    {
      id: 2,
      invoice_number: 'PI-2024-002',
      supplier_id: 2,
      supplier_name: 'Global Supplies Co',
      invoice_date: '2024-02-01',
      due_date: '2024-02-16',
      amount: 3200.00,
      tax_amount: 480.00,
      total_amount: 3680.00,
      currency: 'SAR',
      status: 'approved',
      description: 'Monthly maintenance supplies',
      description_arabic: 'لوازم الصيانة الشهرية',
      payment_terms: 'net_15',
      reference_number: 'REF-002',
      notes: ''
    }
  ]

  let suppliers = [
    { id: 1, company_name: 'Tech Solutions Ltd' },
    { id: 2, company_name: 'Global Supplies Co' },
    { id: 3, company_name: 'Office Depot' }
  ]

  let searchTerm = ''
  let showInvoiceModal = false
  let editingInvoice = null
  let loading = false

  // Filter invoices based on search term
  $: filteredInvoices = purchaseInvoices.filter(invoice => 
    invoice.invoice_number.toLowerCase().includes(searchTerm.toLowerCase()) ||
    invoice.supplier_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    invoice.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
    invoice.reference_number.toLowerCase().includes(searchTerm.toLowerCase())
  )

  function handleAddInvoice() {
    editingInvoice = null
    showInvoiceModal = true
  }

  function handleEditInvoice(invoice) {
    editingInvoice = invoice
    showInvoiceModal = true
  }

  function handleDeleteInvoice(invoice) {
    if (confirm(`Are you sure you want to delete invoice ${invoice.invoice_number}?`)) {
      purchaseInvoices = purchaseInvoices.filter(i => i.id !== invoice.id)
    }
  }

  function handleInvoiceSave(event) {
    loading = true
    const invoiceData = event.detail

    setTimeout(() => {
      if (editingInvoice) {
        // Update existing invoice
        const index = purchaseInvoices.findIndex(i => i.id === editingInvoice.id)
        if (index !== -1) {
          // Get supplier name
          const supplier = suppliers.find(s => s.id == invoiceData.supplier_id)
          purchaseInvoices[index] = { 
            ...purchaseInvoices[index], 
            ...invoiceData,
            supplier_name: supplier ? supplier.company_name : 'Unknown Supplier'
          }
          purchaseInvoices = [...purchaseInvoices]
        }
      } else {
        // Add new invoice
        const supplier = suppliers.find(s => s.id == invoiceData.supplier_id)
        const newInvoice = {
          id: Math.max(...purchaseInvoices.map(i => i.id)) + 1,
          ...invoiceData,
          supplier_name: supplier ? supplier.company_name : 'Unknown Supplier'
        }
        purchaseInvoices = [...purchaseInvoices, newInvoice]
      }
      
      showInvoiceModal = false
      editingInvoice = null
      loading = false
    }, 1000)
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
    searchTerm = event.detail
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
      key: 'supplier_name',
      label: 'Supplier',
      labelArabic: 'المورد',
      sortable: true,
      render: (invoice) => `
        <div class="font-medium text-white">${invoice.supplier_name}</div>
      `
    },
    {
      key: 'invoice_date',
      label: 'Date',
      labelArabic: 'التاريخ',
      sortable: true,
      render: (invoice) => `
        <div>
          <div class="text-white">${formatDate(invoice.invoice_date)}</div>
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
          <div class="text-sm text-white/60">Tax: ${formatCurrency(invoice.tax_amount)}</div>
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
      actions: [
        { key: 'edit', text: 'Edit', icon: 'fa-edit', class: 'btn-primary' },
        { key: 'delete', text: 'Delete', icon: 'fa-trash', class: 'btn-danger' }
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