<script>
  import FormField from './FormField.svelte'
  
  export let productForm = {}
</script>

<div class="space-y-6">
  <!-- Description Information -->
  <div class="glass-card">
    <div class="p-6">
      <h3 class="heading-3 mb-6">Product Description</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- English Description -->
        <FormField
          label="Description (English)"
          type="textarea"
          bind:value={productForm.description}
          placeholder="Enter product description"
          rows={5}
        />

        <!-- Arabic Description -->
        <FormField
          label="Description (Arabic)"
          labelArabic="الوصف (عربي)"
          type="textarea"
          bind:value={productForm.description_arabic}
          placeholder="أدخل وصف المنتج"
          rows={5}
          dir="rtl"
        />
      </div>
    </div>
  </div>

  <!-- Additional Notes -->
  <div class="glass-card">
    <div class="p-6">
      <h3 class="heading-3 mb-6">Additional Information</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- English Notes -->
        <FormField
          label="Additional Notes (English)"
          type="textarea"
          bind:value={productForm.additional_notes}
          placeholder="Any additional notes or specifications"
          rows={3}
        />

        <!-- Arabic Notes -->
        <FormField
          label="Additional Notes (Arabic)"
          labelArabic="ملاحظات إضافية (عربي)"
          type="textarea"
          bind:value={productForm.additional_notes_arabic}
          placeholder="أي ملاحظات أو مواصفات إضافية"
          rows={3}
          dir="rtl"
        />
      </div>
    </div>
  </div>
</div>