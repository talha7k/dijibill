<script>
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'
  import InvoiceSelectionModal from './InvoiceSelectionModal.svelte'
  import {
    showCustomerModal,
    showTableModal,
    showSalesCategoryModal,
    showTransferModal,
    showRefundModal,
    showPaymentModal,
    showInvoiceSelectionModal,
    customers,
    salesCategories,
    selectedSalesCategory,
    openInvoices,
    selectedTransferInvoice,
    currentSale,
    paymentTypes,
    paymentItems,
    remainingAmount,
    refundReason,
    loading
  } from '../stores/posStore.js'
  import {
    loadOpenInvoices,
    addPaymentItem,
    removePaymentItem,
    processPayments
  } from '../services/posService.js'

  export let onTransferItems = () => {}
  export let onRefundSale = () => {}

  // Payment modal state
  let selectedPaymentType = ''
  let paymentAmount = ''

  // Load open invoices when transfer modal opens
  $: if ($showTransferModal) {
    loadOpenInvoices()
  }

  function selectTransferInvoice(invoice) {
    $selectedTransferInvoice = invoice
  }

  function addPayment() {
    if (!selectedPaymentType || !paymentAmount) {
      alert('Please select payment type and enter amount')
      return
    }

    const amount = parseFloat(paymentAmount)
    if (isNaN(amount) || amount <= 0) {
      alert('Please enter a valid amount')
      return
    }

    addPaymentItem(parseInt(selectedPaymentType), amount)
    selectedPaymentType = ''
    paymentAmount = ''
  }

  function setFullAmount() {
    paymentAmount = $remainingAmount.toString()
  }

  function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount)
  }
</script>

<!-- Customer Selection Modal -->
<Modal
  show={$showCustomerModal}
  title="Select Customer"
  size="md"
  on:close={() => $showCustomerModal = false}
  on:secondary={() => $showCustomerModal = false}
>
  <div class="space-y-4">
    <FormField
      label="Customer"
      type="select"
      bind:value={$currentSale.customer_id}
      options={$customers.map(c => ({ value: c.id, label: c.name }))}
      placeholder="Select customer"
    />
    <div class="flex justify-end space-x-2">
      <button
        class="btn btn-secondary"
        on:click={() => $showCustomerModal = false}
      >
        Cancel
      </button>
      <button
        class="btn btn-primary"
        on:click={() => $showCustomerModal = false}
      >
        Select
      </button>
    </div>
  </div>
</Modal>

<!-- Table Selection Modal -->
<Modal
  show={$showTableModal}
  title="Select Table"
  size="md"
  on:close={() => $showTableModal = false}
  on:secondary={() => $showTableModal = false}
>
  <div class="space-y-4">
    <FormField
      label="Table Number"
      type="text"
      bind:value={$currentSale.table_number}
      placeholder="Enter table number"
    />
    <div class="flex justify-end space-x-2">
      <button
        class="btn btn-secondary"
        on:click={() => $showTableModal = false}
      >
        Cancel
      </button>
      <button
        class="btn btn-primary"
        on:click={() => $showTableModal = false}
      >
        Select
      </button>
    </div>
  </div>
</Modal>

<!-- Sales Category Selection Modal -->
<Modal
  show={$showSalesCategoryModal}
  title="Select Sales Category"
  size="md"
  on:close={() => $showSalesCategoryModal = false}
  on:secondary={() => $showSalesCategoryModal = false}
>
  <div class="space-y-4">
    <FormField
      label="Sales Category"
      type="select"
      bind:value={$selectedSalesCategory}
      options={$salesCategories.map(c => ({ value: c, label: c.name }))}
      placeholder="Select sales category"
    />
    <div class="flex justify-end space-x-2">
      <button
        class="btn btn-secondary"
        on:click={() => $showSalesCategoryModal = false}
      >
        Cancel
      </button>
      <button
        class="btn btn-primary"
        on:click={() => $showSalesCategoryModal = false}
      >
        Select
      </button>
    </div>
  </div>
</Modal>

<!-- Transfer Modal with Invoice Grid -->
<Modal
  show={$showTransferModal}
  title="Transfer Items to Invoice"
  size="lg"
  on:close={() => {
    $showTransferModal = false
    $selectedTransferInvoice = null
  }}
  on:secondary={() => {
    $showTransferModal = false
    $selectedTransferInvoice = null
  }}
