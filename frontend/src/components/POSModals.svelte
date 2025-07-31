<script>
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'
  import { 
    showCustomerModal,
    showTableModal,
    showSalesCategoryModal,
    showTransferModal,
    showRefundModal,
    currentSale,
    customers,
    salesCategories,
    selectedSalesCategory,
    transferToInvoiceId,
    refundReason
  } from '../stores/posStore.js'

  export let onTransferItems = () => {}
  export let onRefundSale = () => {}
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

<!-- Transfer Modal -->
<Modal
  show={$showTransferModal}
  title="Transfer Items"
  size="md"
  on:close={() => $showTransferModal = false}
  on:secondary={() => $showTransferModal = false}
>
  <div class="space-y-4">
    <FormField
      label="Transfer to Invoice ID"
      type="text"
      bind:value={$transferToInvoiceId}
      placeholder="Enter invoice ID"
    />
    <div class="flex justify-end space-x-2">
      <button
        class="btn btn-secondary"
        on:click={() => $showTransferModal = false}
      >
        Cancel
      </button>
      <button
        class="btn btn-primary"
        on:click={onTransferItems}
        disabled={!$transferToInvoiceId}
      >
        Transfer
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