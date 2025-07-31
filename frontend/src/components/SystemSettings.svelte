<script>
  import { createEventDispatcher, onMount } from 'svelte'
  import { GetSystemSettings, UpdateSystemSettings } from '../../wailsjs/go/main/App.js'
  import { database } from '../../wailsjs/go/models'
  import FormField from './FormField.svelte'

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

  onMount(async () => {
    try {
      isLoading = true
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
    } catch (error) {
      console.error('Error loading system settings:', error)
      // Fallback to localStorage for backward compatibility
      const savedSettings = localStorage.getItem('systemSettings')
      if (savedSettings) {
        try {
          const parsed = JSON.parse(savedSettings)
          systemSettings = { ...systemSettings, ...parsed }
        } catch (parseError) {
          console.error('Error parsing localStorage settings:', parseError)
        }
      }
    } finally {
      isLoading = false
    }
  })

  async function saveSystemSettings() {
    try {
      isLoading = true
      
      const settings = database.SystemSettings.createFrom({
        id: 1, // Assuming single system settings record
        currency: systemSettings.currency,
        language: systemSettings.language,
        timezone: systemSettings.timezone,
        date_format: systemSettings.date_format,
        invoice_language: systemSettings.invoice_language,
        zatca_enabled: systemSettings.zatca_enabled,
        auto_backup: systemSettings.auto_backup,
        backup_frequency: systemSettings.backup_frequency
      })

      await UpdateSystemSettings(settings)
      
      // Also save to localStorage for immediate access by other components
      localStorage.setItem('systemSettings', JSON.stringify(systemSettings))
      
      dispatch('save', { 
        systemSettings,
        message: 'System settings saved successfully!'
      })
    } catch (error) {
      console.error('Error saving system settings:', error)
      dispatch('error', { 
        message: 'Failed to save system settings. Please try again.'
      })
    } finally {
      isLoading = false
    }
  }
</script>

<div class="space-y-6">
  <div>
    <h3 class="text-lg leading-6 font-medium text-gray-100">System Settings</h3>
    <p class="mt-1 text-sm text-gray-300">
      Configure general system preferences and features
    </p>
  </div>

  <form on:submit|preventDefault={saveSystemSettings} class="space-y-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Currency -->
      <FormField
        label="Currency"
        labelArabic="العملة"
        type="select"
        bind:value={systemSettings.currency}
        options={[
          { value: 'SAR', label: 'Saudi Riyal (SAR)' },
          { value: 'USD', label: 'US Dollar (USD)' },
          { value: 'EUR', label: 'Euro (EUR)' },
          { value: 'AED', label: 'UAE Dirham (AED)' }
        ]}
      />

      <!-- Language -->
      <FormField
        label="Language"
        labelArabic="اللغة"
        type="select"
        bind:value={systemSettings.language}
        options={[
          { value: 'en', label: 'English' },
          { value: 'ar', label: 'Arabic' }
        ]}
      />

      <!-- Timezone -->
      <FormField
        label="Timezone"
        labelArabic="المنطقة الزمنية"
        type="select"
        bind:value={systemSettings.timezone}
        options={[
          { value: 'Asia/Riyadh', label: 'Asia/Riyadh' },
          { value: 'Asia/Dubai', label: 'Asia/Dubai' },
          { value: 'UTC', label: 'UTC' }
        ]}
      />

      <!-- Date Format -->
      <FormField
        label="Date Format"
        labelArabic="تنسيق التاريخ"
        type="select"
        bind:value={systemSettings.date_format}
        options={[
          { value: 'DD/MM/YYYY', label: 'DD/MM/YYYY' },
          { value: 'MM/DD/YYYY', label: 'MM/DD/YYYY' },
          { value: 'YYYY-MM-DD', label: 'YYYY-MM-DD' }
        ]}
      />

      <!-- Invoice Language -->
      <div class="col-span-2">
        <FormField
          label="Invoice Language"
          labelArabic="لغة الفاتورة"
          type="select"
          bind:value={systemSettings.invoice_language}
          options={[
            { value: 'english', label: 'English' },
            { value: 'arabic', label: 'Arabic' },
            { value: 'bilingual', label: 'Bilingual (Arabic + English)' }
          ]}
        />
        <p class="mt-1 text-xs text-gray-300">
          Default language for all invoice printouts and PDFs
        </p>
      </div>
    </div>

    <!-- System Features -->
    <div class="space-y-4">
      <h4 class="text-md font-medium text-gray-100">System Features</h4>
      
      <div class="space-y-4">
        <div class="flex items-start">
          <div class="flex items-center h-5">
            <FormField
              type="checkbox"
              bind:checked={systemSettings.zatca_enabled}
              placeholder="Enable ZATCA Integration"
            />
          </div>
          <div class="ml-3 text-sm">
            <label class="font-medium text-gray-300">
              Enable ZATCA Integration
            </label>
            <p class="text-gray-300">Enable Saudi Arabia's ZATCA e-invoicing compliance</p>
          </div>
        </div>

        <div class="flex items-start">
          <div class="flex items-center h-5">
            <FormField
              type="checkbox"
              bind:checked={systemSettings.auto_backup}
              placeholder="Enable Auto Backup"
            />
          </div>
          <div class="ml-3 text-sm">
            <label class="font-medium text-gray-300">
              Enable Auto Backup
            </label>
            <p class="text-gray-300">Automatically backup your data at regular intervals</p>
          </div>
        </div>

        {#if systemSettings.auto_backup}
          <div class="ml-7">
            <FormField
              label="Backup Frequency"
              labelArabic="تكرار النسخ الاحتياطي"
              type="select"
              bind:value={systemSettings.backup_frequency}
              options={[
                { value: 'daily', label: 'Daily' },
                { value: 'weekly', label: 'Weekly' },
                { value: 'monthly', label: 'Monthly' }
              ]}
            />
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
          <i class="fas fa-spinner fa-spin -ml-1 mr-3 h-5 w-5 text-white"></i>
          Saving...
        {:else}
          Save System Settings
        {/if}
      </button>
    </div>
  </form>
</div>