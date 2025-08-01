<script>
  import { GetCompany, UpdateCompany, UploadCompanyLogo, UploadCompanyLogoFromData, GetCompanyLogoAsBase64 } from '../../wailsjs/go/main/App.js'
  import { createEventDispatcher, onMount } from 'svelte'
  import FormField from './FormField.svelte'
  import { showDbSuccess, showDbError } from '../stores/notificationStore.js'

  const dispatch = createEventDispatcher()

  export let companySettings = {
    id: 0,
    name: '',
    name_arabic: '',
    vat_number: '',
    cr_number: '',
    email: '',
    phone: '',
    address: '',
    address_arabic: '',
    city: '',
    city_arabic: '',
    country: '',
    country_arabic: '',
    logo: ''
  }

  export let isLoading = false

  let logoPreview = companySettings.logo || ''
  let logoError = ''
  let fileInput

  // Allowed file types and size limit
  const ALLOWED_TYPES = ['image/webp', 'image/png', 'image/svg+xml', 'image/jpeg', 'image/jpg']
  const MAX_SIZE = 500 * 1024 // 500KB in bytes

  // Load existing logo on component mount
  onMount(async () => {
    if (companySettings.logo_file_id && companySettings.logo_file_id > 0) {
      try {
        const base64Logo = await GetCompanyLogoAsBase64(companySettings.logo_file_id)
        if (base64Logo) {
          logoPreview = `data:image/jpeg;base64,${base64Logo}`
        }
      } catch (error) {
        console.error('Error loading existing logo:', error)
        // If file not found, clear the file ID
        companySettings.logo_file_id = null
      }
    } else if (companySettings.logo) {
      // Fallback to legacy base64 logo
      logoPreview = companySettings.logo
    }
  })

  async function handleLogoUpload(event) {
    const file = event.target.files?.[0]
    if (!file) return

    logoError = ''

    // Validate file type
    if (!ALLOWED_TYPES.includes(file.type)) {
      logoError = 'Please select a valid image file (WebP, PNG, SVG, or JPG only)'
      if (fileInput) fileInput.value = ''
      return
    }

    // Validate file size
    if (file.size > MAX_SIZE) {
      logoError = 'File size must be less than 500KB'
      if (fileInput) fileInput.value = ''
      return
    }

    try {
      // Convert file to base64
      const reader = new FileReader()
      reader.onload = async (e) => {
        try {
          const result = e.target.result
          if (typeof result !== 'string') return
          
          const dataURL = result
          const base64Data = dataURL.split(',')[1] // Remove data:image/...;base64, prefix
          
          // Upload the file and get the file ID
          const fileID = await UploadCompanyLogoFromData(base64Data, file.name, file.type)
          if (fileID > 0) {
            // Store the file ID in company settings
            companySettings.logo_file_id = fileID
            // Clear the old base64 logo
            companySettings.logo = ''
            
            // Set preview directly from the uploaded file
            logoPreview = dataURL
          }
        } catch (error) {
          console.error('Error uploading logo:', error)
          logoError = 'Error uploading file. Please try again.'
          if (fileInput) fileInput.value = ''
        }
      }
      reader.readAsDataURL(file)
    } catch (error) {
      console.error('Error reading file:', error)
      logoError = 'Error reading file. Please try again.'
      if (fileInput) fileInput.value = ''
    }
  }

  function removeLogo() {
    logoPreview = ''
    companySettings.logo = ''
    companySettings.logo_file_id = null
    logoError = ''
    if (fileInput) fileInput.value = ''
  }

  function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  async function saveCompanySettings() {
    isLoading = true
    try {
      await UpdateCompany(companySettings)
      showDbSuccess('save', 'Company Settings')
    } catch (error) {
      console.error('Error saving company settings:', error)
      showDbError('save', 'Company Settings', error)
    } finally {
      isLoading = false
    }
  }
</script>

