<script>
  import { GetCompany, UpdateCompany } from '../../wailsjs/go/main/App.js'
  import { createEventDispatcher } from 'svelte'
  import FormField from './FormField.svelte'

  const dispatch = createEventDispatcher()

  export let companySettings = {
    id: 1,
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
    country: 'Saudi Arabia',
    country_arabic: 'المملكة العربية السعودية',
    logo: ''
  }

  export let isLoading = false

  async function saveCompanySettings() {
    isLoading = true
    try {
      await UpdateCompany(companySettings)
      dispatch('saved', { message: 'Company settings saved successfully!' })
    } catch (error) {
      console.error('Error saving company settings:', error)
      dispatch('error', { message: 'Error saving company settings: ' + error.message })
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
    <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
      <svg class="mx-auto h-12 w-12 text-gray-300" stroke="currentColor" fill="none" viewBox="0 0 48 48">
        <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
      </svg>
      <div class="mt-4">
        <label for="logo-upload" class="cursor-pointer">
          <span class="mt-2 block text-sm font-medium text-gray-100">
            Upload a logo
          </span>
          <input id="logo-upload" name="logo-upload" type="file" class="sr-only" accept="image/*" />
        </label>
        <p class="mt-1 text-xs text-gray-300">PNG, JPG, GIF up to 10MB</p>
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