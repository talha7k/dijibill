<script>
  import { onMount } from 'svelte'
  import { wailsReady } from './stores/wailsStore.js'
  import PageLayout from './components/PageLayout.svelte'
  import { GetCustomers, GetInvoices, GetProducts, GetTodaysSales, GetTopSellingProducts } from '../wailsjs/go/main/App.js'

  let customers = []
  let invoices = []
  let todaysSales = {}
  let topProducts = []
  let isLoading = true
  let stats = {
    items: 0,
    customers: 0,
    products: 0,
    payments: 0
  }

  // Load data when Wails is ready
  $: if ($wailsReady) {
    loadDashboardData()
  }

  // Mock data for payments chart
  let paymentsData = [
    { month: 'Jan', amount: 0 },
    { month: 'Feb', amount: 0 },
    { month: 'Mar', amount: 0 },
    { month: 'Apr', amount: 0 },
    { month: 'May', amount: 0 },
    { month: 'Jun', amount: 0 },
    { month: 'Jul', amount: 0 },
    { month: 'Aug', amount: 0 },
    { month: 'Sep', amount: 0 },
    { month: 'Oct', amount: 0 },
    { month: 'Nov', amount: 0 },
    { month: 'Dec', amount: 0 }
  ]

  async function loadDashboardData() {
    try {
      isLoading = true
      console.log('📊 Loading dashboard data...')
      
      // Load data from backend
      const [customersData, invoicesData, productsData, todaysSalesData, topProductsData] = await Promise.all([
        GetCustomers(),
        GetInvoices(),
        GetProducts(),
        GetTodaysSales(),
        GetTopSellingProducts()
      ])

      customers = customersData || []
      invoices = invoicesData || []
      todaysSales = todaysSalesData || {}
      topProducts = topProductsData || []
      
      // Calculate stats
      stats.customers = customers.length
      stats.items = todaysSales.items_sold || 0
      stats.products = productsData?.length || 0
      stats.payments = (invoices || []).filter(inv => inv.status === 'paid').length
      
      console.log('✅ Dashboard data loaded successfully')
    } catch (error) {
      console.error('❌ Error loading dashboard data:', error)
      // Fallback to empty data
      customers = []
      invoices = []
      todaysSales = {}
      topProducts = []
      stats = { items: 0, customers: 0, products: 0, payments: 0 }
    } finally {
      isLoading = false
    }
  }

  // Calculate invoice statistics
  $: paidInvoices = (invoices || []).filter(inv => inv.status === 'paid').length
  $: totalInvoices = (invoices || []).length
  $: overdueCustomers = (customers || []).slice(0, 3) // Mock overdue customers
  $: topCustomers = (customers || []).slice(0, 3) // Mock top customers

  // Format currency
  function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'SAR',
      minimumFractionDigits: 2
    }).format(amount || 0)
  }
</script>

<PageLayout 
  title="Dashboard" 
  icon="fa-tachometer-alt" 
  showIndicator={true}
