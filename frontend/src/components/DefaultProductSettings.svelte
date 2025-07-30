<script>
  import { createEventDispatcher, onMount } from 'svelte'
  import { GetDefaultProductSettings, UpdateDefaultProductSettings, GetTaxRates, GetUnitsOfMeasurement, GetPaymentTypes, GetProductCategories } from '../../wailsjs/go/main/App.js'
  import * as main from '../../wailsjs/go/models'

  const dispatch = createEventDispatcher()

  export let isLoading = false

  // Default product settings
  let defaultProductSettings = {
    id: 1,
    default_stock: 0,
    default_tax_rate_id: 1,
    default_unit_id: 1,
    default_payment_type_id: 1,
    default_product_category_id: 1,
    default_product_type: 'product',
    default_markup: 0,
    default_product_status: true,
    default_price_includes_tax: false,
    default_price_change_allowed: true
  }

  // Reference data
  let taxRates = []
  let units = []
  let paymentTypes = []
  let productCategories = []

  onMount(async () => {
    await loadData()
  })

  async function loadData() {
    isLoading = true
    try {
      // Load reference data
      const [
        loadedTaxRates,
        loadedUnits,
        loadedPaymentTypes,
        loadedProductCategories,
        loadedDefaultSettings
      ] = await Promise.all([
        GetTaxRates(),
        GetUnitsOfMeasurement(),
        GetPaymentTypes(),
        GetProductCategories(),
        GetDefaultProductSettings()
      ])

      taxRates = loadedTaxRates || []
      units = loadedUnits || []
      paymentTypes = loadedPaymentTypes || []
      productCategories = loadedProductCategories || []

      if (loadedDefaultSettings) {
        defaultProductSettings = loadedDefaultSettings
      }
    } catch (error) {
      console.error('Error loading default product settings:', error)
      dispatch('error', { message: 'Failed to load default product settings' })
    } finally {
      isLoading = false
    }
  }

  async function saveDefaultProductSettings() {
    isLoading = true
    try {
      // Create a new DefaultProductSettings object
      const settings = new main.main.DefaultProductSettings()
      settings.id = defaultProductSettings.id
      settings.default_stock = defaultProductSettings.default_stock
      settings.default_tax_rate_id = defaultProductSettings.default_tax_rate_id
      settings.default_unit_id = defaultProductSettings.default_unit_id
      settings.default_payment_type_id = defaultProductSettings.default_payment_type_id
      settings.default_product_category_id = defaultProductSettings.default_product_category_id
      settings.default_product_type = defaultProductSettings.default_product_type
      settings.default_markup = defaultProductSettings.default_markup
      settings.default_product_status = defaultProductSettings.default_product_status
      settings.default_price_includes_tax = defaultProductSettings.default_price_includes_tax
      settings.default_price_change_allowed = defaultProductSettings.default_price_change_allowed

      await UpdateDefaultProductSettings(settings)
      dispatch('saved', { message: 'Default product settings saved successfully!' })
    } catch (error) {
      console.error('Error saving default product settings:', error)
      dispatch('error', { message: 'Failed to save default product settings' })
    } finally {
      isLoading = false
    }
  }
</script>

