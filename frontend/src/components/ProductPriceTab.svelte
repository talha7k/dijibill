<script>
  export let productForm = {}
  export let taxRates = []

  // Calculate sale price based on cost and markup
  $: if (productForm.cost && productForm.markup) {
    productForm.unit_price = productForm.cost * (1 + productForm.markup / 100)
  }
</script>

<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
  <!-- Left Column -->
  <div class="space-y-4">
    <!-- Tax Rate -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Tax Rate</span>
      </label>
      <div class="join">
        <select 
          bind:value={productForm.vat_rate}
          class="select select-bordered join-item flex-1 bg-white/10 text-white"
        >
          {#each taxRates as rate}
            <option value={rate.rate}>{rate.name} ({rate.rate}%)</option>
          {/each}
        </select>
        <button type="button" class="btn btn-outline join-item text-white border-white/20">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
        </button>
        <button type="button" class="btn btn-outline join-item text-white border-white/20">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4"></path>
          </svg>
        </button>
      </div>
    </div>

    <!-- Cost -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Cost</span>
      </label>
      <input 
        type="number" 
        bind:value={productForm.cost}
        placeholder="0" 
        class="input input-bordered bg-white/10 text-white placeholder-white/50"
        min="0"
        step="0.01"
      />
    </div>

    <!-- Markup -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Markup</span>
      </label>
      <div class="join">
        <input 
          type="number" 
          bind:value={productForm.markup}
          placeholder="0" 
          class="input input-bordered join-item flex-1 bg-white/10 text-white placeholder-white/50"
          min="0"
          step="0.01"
        />
        <span class="btn btn-outline join-item text-white border-white/20 no-animation">%</span>
      </div>
    </div>

    <!-- Sale Price -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Sale Price</span>
      </label>
      <input 
        type="number" 
        bind:value={productForm.unit_price}
        placeholder="0" 
        class="input input-bordered bg-white/10 text-white placeholder-white/50"
        min="0"
        step="0.01"
      />
    </div>
  </div>

  <!-- Right Column -->
  <div class="space-y-4">
    <!-- Price Settings -->
    <div class="space-y-3">
      <div class="form-control">
        <label class="label cursor-pointer justify-start gap-3">
          <input 
            type="checkbox" 
            bind:checked={productForm.price_includes_tax}
            class="checkbox checkbox-success" 
          />
          <span class="label-text text-white">Price includes tax</span>
        </label>
      </div>

      <div class="form-control">
        <label class="label cursor-pointer justify-start gap-3">
          <input 
            type="checkbox" 
            bind:checked={productForm.price_change_allowed}
            class="checkbox checkbox-primary" 
          />
          <span class="label-text text-white">Price change allowed</span>
        </label>
      </div>
    </div>

    <!-- Price Calculation Info -->
    <div class="bg-white/5 p-4 rounded-lg border border-white/10">
      <h4 class="text-white font-medium mb-2">Price Calculation</h4>
      <div class="space-y-1 text-sm text-white/70">
        <div class="flex justify-between">
          <span>Cost:</span>
          <span>{productForm.cost || 0} SAR</span>
        </div>
        <div class="flex justify-between">
          <span>Markup ({productForm.markup || 0}%):</span>
          <span>{((productForm.cost || 0) * (productForm.markup || 0) / 100).toFixed(2)} SAR</span>
        </div>
        <div class="border-t border-white/20 pt-1 flex justify-between font-medium text-white">
          <span>Sale Price:</span>
          <span>{(productForm.unit_price || 0).toFixed(2)} SAR</span>
        </div>
        {#if productForm.price_includes_tax}
          <div class="flex justify-between text-xs">
            <span>Tax ({productForm.vat_rate}%):</span>
            <span>{((productForm.unit_price || 0) * (productForm.vat_rate || 0) / (100 + (productForm.vat_rate || 0))).toFixed(2)} SAR</span>
          </div>
        {:else}
          <div class="flex justify-between text-xs">
            <span>+ Tax ({productForm.vat_rate}%):</span>
            <span>{((productForm.unit_price || 0) * (productForm.vat_rate || 0) / 100).toFixed(2)} SAR</span>
          </div>
          <div class="flex justify-between font-medium text-white">
            <span>Total with Tax:</span>
            <span>{((productForm.unit_price || 0) * (1 + (productForm.vat_rate || 0) / 100)).toFixed(2)} SAR</span>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>