>
  <svelte:fragment slot="subtitle">
    Welcome back! Here's what's happening with your business.
  </svelte:fragment>

  {#if isLoading}
    <div class="flex justify-center items-center h-64">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else}
    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <!-- Today's Sales Card -->
      <div class="card bg-base-100/20 backdrop-blur-lg border border-white/20 shadow-xl">
        <div class="card-body p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-white/70 mb-1">Today's Sales</p>
              <p class="text-2xl font-bold text-white">{todaysSales.sales_count || 0}</p>
              <p class="text-xs text-primary mt-1">{formatCurrency(todaysSales.total_amount)}</p>
            </div>
            <div class="text-primary">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Items Sold Today Card -->
      <div class="card bg-base-100/20 backdrop-blur-lg border border-white/20 shadow-xl">
        <div class="card-body p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-white/70 mb-1">Items Sold Today</p>
              <p class="text-2xl font-bold text-white">{todaysSales.items_sold || 0}</p>
              <p class="text-xs text-secondary mt-1">Total items</p>
            </div>
            <div class="text-secondary">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Products and Services Card -->
      <div class="card bg-base-100/20 backdrop-blur-lg border border-white/20 shadow-xl">
        <div class="card-body p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-white/70 mb-1">Products and Services</p>
              <p class="text-2xl font-bold text-white">{stats.products}</p>
              <p class="text-xs text-accent mt-1">All products</p>
            </div>
            <div class="text-accent">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Payments Card -->
      <div class="card bg-base-100/20 backdrop-blur-lg border border-white/20 shadow-xl">
        <div class="card-body p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-white/70 mb-1">Payments</p>
              <p class="text-2xl font-bold text-white">{stats.payments}</p>
              <p class="text-xs text-error mt-1">This Year</p>
            </div>
            <div class="text-error">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z"></path>
              </svg>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Charts Section -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
      <!-- Payments Chart -->
      <div class="lg:col-span-2 card bg-base-100/20 backdrop-blur-lg border border-white/20 shadow-xl">
        <div class="card-body p-6">
          <h3 class="text-lg font-semibold mb-4 text-white">Payments</h3>
          <p class="text-sm text-white/70 mb-4">This Year</p>
          <div class="h-64 flex items-end justify-between gap-1 px-2">
            {#each paymentsData as data}
              <div class="flex flex-col items-center flex-1 max-w-[60px]">
                <div 
                  class="bg-error w-full rounded-t transition-all duration-300 hover:bg-error/80 min-h-[4px]"
                  style="height: {data.amount > 0 ? '80%' : '4px'}"
                ></div>
                <span class="text-xs text-white/70 mt-2 truncate">{data.month}</span>
              </div>
            {/each}
          </div>
        </div>
      </div>

      <!-- Invoices Chart -->
      <div class="card bg-base-100/20 backdrop-blur-lg border border-white/20 shadow-xl">
        <div class="card-body p-6">
          <h3 class="text-lg font-semibold mb-4 text-white">Invoices</h3>
          <div class="flex flex-col items-center">
            <div class="radial-progress text-error text-4xl font-bold mb-4" style="--value:{totalInvoices > 0 ? (paidInvoices / totalInvoices) * 100 : 0}; --size:8rem; --thickness: 8px;">
              {paidInvoices}
            </div>
            <div class="space-y-2 w-full">
              <div class="flex items-center gap-2">
                <div class="w-3 h-3 rounded-full bg-info flex-shrink-0"></div>
                <span class="text-sm text-white">Paid</span>
              </div>
              <div class="flex items-center gap-2">
                <div class="w-3 h-3 rounded-full bg-warning flex-shrink-0"></div>
                <span class="text-sm text-white">Partially paid</span>
              </div>
              <div class="flex items-center gap-2">
                <div class="w-3 h-3 rounded-full bg-primary flex-shrink-0"></div>
                <span class="text-sm text-white">Awaiting payment</span>
              </div>
              <div class="flex items-center gap-2">
                <div class="w-3 h-3 rounded-full bg-error flex-shrink-0"></div>
                <span class="text-sm text-white">Overdue</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Customer Tables -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Top Selling Products -->
      <div class="card bg-base-100/20 backdrop-blur-lg border border-white/20 shadow-xl">
        <div class="card-body p-6">
          <h3 class="text-lg font-semibold mb-4 text-white">Top Selling Products</h3>
          <div class="overflow-x-auto">
            <table class="table table-sm w-full">
              <thead>
                <tr class="border-white/20">
                  <th class="text-white/70">Product</th>
                  <th class="text-white/70">Qty Sold</th>
                  <th class="text-white/70">Revenue</th>
                </tr>
              </thead>
              <tbody>
                {#each topProducts as product}
                  <tr class="border-white/20">
                    <td class="text-white">{product.name}</td>
                    <td class="text-white">{product.total_quantity}</td>
                    <td class="text-primary font-semibold">{formatCurrency(product.total_revenue)}</td>
                  </tr>
                {:else}
                  <tr class="border-white/20">
                    <td colspan="3" class="text-center text-white/50">No products found</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Overdue Customers -->
      <div class="card bg-base-100/20 backdrop-blur-lg border border-white/20 shadow-xl">
        <div class="card-body p-6">
          <h3 class="text-lg font-semibold mb-4 text-white">Overdue Customers</h3>
          <div class="overflow-x-auto">
            <table class="table table-sm w-full">
              <thead>
                <tr class="border-white/20">
                  <th class="text-white/70">Customer Name</th>
                  <th class="text-white/70">Due Date</th>
                  <th class="text-white/70">Amount</th>
                </tr>
              </thead>
              <tbody>
                {#each overdueCustomers as customer}
                  <tr class="border-white/20">
                    <td class="text-white">{customer.name || 'Sample Customer'}</td>
                    <td class="text-white/70">2024-01-15</td>
                    <td class="text-error font-semibold">€ 4,918.00</td>
                  </tr>
                {:else}
                  <tr class="border-white/20">
                    <td class="text-white">Edenta Arabia LTD</td>
                    <td class="text-white/70">2024-01-15</td>
                    <td class="text-error font-semibold">€ 4,918.00</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  {/if}
</PageLayout>