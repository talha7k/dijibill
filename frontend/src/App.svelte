<script>
  import Dashboard from './Dashboard.svelte'
  import QRValidation from './QRValidation.svelte'
  import SalesInvoices from './SalesInvoices.svelte'
  import Customer from './Customer.svelte'
  import Payments from './Payments.svelte'
  import PurchaseInvoices from './PurchaseInvoices.svelte'
  import Suppliers from './Suppliers.svelte'
  import Products from './Products.svelte'
  import Users from './Users.svelte'
  import GeneralSettings from './GeneralSettings.svelte'
  import Administration from './Administration.svelte'
  import {Greet} from '../wailsjs/go/main/App.js'

  let currentView = 'dashboard' // Default to dashboard
  let resultText = "Welcome to DijiBill - ZATCA Compliant Invoicing"
  let name = ''

  function greet() {
    Greet(name).then(result => resultText = result)
  }

  function switchView(view) {
    currentView = view
  }

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

<div class="min-h-screen bg-gradient-to-br from-primary to-secondary flex" data-theme="dijibill">
  <!-- Sidebar -->
  <div class="w-64 flex-shrink-0 bg-base-100/10 backdrop-blur-lg border-r border-white/20 flex flex-col">
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

  <!-- Main Content Area -->
  <div class="flex-1 min-w-0 flex flex-col">
    <!-- Top Header -->
    <header class="bg-base-100/10 backdrop-blur-lg border-b border-white/20 p-4">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-2xl font-bold text-white capitalize">
            {#if currentView === 'dashboard'}
              Dashboard
            {:else if currentView === 'sales-invoices'}
              Sales Invoices
            {:else if currentView === 'customer'}
              Customer
            {:else if currentView === 'payments'}
              Payments
            {:else if currentView === 'purchase-invoices'}
              Purchase Invoices
            {:else if currentView === 'suppliers'}
              Suppliers
            {:else if currentView === 'products'}
              Products
            {:else if currentView === 'users'}
              Users
            {:else if currentView === 'general-settings'}
              General Settings
            {:else if currentView === 'administration'}
              Administration
            {:else if currentView === 'qr-validation'}
              QR Validation
            {/if}
          </h2>
        </div>
        
        <div class="flex items-center gap-4">
          <!-- Search -->
          <div class="form-control">
            <input type="text" placeholder="Search..." class="input input-sm bg-white/10 border-white/20 text-white placeholder-white/60" />
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
      {#if currentView === 'dashboard'}
        <Dashboard />
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
    </main>
  </div>
</div>