<div class="space-y-6">
  <div>
    <h3 class="text-lg leading-6 font-medium text-gray-900">Default Product Settings</h3>
    <p class="mt-1 text-sm text-gray-500">
      Configure default values that will be pre-filled when creating new products
    </p>
  </div>

  <form on:submit|preventDefault={saveDefaultProductSettings} class="space-y-6">
    <div class="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-2">
      <!-- Default Stock -->
      <div>
        <label for="default_stock" class="block text-sm font-medium text-gray-700">
          Default Stock Quantity
        </label>
        <div class="mt-1">
          <input
            type="number"
            id="default_stock"
            bind:value={defaultProductSettings.default_stock}
            min="0"
            step="0.01"
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
            placeholder="0"
          />
        </div>
        <p class="mt-1 text-xs text-gray-500">Initial stock quantity for new products</p>
      </div>

      <!-- Default Tax Rate -->
      <div>
        <label for="default_tax_rate" class="block text-sm font-medium text-gray-700">
          Default Tax Rate
        </label>
        <div class="mt-1">
          <select
            id="default_tax_rate"
            bind:value={defaultProductSettings.default_tax_rate_id}
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
          >
            {#each taxRates as taxRate}
              <option value={taxRate.id}>{taxRate.name} ({taxRate.rate}%)</option>
            {/each}
          </select>
        </div>
      </div>

      <!-- Default Unit -->
      <div>
        <label for="default_unit" class="block text-sm font-medium text-gray-700">
          Default Unit of Measurement
        </label>
        <div class="mt-1">
          <select
            id="default_unit"
            bind:value={defaultProductSettings.default_unit_id}
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
          >
            {#each units as unit}
              <option value={unit.id}>{unit.name} ({unit.symbol})</option>
            {/each}
          </select>
        </div>
      </div>

      <!-- Default Payment Type -->
      <div>
        <label for="default_payment_type" class="block text-sm font-medium text-gray-700">
          Default Payment Type
        </label>
        <div class="mt-1">
          <select
            id="default_payment_type"
            bind:value={defaultProductSettings.default_payment_type_id}
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
          >
            {#each paymentTypes as paymentType}
              <option value={paymentType.id}>{paymentType.name}</option>
            {/each}
          </select>
        </div>
      </div>

      <!-- Default Product Category -->
      <div>
        <label for="default_product_category" class="block text-sm font-medium text-gray-700">
          Default Product Category
        </label>
        <div class="mt-1">
          <select
            id="default_product_category"
            bind:value={defaultProductSettings.default_product_category_id}
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
          >
            {#each productCategories as category}
              <option value={category.id}>{category.name}</option>
            {/each}
          </select>
        </div>
      </div>

      <!-- Default Product Type -->
      <div>
        <label for="default_product_type" class="block text-sm font-medium text-gray-700">
          Default Product Type
        </label>
        <div class="mt-1">
          <select
            id="default_product_type"
            bind:value={defaultProductSettings.default_product_type}
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
          >
            <option value="product">Product</option>
            <option value="service">Service</option>
          </select>
        </div>
      </div>

      <!-- Default Markup -->
      <div>
        <label for="default_markup" class="block text-sm font-medium text-gray-700">
          Default Markup (%)
        </label>
        <div class="mt-1">
          <input
            type="number"
            id="default_markup"
            bind:value={defaultProductSettings.default_markup}
            min="0"
            step="0.01"
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
            placeholder="0"
          />
        </div>
        <p class="mt-1 text-xs text-gray-500">Default markup percentage for pricing</p>
      </div>

      <!-- Default Status -->
      <div>
        <label for="default_status" class="block text-sm font-medium text-gray-700">
          Default Product Status
        </label>
        <div class="mt-1">
          <select
            id="default_status"
            bind:value={defaultProductSettings.default_product_status}
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
          >
            <option value={true}>Active</option>
            <option value={false}>Inactive</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Checkboxes -->
    <div class="space-y-4">
      <div class="flex items-start">
        <div class="flex items-center h-5">
          <input
            id="price_includes_tax"
            type="checkbox"
            bind:checked={defaultProductSettings.default_price_includes_tax}
            class="focus:ring-blue-500 h-4 w-4 text-blue-600 border-gray-300 rounded"
          />
        </div>
        <div class="ml-3 text-sm">
          <label for="price_includes_tax" class="font-medium text-gray-700">
            Price includes tax by default
          </label>
          <p class="text-gray-500">When enabled, product prices will include tax by default</p>
        </div>
      </div>

      <div class="flex items-start">
        <div class="flex items-center h-5">
          <input
            id="allow_price_change"
            type="checkbox"
            bind:checked={defaultProductSettings.default_price_change_allowed}
            class="focus:ring-blue-500 h-4 w-4 text-blue-600 border-gray-300 rounded"
          />
        </div>
        <div class="ml-3 text-sm">
          <label for="allow_price_change" class="font-medium text-gray-700">
            Allow price changes by default
          </label>
          <p class="text-gray-500">When enabled, product prices can be modified during sales</p>
        </div>
      </div>
    </div>

    <!-- Save Button -->
    <div class="flex justify-end">
      <button
        type="submit"
        disabled={isLoading}
        class="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {#if isLoading}
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Saving...
        {:else}
          Save Default Product Settings
        {/if}
      </button>
    </div>
  </form>
</div>