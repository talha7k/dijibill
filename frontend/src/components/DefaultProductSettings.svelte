<script>
  import { createEventDispatcher, onMount } from 'svelte'
  import { GetDefaultProductSettings, UpdateDefaultProductSettings, GetTaxRates, GetUnitsOfMeasurement } from '../../wailsjs/go/main/App.js'
   import { database } from '../../wailsjs/go/models'
  import FormField from './FormField.svelte'

  const dispatch = createEventDispatcher()

  export let isLoading = false

  // Default product settings
  let defaultProductSettings = {
    id: 1,
    default_stock: 0,
    default_tax_rate_id: 1,
    default_unit_id: 1,
    default_product_type: 'product',
    default_markup: 0,
    default_product_status: true,
    default_price_includes_tax: false,
    default_price_change_allowed: true
  }

  // Reference data
  let taxRates = []
  let units = []

  // Computed properties for string conversion
  $: stockValue = defaultProductSettings.default_stock?.toString() || '0'
  $: markupValue = defaultProductSettings.default_markup?.toString() || '0'
  $: taxRateValue = defaultProductSettings.default_tax_rate_id?.toString() || '1'
  $: unitValue = defaultProductSettings.default_unit_id?.toString() || '1'
  $: statusValue = defaultProductSettings.default_product_status?.toString() || 'true'

  // Handle input changes
  function handleStockChange(event) {
    defaultProductSettings.default_stock = parseFloat(event.target.value) || 0
  }

  function handleMarkupChange(event) {
    defaultProductSettings.default_markup = parseFloat(event.target.value) || 0
  }

  function handleTaxRateChange(event) {
    defaultProductSettings.default_tax_rate_id = parseInt(event.target.value) || 1
  }

  function handleUnitChange(event) {
    defaultProductSettings.default_unit_id = parseInt(event.target.value) || 1
  }

  function handleStatusChange(event) {
    defaultProductSettings.default_product_status = event.target.value === 'true'
  }

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
        loadedDefaultSettings
      ] = await Promise.all([
        GetTaxRates(),
        GetUnitsOfMeasurement(),
        GetDefaultProductSettings()
      ])

      taxRates = loadedTaxRates || []
      units = loadedUnits || []

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
      const settings = new database.DefaultProductSettings()
      settings.id = defaultProductSettings.id
      settings.default_stock = defaultProductSettings.default_stock
      settings.default_tax_rate_id = defaultProductSettings.default_tax_rate_id
      settings.default_unit_id = defaultProductSettings.default_unit_id
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
    <h3 class="text-lg leading-6 font-medium text-gray-100">Default Product Settings</h3>
    <p class="mt-1 text-sm text-gray-300">
      Configure default values that will be pre-filled when creating new products
    </p>
  </div>

  <form on:submit|preventDefault={saveDefaultProductSettings} class="space-y-6">
    <div class="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-2">
      <!-- Default Stock -->
      <FormField
        label="Default Stock Quantity"
        type="number"
        value={stockValue}
        on:input={handleStockChange}
        min="0"
        step="0.01"
        placeholder="0"
      />
      <p class="text-xs text-gray-300 -mt-4">Initial stock quantity for new products</p>

      <!-- Default Tax Rate -->
      <FormField
        label="Default Tax Rate"
        type="select"
        value={taxRateValue}
        on:input={handleTaxRateChange}
        options={taxRates.map(rate => ({ value: rate.id.toString(), label: `${rate.name} (${rate.rate}%)` }))}
      />

      <!-- Default Unit -->
      <FormField
        label="Default Unit of Measurement"
        type="select"
        value={unitValue}
        on:input={handleUnitChange}
        options={units.map(unit => ({ value: unit.id.toString(), label: `${unit.label} (${unit.value})` }))}
      />

      <!-- Default Product Type -->
      <FormField
        label="Default Product Type"
        type="select"
        bind:value={defaultProductSettings.default_product_type}
        options={[
          { value: 'product', label: 'Product' },
          { value: 'service', label: 'Service' }
        ]}
      />

      <!-- Default Markup -->
      <FormField
        label="Default Markup (%)"
        type="number"
        value={markupValue}
        on:input={handleMarkupChange}
        min="0"
        step="0.01"
        placeholder="0"
      />
      <p class="text-xs text-gray-300 -mt-4">Default markup percentage for pricing</p>

      <!-- Default Status -->
      <FormField
        label="Default Product Status"
        type="select"
        value={statusValue}
        on:input={handleStatusChange}
        options={[
          { value: 'true', label: 'Active' },
          { value: 'false', label: 'Inactive' }
        ]}
      />
    </div>

    <!-- Checkboxes -->
    <div class="space-y-4">
      <FormField
        label="Price includes tax by default"
        type="checkbox"
        bind:checked={defaultProductSettings.default_price_includes_tax}
      />
      <p class="text-xs text-gray-300 -mt-2">When enabled, product prices will include tax by default</p>

      <FormField
        label="Allow price changes by default"
        type="checkbox"
        bind:checked={defaultProductSettings.default_price_change_allowed}
      />
      <p class="text-xs text-gray-300 -mt-2">When enabled, product prices can be modified during sales</p>
    </div>

    <!-- Save Button -->
    <div class="flex justify-end">
      <button
        type="submit"
        disabled={isLoading}
        class="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {#if isLoading}
          <i class="fas fa-spinner fa-spin -ml-1 mr-3 h-5 w-5 text-white"></i>
          Saving...
        {:else}
          Save Default Product Settings
        {/if}
      </button>
    </div>
  </form>
</div>