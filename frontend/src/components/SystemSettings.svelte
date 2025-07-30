<script>
  import { createEventDispatcher, onMount } from 'svelte'

  const dispatch = createEventDispatcher()

  export let systemSettings = {
    currency: 'SAR',
    language: 'en',
    timezone: 'Asia/Riyadh',
    date_format: 'DD/MM/YYYY',
    invoice_language: 'english', // Default language for invoices/printouts
    zatca_enabled: true,
    auto_backup: true,
    backup_frequency: 'daily'
  }

  export let isLoading = false

  onMount(() => {
    // Load existing settings from localStorage
    const savedSettings = localStorage.getItem('systemSettings')
    if (savedSettings) {
      try {
        const parsed = JSON.parse(savedSettings)
        systemSettings = { ...systemSettings, ...parsed }
      } catch (error) {
        console.error('Error loading system settings from localStorage:', error)
      }
    }
  })

  function saveSystemSettings() {
    // Save to localStorage for immediate access by other components
    localStorage.setItem('systemSettings', JSON.stringify(systemSettings))
    
    dispatch('save', { systemSettings })
  }
</script>

<div class="space-y-6">
  <div>
    <h3 class="text-lg leading-6 font-medium text-gray-900">System Settings</h3>
    <p class="mt-1 text-sm text-gray-500">
      Configure general system preferences and features
    </p>
  </div>

  <form on:submit|preventDefault={saveSystemSettings} class="space-y-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Currency -->
      <div>
        <label for="currency" class="block text-sm font-medium text-gray-700">
          Currency
        </label>
        <div class="mt-1">
          <select
            id="currency"
            bind:value={systemSettings.currency}
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
          >
            <option value="SAR">Saudi Riyal (SAR)</option>
            <option value="USD">US Dollar (USD)</option>
            <option value="EUR">Euro (EUR)</option>
            <option value="AED">UAE Dirham (AED)</option>
          </select>
        </div>
      </div>

      <!-- Language -->
      <div>
        <label for="language" class="block text-sm font-medium text-gray-700">
          Language
        </label>
        <div class="mt-1">
          <select
            id="language"
            bind:value={systemSettings.language}
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
          >
            <option value="en">English</option>
            <option value="ar">Arabic</option>
          </select>
        </div>
      </div>

      <!-- Timezone -->
      <div>
        <label for="timezone" class="block text-sm font-medium text-gray-700">
          Timezone
        </label>
        <div class="mt-1">
          <select
            id="timezone"
            bind:value={systemSettings.timezone}
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
          >
            <option value="Asia/Riyadh">Asia/Riyadh</option>
            <option value="Asia/Dubai">Asia/Dubai</option>
            <option value="UTC">UTC</option>
          </select>
        </div>
      </div>

      <!-- Date Format -->
      <div>
        <label for="date-format" class="block text-sm font-medium text-gray-700">
          Date Format
        </label>
        <div class="mt-1">
          <select
            id="date-format"
            bind:value={systemSettings.date_format}
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
          >
            <option value="DD/MM/YYYY">DD/MM/YYYY</option>
            <option value="MM/DD/YYYY">MM/DD/YYYY</option>
            <option value="YYYY-MM-DD">YYYY-MM-DD</option>
          </select>
        </div>
      </div>

      <!-- Invoice Language -->
      <div>
        <label for="invoice-language" class="block text-sm font-medium text-gray-700">
          Invoice Language
        </label>
        <div class="mt-1">
          <select
            id="invoice-language"
            bind:value={systemSettings.invoice_language}
            class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
          >
            <option value="english">English</option>
            <option value="arabic">Arabic</option>
            <option value="bilingual">Bilingual (Arabic + English)</option>
          </select>
        </div>
        <p class="mt-1 text-xs text-gray-500">
          Default language for all invoice printouts and PDFs
        </p>
      </div>
    </div>

    <!-- System Features -->
    <div class="space-y-4">
      <h4 class="text-md font-medium text-gray-900">System Features</h4>
      
      <div class="space-y-4">
        <div class="flex items-start">
          <div class="flex items-center h-5">
            <input
              id="zatca_enabled"
              type="checkbox"
              bind:checked={systemSettings.zatca_enabled}
              class="focus:ring-blue-500 h-4 w-4 text-blue-600 border-gray-300 rounded"
            />
          </div>
          <div class="ml-3 text-sm">
            <label for="zatca_enabled" class="font-medium text-gray-700">
              Enable ZATCA Integration
            </label>
            <p class="text-gray-500">Enable Saudi Arabia's ZATCA e-invoicing compliance</p>
          </div>
        </div>

        <div class="flex items-start">
          <div class="flex items-center h-5">
            <input
              id="auto_backup"
              type="checkbox"
              bind:checked={systemSettings.auto_backup}
              class="focus:ring-blue-500 h-4 w-4 text-blue-600 border-gray-300 rounded"
            />
          </div>
          <div class="ml-3 text-sm">
            <label for="auto_backup" class="font-medium text-gray-700">
              Enable Auto Backup
            </label>
            <p class="text-gray-500">Automatically backup your data at regular intervals</p>
          </div>
        </div>

        {#if systemSettings.auto_backup}
          <div class="ml-7">
            <label for="backup_frequency" class="block text-sm font-medium text-gray-700">
              Backup Frequency
            </label>
            <div class="mt-1">
              <select
                id="backup_frequency"
                bind:value={systemSettings.backup_frequency}
                class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 rounded-md"
              >
                <option value="daily">Daily</option>
                <option value="weekly">Weekly</option>
                <option value="monthly">Monthly</option>
              </select>
            </div>
          </div>
        {/if}
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
          Save System Settings
        {/if}
      </button>
    </div>
  </form>
</div>