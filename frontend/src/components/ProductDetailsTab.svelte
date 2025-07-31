<script>
  import { createEventDispatcher } from 'svelte'
  import FormField from './FormField.svelte'
  
  export let productForm = {}
  export let categories = []
  export let units = []
  export let generateSKU = () => {}
  export let generateBarcode = () => {}

  const dispatch = createEventDispatcher()

  // Transform categories and units for FormField options
  $: categoryOptions = (categories || []).map(cat => ({ value: cat.id, label: cat.name }))
  $: unitOptions = (units || []).map(unit => ({ value: unit.value, label: unit.label }))
</script>

<div class="space-y-6">
  <!-- Basic Information -->
  <div class="glass-card">
    <div class="p-6">
      <h3 class="heading-3 mb-6">Basic Information</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Product Name -->
        <FormField
          label="Product Name"
          type="text"
          bind:value={productForm.name}
          placeholder="Enter product name"
          required={true}
        />

        <!-- Product Name Arabic -->
        <FormField
          label="Product Name"
          labelArabic="اسم المنتج"
          type="text"
          bind:value={productForm.name_arabic}
          placeholder="أدخل اسم المنتج"
          dir="rtl"
        />

        <!-- SKU -->
        <div class="form-control">
          <label for="product-sku" class="label-glass">SKU</label>
          <div class="flex space-x-2">
            <input 
              id="product-sku"
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
              <i class="fas fa-sync-alt w-4 h-4"></i>
            </button>
          </div>
        </div>

        <!-- Barcode -->
        <div class="form-control">
          <label for="product-barcode" class="label-glass">Barcode</label>
          <div class="flex space-x-2">
            <input 
              id="product-barcode"
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
              <i class="fas fa-sync-alt w-4 h-4"></i>
            </button>
          </div>
        </div>

        <!-- Category -->
        <div class="form-control">
          <label for="product-category" class="label-glass">Category</label>
          <div class="flex space-x-2">
            <select id="product-category" class="select-glass flex-1" bind:value={productForm.category_id}>
              <option value={0}>Select Category</option>
              {#each (categories || []) as category}
                <option value={category.id}>{category.name}</option>
              {/each}
            </select>
            <button 
              type="button" 
              class="btn-glass"
              title="Add Category"
              on:click={() => dispatch('addCategory')}
            >
              <i class="fas fa-plus w-4 h-4"></i>
            </button>
          </div>
        </div>

        <!-- Unit of Measurement -->
        <FormField
          label="Unit of Measurement"
          labelArabic="وحدة القياس"
          type="select"
          bind:value={productForm.unit}
          options={unitOptions}
          placeholder="Select unit"
        />
       </div>
     </div>
   </div>

   <!-- Inventory & Status -->
   <div class="glass-card">
     <div class="p-6">
       <h3 class="heading-3 mb-6">Inventory & Status</h3>
       
       <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
         <!-- Current Stock -->
         <FormField
           label="Current Stock"
           labelArabic="المخزون الحالي"
           type="number"
           bind:value={productForm.stock}
           placeholder="0"
           min={0}
         />

         <!-- Minimum Stock Alert -->
         <FormField
           label="Minimum Stock Alert"
           labelArabic="تنبيه الحد الأدنى للمخزون"
           type="number"
           bind:value={productForm.min_stock}
           placeholder="0"
           min={0}
         />
       </div>

       <!-- Status Checkboxes -->
       <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-6">
         <div class="form-control">
           <label id="product-status-label" class="label-glass">Product Status</label>
           <div class="space-y-3" role="group" aria-labelledby="product-status-label">
             <FormField
               label="Active Product"
               type="checkbox"
               bind:checked={productForm.is_active}
               placeholder="Active Product"
             />
             
             <FormField
               label="Service (not using stock)"
               type="checkbox"
               bind:checked={productForm.service_not_using_stock}
               placeholder="Service (not using stock)"
             />
           </div>
         </div>
       </div>
    </div>
  </div>
</div>