<div class="space-y-6">
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <!-- Company Name -->
    <FormField
      label="Company Name"
      labelArabic="اسم الشركة"
      type="text"
      bind:value={companySettings.name}
      placeholder="Enter company name"
      required
    />

    <!-- Company Name Arabic -->
    <FormField
      label="Company Name (Arabic)"
      labelArabic="اسم الشركة بالعربية"
      type="text"
      bind:value={companySettings.name_arabic}
      placeholder="اسم الشركة"
      dir="rtl"
    />

    <!-- VAT Number -->
    <FormField
      label="VAT Number"
      labelArabic="رقم ضريبة القيمة المضافة"
      type="text"
      bind:value={companySettings.vat_number}
      placeholder="Enter VAT number"
    />

    <!-- CR Number -->
    <FormField
      label="CR Number"
      labelArabic="رقم السجل التجاري"
      type="text"
      bind:value={companySettings.cr_number}
      placeholder="Enter CR number"
    />

    <!-- Email -->
    <FormField
      label="Email"
      labelArabic="البريد الإلكتروني"
      type="email"
      bind:value={companySettings.email}
      placeholder="Enter email address"
    />

    <!-- Phone -->
    <FormField
      label="Phone"
      labelArabic="رقم الهاتف"
      type="tel"
      bind:value={companySettings.phone}
      placeholder="Enter phone number"
    />

    <!-- Address -->
    <FormField
      label="Address"
      labelArabic="العنوان"
      type="textarea"
      bind:value={companySettings.address}
      placeholder="Enter address"
      rows={3}
    />

    <!-- Address Arabic -->
    <FormField
      label="Address (Arabic)"
      labelArabic="العنوان بالعربية"
      type="textarea"
      bind:value={companySettings.address_arabic}
      placeholder="العنوان"
      dir="rtl"
      rows={3}
    />

    <!-- City -->
    <FormField
      label="City"
      labelArabic="المدينة"
      type="text"
      bind:value={companySettings.city}
      placeholder="Enter city"
    />

    <!-- City Arabic -->
    <FormField
      label="City (Arabic)"
      labelArabic="المدينة بالعربية"
      type="text"
      bind:value={companySettings.city_arabic}
      placeholder="المدينة"
      dir="rtl"
    />

    <!-- Country -->
    <FormField
      label="Country"
      labelArabic="البلد"
      type="text"
      bind:value={companySettings.country}
      placeholder="Enter country"
    />

    <!-- Country Arabic -->
    <FormField
      label="Country (Arabic)"
      labelArabic="البلد بالعربية"
      type="text"
      bind:value={companySettings.country_arabic}
      placeholder="البلد"
      dir="rtl"
    />
  </div>

  <!-- Logo Upload -->
  <div>
    <label for="logo" class="block text-sm font-medium text-gray-300 mb-2">
      Company Logo
    </label>
    
    {#if logoPreview}
      <!-- Logo Preview -->
      <div class="mb-4">
        <div class="relative inline-block">
          <img 
            src={logoPreview} 
            alt="Company Logo Preview" 
            class="max-w-xs max-h-32 object-contain border border-gray-300 rounded-lg"
          />
          <button
            type="button"
            on:click={removeLogo}
            class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center text-xs hover:bg-red-600"
          >
            ×
          </button>
        </div>
      </div>
    {/if}

    <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
      <i class="fas fa-image mx-auto h-12 w-12 text-gray-300"></i>
      <div class="mt-4">
        <label for="logo-upload" class="cursor-pointer">
          <span class="mt-2 block text-sm font-medium text-gray-100">
            {logoPreview ? 'Change logo' : 'Upload a logo'}
          </span>
          <input 
            id="logo-upload" 
            name="logo-upload" 
            type="file" 
            class="sr-only" 
            accept=".webp,.png,.svg,.jpg,.jpeg"
            on:change={handleLogoUpload}
            bind:this={fileInput}
          />
        </label>
        <p class="mt-1 text-xs text-gray-300">WebP, PNG, SVG, JPG up to 500KB</p>
        
        {#if logoError}
          <p class="mt-2 text-xs text-red-400">{logoError}</p>
        {/if}
      </div>
    </div>
  </div>

  <!-- Save Button -->
  <div class="flex justify-end">
    <button
      on:click={saveCompanySettings}
      disabled={isLoading}
      class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
    >
      {isLoading ? 'Saving...' : 'Save Company Settings'}
    </button>
  </div>
</div>