<script>
  import { onMount } from 'svelte'
  import { wailsReady, wailsError, initializeWailsRuntime } from './stores/wailsStore.js'
  import { GetSystemSettings } from '../wailsjs/go/main/App'
  import Dashboard from './Dashboard.svelte'
  import QRValidation from './QRValidation.svelte'
  import SalesInvoices from './SalesInvoices.svelte'
  import POS from './POS.svelte'
  import Customer from './Customer.svelte'
  import Payments from './Payments.svelte'
  import PurchaseInvoices from './PurchaseInvoices.svelte'
  import PurchaseProducts from './PurchaseProducts.svelte'
  import Suppliers from './Suppliers.svelte'
  import Products from './Products.svelte'
  import Users from './Users.svelte'
  import GeneralSettings from './GeneralSettings.svelte'
  import Administration from './Administration.svelte'
  import {Greet} from '../wailsjs/go/main/App.js'

  let currentView = 'dashboard' // Default to dashboard
  let resultText = "Welcome to DijiBill - ZATCA Compliant Invoicing"
  let sidebarVisible = true // Sidebar visibility state
  let sidebarAsOverlay = false // Whether sidebar should appear as overlay
  
  // Get runtime status for display
  $: runtimeStatus = $wailsReady ? 'Runtime Ready' : $wailsError ? 'Runtime Error' : 'Initializing...'
  let name = ''
  let backupStatus = 'unknown' // 'recent', 'overdue', 'unknown'
  let lastBackupTime = null
  let autoBackupEnabled = false

  // Check backup status
  async function checkBackupStatus() {
    try {
      const settings = await GetSystemSettings()
      if (settings) {
        autoBackupEnabled = settings.auto_backup
        lastBackupTime = settings.last_backup_time
        
        if (!autoBackupEnabled) {
          backupStatus = 'disabled'
          return
        }
        
        if (!lastBackupTime) {
          backupStatus = 'never'
          return
        }
        
        // Convert Go time to JavaScript Date
        const lastBackup = new Date(lastBackupTime.toString())
        const now = new Date()
        const hoursSinceBackup = (now.getTime() - lastBackup.getTime()) / (1000 * 60 * 60)
        
        // Check backup frequency from settings
        const backupFrequency = Number(settings.backup_frequency) || 24 // default 24 hours
        
        if (hoursSinceBackup <= backupFrequency) {
          backupStatus = 'recent'
        } else {
          backupStatus = 'overdue'
        }
      }
    } catch (error) {
      console.error('Error checking backup status:', error)
      backupStatus = 'unknown'
    }
  }

  // Auto-hide sidebar when switching to POS only
  $: {
    if (currentView === 'pos' && sidebarVisible && !sidebarAsOverlay) {
      sidebarVisible = false
      sidebarAsOverlay = false
    }
  }

  function toggleSidebar() {
    if (sidebarVisible) {
      // If sidebar is visible, hide it
      sidebarVisible = false
      sidebarAsOverlay = false
    } else {
      // If sidebar is hidden, show it as overlay
      sidebarVisible = true
      sidebarAsOverlay = true
    }
  }

  // Initialize Wails runtime on mount
  onMount(async () => {
    console.log('🎯 App component mounted, initializing Wails runtime...');
    await initializeWailsRuntime();
    
    // Check backup status when app loads
    checkBackupStatus();
    
    // Check backup status every 5 minutes
    const backupStatusInterval = setInterval(checkBackupStatus, 5 * 60 * 1000);
    
    return () => {
      clearInterval(backupStatusInterval);
    };
  });

  // Retry function for manual retry
  async function retryInitialization() {
    console.log('🔄 Retrying Wails runtime initialization...');
    await initializeWailsRuntime();
  }

  function greet() {
    Greet(name).then(result => resultText = result)
  }

  function switchView(view) {
    currentView = view
    // Close overlay sidebar when navigating
    if (sidebarAsOverlay) {
      sidebarVisible = false
      sidebarAsOverlay = false
    }
  }

  // Reactive functions to get the current page icon and title
  $: currentPageIcon = (() => {
    // Dashboard icon
    if (currentView === 'dashboard') {
      return menuItems[0].icon
    }
    
    // Search through categories for the current view
    for (const category of menuItems.slice(1)) {
      const item = category.items.find(item => item.id === currentView)
      if (item) {
        return item.icon
      }
    }
    
    // Default icon if not found
    return 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z'
  })()

  $: currentPageTitle = (() => {
    if (currentView === 'dashboard') return 'Dashboard'
    if (currentView === 'pos') return 'Point of Sale'
    if (currentView === 'sales-invoices') return 'Sales Invoices'
    if (currentView === 'customer') return 'Customer'
    if (currentView === 'payments') return 'Payments'
    if (currentView === 'purchase-invoices') return 'Purchase Invoices'
    if (currentView === 'purchase-products') return 'Purchase Products'
    if (currentView === 'suppliers') return 'Suppliers'
    if (currentView === 'products') return 'Products'
    if (currentView === 'users') return 'Users'
    if (currentView === 'general-settings') return 'General Settings'
    if (currentView === 'administration') return 'Administration'
    if (currentView === 'qr-validation') return 'QR Validation'
    return 'Page'
  })()

  const menuItems = [
    {
      id: 'dashboard',
      label: 'Dashboard',
      icon: 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z M8 5a2 2 0 012-2h4a2 2 0 012 2v6H8V5z',
      category: null
    },
    {
      category: 'Sales',
      items: [
        {
          id: 'pos',
          label: 'POS',
          icon: 'M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z'
        },
        {
          id: 'sales-invoices',
          label: 'Sales Invoices',
          icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z'
        },
        {
          id: 'customer',
          label: 'Customer',
          icon: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z'
        },
        {
          id: 'payments',
          label: 'Payments',
          icon: 'M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z'
        }
      ]
    },
    {
      category: 'Purchases',
      items: [
        {
          id: 'purchase-invoices',
          label: 'Purchase Invoices',
          icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z'
        },
        {
          id: 'purchase-products',
          label: 'Purchase Products',
          icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4'
        },
        {
          id: 'suppliers',
          label: 'Suppliers',
          icon: 'M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4'
        }
      ]
    },
    {
      category: 'Other',
      items: [
        {
          id: 'products',
          label: 'Products',
          icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4'
        },
        {
          id: 'users',
          label: 'Users',
          icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z'
        },
        {
          id: 'general-settings',
          label: 'General settings',
          icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z'
        },
        {
          id: 'administration',
          label: 'Administration',
          icon: 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z'
        },
        {
          id: 'qr-validation',
          label: 'QR Validation',
          icon: 'M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z'
        }
      ]
    }
  ]
