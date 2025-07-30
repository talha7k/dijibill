<script>
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import HorizontalTabs from './HorizontalTabs.svelte'
  import ProductDetailsTab from './ProductDetailsTab.svelte'
  import ProductPriceTab from './ProductPriceTab.svelte'
  import ProductDescriptionTab from './ProductDescriptionTab.svelte'
  import ProductImageTab from './ProductImageTab.svelte'

  export let show = false
  export let editingProduct = null
  export let productForm = {}
  export let categories = []
  export let taxRates = []
  export let units = []
  export let isLoading = false
  export let generateSKU
  export let generateBarcode

  let activeTab = 'details'

  const dispatch = createEventDispatcher()

  // Tab configuration for HorizontalTabs
  const tabs = [
    { id: 'details', label: 'Details', name: 'Details' },
    { id: 'price', label: 'Price & Tax', name: 'Price & Tax' },
    { id: 'description', label: 'Description', name: 'Description' },
    { id: 'image', label: 'Image & Color', name: 'Image & Color' }
  ]

  function closeModal() {
    dispatch('close')
  }

  function saveProduct() {
    dispatch('save')
  }

  function handleTabChange(event) {
    activeTab = event.detail.tabId
  }

  // Reset to first tab when modal opens
  $: if (show) {
    activeTab = 'details'
  }
</script>

<Modal
  {show}
  title={editingProduct ? 'Edit Product' : 'New Product'}
  size="xl"
  loading={isLoading}
  primaryButtonText={editingProduct ? 'Update Product' : 'Create Product'}
  secondaryButtonText="Cancel"
  primaryButtonDisabled={isLoading}
  on:close={closeModal}
  on:primary={saveProduct}
  on:secondary={closeModal}
>
  <!-- Tabs -->
  <div class="mb-6">
    <HorizontalTabs
      {tabs}
      {activeTab}
      variant="glass"
      on:tabChange={handleTabChange}
    />
  </div>

  <!-- Tab Content -->
  <div class="min-h-[400px]">
    {#if activeTab === 'details'}
      <ProductDetailsTab 
        bind:productForm 
        {categories} 
        {units}
        {generateSKU}
        {generateBarcode}
      />
    {:else if activeTab === 'price'}
      <ProductPriceTab 
        bind:productForm 
        {taxRates}
      />
    {:else if activeTab === 'description'}
      <ProductDescriptionTab 
        bind:productForm 
      />
    {:else if activeTab === 'image'}
      <ProductImageTab 
        bind:productForm 
      />
    {/if}
  </div>
</Modal>