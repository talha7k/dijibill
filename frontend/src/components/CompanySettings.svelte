<script>
  import { GetCompany, UpdateCompany } from '../../wailsjs/go/main/App.js'
  import { createEventDispatcher } from 'svelte'

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
    <div>
      <label for="company-name" class="block text-sm font-medium text-gray-700 mb-2">
        Company Name
      </label>
      <input
        id="company-name"
        type="text"
        bind:value={companySettings.name}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="Enter company name"
      />
    </div>

    <!-- Company Name Arabic -->
    <div>
      <label for="company-name-arabic" class="block text-sm font-medium text-gray-700 mb-2">
        Company Name (Arabic)
      </label>
      <input
        id="company-name-arabic"
        type="text"
        bind:value={companySettings.name_arabic}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="اسم الشركة"
        dir="rtl"
      />
    </div>

    <!-- VAT Number -->
    <div>
      <label for="vat-number" class="block text-sm font-medium text-gray-700 mb-2">
        VAT Number
      </label>
      <input
        id="vat-number"
        type="text"
        bind:value={companySettings.vat_number}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="Enter VAT number"
      />
    </div>

    <!-- CR Number -->
    <div>
      <label for="cr-number" class="block text-sm font-medium text-gray-700 mb-2">
        CR Number
      </label>
      <input
        id="cr-number"
        type="text"
        bind:value={companySettings.cr_number}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="Enter CR number"
      />
    </div>

    <!-- Email -->
    <div>
      <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
        Email
      </label>
      <input
        id="email"
        type="email"
        bind:value={companySettings.email}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="Enter email address"
      />
    </div>

    <!-- Phone -->
    <div>
      <label for="phone" class="block text-sm font-medium text-gray-700 mb-2">
        Phone
      </label>
      <input
        id="phone"
        type="tel"
        bind:value={companySettings.phone}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="Enter phone number"
      />
    </div>

    <!-- Address -->
    <div>
      <label for="address" class="block text-sm font-medium text-gray-700 mb-2">
        Address
      </label>
      <textarea
        id="address"
        bind:value={companySettings.address}
        rows="3"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="Enter address"
      ></textarea>
    </div>

    <!-- Address Arabic -->
    <div>
      <label for="address-arabic" class="block text-sm font-medium text-gray-700 mb-2">
        Address (Arabic)
      </label>
      <textarea
        id="address-arabic"
        bind:value={companySettings.address_arabic}
        rows="3"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="العنوان"
        dir="rtl"
      ></textarea>
    </div>

    <!-- City -->
    <div>
      <label for="city" class="block text-sm font-medium text-gray-700 mb-2">
        City
      </label>
      <input
        id="city"
        type="text"
        bind:value={companySettings.city}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="Enter city"
      />
    </div>

    <!-- City Arabic -->
    <div>
      <label for="city-arabic" class="block text-sm font-medium text-gray-700 mb-2">
        City (Arabic)
      </label>
      <input
        id="city-arabic"
        type="text"
        bind:value={companySettings.city_arabic}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="المدينة"
        dir="rtl"
      />
    </div>

    <!-- Country -->
    <div>
      <label for="country" class="block text-sm font-medium text-gray-700 mb-2">
        Country
      </label>
      <input
        id="country"
        type="text"
        bind:value={companySettings.country}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="Enter country"
      />
    </div>

    <!-- Country Arabic -->
    <div>
      <label for="country-arabic" class="block text-sm font-medium text-gray-700 mb-2">
        Country (Arabic)
      </label>
      <input
        id="country-arabic"
        type="text"
        bind:value={companySettings.country_arabic}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="البلد"
        dir="rtl"
      />
    </div>
  </div>

  <!-- Logo Upload -->
  <div>
    <label for="logo" class="block text-sm font-medium text-gray-700 mb-2">
      Company Logo
    </label>
    <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
      <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48">
        <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
      </svg>
      <div class="mt-4">
        <label for="logo-upload" class="cursor-pointer">
          <span class="mt-2 block text-sm font-medium text-gray-900">
            Upload a logo
          </span>
          <input id="logo-upload" name="logo-upload" type="file" class="sr-only" accept="image/*" />
        </label>
        <p class="mt-1 text-xs text-gray-500">PNG, JPG, GIF up to 10MB</p>
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