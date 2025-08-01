<script>
  import { onMount } from 'svelte'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'
  import StatusBadge from './StatusBadge.svelte'
  import { showDbSuccess, showDbError } from '../stores/notificationStore.js'
  import {
    showInvoiceSelectionModal,
    openInvoices,
    pastInvoices,
    selectedInvoice,
    invoiceSelectionTab,
    invoiceDateFilter,
    loading
  } from '../stores/posStore.js'
  import {
    loadOpenInvoices,
    loadPastInvoices,
    filterPastInvoicesByDate,
    selectInvoice,
    refundInvoice,
    reprintReceipt
  } from '../services/posService.js'

  export let onClose = () => {}
  export let onInvoiceSelected = () => {}

  let filteredPastInvoices = []
  let refundReason = ''
  let showRefundConfirm = false
  let invoiceToRefund = null

  // Load data when modal opens
  $: if ($showInvoiceSelectionModal) {
    loadData()
  }

  // Filter past invoices when date filter changes
  $: if ($pastInvoices && $invoiceDateFilter) {
    filteredPastInvoices = filterPastInvoicesByDate(
      $invoiceDateFilter.start_date,
      $invoiceDateFilter.end_date
    )
  }

  async function loadData() {
    await loadOpenInvoices()
    await loadPastInvoices()
  }

  function handleTabChange(tab) {
    invoiceSelectionTab.set(tab)
    selectedInvoice.set(null)
  }

  function handleInvoiceSelect(invoice) {
    selectInvoice(invoice)
  }

  function handleDateFilterChange() {
    filteredPastInvoices = filterPastInvoicesByDate(
      $invoiceDateFilter.start_date,
      $invoiceDateFilter.end_date
    )
  }

  function formatCurrency(amount) {
    return new Intl.NumberFormat('en-SA', {
      style: 'currency',
      currency: 'SAR'
    }).format(amount || 0)
  }

  function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString('en-SA', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  async function handleRefund(invoice) {
    invoiceToRefund = invoice
    showRefundConfirm = true
  }

  async function confirmRefund() {
    if (!invoiceToRefund || !refundReason.trim()) {
      showDbError('Please enter a refund reason')
      return
    }

    await refundInvoice(invoiceToRefund.id, refundReason)
    
    // Refresh data
    await loadData()
    
    // Reset refund state
    showRefundConfirm = false
    invoiceToRefund = null
    refundReason = ''
    selectedInvoice.set(null)
  }

  function cancelRefund() {
    showRefundConfirm = false
    invoiceToRefund = null
    refundReason = ''
  }

  async function handleReprint(invoice) {
    await reprintReceipt(invoice.id)
  }

  function handleClose() {
    selectedInvoice.set(null)
    showInvoiceSelectionModal.set(false)
    onClose()
  }
</script>

<Modal
  show={$showInvoiceSelectionModal}
  title="Select Invoice"
  size="xl"
  on:close={handleClose}
>
  <div class="space-y-6">
    <!-- Tab Navigation -->
    <div class="flex space-x-1 p-1 rounded-lg">
      <button
        class="flex-1 py-2 px-4 rounded-md text-sm font-medium transition-colors {$invoiceSelectionTab === 'open' ? 'bg-white/20 text-blue-300 shadow-sm' : 'text-gray-300 hover:text-gray-200'}"
        on:click={() => handleTabChange('open')}
      >
        Open Invoices ({$openInvoices.length})
      </button>
      <button
        class="flex-1 py-2 px-4 rounded-md text-sm font-medium transition-colors {$invoiceSelectionTab === 'past' ? 'bg-white/20 text-blue-300 shadow-sm' : 'text-gray-300 hover:text-gray-200'}"
        on:click={() => handleTabChange('past')}
      >
        Past Invoices ({$pastInvoices.length})
      </button>
    </div>

    <!-- Date Filter for Past Invoices -->
    {#if $invoiceSelectionTab === 'past'}
      <div class=" p-4 rounded-lg">
        <h4 class="font-medium mb-3">Date Filter</h4>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField
            label="Start Date"
            type="date"
            bind:value={$invoiceDateFilter.start_date}
            on:change={handleDateFilterChange}
          />
          <FormField
            label="End Date"
            type="date"
            bind:value={$invoiceDateFilter.end_date}
            on:change={handleDateFilterChange}
          />
        </div>
        <div class="text-sm text-gray-300 mt-2">
          Showing {filteredPastInvoices.length} invoices
        </div>
      </div>
    {/if}

    <!-- Invoice List -->
    <div class="max-h-96 overflow-y-auto">
      {#if $loading}
        <div class="text-center py-8">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-300"></div>
          <p class="mt-2 text-gray-300">Loading invoices...</p>
        </div>
      {:else}
        {#if $invoiceSelectionTab === 'open'}
          {#if $openInvoices.length === 0}
            <div class="text-center py-8 text-gray-300">
              No open invoices available
            </div>
          {:else}
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {#each $openInvoices as invoice}
                <div 
                  class="border rounded-lg p-4 cursor-pointer transition-all hover:shadow-md {$selectedInvoice?.id === invoice.id ? 'border-blue-300 bg-blue-50' : 'border-gray-200 hover:border-gray-300'}"
                  on:click={() => handleInvoiceSelect(invoice)}
                >
                  <div class="flex justify-between items-start mb-2">
                    <div class="font-medium text-gray-200">
                      {invoice.invoice_number || `Invoice #${invoice.id}`}
                    </div>
                    <StatusBadge status={invoice.status} />
                  </div>
                  <div class="text-sm text-gray-300 mb-1">
                    {invoice.customer?.name || 'Walk-in Customer'}
                  </div>
                  <div class="text-sm text-gray-300 mb-2">
                    Table: {invoice.table_number || 'N/A'}
                  </div>
                  <div class="text-sm font-medium text-green-600 mb-1">
                    {formatCurrency(invoice.total_amount)}
                  </div>
                  <div class="text-xs text-gray-400">
                    {formatDate(invoice.created_at)}
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        {:else}
          {#if filteredPastInvoices.length === 0}
            <div class="text-center py-8 text-gray-300">
              No past invoices found for the selected date range
            </div>
          {:else}
            <div class="space-y-3">
              {#each filteredPastInvoices as invoice}
                <div 
                  class="border rounded-lg p-4 cursor-pointer transition-all hover:shadow-md {$selectedInvoice?.id === invoice.id ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-gray-300'}"
                  on:click={() => handleInvoiceSelect(invoice)}
                >
                  <div class="flex justify-between items-start">
                    <div class="flex-1">
                      <div class="flex justify-between items-center mb-2">
                        <div class="font-medium text-gray-200">
                          {invoice.invoice_number || `Invoice #${invoice.id}`}
                        </div>
                        <StatusBadge status={invoice.status} />
                      </div>
                      <div class="grid grid-cols-2 gap-4 text-sm">
                        <div>
                          <div class="text-gray-300">Customer:</div>
                          <div class="font-medium">{invoice.customer?.name || 'Walk-in Customer'}</div>
                        </div>
                        <div>
                          <div class="text-gray-300">Table:</div>
                          <div class="font-medium">{invoice.table_number || 'N/A'}</div>
                        </div>
                        <div>
                          <div class="text-gray-300">Amount:</div>
                          <div class="font-medium text-green-600">{formatCurrency(invoice.total_amount)}</div>
                        </div>
                        <div>
                          <div class="text-gray-300">Date:</div>
                          <div class="font-medium">{formatDate(invoice.created_at)}</div>
                        </div>
                      </div>
                    </div>
                    <div class="flex flex-col space-y-2 ml-4">
                      <button
                        class="btn btn-sm btn-primary"
                        on:click|stopPropagation={() => handleReprint(invoice)}
                        title="Reprint Receipt"
                      >
                        <i class="fas fa-print"></i>
                      </button>
                      {#if invoice.status !== 'refunded'}
                        <button
                          class="btn btn-sm btn-warning"
                          on:click|stopPropagation={() => handleRefund(invoice)}
                          title="Refund Invoice"
                        >
                          <i class="fas fa-undo"></i>
                        </button>
                      {/if}
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        {/if}
      {/if}
    </div>

    <!-- Action Buttons -->
    <div class="flex justify-end space-x-2 pt-4 border-t">
      <button
        class="btn btn-secondary"
        on:click={handleClose}
      >
        Close
      </button>
      {#if $selectedInvoice && $invoiceSelectionTab === 'open'}
        <button
          class="btn btn-primary"
          on:click={() => {
            onInvoiceSelected()
            handleClose()
          }}
        >
          Select Invoice
        </button>
      {/if}
    </div>
  </div>
</Modal>

<!-- Refund Confirmation Modal -->
<Modal
  show={showRefundConfirm}
  title="Confirm Refund"
  size="md"
  on:close={cancelRefund}
>
  <div class="space-y-4">
    {#if invoiceToRefund}
      <div class=" p-4 rounded-lg">
        <div class="font-medium">Invoice: {invoiceToRefund.invoice_number || `#${invoiceToRefund.id}`}</div>
        <div class="text-sm text-gray-300">Amount: {formatCurrency(invoiceToRefund.total_amount)}</div>
        <div class="text-sm text-gray-300">Customer: {invoiceToRefund.customer?.name || 'Walk-in Customer'}</div>
      </div>
    {/if}
    
    <FormField
      label="Refund Reason"
      type="text"
      bind:value={refundReason}
      placeholder="Enter reason for refund"
      required
    />
    
    <div class="flex justify-end space-x-2">
      <button
        class="btn btn-secondary"
        on:click={cancelRefund}
      >
        Cancel
      </button>
      <button
        class="btn btn-danger"
        on:click={confirmRefund}
        disabled={!refundReason.trim()}
      >
        Confirm Refund
      </button>
    </div>
  </div>
</Modal>