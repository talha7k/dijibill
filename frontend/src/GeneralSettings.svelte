<script>
  import { onMount } from 'svelte'
  
  // Import backend functions
  import { 
    GetCompany, 
    GetTaxRates, 
    GetPaymentTypes, 
    GetSalesCategories, 
    GetUnitsOfMeasurement,
    GetSystemSettings
  } from '../wailsjs/go/main/App.js'

  // Import components
  import PageLayout from './components/PageLayout.svelte'
  import CompanySettings from './components/CompanySettings.svelte'
  import TaxRateSettings from './components/TaxRateSettings.svelte'
  import PaymentTypeSettings from './components/PaymentTypeSettings.svelte'
  import SalesCategorySettings from './components/SalesCategorySettings.svelte'
  import UnitSettings from './components/UnitSettings.svelte'
  import SystemSettings from './components/SystemSettings.svelte'
  import DefaultProductSettings from './components/DefaultProductSettings.svelte'
  import HorizontalTabs from './components/HorizontalTabs.svelte'

  // State variables
  let activeTab = 'company'
  let isLoading = false

  // Handle tab change
  function handleTabChange(event) {
    activeTab = event.detail.tabId
  }

  // Company settings
  let companySettings = {
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
    logo: '',
    created_at: '',
    updated_at: ''
  }

  // System settings
  let systemSettings = {
    currency: 'SAR',
    language: 'english',
    timezone: 'Asia/Riyadh',
    date_format: 'DD/MM/YYYY',
    invoice_language: 'english',
    zatca_enabled: false,
    auto_backup: true,
    backup_frequency: 'daily'
  }

  // Data arrays
  let taxRates = []
  let paymentTypes = []
  let salesCategories = []
  let units = []

  onMount(() => {
    loadSettings()
  })

  async function loadSettings() {
    isLoading = true
    try {
      // Load company settings
      const company = await GetCompany()
      if (company) {
        companySettings = company
      }

      // Load system settings
      try {
        const settings = await GetSystemSettings()
        if (settings) {
          systemSettings = {
            currency: settings.currency,
            language: settings.language,
            timezone: settings.timezone,
            date_format: settings.date_format,
            invoice_language: settings.invoice_language,
            zatca_enabled: settings.zatca_enabled,
            auto_backup: settings.auto_backup,
            backup_frequency: settings.backup_frequency
          }
        }
      } catch (systemError) {
        console.error('Error loading system settings:', systemError)
        // Keep default system settings if loading fails
      }

      // Load all settings data from backend
      const [
        loadedTaxRates,
        loadedPaymentTypes, 
        loadedSalesCategories,
        loadedUnits
      ] = await Promise.all([
        GetTaxRates(),
        GetPaymentTypes(),
        GetSalesCategories(),
        GetUnitsOfMeasurement()
      ])

      taxRates = loadedTaxRates || []
      paymentTypes = loadedPaymentTypes || []
      salesCategories = loadedSalesCategories || []
      units = loadedUnits || []

      console.log('Settings loaded successfully')
    } catch (error) {
      console.error('Error loading settings:', error)
      // Keep the default mock data if loading fails
    } finally {
      isLoading = false
    }
  }

  // Event handlers
  function handleCompanySaved(event) {
    alert(event.detail.message)
  }

  function handleError(event) {
    alert(event.detail.message)
  }

  function handleReload() {
    loadSettings()
  }

  function handleSystemSave(event) {
    systemSettings = event.detail.systemSettings
    if (event.detail.message) {
      alert(event.detail.message)
    }
  }

  function handleSystemError(event) {
    alert(event.detail.message)
  }

  // Tab configuration
  const tabs = [
    { id: 'company', name: 'Company', icon: '🏢' },
    { id: 'system', name: 'System', icon: '⚙️' },
    { id: 'default-products', name: 'Default Products', icon: '📦' },
    { id: 'tax-rates', name: 'Tax Rates', icon: '💰' },
    { id: 'payment-types', name: 'Payment Types', icon: '💳' },
    { id: 'sales-categories', name: 'Sales Categories', icon: '📊' },
    { id: 'units', name: 'Units', icon: '📏' }
  ]
</script>

<PageLayout 
  title="Settings"
  icon="settings"
>
  <div slot="actions">
    {#if isLoading}
      <div class="flex items-center text-sm text-gray-300">
        <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-gray-300" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        Loading...
      </div>
    {/if}
  </div>
  <!-- Horizontal Tabs -->
  <HorizontalTabs 
    {tabs}
    {activeTab}
    variant="standard"
    showScrollButtons={true}
    on:tabChange={handleTabChange}
  />

  <!-- Main content -->
  <div class="mt-6">
    <div class="backdrop-blur-sm shadow-lg rounded-xl border border-white/20">
      <div class="px-6 py-8">
        {#if activeTab === 'company'}
          <CompanySettings 
            bind:companySettings={companySettings}
            bind:isLoading={isLoading}
            on:saved={handleCompanySaved}
            on:error={handleError}
          />
        {:else if activeTab === 'system'}
          <SystemSettings 
            bind:systemSettings={systemSettings}
            bind:isLoading={isLoading}
            on:save={handleSystemSave}
            on:error={handleSystemError}
          />
        {:else if activeTab === 'default-products'}
          <DefaultProductSettings />
        {:else if activeTab === 'tax-rates'}
          <TaxRateSettings 
            bind:taxRates={taxRates}
            bind:isLoading={isLoading}
            on:reload={handleReload}
            on:error={handleError}
          />
        {:else if activeTab === 'payment-types'}
          <PaymentTypeSettings 
            bind:paymentTypes={paymentTypes}
            bind:isLoading={isLoading}
            on:reload={handleReload}
            on:error={handleError}
          />
        {:else if activeTab === 'sales-categories'}
          <SalesCategorySettings 
            bind:salesCategories={salesCategories}
            bind:isLoading={isLoading}
            on:reload={handleReload}
            on:error={handleError}
          />
        {:else if activeTab === 'units'}
          <UnitSettings 
            bind:units={units}
            bind:isLoading={isLoading}
            on:reload={handleReload}
            on:error={handleError}
          />
        {/if}
      </div>
    </div>
  </div>
</PageLayout>