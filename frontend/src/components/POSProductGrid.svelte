<script>
  import { 
    currentView,
    selectedCategory,
    categories,
    filteredProducts,
    searchQuery,
    loading,
    addProductToSale
  } from '../stores/posStore.js'

  export let onSelectCategory = (category) => {}
  export let onBackToCategories = () => {}

  function selectCategory(category) {
    onSelectCategory(category)
  }

  function backToCategories() {
    currentView.set('categories')
    selectedCategory.set(null)
    searchQuery.set('')
    onBackToCategories()
  }
</script>

<div class="flex-1 bg-gray-900 flex flex-col">
  <!-- Header -->
  <div class="p-4 border-b border-gray-700">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-4">
        {#if $currentView === 'products'}
          <button
            class="p-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
            on:click={backToCategories}
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
            </svg>
          </button>
          <h2 class="text-xl font-bold">{$selectedCategory?.name || 'Products'}</h2>
        {:else}
          <h2 class="text-xl font-bold">Product Categories</h2>
        {/if}
      </div>

      <!-- Search (only show in products view) -->
      {#if $currentView === 'products'}
        <div class="relative">
          <input
            type="text"
            bind:value={$searchQuery}
            placeholder="Search products..."
            class="w-64 pl-4 pr-10 py-2 bg-gray-800 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
          </div>
        </div>
      {/if}

      <!-- Placeholder for language change -->
      <button
        class="p-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors opacity-50 cursor-not-allowed"
        disabled
        title="Change Language (Coming Soon)"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129"></path>
        </svg>
      </button>
    </div>
  </div>

  <!-- Content Area -->
  <div class="flex-1 p-6 overflow-y-auto">
    {#if $loading}
      <div class="flex items-center justify-center h-64">
        <div class="text-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
          <p>Loading...</p>
        </div>
      </div>
    {:else if $currentView === 'categories'}
      <!-- Categories Grid -->
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {#each $categories as category}
          <button
            class="aspect-square bg-gradient-to-br from-blue-500 to-blue-700 hover:from-blue-600 hover:to-blue-800 rounded-lg p-4 text-white font-medium text-lg transition-all transform hover:scale-105 shadow-lg"
            on:click={() => selectCategory(category)}
          >
            <div class="h-full flex flex-col items-center justify-center">
              <div class="text-3xl mb-2">📦</div>
              <div class="text-center">{category.name}</div>
            </div>
          </button>
        {/each}
      </div>
    {:else}
      <!-- Products Grid -->
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {#each $filteredProducts as product}
          <button
            class="aspect-square bg-gradient-to-br from-green-500 to-green-700 hover:from-green-600 hover:to-green-800 rounded-lg p-4 text-white font-medium transition-all transform hover:scale-105 shadow-lg"
            on:click={() => addProductToSale(product)}
          >
            <div class="h-full flex flex-col items-center justify-center">
              <div class="text-2xl mb-2">🛍️</div>
              <div class="text-center text-sm mb-1">{product.name}</div>
              <div class="text-lg font-bold">{(product.unit_price || 0).toFixed(2)} SAR</div>
            </div>
          </button>
        {/each}
      </div>
    {/if}
  </div>
</div>