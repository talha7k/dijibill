<script>
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'

  export let show = false
  export let editingPurchaseProduct = null
  export let purchaseProductForm = {}
  export let categories = []
  export let units = []
  export let isLoading = false
  export let generateSKU
  export let generateBarcode

  const dispatch = createEventDispatcher()

  function closeModal() {
    dispatch('close')
  }

  function savePurchaseProduct() {
    dispatch('save')
  }

  // Get category name for display
  $: selectedCategory = categories.find(cat => cat.id === purchaseProductForm.category_id)
  
  // Get unit details for display
  $: selectedUnit = units.find(unit => unit.value === purchaseProductForm.unit)
</script>

<Modal
  {show}
  title={editingPurchaseProduct ? 'Edit Purchase Product' : 'New Purchase Product'}
  size="lg"
  loading={isLoading}
  primaryButtonText={editingPurchaseProduct ? 'Update Purchase Product' : 'Create Purchase Product'}
  secondaryButtonText="Cancel"
  primaryButtonDisabled={isLoading}
  on:close={closeModal}
  on:primary={savePurchaseProduct}
  on:secondary={closeModal}
>
  <div class="space-y-6">
    <!-- Basic Information -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Product Name -->
      <FormField
        label="Product Name"
        type="text"
        bind:value={purchaseProductForm.name}
        placeholder="Enter product name"
        required={true}
      />

      <!-- Product Name Arabic -->
      <FormField
        label="Product Name (Arabic)"
        labelArabic="اسم المنتج"
        type="text"
        bind:value={purchaseProductForm.name_arabic}
        placeholder="أدخل اسم المنتج"
        dir="rtl"
      />
    </div>

    <!-- Category and Unit -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Category -->
      <FormField
        label="Category"
        type="select"
        bind:value={purchaseProductForm.category_id}
        options={[
          { value: 0, label: 'Select Category' },
          ...categories.map(cat => ({ value: cat.id, label: cat.name }))
        ]}
      />

      <!-- Unit -->
      <FormField
        label="Unit of Measurement"
        type="select"
        bind:value={purchaseProductForm.unit}
        options={units.map(unit => ({ value: unit.value, label: unit.label }))}
      />
    </div>

    <!-- Price and VAT -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Unit Price -->
      <FormField
        label="Unit Price (SAR)"
        type="number"
        bind:value={purchaseProductForm.unit_price}
        placeholder="0.00"
        step="0.01"
        min="0"
        required={true}
      />

      <!-- VAT Rate -->
      <FormField
        label="VAT Rate (%)"
        type="number"
        bind:value={purchaseProductForm.vat_rate}
        placeholder="15.0"
        step="0.1"
        min="0"
        max="100"
      />
    </div>

    <!-- SKU and Barcode -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- SKU -->
      <div class="space-y-2">
        <FormField
          label="SKU"
          type="text"
          bind:value={purchaseProductForm.sku}
          placeholder="Enter SKU"
        />
        <button 
          type="button" 
          class="btn btn-sm btn-outline btn-secondary"
          on:click={generateSKU}
        >
          Generate SKU
        </button>
      </div>

      <!-- Barcode -->
      <div class="space-y-2">
        <FormField
          label="Barcode"
          type="text"
          bind:value={purchaseProductForm.barcode}
          placeholder="Enter barcode"
        />
        <button 
          type="button" 
          class="btn btn-sm btn-outline btn-secondary"
          on:click={generateBarcode}
        >
          Generate Barcode
        </button>
      </div>
    </div>

    <!-- Description -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Description -->
      <FormField
        label="Description"
        type="textarea"
        bind:value={purchaseProductForm.description}
        placeholder="Enter product description"
        rows={3}
      />

      <!-- Description Arabic -->
      <FormField
        label="Description (Arabic)"
        labelArabic="الوصف"
        type="textarea"
        bind:value={purchaseProductForm.description_arabic}
        placeholder="أدخل وصف المنتج"
        dir="rtl"
        rows={3}
      />
    </div>

    <!-- Notes -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Notes -->
      <FormField
        label="Notes"
        type="textarea"
        bind:value={purchaseProductForm.notes}
        placeholder="Enter any additional notes"
        rows={2}
      />

      <!-- Notes Arabic -->
      <FormField
        label="Notes (Arabic)"
        labelArabic="ملاحظات"
        type="textarea"
        bind:value={purchaseProductForm.notes_arabic}
        placeholder="أدخل أي ملاحظات إضافية"
        dir="rtl"
        rows={2}
      />
    </div>

    <!-- Status -->
    <div class="form-control">
      <label class="label cursor-pointer justify-start gap-3">
        <input 
          type="checkbox" 
          class="checkbox checkbox-primary" 
          bind:checked={purchaseProductForm.is_active}
        />
        <span class="label-text text-white">Active</span>
      </label>
    </div>
  </div>
</Modal>