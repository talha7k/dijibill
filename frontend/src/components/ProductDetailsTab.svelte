<script>
  export let productForm = {}
  export let categories = []
  export let units = []
  export let generateSKU = () => {}
  export let generateBarcode = () => {}
</script>

<div class="space-y-6">
  <!-- Basic Information -->
  <div class="glass-card">
    <div class="p-6">
      <h3 class="heading-3 mb-6">Basic Information</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Product Name -->
        <div class="form-control">
          <label class="label-glass">Product Name</label>
          <input 
            type="text" 
            class="input-glass"
            bind:value={productForm.name}
            placeholder="Enter product name"
            required
          />
        </div>

        <!-- Product Name Arabic -->
        <div class="form-control">
          <label class="label-glass">Product Name (Arabic)</label>
          <input 
            type="text" 
            class="input-glass"
            bind:value={productForm.name_arabic}
            placeholder="أدخل اسم المنتج"
            dir="rtl"
          />
        </div>

        <!-- SKU -->
        <div class="form-control">
          <label class="label-glass">SKU</label>
          <div class="flex space-x-2">
            <input 
              type="text" 
              class="input-glass flex-1"
              bind:value={productForm.sku}
              placeholder="Product SKU"
            />
            <button 
              type="button" 
              class="btn-glass"
              on:click={generateSKU}
              title="Generate SKU"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
              </svg>
            </button>
          </div>
        </div>

        <!-- Barcode -->
        <div class="form-control">
          <label class="label-glass">Barcode</label>
          <div class="flex space-x-2">
            <input 
              type="text" 
              class="input-glass flex-1"
              bind:value={productForm.barcode}
              placeholder="Product barcode"
            />
            <button 
              type="button" 
              class="btn-glass"
              on:click={generateBarcode}
              title="Generate Barcode"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
              </svg>
            </button>
          </div>
        </div>

        <!-- Category -->
        <div class="form-control">
          <label class="label-glass">Category</label>
          <select class="select-glass" bind:value={productForm.category_id}>
            <option value={0}>Select Category</option>
            {#each categories as category}
              <option value={category.id}>{category.name}</option>
            {/each}
          </select>
        </div>

        <!-- Unit of Measurement -->
         <div class="form-control">
           <label class="label-glass">Unit of Measurement</label>
           <select class="select-glass" bind:value={productForm.unit}>
             {#each units as unit}
               <option value={unit.value}>{unit.label}</option>
             {/each}
           </select>
         </div>
       </div>
     </div>
   </div>

   <!-- Inventory & Status -->
   <div class="glass-card">
     <div class="p-6">
       <h3 class="heading-3 mb-6">Inventory & Status</h3>
       
       <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
         <!-- Current Stock -->
         <div class="form-control">
           <label class="label-glass">Current Stock</label>
           <input 
             type="number" 
             class="input-glass"
             bind:value={productForm.stock}
             placeholder="0"
             min="0"
           />
         </div>

         <!-- Minimum Stock Alert -->
         <div class="form-control">
           <label class="label-glass">Minimum Stock Alert</label>
           <input 
             type="number" 
             class="input-glass"
             bind:value={productForm.min_stock}
             placeholder="0"
             min="0"
           />
         </div>
       </div>

       <!-- Status Checkboxes -->
       <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-6">
         <div class="form-control">
           <label class="label-glass">Status Options</label>
           <div class="space-y-3">
             <label class="checkbox-glass">
               <input 
                 type="checkbox" 
                 class="checkbox-primary"
                 bind:checked={productForm.is_active}
               />
               <span>Active</span>
             </label>
             
             <label class="checkbox-glass">
               <input 
                 type="checkbox" 
                 class="checkbox-primary"
                 bind:checked={productForm.default_quantity}
               />
               <span>Default quantity</span>
             </label>
             
             <label class="checkbox-glass">
               <input 
                 type="checkbox" 
                 class="checkbox-primary"
                 bind:checked={productForm.service_not_using_stock}
               />
               <span>Service (not using stock)</span>
             </label>
           </div>
         </div>
       </div>
    </div>
  </div>
</div>