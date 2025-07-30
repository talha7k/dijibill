<script>
  import { createEventDispatcher } from 'svelte'
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

  function closeModal() {
    dispatch('close')
  }

  function saveProduct() {
    dispatch('save')
  }



  // Reset to first tab when modal opens
  $: if (show) {
    activeTab = 'details'
  }
</script>

{#if show}
  <div class="modal modal-open">
    <div class="modal-box w-11/12 max-w-5xl bg-base-100/90 backdrop-blur-lg border border-white/20">
      <div class="flex items-center justify-between mb-6">
        <h3 class="font-bold text-lg text-white">
          {editingProduct ? 'Edit Product' : 'New Product'}
        </h3>
        <button class="btn btn-sm btn-circle btn-ghost text-white" on:click={closeModal}>
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>
      
      <!-- Tabs -->
      <div class="tabs tabs-boxed bg-white/10 mb-6">
        <button 
          class="tab {activeTab === 'details' ? 'tab-active' : ''}" 
          on:click={() => activeTab = 'details'}
        >
          Details
        </button>
        <button 
          class="tab {activeTab === 'price' ? 'tab-active' : ''}" 
          on:click={() => activeTab = 'price'}
        >
          Price & tax
        </button>
        <button 
          class="tab {activeTab === 'description' ? 'tab-active' : ''}" 
          on:click={() => activeTab = 'description'}
        >
          Description
        </button>
        <button 
          class="tab {activeTab === 'image' ? 'tab-active' : ''}" 
          on:click={() => activeTab = 'image'}
        >
          Image & color
        </button>
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

      <div class="modal-action">
        <button class="btn btn-ghost" on:click={closeModal}>Cancel</button>
        <button class="btn btn-primary" on:click={saveProduct} disabled={isLoading}>
          {#if isLoading}
            <span class="loading loading-spinner loading-sm"></span>
          {/if}
          {editingProduct ? 'Update' : 'Create'} Product
        </button>
      </div>
    </div>
  </div>
{/if}