<script>
  import { 
    showTableModal, 
    showSalesCategoryModal, 
    showCustomerModal,
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
      <!-- Table Selection -->
      <button
        class="p-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
        on:click={() => { $showTableModal = true; onTableClick(); }}
        title="Select Table"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M3 14h18m-9-4v8m-7 0h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
        </svg>
      </button>
      
      <!-- Sales Category Selection -->
      <button
        class="p-2 bg-purple-600 hover:bg-purple-700 rounded-lg transition-colors"
        on:click={() => { $showSalesCategoryModal = true; onSalesCategoryClick(); }}
        title="Select Sales Category"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"></path>
        </svg>
      </button>
      
      <!-- Customer Selection -->
      <button
        class="p-2 bg-green-600 hover:bg-green-700 rounded-lg transition-colors"
        on:click={() => { $showCustomerModal = true; onCustomerClick(); }}
        title="Select Customer"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
        </svg>
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