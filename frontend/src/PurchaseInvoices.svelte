<script>
  import PageLayout from './components/PageLayout.svelte'
  import DataTable from './components/DataTable.svelte'
  import StatusBadge from './components/StatusBadge.svelte'
  
  // Purchase Invoices functionality will be implemented here
  let purchaseInvoices = []
  let loading = false
  let searchTerm = ''

  // Mock data for demonstration
  purchaseInvoices = [
    // Add some mock data when backend is ready
  ]

  function formatDate(dateString) {
    if (!dateString) return 'N/A'
    try {
      return new Date(dateString).toLocaleDateString()
    } catch (error) {
      return 'Invalid Date'
    }
  }

  function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'SAR'
    }).format(amount)
  }

  $: filteredPurchases = (purchaseInvoices || []).filter(purchase => 
    purchase.invoice_number?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (purchase.supplier && purchase.supplier.name?.toLowerCase().includes(searchTerm.toLowerCase()))
  )

  // Table configuration
  const columns = [
    { key: 'invoice_number', label: 'Invoice #', render: (item) => `<span class="font-medium">${item.invoice_number}</span>` },
    { 
      key: 'supplier', 
      label: 'Supplier', 
      render: (item) => {
        if (item.supplier) {
          return `<div><div class="font-medium">${item.supplier.name}</div>${item.supplier.email ? `<div class="text-sm opacity-70">${item.supplier.email}</div>` : ''}</div>`
        }
        return '<span class="opacity-40">No supplier</span>'
      }
    },
    { key: 'date', label: 'Date', render: (item) => formatDate(item.date) },
    { key: 'total_amount', label: 'Amount', render: (item) => `<span class="font-medium">${formatCurrency(item.total_amount)}</span>` },
    { 
      key: 'status', 
      label: 'Status',
      actions: [
        { key: 'edit', icon: 'fa-edit', class: 'btn-warning', title: 'Edit Purchase' },
        { key: 'delete', icon: 'fa-trash', class: 'btn-error', title: 'Delete Purchase' }
      ]
    }
  ]

  const primaryAction = {
    text: 'New Purchase',
    icon: 'fa-plus'
  }

  const secondaryActions = [
    {
      text: 'Import',
      icon: 'fa-upload'
    }
  ]

  function handlePrimaryAction() {
    console.log('New Purchase clicked')
    // TODO: Implement new purchase functionality
  }

  function handleSecondaryAction(action) {
    if (action.text === 'Import') {
      console.log('Import clicked')
      // TODO: Implement import functionality
    }
  }

  function handleRowAction(event) {
    const { action, item } = event.detail
    if (action === 'edit') {
      console.log('Edit purchase:', item)
      // TODO: Implement edit functionality
    } else if (action === 'delete') {
      console.log('Delete purchase:', item)
      // TODO: Implement delete functionality
    }
  }

  function handleSearch(event) {
    searchTerm = event.detail.searchTerm
  }
</script>

<PageLayout 
  title="Purchase Invoices" 
  icon="fa-file-invoice-dollar" 
  showIndicator={true}
>
  <svelte:fragment slot="actions">
    <!-- Actions are handled by DataTable -->
  </svelte:fragment>

  <DataTable
    data={filteredPurchases}
    {columns}
    {loading}
    {searchTerm}
    searchPlaceholder="Search purchases..."
    emptyStateTitle="No purchase invoices found"
    emptyStateMessage="Create your first purchase to get started"
    emptyStateIcon="fa-file-invoice-dollar"
    {primaryAction}
    {secondaryActions}
    on:primary-action={handlePrimaryAction}
    on:secondary-action={handleSecondaryAction}
    on:row-action={handleRowAction}
    on:search={handleSearch}
  >
    <svelte:fragment slot="cell" let:item let:column>
      {#if column.key === 'status'}
        <StatusBadge status={item.status} />
      {:else if column.render}
        {@html column.render(item)}
      {:else}
        {item[column.key] || '-'}
      {/if}
    </svelte:fragment>
  </DataTable>
</PageLayout>