>
  <div class="space-y-4">
    <div class="text-sm text-gray-600 mb-4">
      Select an open invoice to transfer items to:
    </div>
    
    {#if $openInvoices.length === 0}
      <div class="text-center py-8 text-gray-500">
        No open invoices available for transfer
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 max-h-96 overflow-y-auto">
        {#each $openInvoices as invoice}
          <div 
            class="border rounded-lg p-4 cursor-pointer transition-all hover:shadow-md {$selectedTransferInvoice?.id === invoice.id ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-gray-300'}"
            on:click={() => selectTransferInvoice(invoice)}
          >
            <div class="font-medium text-gray-900">
              {invoice.invoice_number || `Invoice #${invoice.id}`}
            </div>
            <div class="text-sm text-gray-600 mt-1">
              {invoice.customer?.name || 'Walk-in Customer'}
            </div>
            <div class="text-sm text-gray-500 mt-1">
              Table: {invoice.table_number || 'N/A'}
            </div>
            <div class="text-sm font-medium text-green-600 mt-2">
              {formatCurrency(invoice.total_amount)}
            </div>
            <div class="text-xs text-gray-400 mt-1">
              Status: {invoice.status}
            </div>
          </div>
        {/each}
      </div>
    {/if}
    
    <div class="flex justify-end space-x-2 pt-4 border-t">
      <button
        class="btn btn-secondary"
        on:click={() => {
          $showTransferModal = false
          $selectedTransferInvoice = null
        }}
      >
        Cancel
      </button>
      <button
        class="btn btn-primary"
        on:click={onTransferItems}
        disabled={!$selectedTransferInvoice}
      >
        Transfer to Selected Invoice
      </button>
    </div>
  </div>
</Modal>

<!-- Payment Modal -->
<Modal
  show={$showPaymentModal}
  title="Process Payment"
  size="lg"
  on:close={() => $showPaymentModal = false}
  on:secondary={() => $showPaymentModal = false}
>
  <div class="space-y-6">
    <!-- Payment Summary -->
    <div class="p-4 rounded-lg">
      <div class="flex justify-between items-center mb-2">
        <span class="font-medium">Total Amount:</span>
        <span class="text-lg font-bold">{formatCurrency($currentSale.total)}</span>
      </div>
      <div class="flex justify-between items-center mb-2">
        <span class="font-medium">Paid Amount:</span>
        <span class="text-green-600">{formatCurrency($currentSale.total - $remainingAmount)}</span>
      </div>
      <div class="flex justify-between items-center">
        <span class="font-medium">Remaining:</span>
        <span class="text-red-600 font-bold">{formatCurrency($remainingAmount)}</span>
      </div>
    </div>

    <!-- Add Payment -->
    <div class="border rounded-lg p-4">
      <h4 class="font-medium mb-4">Add Payment</h4>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <FormField
          label="Payment Type"
          type="select"
          bind:value={selectedPaymentType}
          options={$paymentTypes.map(pt => ({ value: pt.id, label: pt.name }))}
          placeholder="Select payment type"
        />
        <div class="space-y-2">
          <FormField
            label="Amount"
            type="number"
            bind:value={paymentAmount}
            placeholder="0.00"
            step="0.01"
            min="0"
            max={$remainingAmount}
          />
          <button
            class="btn btn-sm btn-outline w-full"
            on:click={setFullAmount}
            disabled={$remainingAmount <= 0}
          >
            Full Amount
          </button>
        </div>
        <div class="flex items-end">
          <button
            class="btn btn-primary w-full"
            on:click={addPayment}
            disabled={!selectedPaymentType || !paymentAmount || $remainingAmount <= 0}
          >
            Add Payment
          </button>
        </div>
      </div>
    </div>

    <!-- Payment Items -->
    {#if $paymentItems.length > 0}
      <div class="border rounded-lg p-4">
        <h4 class="font-medium mb-4">Payment Breakdown</h4>
        <div class="space-y-2">
          {#each $paymentItems as item}
            <div class="flex justify-between items-center py-2 border-b last:border-b-0">
              <div>
                <span class="font-medium">{item.payment_type_name}</span>
              </div>
              <div class="flex items-center space-x-2">
                <span class="font-medium">{formatCurrency(item.amount)}</span>
                <button
                  class="btn btn-sm btn-danger"
                  on:click={() => removePaymentItem(item.id)}
                >
                  Remove
                </button>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <div class="flex justify-end space-x-2 pt-4 border-t">
      <button
        class="btn btn-secondary"
        on:click={() => $showPaymentModal = false}
      >
        Cancel
      </button>
      <button
        class="btn btn-success"
        on:click={processPayments}
        disabled={$paymentItems.length === 0}
      >
        Complete Payment
      </button>
    </div>
  </div>
</Modal>

<!-- Refund Modal -->
<Modal
  show={$showRefundModal}
  title="Refund Sale"
  size="md"
  on:close={() => $showRefundModal = false}
  on:secondary={() => $showRefundModal = false}
>
  <div class="space-y-4">
    <FormField
      label="Refund Reason"
      type="text"
      bind:value={$refundReason}
      placeholder="Enter refund reason"
    />
    <div class="flex justify-end space-x-2">
      <button
        class="btn btn-secondary"
        on:click={() => $showRefundModal = false}
      >
        Cancel
      </button>
      <button
        class="btn btn-danger"
        on:click={onRefundSale}
        disabled={!$refundReason}
      >
        Process Refund
      </button>
    </div>
  </div>
</Modal>

<!-- Invoice Selection Modal -->
<InvoiceSelectionModal />