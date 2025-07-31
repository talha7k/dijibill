<script>
  import { onMount } from 'svelte'
  import POSCart from './components/POSCart.svelte'
  import POSProductGrid from './components/POSProductGrid.svelte'
  import POSModals from './components/POSModals.svelte'
  import { 
    loadInitialData,
    selectCategory,
    saveSale,
    refundSale,
    transferItems
  } from './services/posService.js'

  // Load initial data when component mounts
  onMount(() => {
    loadInitialData()
  })

  // Event handlers
  function handleSelectCategory(category) {
    selectCategory(category)
  }

  function handleSaveSale() {
    saveSale()
  }

  function handleRefund() {
    refundSale()
  }

  function handleTransfer() {
    transferItems()
  }
</script>

<div class="h-screen bg-gray-900 text-white flex">
  <!-- Left Panel - Current Sale -->
  <POSCart 
    onSaveSale={handleSaveSale}
    onRefund={handleRefund}
    onTransfer={handleTransfer}
  />

  <!-- Right Panel - Products -->
  <POSProductGrid 
    onSelectCategory={handleSelectCategory}
  />
</div>

<!-- All Modals -->
<POSModals 
  onTransferItems={handleTransfer}
  onRefundSale={handleRefund}
/>