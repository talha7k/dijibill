<script>
  import FormField from './FormField.svelte'
  import IconButton from './IconButton.svelte'
  
  export let productForm = {}
  export let taxRates = []

  // Transform tax rates for FormField options
  $: taxRateOptions = (taxRates || []).map(rate => ({ 
    value: rate.rate, 
    label: `${rate.name} (${rate.rate}%)` 
  }))

  // Calculate sale price based on cost and markup
  $: if (productForm.cost && productForm.markup) {
    productForm.unit_price = productForm.cost * (1 + productForm.markup / 100)
  }
</script>

<div class="space-y-6">
  <!-- Pricing Information -->
  <div class="glass-card">
    <div class="p-6">
      <h3 class="heading-3 mb-6">Pricing Information</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Tax Rate -->
        <div class="form-control">
          <label for="tax-rate-select" class="label-glass">Tax Rate</label>
          <div class="flex space-x-2">
            <select 
              id="tax-rate-select"
              bind:value={productForm.vat_rate}
              class="select-glass flex-1"
            >
              {#each (taxRates || []) as rate}
                <option value={rate.rate}>{rate.name} ({rate.rate}%)</option>
              {/each}
            </select>
            <IconButton
              icon="plus"
              variant="outline"
              title="Add Tax Rate"
            />
          </div>
        </div>

        <!-- Cost -->
        <FormField
          label="Cost"
          labelArabic="التكلفة"
          type="number"
          bind:value={productForm.cost}
          placeholder="0.00"
          min={0}
          step={0.01}
        />

        <!-- Markup -->
        <FormField
          label="Markup (%)"
          labelArabic="هامش الربح (%)"
          type="number"
          bind:value={productForm.markup}
          placeholder="0"
          min={0}
          step={0.01}
        />

        <!-- Sale Price -->
        <FormField
          label="Sale Price"
          labelArabic="سعر البيع"
          type="number"
          bind:value={productForm.unit_price}
          placeholder="0.00"
          min={0}
          step={0.01}
        />
      </div>
    </div>
  </div>

  <!-- Price Settings -->
  <div class="glass-card">
    <div class="p-6">
      <h3 class="heading-3 mb-6">Price Settings</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="form-control">
          <label id="price-options-label" class="label-glass">Price Options</label>
          <div class="space-y-3" role="group" aria-labelledby="price-options-label">
            <FormField
              label="Price includes tax"
              type="checkbox"
              bind:checked={productForm.price_includes_tax}
              placeholder="Price includes tax"
            />
            
            <FormField
              label="Price change allowed"
              type="checkbox"
              bind:checked={productForm.price_change_allowed}
              placeholder="Price change allowed"
            />
          </div>
        </div>

        <!-- Price Calculation Info -->
        <div class="form-control">
          <label id="price-calculation-label" class="label-glass">Price Calculation</label>
          <div class="bg-white/5 p-4 rounded-lg border border-white/10" role="region" aria-labelledby="price-calculation-label">
            <div class="space-y-2 text-sm text-white/70">
              <div class="flex justify-between">
                <span>Cost:</span>
                <span>{(productForm.cost || 0).toFixed(2)} SAR</span>
              </div>
              <div class="flex justify-between">
                <span>Markup ({productForm.markup || 0}%):</span>
                <span>{((productForm.cost || 0) * (productForm.markup || 0) / 100).toFixed(2)} SAR</span>
              </div>
              <div class="border-t border-white/20 pt-2 flex justify-between font-medium text-white">
                <span>Sale Price:</span>
                <span>{(productForm.unit_price || 0).toFixed(2)} SAR</span>
              </div>
              {#if productForm.price_includes_tax}
                <div class="flex justify-between text-xs">
                  <span>Tax ({productForm.vat_rate || 0}%):</span>
                  <span>{((productForm.unit_price || 0) * (productForm.vat_rate || 0) / (100 + (productForm.vat_rate || 0))).toFixed(2)} SAR</span>
                </div>
              {:else}
                <div class="flex justify-between text-xs">
                  <span>+ Tax ({productForm.vat_rate || 0}%):</span>
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
    </div>
  </div>
</div>