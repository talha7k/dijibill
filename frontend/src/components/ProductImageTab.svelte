<script>
  import FormField from './FormField.svelte'
  
  export let productForm = {}

  let fileInput
  let imagePreview = productForm.image_url || ''

  function handleImageUpload(event) {
    const file = event.target.files[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        imagePreview = e.target.result
        productForm.image_url = e.target.result
      }
      reader.readAsDataURL(file)
    }
  }

  function handleImageUrlChange() {
    imagePreview = productForm.image_url
  }

  function selectColor(color) {
    productForm.color = color
  }

  // Predefined colors
  const quickColors = [
    '#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7',
    '#DDA0DD', '#98D8C8', '#F7DC6F', '#BB8FCE', '#85C1E9'
  ]

  function removeImage() {
    imagePreview = null
    productForm.image_url = ''
  }
</script>

<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
  <!-- Left Column - Image Upload -->
  <div class="space-y-6">
    <!-- Image URL Input -->
    <FormField
      label="Image URL"
      labelArabic="رابط الصورة"
      type="url"
      bind:value={productForm.image_url}
      placeholder="https://example.com/image.jpg"
      on:input={handleImageUrlChange}
    />

    <!-- File Upload -->
    <div class="space-y-3">
      <label for="product-image-upload" class="label-standard">
        Upload Image
        <span class="text-white/50 text-sm mr-2">رفع صورة</span>
      </label>
      <div class="flex gap-3">
        <input
          id="product-image-upload"
          bind:this={fileInput}
          type="file"
          accept="image/*"
          on:change={handleImageUpload}
          class="file-input file-input-bordered file-input-primary w-full bg-white/10 border-white/20 text-white"
        />
        {#if imagePreview}
          <button
            type="button"
            on:click={() => { imagePreview = ''; productForm.image_url = '' }}
            class="btn btn-outline btn-error btn-sm"
          >
            Remove
          </button>
        {/if}
      </div>
    </div>

    <!-- Image Preview -->
    {#if imagePreview}
      <div class="space-y-3">
        <label for="product-image-preview" class="label-standard">
          Preview
          <span class="text-white/50 text-sm mr-2">معاينة</span>
        </label>
        <div id="product-image-preview" class="relative w-full h-48 bg-white/5 rounded-lg border border-white/20 overflow-hidden">
          <img
            src={imagePreview}
            alt="Product preview"
            class="w-full h-full object-cover"
          />
        </div>
      </div>
    {/if}
  </div>

  <!-- Right Column - Color Selection -->
  <div class="space-y-6">
    <!-- Color Picker -->
    <FormField
      label="Product Color"
      labelArabic="لون المنتج"
      type="color"
      bind:value={productForm.color}
      placeholder="#000000"
    />

  

    <!-- Quick Color Selection -->
    <div class="space-y-3">
      <label id="quick-colors-label" class="label-standard">
        Quick Colors
        <span class="text-white/50 text-sm mr-2">الألوان السريعة</span>
      </label>
      <div class="grid grid-cols-5 gap-2" role="group" aria-labelledby="quick-colors-label">
        {#each quickColors as color}
          <button
            type="button"
            class="w-10 h-10 rounded-lg border-2 transition-all duration-200 hover:scale-110 {productForm.color === color ? 'border-white shadow-lg' : 'border-white/30'}"
            style="background-color: {color}"
            on:click={() => selectColor(color)}
            title={color}
            aria-label="Select color {color}"
          ></button>
        {/each}
      </div>
    </div>
  </div>
</div>