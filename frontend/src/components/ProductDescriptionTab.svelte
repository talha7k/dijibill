<script>
  export let productForm = {}
</script>

<div class="space-y-6">
  <!-- Description Information -->
  <div class="glass-card">
    <div class="p-6">
      <h3 class="heading-3 mb-6">Product Description</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- English Description -->
        <div class="form-control">
          <label class="label-glass">Description (English)</label>
          <textarea 
            bind:value={productForm.description}
            placeholder="Enter product description" 
            class="textarea-glass h-32"
          ></textarea>
        </div>

        <!-- Arabic Description -->
        <div class="form-control">
          <label class="label-glass">Description (Arabic)</label>
          <textarea 
            bind:value={productForm.description_arabic}
            placeholder="أدخل وصف المنتج" 
            class="textarea-glass h-32"
            dir="rtl"
          ></textarea>
        </div>
      </div>
    </div>
  </div>

  <!-- Additional Notes -->
  <div class="glass-card">
    <div class="p-6">
      <h3 class="heading-3 mb-6">Additional Information</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- English Notes -->
        <div class="form-control">
          <label class="label-glass">Additional Notes (English)</label>
          <textarea 
            bind:value={productForm.additional_notes}
            placeholder="Any additional notes or specifications" 
            class="textarea-glass h-24"
          ></textarea>
        </div>

        <!-- Arabic Notes -->
        <div class="form-control">
          <label class="label-glass">Additional Notes (Arabic)</label>
          <textarea 
            bind:value={productForm.additional_notes_arabic}
            placeholder="أي ملاحظات أو مواصفات إضافية" 
            class="textarea-glass h-24"
            dir="rtl"
          ></textarea>
        </div>
      </div>
    </div>
  </div>

  <!-- Description Preview -->
  <div class="glass-card">
    <div class="p-6">
      <h3 class="heading-3 mb-6">Description Preview</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="form-control">
          <label class="label-glass">English Preview</label>
          <div class="bg-white/5 p-4 rounded-lg border border-white/10 min-h-[100px]">
            <div class="text-white/80 text-sm whitespace-pre-wrap">
              {productForm.description || 'No description provided'}
            </div>
          </div>
        </div>
        
        <div class="form-control">
          <label class="label-glass">Arabic Preview</label>
          <div class="bg-white/5 p-4 rounded-lg border border-white/10 min-h-[100px]" dir="rtl">
            <div class="text-white/80 text-sm whitespace-pre-wrap">
              {productForm.description_arabic || 'لم يتم تقديم وصف'}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>