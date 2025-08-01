<script>
  import { onMount } from 'svelte'
  import { 
    GetSalesInvoices, 
    GetSalesInvoiceByID,
    GenerateInvoiceHTML, 
    ViewInvoiceHTML, 
    SaveInvoiceHTML,
    ViewInvoiceHTMLEnglish,
    ViewInvoiceHTMLArabic,
    ViewInvoiceHTMLBilingual,
    SaveInvoiceHTMLEnglish,
    SaveInvoiceHTMLArabic,
    SaveInvoiceHTMLBilingual
  } from '../wailsjs/go/main/App.js'
  import { refundInvoice } from './services/posService.js'
  import SalesInvoiceModal from './SalesInvoiceModal.svelte'
  import PageLayout from './components/PageLayout.svelte'
  import DataTable from './components/DataTable.svelte'
  import StatusBadge from './components/StatusBadge.svelte'
  
  let invoices = []
  let isLoading = false
  let showInvoiceModal = false
  let editingInvoice = null
  let searchTerm = ''
  
  // Get invoice language preference from localStorage (set by system settings)
  function getInvoiceLanguagePreference() {
    const systemSettings = localStorage.getItem('systemSettings')
    if (systemSettings) {
      try {
        const settings = JSON.parse(systemSettings)
        return settings.invoice_language || 'english'
      } catch (error) {
        console.error('Error parsing system settings:', error)
      }
    }
    return 'english' // default to English
  }
  
  onMount(async () => {
    await loadInvoices()
  })
  
  async function loadInvoices() {
    isLoading = true
    try {
      console.log('Loading invoices...')
      const result = await GetSalesInvoices()
      console.log('Invoices loaded:', result)
      invoices = result || []
      console.log('Total invoices:', invoices.length)
    } catch (error) {
      console.error('Error loading invoices:', error)
      invoices = []
    } finally {
      isLoading = false
    }
  }
  
  function openInvoiceModal() {
    editingInvoice = null
    showInvoiceModal = true
  }
  
  function closeInvoiceModal() {
    editingInvoice = null
    showInvoiceModal = false
  }
  
  function handleInvoiceSaved() {
    console.log('Invoice saved event received, refreshing table...')
    editingInvoice = null
    showInvoiceModal = false
    loadInvoices() // Refresh the list
  }
  
  function formatDate(dateString) {
    if (!dateString) return 'N/A'
    try {
      return new Date(dateString).toLocaleDateString()
    } catch (error) {
      console.error('Error formatting date:', error)
      return 'Invalid Date'
    }
  }
  
  function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'SAR'
    }).format(amount)
  }
  
  function getStatusBadgeClass(status) {
    switch (status) {
      case 'paid': return 'badge-success'
      case 'sent': return 'badge-info'
      case 'draft': return 'badge-warning'
      case 'cancelled': return 'badge-error'
      case 'refunded': return 'badge-secondary'
      case 'partially_paid': return 'badge-warning'
      default: return 'badge-neutral'
    }
  }

  async function saveOrPrintInvoice(invoiceId) {
    console.log('saveOrPrintInvoice called with ID:', invoiceId)
    const language = getInvoiceLanguagePreference()
    
    try {
      if (language === 'arabic') {
        await SaveInvoiceHTMLArabic(invoiceId)
      } else if (language === 'bilingual') {
        await SaveInvoiceHTMLBilingual(invoiceId)
      } else {
        await SaveInvoiceHTMLEnglish(invoiceId)
      }
      console.log(`Invoice ${invoiceId} saved in ${language}`)
    } catch (error) {
      console.error('Error saving invoice:', error)
      alert('Error saving invoice: ' + error.message)
    }
  }

  async function viewInvoiceHTML(invoiceId) {
    console.log('viewInvoiceHTML called with ID:', invoiceId)
    const language = getInvoiceLanguagePreference()
    
    try {
      if (language === 'arabic') {
        await ViewInvoiceHTMLArabic(invoiceId)
      } else if (language === 'bilingual') {
        await ViewInvoiceHTMLBilingual(invoiceId)
      } else {
        await ViewInvoiceHTMLEnglish(invoiceId)
      }
      console.log(`Invoice ${invoiceId} viewed in ${language}`)
    } catch (error) {
      console.error('Error viewing invoice:', error)
      alert('Error viewing invoice: ' + error.message)
    }
  }

  async function editInvoice(invoice) {
    console.log('editInvoice called with invoice:', invoice)
    try {
      // Fetch complete invoice data including items and customer
      const completeInvoice = await GetSalesInvoiceByID(invoice.id)
      console.log('Complete invoice data loaded:', completeInvoice)
      editingInvoice = completeInvoice
      showInvoiceModal = true
    } catch (error) {
      console.error('Error loading invoice for editing:', error)
      alert('Error loading invoice: ' + error.message)
    }
  }

  async function handleRefundInvoice(invoice) {
    console.log('handleRefundInvoice called with invoice:', invoice)
    
    // Confirm refund action
    const confirmed = confirm(`Are you sure you want to refund invoice ${invoice.invoice_number}? This will mark the invoice and all its payments as refunded.`)
    if (!confirmed) return
    
    // Get refund reason
    const reason = prompt('Please enter a reason for the refund (optional):') || 'Refund requested'
    
    try {
      await refundInvoice(invoice.id, reason)
      alert(`Invoice ${invoice.invoice_number} has been successfully refunded.`)
      await loadInvoices() // Refresh the list
    } catch (error) {
      console.error('Error refunding invoice:', error)
      alert('Error refunding invoice: ' + error.message)
    }
  }
  
  $: filteredInvoices = (invoices || []).filter(invoice => 
    invoice.invoice_number?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (invoice.customer && invoice.customer.name?.toLowerCase().includes(searchTerm.toLowerCase()))
  )

  // Table configuration
  const columns = [
    { key: 'invoice_number', label: 'Invoice #', render: (item) => `<span class="font-medium">${item.invoice_number}</span>` },
    { 
      key: 'customer', 
      label: 'Customer', 
      render: (item) => {
        if (item.customer) {
          return `<div><div class="font-medium">${item.customer.name}</div>${item.customer.email ? `<div class="text-sm opacity-70">${item.customer.email}</div>` : ''}</div>`
        }
        return '<span class="opacity-40">No customer</span>'
      }
    },
    { key: 'issue_date', label: 'Issue Date', render: (item) => formatDate(item.issue_date) },
    { key: 'due_date', label: 'Due Date', render: (item) => formatDate(item.due_date) },
    { key: 'total_amount', label: 'Amount', render: (item) => `<span class="font-medium">${formatCurrency(item.total_amount)}</span>` },
    { key: 'status', label: 'Status' },
    { 
      label: 'Actions',
      actions: [
        { key: 'view', text: 'View', icon: 'view', class: 'btn-info', title: 'View HTML' },
        { key: 'edit', text: 'Edit', icon: 'edit', class: 'btn-warning', title: 'Edit Invoice' },
        { key: 'save', text: 'Save', icon: 'download', class: 'btn-success', title: 'Save/Print' },
        { key: 'refund', text: 'Refund', icon: 'undo', class: 'btn-error', title: 'Refund Invoice' }
      ]
    }
  ]

  const primaryAction = {
    text: 'New Invoice',
    icon: 'add'
  }

  const secondaryActions = [
    {
      text: 'Filter',
      icon: 'filter'
    }
  ]

  function handlePrimaryAction() {
    openInvoiceModal()
  }

  function handleSecondaryAction(action) {
    if (action.text === 'Filter') {
      console.log('Filter clicked')
    }
  }

  function handleRowAction(event) {
    const { action, item } = event.detail
    if (action === 'view') {
      viewInvoiceHTML(item.id)
    } else if (action === 'edit') {
      editInvoice(item)
    } else if (action === 'save') {
      saveOrPrintInvoice(item.id)
    } else if (action === 'refund') {
      handleRefundInvoice(item)
    }
  }

  function handleSearch(event) {
    searchTerm = event.detail.searchTerm
  }
</script>

<PageLayout 
  title="Sales Invoices" 
  icon="invoice" 
  showIndicator={true}
>
  <svelte:fragment slot="actions">
    <!-- Actions are handled by DataTable -->
  </svelte:fragment>

  <DataTable
    data={filteredInvoices}
    {columns}
    loading={isLoading}
    {searchTerm}
    searchPlaceholder="Search invoices..."
    emptyStateTitle="No invoices found"
    emptyStateMessage={searchTerm ? 'Try adjusting your search terms' : 'Create your first invoice to get started'}
    emptyStateIcon="invoice"
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

<!-- Invoice Modal -->
<SalesInvoiceModal
  bind:isOpen={showInvoiceModal}
  {editingInvoice}
  on:close={closeInvoiceModal}
  on:saved={handleInvoiceSaved}
/>