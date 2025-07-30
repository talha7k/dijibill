<script>
  export let productForm = {}
</script>

<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
  <!-- Left Column - English -->
  <div class="space-y-4">
    <h4 class="text-white font-medium mb-4">English Description</h4>
    
    <!-- Description -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Description</span>
      </label>
      <textarea 
        bind:value={productForm.description}
        placeholder="Enter product description" 
        class="textarea textarea-bordered bg-white/10 text-white placeholder-white/50 h-32"
      ></textarea>
    </div>

    <!-- Additional Notes -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Additional Notes</span>
      </label>
      <textarea 
        placeholder="Any additional notes or specifications" 
        class="textarea textarea-bordered bg-white/10 text-white placeholder-white/50 h-24"
      ></textarea>
    </div>
  </div>

  <!-- Right Column - Arabic -->
  <div class="space-y-4">
    <h4 class="text-white font-medium mb-4">Arabic Description</h4>
    
    <!-- Description Arabic -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Description (Arabic)</span>
      </label>
      <textarea 
        bind:value={productForm.description_arabic}
        placeholder="أدخل وصف المنتج" 
        class="textarea textarea-bordered bg-white/10 text-white placeholder-white/50 h-32"
        dir="rtl"
      ></textarea>
    </div>

    <!-- Additional Notes Arabic -->
    <div class="form-control">
      <label class="label">
        <span class="label-text text-white">Additional Notes (Arabic)</span>
      </label>
      <textarea 
        placeholder="أي ملاحظات أو مواصفات إضافية" 
        class="textarea textarea-bordered bg-white/10 text-white placeholder-white/50 h-24"
        dir="rtl"
      ></textarea>
    </div>
  </div>
</div>

<!-- Description Preview -->
<div class="mt-6 bg-white/5 p-4 rounded-lg border border-white/10">
  <h4 class="text-white font-medium mb-3">Description Preview</h4>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <div>
      <h5 class="text-white/70 text-sm mb-2">English</h5>
      <div class="text-white/80 text-sm whitespace-pre-wrap">
        {productForm.description || 'No description provided'}
      </div>
    </div>
    <div>
      <h5 class="text-white/70 text-sm mb-2">Arabic</h5>
      <div class="text-white/80 text-sm whitespace-pre-wrap" dir="rtl">
        {productForm.description_arabic || 'لم يتم تقديم وصف'}
      </div>
    </div>
  </div>
</div>