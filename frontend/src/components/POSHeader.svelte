<script>
  import { 
    showTableModal, 
    showSalesCategoryModal, 
    showCustomerModal,
    showInvoiceSelectionModal,
    currentSale,
    selectedSalesCategory,
    customers,
    getCustomerName
  } from '../stores/posStore.js'

  export let onTableClick = () => {}
  export let onSalesCategoryClick = () => {}
  export let onCustomerClick = () => {}

  $: customerName = getCustomerName($currentSale.customer_id, $customers)
</script>

<div class="p-4 border-b border-gray-700">
  <div class="flex justify-between items-center mb-4">
    <h2 class="text-xl font-bold">Current Sale</h2>
    <div class="flex space-x-2">
      <!-- Select Invoice -->
      <button
        class="p-2 bg-orange-600 hover:bg-orange-700 rounded-lg transition-colors"
        on:click={() => $showInvoiceSelectionModal = true}
        title="Select Previous Invoice"
      >
        <i class="fas fa-file-invoice w-5 h-5"></i>
      </button>
      
      <!-- Table Selection -->
      <button
        class="p-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
        on:click={() => { $showTableModal = true; onTableClick(); }}
        title="Select Table"
      >
        <i class="fas fa-table w-5 h-5"></i>
      </button>
      
      <!-- Sales Category Selection -->
      <button
        class="p-2 bg-purple-600 hover:bg-purple-700 rounded-lg transition-colors"
        on:click={() => { $showSalesCategoryModal = true; onSalesCategoryClick(); }}
        title="Select Sales Category"
      >
        <i class="fas fa-tag w-5 h-5"></i>
      </button>
      
      <!-- Customer Selection -->
      <button
        class="p-2 bg-green-600 hover:bg-green-700 rounded-lg transition-colors"
        on:click={() => { $showCustomerModal = true; onCustomerClick(); }}
        title="Select Customer"
      >
        <i class="fas fa-user w-5 h-5"></i>
      </button>
    </div>
  </div>

  <!-- Customer, Table, and Sales Category Info -->
  <div class="text-sm text-gray-300">
    <div>Customer: {customerName}</div>
    {#if $currentSale.table_number}
      <div>Table: {$currentSale.table_number}</div>
    {/if}
    <div>Category: {$selectedSalesCategory?.name || 'None'}</div>
  </div>
</div>