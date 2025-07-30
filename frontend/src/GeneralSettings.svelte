<script>
  import { onMount } from 'svelte'
  
  // Import backend functions
  import { 
    GetCompany, 
    GetTaxRates, 
    GetPaymentTypes, 
    GetSalesCategories, 
    GetUnitsOfMeasurement 
  } from '../wailsjs/go/main/App.js'

  // Import components
  import CompanySettings from './components/CompanySettings.svelte'
  import TaxRateSettings from './components/TaxRateSettings.svelte'
  import PaymentTypeSettings from './components/PaymentTypeSettings.svelte'
  import SalesCategorySettings from './components/SalesCategorySettings.svelte'
  import UnitSettings from './components/UnitSettings.svelte'
  import SystemSettings from './components/SystemSettings.svelte'
  import DefaultProductSettings from './components/DefaultProductSettings.svelte'

  // State variables
  let activeTab = 'company'
  let isLoading = false
  let tabsContainer

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
    logo: ''
  }

  // System settings
  let systemSettings = {
    currency: 'SAR',
    language: 'en',
    timezone: 'Asia/Riyadh',
    date_format: 'DD/MM/YYYY',
    zatca_enabled: true,
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
    alert('System settings saved successfully!')
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

  // Tab scrolling functions
  function scrollLeft() {
    if (tabsContainer) {
      tabsContainer.scrollBy({ left: -200, behavior: 'smooth' })
    }
  }

  function scrollRight() {
    if (tabsContainer) {
      tabsContainer.scrollBy({ left: 200, behavior: 'smooth' })
    }
  }

  // Check if scrolling is needed
  let canScrollLeft = false
  let canScrollRight = false

  function updateScrollButtons() {
    if (tabsContainer) {
      canScrollLeft = tabsContainer.scrollLeft > 0
      canScrollRight = tabsContainer.scrollLeft < (tabsContainer.scrollWidth - tabsContainer.clientWidth)
    }
  }

  onMount(() => {
    if (tabsContainer) {
      tabsContainer.addEventListener('scroll', updateScrollButtons)
      // Initial check
      setTimeout(updateScrollButtons, 100)
    }
  })
</script>

<div class="min-h-screen bg-gray-50">
  <!-- Header -->
  <div class="bg-white shadow">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">General Settings</h1>
            <p class="mt-1 text-sm text-gray-500">
              Manage your company information, system preferences, and business settings
            </p>
          </div>
          {#if isLoading}
            <div class="flex items-center text-sm text-gray-500">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Loading...
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>

  <!-- Horizontal Tabs -->
  <div class="bg-white border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="relative flex items-center">
        <!-- Left scroll button -->
        <button
          on:click={scrollLeft}
          class="flex-shrink-0 p-2 text-gray-400 hover:text-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={!canScrollLeft}
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>

        <!-- Scrollable tabs container -->
        <div 
          bind:this={tabsContainer}
          class="flex-1 overflow-x-auto scrollbar-hide"
          style="scrollbar-width: none; -ms-overflow-style: none;"
          on:scroll={updateScrollButtons}
        >
          <nav class="flex space-x-8 py-4" style="min-width: max-content;">
            {#each tabs as tab}
              <button
                on:click={() => activeTab = tab.id}
                class="flex items-center px-3 py-2 text-sm font-medium rounded-md whitespace-nowrap transition-colors duration-200 {activeTab === tab.id 
                  ? 'bg-blue-100 text-blue-700 border-b-2 border-blue-500' 
                  : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'}"
              >
                <span class="mr-2 text-lg">{tab.icon}</span>
                {tab.name}
              </button>
            {/each}
          </nav>
        </div>

        <!-- Right scroll button -->
        <button
          on:click={scrollRight}
          class="flex-shrink-0 p-2 text-gray-400 hover:text-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={!canScrollRight}
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>
    </div>
  </div>

  <!-- Main content -->
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="bg-white shadow sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6">
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
</div>

<style>
  /* Hide scrollbar for Chrome, Safari and Opera */
  .scrollbar-hide::-webkit-scrollbar {
    display: none;
  }

  /* Hide scrollbar for IE, Edge and Firefox */
  .scrollbar-hide {
    -ms-overflow-style: none;  /* IE and Edge */
    scrollbar-width: none;  /* Firefox */
  }
</style>