</script>

<div class="min-h-screen bg-gradient-to-br from-primary to-secondary flex relative" data-theme="dijibill">
  <!-- Sidebar Overlay (when toggled to show as overlay) -->
  {#if sidebarVisible && sidebarAsOverlay}
    <!-- Backdrop -->
    <div 
      class="fixed inset-0 bg-black/50 z-40 transition-opacity duration-300"
      on:click={toggleSidebar}
      on:keydown={(e) => e.key === 'Escape' && toggleSidebar()}
      role="button"
      tabindex="0"
      aria-label="Close sidebar"
    ></div>
    
    <!-- Sidebar as overlay -->
    <div class="fixed left-0 top-0 h-full w-64 bg-base-100/10 backdrop-blur-lg border-r border-white/20 flex flex-col z-50 transition-transform duration-300">
      <!-- Logo/Header -->
      <div class="p-4 border-b border-white/20">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 24 24">
              <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
            </svg>
          </div>
          <div>
            <h1 class="text-lg font-bold text-white">E-Invoice</h1>
            <p class="text-xs text-white/60">الفاتورة الإلكترونية</p>
          </div>
        </div>
      </div>

      <!-- Navigation Menu -->
      <nav class="flex-1 p-4 overflow-y-auto">
        <ul class="space-y-2">
          <!-- Dashboard (standalone) -->
          <li>
            <button
              class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors {currentView === 'dashboard' ? 'bg-primary text-white' : 'text-white/80 hover:bg-white/10'}"
              on:click={() => switchView('dashboard')}
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="{menuItems[0].icon}"></path>
              </svg>
              <span class="text-sm font-medium">{menuItems[0].label}</span>
            </button>
          </li>

          <!-- Categories -->
          {#each menuItems.slice(1) as category}
            <li class="mt-6">
              <div class="text-xs font-semibold text-white/60 uppercase tracking-wider mb-2 px-3">
                {category.category}
              </div>
              <ul class="space-y-1">
                {#each category.items as item}
                  <li>
                    <button
                      class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors {currentView === item.id ? 'bg-primary text-white' : 'text-white/80 hover:bg-white/10'}"
                      on:click={() => switchView(item.id)}
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="{item.icon}"></path>
                      </svg>
                      <span class="text-sm">{item.label}</span>
                      {#if item.id === 'payments'}
                        <div class="ml-auto">
                          <div class="w-4 h-4 bg-accent rounded-full flex items-center justify-center">
                            <svg class="w-2 h-2 text-white" fill="currentColor" viewBox="0 0 8 8">
                              <circle cx="4" cy="4" r="3"/>
                            </svg>
                          </div>
                        </div>
                      {/if}
                    </button>
                  </li>
                {/each}
              </ul>
            </li>
          {/each}
        </ul>
      </nav>
    </div>
  {:else if sidebarVisible && !sidebarAsOverlay}
    <!-- Regular sidebar for when not in overlay mode -->
    <div class="w-64 flex-shrink-0 bg-base-100/10 backdrop-blur-lg border-r border-white/20 flex flex-col transition-all duration-300">
      <!-- Logo/Header -->
      <div class="p-4 border-b border-white/20">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 24 24">
              <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
            </svg>
          </div>
          <div>
            <h1 class="text-lg font-bold text-white">E-Invoice</h1>
            <p class="text-xs text-white/60">الفاتورة الإلكترونية</p>
          </div>
        </div>
      </div>

      <!-- Navigation Menu -->
      <nav class="flex-1 p-4 overflow-y-auto">
        <ul class="space-y-2">
          <!-- Dashboard (standalone) -->
          <li>
            <button
              class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors {currentView === 'dashboard' ? 'bg-primary text-white' : 'text-white/80 hover:bg-white/10'}"
              on:click={() => switchView('dashboard')}
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="{menuItems[0].icon}"></path>
              </svg>
              <span class="text-sm font-medium">{menuItems[0].label}</span>
            </button>
          </li>

          <!-- Categories -->
          {#each menuItems.slice(1) as category}
            <li class="mt-6">
              <div class="text-xs font-semibold text-white/60 uppercase tracking-wider mb-2 px-3">
                {category.category}
              </div>
              <ul class="space-y-1">
                {#each category.items as item}
                  <li>
                    <button
                      class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors {currentView === item.id ? 'bg-primary text-white' : 'text-white/80 hover:bg-white/10'}"
                      on:click={() => switchView(item.id)}
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="{item.icon}"></path>
                      </svg>
                      <span class="text-sm">{item.label}</span>
                      {#if item.id === 'payments'}
                        <div class="ml-auto">
                          <div class="w-4 h-4 bg-accent rounded-full flex items-center justify-center">
                            <svg class="w-2 h-2 text-white" fill="currentColor" viewBox="0 0 8 8">
                              <circle cx="4" cy="4" r="3"/>
                            </svg>
                          </div>
                        </div>
                      {/if}
                    </button>
                  </li>
                {/each}
              </ul>
            </li>
          {/each}
        </ul>
      </nav>
    </div>
  {/if}

  <!-- Main Content Area -->
  <div class="flex-1 min-w-0 flex flex-col">
    <!-- Top Header -->
    <header class="bg-base-100/10 backdrop-blur-lg border-b border-white/20 p-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-4">
          <!-- Sidebar Toggle Button - Always visible -->
          <button
            class="btn btn-ghost btn-circle text-white hover:bg-white/10"
            on:click={toggleSidebar}
            title="Toggle Menu"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
            </svg>
          </button>
          
          <div class="flex items-center gap-3">
            <!-- Page Icon -->
            <div class="w-10 h-10 bg-primary/20 rounded-lg flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="{currentPageIcon}"></path>
              </svg>
            </div>
            
            <!-- Page Title -->
            <h2 class="text-2xl font-bold text-white">
              {currentPageTitle}
            </h2>
          </div>
        </div>
        
        <div class="flex items-center gap-4">
          <!-- Runtime Status -->
          <div class="flex items-center gap-2">
            {#if $wailsReady}
              <div class="badge badge-success badge-sm">
                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                </svg>
                Connected
              </div>
            {:else if $wailsError}
              <div class="badge badge-error badge-sm">
                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
                </svg>
                Error
              </div>
            {:else}
              <div class="badge badge-warning badge-sm">
                <div class="loading loading-spinner loading-xs mr-1"></div>
                Connecting...
              </div>
            {/if}
          </div>
          
          <!-- Backup Status -->
          <div class="flex items-center gap-2">
            {#if backupStatus === 'recent'}
              <div class="badge badge-success badge-sm">
                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z"></path>
                </svg>
                Backup OK
              </div>
            {:else if backupStatus === 'overdue'}
              <div class="badge badge-warning badge-sm">
                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
                </svg>
                Backup Overdue
              </div>
            {:else if backupStatus === 'never'}
              <div class="badge badge-error badge-sm">
                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
                </svg>
                No Backup
              </div>
            {:else if backupStatus === 'disabled'}
              <div class="badge badge-neutral badge-sm">
                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M13.477 14.89A6 6 0 015.11 6.524l8.367 8.368zm1.414-1.414L6.524 5.11a6 6 0 018.367 8.367zM18 10a8 8 0 11-16 0 8 8 0 0116 0z" clip-rule="evenodd"></path>
                </svg>
                Backup Disabled
              </div>
            {:else}
              <div class="badge badge-ghost badge-sm">
                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9 12a1 1 0 002 0V8a1 1 0 10-2 0v4zm1-7a1 1 0 100 2 1 1 0 000-2z"></path>
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm0-2a6 6 0 100-12 6 6 0 000 12z" clip-rule="evenodd"></path>
                </svg>
                Backup Status Unknown
              </div>
            {/if}
          </div>
          
          <!-- User Menu -->
          <div class="dropdown dropdown-end">
            <div role="button" class="btn btn-ghost btn-circle text-white">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
              </svg>
            </div>
            <ul class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52">
              <li><button class="btn btn-ghost">Profile</button></li>
              <li><button class="btn btn-ghost">Settings</button></li>
              <li><button class="btn btn-ghost">Logout</button></li>
            </ul>
          </div>
        </div>
      </div>
    </header>

    <!-- Page Content -->
    <main class="flex-1 overflow-y-auto">
      {#if !$wailsReady}
        <!-- Loading Screen -->
        <div class="flex items-center justify-center h-full">
          <div class="text-center">
            {#if $wailsError}
              <!-- Error State -->
              <div class="alert alert-error max-w-md mx-auto">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
                </svg>
                <div>
                  <h3 class="font-bold">Initialization Error</h3>
                  <div class="text-xs">{$wailsError}</div>
                </div>
                <button class="btn btn-sm" on:click={retryInitialization}>
                  Retry
                </button>
              </div>
            {:else}
              <!-- Loading State -->
              <div class="flex flex-col items-center gap-4">
                <div class="loading loading-spinner loading-lg text-primary"></div>
                <div class="text-white">
                  <h3 class="text-lg font-semibold">Initializing Application</h3>
                  <p class="text-sm opacity-70">Please wait while we connect to the backend...</p>
                </div>
              </div>
            {/if}
          </div>
        </div>
      {:else}
        <!-- Application Content -->
        {#if currentView === 'dashboard'}
          <Dashboard />
        {:else if currentView === 'pos'}
          <POS />
        {:else if currentView === 'qr-validation'}
          <div class="p-6">
            <QRValidation />
          </div>
        {:else if currentView === 'sales-invoices'}
          <SalesInvoices />
        {:else if currentView === 'customer'}
          <Customer />
        {:else if currentView === 'payments'}
          <Payments />
        {:else if currentView === 'purchase-invoices'}
          <PurchaseInvoices />
        {:else if currentView === 'purchase-products'}
          <PurchaseProducts />
        {:else if currentView === 'suppliers'}
          <Suppliers />
        {:else if currentView === 'products'}
          <Products />
        {:else if currentView === 'users'}
          <Users />
        {:else if currentView === 'general-settings'}
          <GeneralSettings />
        {:else if currentView === 'administration'}
          <Administration />
        {/if}
      {/if}
    </main>
  </div>
</div>
