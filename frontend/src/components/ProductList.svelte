<script>
  import { createEventDispatcher } from 'svelte'
  
  export let products = []
  export let isLoading = false
  export let searchTerm = ''

  const dispatch = createEventDispatcher()

  function editProduct(product) {
    dispatch('edit', product)
  }

  function deleteProduct(product) {
    if (confirm(`Are you sure you want to delete "${product.name}"?`)) {
      dispatch('delete', product)
    }
  }

  $: filteredProducts = (products || []).filter(product => 
    product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.sku.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.barcode.toLowerCase().includes(searchTerm.toLowerCase())
  )
</script>

<div class="bg-white/5 backdrop-blur-lg rounded-lg border border-white/20 overflow-hidden">
  <!-- Search Bar -->
  <div class="p-4 border-b border-white/10">
    <div class="form-control">
      <div class="input-group">
        <span class="bg-white/10 border-white/20 text-white">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
          </svg>
        </span>
        <input 
          type="text" 
          bind:value={searchTerm}
          placeholder="Search products by name, SKU, or barcode..." 
          class="input input-bordered flex-1 bg-white/10 text-white placeholder-white/50 border-white/20"
        />
      </div>
    </div>
  </div>

  <!-- Products Table -->
  <div class="overflow-x-auto">
    {#if isLoading}
      <div class="flex justify-center items-center py-12">
        <span class="loading loading-spinner loading-lg text-primary"></span>
      </div>
    {:else if filteredProducts.length === 0}
      <div class="text-center py-12 text-white/60">
        {#if searchTerm}
          <svg class="w-16 h-16 mx-auto mb-4 text-white/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
          </svg>
          <h3 class="text-lg font-medium text-white/70 mb-2">No products found</h3>
          <p>No products match your search criteria.</p>
        {:else}
          <svg class="w-16 h-16 mx-auto mb-4 text-white/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path>
          </svg>
          <h3 class="text-lg font-medium text-white/70 mb-2">No products yet</h3>
          <p>Get started by adding your first product.</p>
        {/if}
      </div>
    {:else}
      <table class="table table-zebra">
        <thead>
          <tr class="border-white/10">
            <th class="text-white/70">Image</th>
            <th class="text-white/70">Product</th>
            <th class="text-white/70">SKU</th>
            <th class="text-white/70">Category</th>
            <th class="text-white/70">Price</th>
            <th class="text-white/70">Stock</th>
            <th class="text-white/70">Status</th>
            <th class="text-white/70">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each filteredProducts as product (product.id)}
            <tr class="hover:bg-white/5 border-white/10">
              <!-- Image -->
              <td>
                {#if product.image_url}
                  <div class="avatar">
                    <div class="w-12 h-12 rounded-lg">
                      <img src={product.image_url} alt={product.name} />
                    </div>
                  </div>
                {:else}
                  <div class="w-12 h-12 bg-white/10 rounded-lg flex items-center justify-center">
                    <svg class="w-6 h-6 text-white/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                    </svg>
                  </div>
                {/if}
              </td>
              
              <!-- Product Info -->
              <td>
                <div class="flex items-center gap-3">
                  {#if product.color}
                    <div 
                      class="w-4 h-4 rounded-full border border-white/20"
                      style="background-color: {product.color}"
                    ></div>
                  {/if}
                  <div>
                    <div class="font-bold text-white">{product.name}</div>
                    {#if product.name_arabic}
                      <div class="text-sm text-white/60" dir="rtl">{product.name_arabic}</div>
                    {/if}
                  </div>
                </div>
              </td>
              
              <!-- SKU -->
              <td class="text-white/80">{product.sku || '-'}</td>
              
              <!-- Category -->
              <td class="text-white/80">{product.category_name || 'Uncategorized'}</td>
              
              <!-- Price -->
              <td class="text-white/80">{product.unit_price.toFixed(2)} SAR</td>
              
              <!-- Stock -->
              <td>
                <div class="flex items-center gap-2">
                  <span class="text-white/80">{product.stock}</span>
                  {#if product.stock <= product.min_stock}
                    <div class="badge badge-error badge-sm">Low</div>
                  {/if}
                  {#if product.service_not_using_stock}
                    <div class="badge badge-info badge-sm">Service</div>
                  {/if}
                </div>
              </td>
              
              <!-- Status -->
              <td>
                <div class="badge {product.is_active ? 'badge-success' : 'badge-error'} badge-sm">
                  {product.is_active ? 'Active' : 'Inactive'}
                </div>
              </td>
              
              <!-- Actions -->
              <td>
                <div class="flex gap-2">
                  <button 
                    class="btn btn-ghost btn-sm text-white/70 hover:text-white"
                    on:click={() => editProduct(product)}
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                    </svg>
                  </button>
                  <button 
                    class="btn btn-ghost btn-sm text-red-400 hover:text-red-300"
                    on:click={() => deleteProduct(product)}
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}
  </div>

  <!-- Results Summary -->
  {#if !isLoading && filteredProducts.length > 0}
    <div class="p-4 border-t border-white/10 text-white/60 text-sm">
      Showing {filteredProducts.length} of {products.length} products
    </div>
  {/if}
</div>