<script>
  import { onMount } from 'svelte'
  import { GetInvoices, GenerateInvoiceHTML, ViewInvoiceHTML, SaveInvoiceHTML } from '../wailsjs/go/main/App.js'
  import InvoiceModal from './InvoiceModal.svelte'
  
  let invoices = []
  let isLoading = false
  let showInvoiceModal = false
  let searchTerm = ''
  
  onMount(async () => {
    await loadInvoices()
  })
  
  async function loadInvoices() {
    isLoading = true
    try {
      console.log('Loading invoices...')
      const result = await GetInvoices()
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
    showInvoiceModal = true
  }
  
  function closeInvoiceModal() {
    showInvoiceModal = false
  }
  
  function handleInvoiceSaved() {
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
      default: return 'badge-neutral'
    }
  }

  async function saveOrPrintInvoice(invoiceId) {
    console.log('saveOrPrintInvoice called with ID:', invoiceId)
    try {
      console.log('Opening print dialog for invoice:', invoiceId)
      await SaveInvoiceHTML(invoiceId)
      console.log('Print dialog opened')
    } catch (error) {
      console.error('Error opening print dialog:', error)
      alert('Failed to open print dialog. Please try again.')
    }
  }

  async function viewInvoiceHTML(invoiceId) {
    console.log('viewInvoiceHTML called with ID:', invoiceId)
    try {
      console.log('Opening HTML preview for invoice:', invoiceId)
      await ViewInvoiceHTML(invoiceId)
      console.log('HTML preview opened')
    } catch (error) {
      console.error('Error viewing HTML:', error)
      alert('Failed to open HTML preview. Please try again.')
    }
  }

  function editInvoice(invoice) {
    console.log('editInvoice called with invoice:', invoice)
    alert('Edit button clicked! Invoice ID: ' + invoice.id)
    // TODO: Implement edit functionality
    // This would typically open the invoice modal with the invoice data pre-filled
    console.log('Edit invoice:', invoice)
  }
  
  $: filteredInvoices = (invoices || []).filter(invoice => 
    invoice.invoice_number?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (invoice.customer && invoice.customer.name?.toLowerCase().includes(searchTerm.toLowerCase()))
  )
</script>

<div class="p-6">
  <div class="glass-card">
    <div class="card-body">
      <h2 class="text-2xl font-semibold text-white mb-4 flex items-center gap-3">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
        </svg>
        Sales Invoices
      </h2>
      
      <div class="flex flex-col lg:flex-row justify-between items-start lg:items-center gap-4 mb-6">
        <div class="flex gap-2">
          <button 
            class="btn-primary"
            on:click={openInvoiceModal}
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
            New Invoice
          </button>
          <button class="btn-secondary">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"></path>
            </svg>
            Filter
          </button>
        </div>
        <div class="form-group w-full lg:w-auto">
          <input 
            type="text" 
            placeholder="Search invoices..." 
            class="input-glass w-full lg:w-64"
            bind:value={searchTerm}
          />
        </div>
      </div>

      {#if isLoading}
        <div class="flex items-center justify-center py-12">
          <div class="loading loading-spinner loading-lg text-primary"></div>
        </div>
      {:else if filteredInvoices.length === 0}
        <div class="text-center py-12">
          <div class="flex flex-col items-center gap-4">
            <svg class="w-16 h-16 text-white/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
            </svg>
            <div>
              <p class="text-lg text-white/60 mb-2">
                {searchTerm ? 'No invoices found matching your search' : 'No invoices found'}
              </p>
              <p class="text-white/40 mb-4">
                {searchTerm ? 'Try adjusting your search terms' : 'Create your first invoice to get started'}
              </p>
              {#if !searchTerm}
                <button 
                  class="btn-primary btn-sm"
                  on:click={openInvoiceModal}
                >
                  Create your first invoice
                </button>
              {/if}
            </div>
          </div>
        </div>
      {:else}
        <div class="overflow-x-auto">
          <table class="table-glass w-full">
            <thead>
              <tr>
                <th>Invoice #</th>
                <th>Customer</th>
                <th>Issue Date</th>
                <th>Due Date</th>
                <th>Amount</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {#each filteredInvoices as invoice}
                <tr>
                  <td class="font-medium">{invoice.invoice_number}</td>
                  <td>
                    {#if invoice.customer}
                      <div>
                        <div class="font-medium">{invoice.customer.name}</div>
                        {#if invoice.customer.email}
                          <div class="text-sm text-white/60">{invoice.customer.email}</div>
                        {/if}
                      </div>
                    {:else}
                      <span class="text-white/40">No customer</span>
                    {/if}
                  </td>
                  <td>{formatDate(invoice.issue_date)}</td>
                  <td>{formatDate(invoice.due_date)}</td>
                  <td class="font-medium">{formatCurrency(invoice.total_amount)}</td>
                  <td>
                    <span class="badge {getStatusBadgeClass(invoice.status)} capitalize">
                      {invoice.status}
                    </span>
                  </td>
                  <td>
                    <div class="flex gap-2">
                      <button 
                        class="btn-icon btn-sm text-info hover:bg-info/20" 
                        title="View HTML"
                        on:click={() => viewInvoiceHTML(invoice.id)}
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                        </svg>
                      </button>
                      <button 
                        class="btn-icon btn-sm text-warning hover:bg-warning/20" 
                        title="Edit"
                        on:click={() => editInvoice(invoice)}
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                        </svg>
                      </button>
                      <button 
                        class="btn-icon btn-sm text-success hover:bg-success/20" 
                        title="Save/Print"
                        on:click={() => saveOrPrintInvoice(invoice.id)}
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                        </svg>
                      </button>
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>
  </div>
</div>

<!-- Invoice Modal -->
<InvoiceModal 
  bind:isOpen={showInvoiceModal}
  on:close={closeInvoiceModal}
  on:saved={handleInvoiceSaved}
/>