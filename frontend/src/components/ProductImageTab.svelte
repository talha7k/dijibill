<script>
  export let productForm = {}

  let imagePreview = null

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

  function removeImage() {
    imagePreview = null
    productForm.image_url = ''
  }
</script>

<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
  <!-- Left Column - Image -->
  <div class="space-y-4">
    <h4 class="text-white font-medium mb-4">Product Image</h4>
    
    <!-- Image Upload -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Image URL</span>
      </label>
      <input 
        type="url" 
        bind:value={productForm.image_url}
        placeholder="https://example.com/image.jpg" 
        class="input input-bordered bg-white/10 text-white placeholder-white/50"
      />
    </div>

    <!-- File Upload -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Upload Image</span>
      </label>
      <input 
        type="file" 
        accept="image/*"
        class="file-input file-input-bordered bg-white/10 text-white"
        on:change={handleImageUpload}
      />
    </div>

    <!-- Image Preview -->
    <div class="bg-white/5 p-4 rounded-lg border border-white/10">
      <h5 class="text-white/70 text-sm mb-3">Image Preview</h5>
      {#if productForm.image_url || imagePreview}
        <div class="relative">
          <img 
            src={imagePreview || productForm.image_url} 
            alt="Product preview"
            class="w-full h-48 object-cover rounded-lg border border-white/20"
          />
          <button 
            type="button"
            class="absolute top-2 right-2 btn btn-sm btn-circle btn-error"
            on:click={removeImage}
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
      {:else}
        <div class="w-full h-48 bg-white/5 border-2 border-dashed border-white/20 rounded-lg flex items-center justify-center">
          <div class="text-center text-white/50">
            <svg class="w-12 h-12 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
            </svg>
            <p>No image selected</p>
          </div>
        </div>
      {/if}
    </div>
  </div>

  <!-- Right Column - Color & Visual -->
  <div class="space-y-4">
    <h4 class="text-white font-medium mb-4">Visual Properties</h4>
    
    <!-- Color Picker -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Product Color</span>
      </label>
      <div class="flex gap-3 items-center">
        <input 
          type="color" 
          bind:value={productForm.color}
          class="w-16 h-12 rounded border border-white/20 bg-transparent cursor-pointer"
        />
        <input 
          type="text" 
          bind:value={productForm.color}
          placeholder="#000000" 
          class="input input-bordered flex-1 bg-white/10 text-white placeholder-white/50"
        />
      </div>
    </div>

    <!-- Color Preview -->
    <div class="bg-white/5 p-4 rounded-lg border border-white/10">
      <h5 class="text-white/70 text-sm mb-3">Color Preview</h5>
      <div class="flex items-center gap-4">
        <div 
          class="w-16 h-16 rounded-lg border border-white/20"
          style="background-color: {productForm.color}"
        ></div>
        <div class="text-white/80">
          <div class="text-sm">Selected Color</div>
          <div class="text-xs text-white/60">{productForm.color}</div>
        </div>
      </div>
    </div>

    <!-- Quick Color Palette -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Quick Colors</span>
      </label>
      <div class="grid grid-cols-6 gap-2">
        {#each ['#000000', '#FFFFFF', '#FF0000', '#00FF00', '#0000FF', '#FFFF00', '#FF00FF', '#00FFFF', '#FFA500', '#800080', '#FFC0CB', '#A52A2A'] as color}
          <button 
            type="button"
            class="w-8 h-8 rounded border border-white/20 cursor-pointer hover:scale-110 transition-transform"
            style="background-color: {color}"
            on:click={() => productForm.color = color}
          ></button>
        {/each}
      </div>
    </div>

    <!-- Image Guidelines -->
    <div class="bg-blue-500/10 p-4 rounded-lg border border-blue-500/20">
      <h5 class="text-blue-300 text-sm font-medium mb-2">Image Guidelines</h5>
      <ul class="text-blue-200/80 text-xs space-y-1">
        <li>• Recommended size: 800x800 pixels</li>
        <li>• Supported formats: JPG, PNG, WebP</li>
        <li>• Maximum file size: 5MB</li>
        <li>• Use high-quality images for better presentation</li>
      </ul>
    </div>
  </div>
</div>