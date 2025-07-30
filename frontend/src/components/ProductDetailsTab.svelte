<script>
  export let productForm = {}
  export let categories = []
  export let units = []
  export let generateSKU = () => {}
  export let generateBarcode = () => {}
</script>

<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
  <!-- Left Column -->
  <div class="space-y-4">
    <!-- Product Name -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Product Name</span>
      </label>
      <input 
        type="text" 
        bind:value={productForm.name}
        placeholder="Enter product name" 
        class="input input-bordered bg-white/10 text-white placeholder-white/50"
        required
      />
    </div>

    <!-- Product Name Arabic -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Product Name (Arabic)</span>
      </label>
      <input 
        type="text" 
        bind:value={productForm.name_arabic}
        placeholder="أدخل اسم المنتج" 
        class="input input-bordered bg-white/10 text-white placeholder-white/50"
        dir="rtl"
      />
    </div>

    <!-- SKU -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">SKU</span>
      </label>
      <div class="join">
        <input 
          type="text" 
          bind:value={productForm.sku}
          placeholder="Product SKU" 
          class="input input-bordered join-item flex-1 bg-white/10 text-white placeholder-white/50"
        />
        <button 
          type="button"
          class="btn btn-outline join-item text-white border-white/20"
          on:click={generateSKU}
        >
          Generate
        </button>
      </div>
    </div>

    <!-- Barcode -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Barcode</span>
      </label>
      <div class="join">
        <input 
          type="text" 
          bind:value={productForm.barcode}
          placeholder="Product barcode" 
          class="input input-bordered join-item flex-1 bg-white/10 text-white placeholder-white/50"
        />
        <button 
          type="button"
          class="btn btn-outline join-item text-white border-white/20"
          on:click={generateBarcode}
        >
          Generate
        </button>
      </div>
    </div>

    <!-- Unit of Measurement -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Unit of Measurement</span>
      </label>
      <select 
        bind:value={productForm.unit}
        class="select select-bordered bg-white/10 text-white"
      >
        {#each units as unit}
          <option value={unit.value}>{unit.label}</option>
        {/each}
      </select>
    </div>
  </div>

  <!-- Right Column -->
  <div class="space-y-4">
    <!-- Category -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Category</span>
      </label>
      <select 
        bind:value={productForm.category_id}
        class="select select-bordered bg-white/10 text-white"
      >
        <option value={0}>Select Category</option>
        {#each categories as category}
          <option value={category.id}>{category.name}</option>
        {/each}
      </select>
    </div>

    <!-- Stock -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Current Stock</span>
      </label>
      <input 
        type="number" 
        bind:value={productForm.stock}
        placeholder="0" 
        class="input input-bordered bg-white/10 text-white placeholder-white/50"
        min="0"
      />
    </div>

    <!-- Minimum Stock -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Minimum Stock Alert</span>
      </label>
      <input 
        type="number" 
        bind:value={productForm.min_stock}
        placeholder="0" 
        class="input input-bordered bg-white/10 text-white placeholder-white/50"
        min="0"
      />
    </div>

    <!-- Checkboxes -->
    <div class="space-y-3">
      <div class="form-control">
        <label class="label cursor-pointer justify-start gap-3">
          <input 
            type="checkbox" 
            bind:checked={productForm.is_active}
            class="checkbox checkbox-primary" 
          />
          <span class="label-text text-white">Active</span>
        </label>
      </div>

      <div class="form-control">
        <label class="label cursor-pointer justify-start gap-3">
          <input 
            type="checkbox" 
            bind:checked={productForm.default_quantity}
            class="checkbox checkbox-primary" 
          />
          <span class="label-text text-white">Default quantity</span>
        </label>
      </div>

      <div class="form-control">
        <label class="label cursor-pointer justify-start gap-3">
          <input 
            type="checkbox" 
            bind:checked={productForm.service_not_using_stock}
            class="checkbox checkbox-primary" 
          />
          <span class="label-text text-white">Service (not using stock)</span>
        </label>
      </div>
    </div>
  </div>
</div>