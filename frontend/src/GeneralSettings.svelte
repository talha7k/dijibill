<script>
  import { onMount } from 'svelte'

  // State variables
  let activeTab = 'company'
  let isLoading = false

  // Company settings
  let companySettings = {
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

  // Tax rates
  let taxRates = [
    { id: 1, name: 'Standard VAT', name_arabic: 'ضريبة القيمة المضافة', rate: 15.0, is_default: true, is_active: true },
    { id: 2, name: 'Zero VAT', name_arabic: 'معفى من الضريبة', rate: 0.0, is_default: false, is_active: true },
    { id: 3, name: 'Exempt', name_arabic: 'خارج نطاق الضريبة', rate: 0.0, is_default: false, is_active: true }
  ]

  // Payment Types
  let paymentTypes = [
    { id: 1, name: 'Cash', name_arabic: 'نقدي', code: 'CASH', description: 'Cash payment', is_default: true, is_active: true },
    { id: 2, name: 'Credit Card', name_arabic: 'بطاقة ائتمان', code: 'CARD', description: 'Credit/Debit card payment', is_default: false, is_active: true },
    { id: 3, name: 'Bank Transfer', name_arabic: 'تحويل بنكي', code: 'TRANSFER', description: 'Bank transfer payment', is_default: false, is_active: true },
    { id: 4, name: 'Check', name_arabic: 'شيك', code: 'CHECK', description: 'Check payment', is_default: false, is_active: true }
  ]

  // Sales Categories
  let salesCategories = [
    { id: 1, name: 'Retail', name_arabic: 'تجزئة', code: 'RETAIL', description: 'Retail sales', description_arabic: 'مبيعات التجزئة', is_default: true, is_active: true },
    { id: 2, name: 'Wholesale', name_arabic: 'جملة', code: 'WHOLESALE', description: 'Wholesale sales', description_arabic: 'مبيعات الجملة', is_default: false, is_active: true },
    { id: 3, name: 'Online', name_arabic: 'إلكتروني', code: 'ONLINE', description: 'Online sales', description_arabic: 'مبيعات إلكترونية', is_default: false, is_active: true }
  ]

  // Units of measurement
  let units = [
    { id: 1, value: 'pcs', label: 'Pieces', arabic: 'قطعة', is_default: true, is_active: true },
    { id: 2, value: 'kg', label: 'Kilograms', arabic: 'كيلوغرام', is_default: false, is_active: true },
    { id: 3, value: 'g', label: 'Grams', arabic: 'غرام', is_default: false, is_active: true },
    { id: 4, value: 'l', label: 'Liters', arabic: 'لتر', is_default: false, is_active: true },
    { id: 5, value: 'ml', label: 'Milliliters', arabic: 'مليلتر', is_default: false, is_active: true },
    { id: 6, value: 'm', label: 'Meters', arabic: 'متر', is_default: false, is_active: true },
    { id: 7, value: 'cm', label: 'Centimeters', arabic: 'سنتيمتر', is_default: false, is_active: true },
    { id: 8, value: 'box', label: 'Box', arabic: 'صندوق', is_default: false, is_active: true },
    { id: 9, value: 'pack', label: 'Pack', arabic: 'عبوة', is_default: false, is_active: true }
  ]

  // Default quantity settings
  let defaultQuantitySettings = {
    enable_default_quantity: true,
    default_quantity_value: 1,
    allow_quantity_override: true,
    min_quantity: 0.01,
    max_quantity: 9999,
    quantity_decimal_places: 2
  }

  // Form states for adding new items
  let newTaxRate = { name: '', name_arabic: '', rate: 0, is_default: false, is_active: true }
  let newPaymentType = { name: '', name_arabic: '', code: '', description: '', is_default: false, is_active: true }
  let newSalesCategory = { name: '', name_arabic: '', code: '', description: '', description_arabic: '', is_default: false, is_active: true }
  let newUnit = { value: '', label: '', arabic: '', is_default: false, is_active: true }
  let showTaxRateForm = false
  let showPaymentTypeForm = false
  let showSalesCategoryForm = false
  let showUnitForm = false
  let editingTaxRate = null
  let editingPaymentType = null
  let editingSalesCategory = null
  let editingUnit = null

  onMount(() => {
    loadSettings()
  })

  async function loadSettings() {
    // Load settings from backend
    isLoading = true
    try {
      // TODO: Implement actual API calls
      console.log('Loading settings...')
    } catch (error) {
      console.error('Error loading settings:', error)
    } finally {
      isLoading = false
    }
  }

  async function saveSettings() {
    isLoading = true
    try {
      // TODO: Implement actual API calls
      console.log('Saving settings...')
      alert('Settings saved successfully!')
    } catch (error) {
      console.error('Error saving settings:', error)
      alert('Error saving settings: ' + error.message)
    } finally {
      isLoading = false
    }
  }

  // Tax rate functions
  function addTaxRate() {
    if (newTaxRate.name && newTaxRate.rate >= 0) {
      const id = Math.max(...taxRates.map(t => t.id), 0) + 1
      taxRates = [...taxRates, { ...newTaxRate, id }]
      resetTaxRateForm()
    }
  }

  function editTaxRate(taxRate) {
    editingTaxRate = taxRate.id
    newTaxRate = { ...taxRate }
    showTaxRateForm = true
  }

  function updateTaxRate() {
    taxRates = taxRates.map(t => t.id === editingTaxRate ? { ...newTaxRate, id: editingTaxRate } : t)
    resetTaxRateForm()
  }

  function deleteTaxRate(id) {
    if (confirm('Are you sure you want to delete this tax rate?')) {
      taxRates = taxRates.filter(t => t.id !== id)
    }
  }

  function resetTaxRateForm() {
    newTaxRate = { name: '', name_arabic: '', rate: 0, is_default: false, is_active: true }
    showTaxRateForm = false
    editingTaxRate = null
  }

  // Payment type functions
  function addPaymentType() {
    if (newPaymentType.name && newPaymentType.code) {
      const id = Math.max(...paymentTypes.map(p => p.id), 0) + 1
      paymentTypes = [...paymentTypes, { ...newPaymentType, id }]
      resetPaymentTypeForm()
    }
  }

  function editPaymentType(paymentType) {
    editingPaymentType = paymentType.id
    newPaymentType = { ...paymentType }
    showPaymentTypeForm = true
  }

  function updatePaymentType() {
    paymentTypes = paymentTypes.map(p => p.id === editingPaymentType ? { ...newPaymentType, id: editingPaymentType } : p)
    resetPaymentTypeForm()
  }

  function deletePaymentType(id) {
    if (confirm('Are you sure you want to delete this payment type?')) {
      paymentTypes = paymentTypes.filter(p => p.id !== id)
    }
  }

  function resetPaymentTypeForm() {
    newPaymentType = { name: '', name_arabic: '', code: '', description: '', is_default: false, is_active: true }
    showPaymentTypeForm = false
    editingPaymentType = null
  }

  function setDefaultPaymentType(id) {
    paymentTypes = paymentTypes.map(p => ({ ...p, is_default: p.id === id }))
  }

  // Sales category functions
  function addSalesCategory() {
    if (newSalesCategory.name && newSalesCategory.code) {
      const id = Math.max(...salesCategories.map(s => s.id), 0) + 1
      salesCategories = [...salesCategories, { ...newSalesCategory, id }]
      resetSalesCategoryForm()
    }
  }

  function editSalesCategory(salesCategory) {
    editingSalesCategory = salesCategory.id
    newSalesCategory = { ...salesCategory }
    showSalesCategoryForm = true
  }

  function updateSalesCategory() {
    salesCategories = salesCategories.map(s => s.id === editingSalesCategory ? { ...newSalesCategory, id: editingSalesCategory } : s)
    resetSalesCategoryForm()
  }

  function deleteSalesCategory(id) {
    if (confirm('Are you sure you want to delete this sales category?')) {
      salesCategories = salesCategories.filter(s => s.id !== id)
    }
  }

  function resetSalesCategoryForm() {
    newSalesCategory = { name: '', name_arabic: '', code: '', description: '', description_arabic: '', is_default: false, is_active: true }
    showSalesCategoryForm = false
    editingSalesCategory = null
  }

  function setDefaultSalesCategory(id) {
    salesCategories = salesCategories.map(s => ({ ...s, is_default: s.id === id }))
  }

  // Unit functions
  function addUnit() {
    if (newUnit.value && newUnit.label) {
      const id = Math.max(...units.map(u => u.id), 0) + 1
      units = [...units, { ...newUnit, id }]
      resetUnitForm()
    }
  }

  function editUnit(unit) {
    editingUnit = unit.id
    newUnit = { ...unit }
    showUnitForm = true
  }

  function updateUnit() {
    units = units.map(u => u.id === editingUnit ? { ...newUnit, id: editingUnit } : u)
    resetUnitForm()
  }

  function deleteUnit(id) {
    if (confirm('Are you sure you want to delete this unit?')) {
      units = units.filter(u => u.id !== id)
    }
  }

  function resetUnitForm() {
    newUnit = { value: '', label: '', arabic: '', is_default: false, is_active: true }
    showUnitForm = false
    editingUnit = null
  }

  function setDefaultTaxRate(id) {
    taxRates = taxRates.map(t => ({ ...t, is_default: t.id === id }))
  }

  function setDefaultUnit(id) {
    units = units.map(u => ({ ...u, is_default: u.id === id }))
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-blue-900 to-purple-900 p-6">
  <div class="max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-bold text-white mb-2">General Settings</h1>
        <p class="text-white/60">Configure your system preferences and business settings</p>
      </div>
      <button 
        class="btn btn-primary"
        on:click={saveSettings}
        disabled={isLoading}
      >
        {#if isLoading}
          <span class="loading loading-spinner loading-sm"></span>
        {/if}
        Save All Settings
      </button>
    </div>

    <!-- Tabs -->
    <div class="tabs tabs-boxed bg-white/10 mb-6">
      <button 
        class="tab {activeTab === 'company' ? 'tab-active' : ''}" 
        on:click={() => activeTab = 'company'}
      >
        Company Info
      </button>
      <button 
        class="tab {activeTab === 'system' ? 'tab-active' : ''}" 
        on:click={() => activeTab = 'system'}
      >
        System Settings
      </button>
      <button 
        class="tab {activeTab === 'payment-types' ? 'tab-active' : ''}" 
        on:click={() => activeTab = 'payment-types'}
      >
        Payment Types
      </button>
      <button 
        class="tab {activeTab === 'sales-categories' ? 'tab-active' : ''}" 
        on:click={() => activeTab = 'sales-categories'}
      >
        Sales Categories
      </button>
      <button 
        class="tab {activeTab === 'taxrates' ? 'tab-active' : ''}" 
        on:click={() => activeTab = 'taxrates'}
      >
        Tax Rates
      </button>
      <button 
        class="tab {activeTab === 'units' ? 'tab-active' : ''}" 
        on:click={() => activeTab = 'units'}
      >
        Units
      </button>
      <button 
        class="tab {activeTab === 'quantity' ? 'tab-active' : ''}" 
        on:click={() => activeTab = 'quantity'}
      >
        Default Quantity
      </button>
    </div>

    <!-- Tab Content -->
    <div class="card bg-base-100/20 backdrop-blur-lg border border-white/20 shadow-xl">
      <div class="card-body">
        
        {#if activeTab === 'company'}
          <!-- Company Information Tab -->
          <h2 class="text-xl font-bold text-white mb-6">Company Information</h2>
          
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div class="space-y-4">
              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Company Name</label>
                <input 
                  type="text" 
                  bind:value={companySettings.name}
                  placeholder="Enter company name" 
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                />
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Company Name (Arabic)</label>
                <input 
                  type="text" 
                  bind:value={companySettings.name_arabic}
                  placeholder="أدخل اسم الشركة" 
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                  dir="rtl"
                />
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">VAT Number</label>
                <input 
                  type="text" 
                  bind:value={companySettings.vat_number}
                  placeholder="Enter VAT number" 
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                />
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">CR Number</label>
                <input 
                  type="text" 
                  bind:value={companySettings.cr_number}
                  placeholder="Enter CR number" 
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                />
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Email</label>
                <input 
                  type="email" 
                  bind:value={companySettings.email}
                  placeholder="Enter email address" 
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                />
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Phone</label>
                <input 
                  type="tel" 
                  bind:value={companySettings.phone}
                  placeholder="Enter phone number" 
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                />
              </div>
            </div>

            <div class="space-y-4">
              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Address</label>
                <textarea 
                  bind:value={companySettings.address}
                  placeholder="Enter company address" 
                  class="textarea textarea-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                  rows="3"
                ></textarea>
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Address (Arabic)</label>
                <textarea 
                  bind:value={companySettings.address_arabic}
                  placeholder="أدخل عنوان الشركة" 
                  class="textarea textarea-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                  rows="3"
                  dir="rtl"
                ></textarea>
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">City</label>
                <input 
                  type="text" 
                  bind:value={companySettings.city}
                  placeholder="Enter city" 
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                />
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">City (Arabic)</label>
                <input 
                  type="text" 
                  bind:value={companySettings.city_arabic}
                  placeholder="أدخل المدينة" 
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                  dir="rtl"
                />
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Country</label>
                <input 
                  type="text" 
                  bind:value={companySettings.country}
                  placeholder="Enter country" 
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                />
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Country (Arabic)</label>
                <input 
                  type="text" 
                  bind:value={companySettings.country_arabic}
                  placeholder="أدخل البلد" 
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                  dir="rtl"
                />
              </div>
            </div>
          </div>

        {:else if activeTab === 'payment-types'}
          <!-- Payment Types Tab -->
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-xl font-bold text-white">Payment Types</h2>
            <button 
              class="btn btn-primary btn-sm"
              on:click={() => showPaymentTypeForm = true}
            >
              Add Payment Type
            </button>
          </div>

          {#if showPaymentTypeForm}
            <div class="card bg-white/10 border border-white/20 mb-6">
              <div class="card-body">
                <h3 class="text-lg font-semibold text-white mb-4">
                  {editingPaymentType ? 'Edit Payment Type' : 'Add New Payment Type'}
                </h3>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Payment Type Name</label>
                    <input 
                      type="text" 
                      bind:value={newPaymentType.name}
                      placeholder="Enter payment type name" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                    />
                  </div>

                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Payment Type Name (Arabic)</label>
                    <input 
                      type="text" 
                      bind:value={newPaymentType.name_arabic}
                      placeholder="أدخل اسم نوع الدفع" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                      dir="rtl"
                    />
                  </div>

                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Code</label>
                    <input 
                      type="text" 
                      bind:value={newPaymentType.code}
                      placeholder="Enter payment type code" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                    />
                  </div>

                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Description</label>
                    <input 
                      type="text" 
                      bind:value={newPaymentType.description}
                      placeholder="Enter description" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                    />
                  </div>

                  <div class="form-control">
                    <label class="cursor-pointer flex items-center gap-3 mt-8">
                      <input 
                        type="checkbox" 
                        bind:checked={newPaymentType.is_default}
                        class="checkbox checkbox-primary" 
                      />
                      <span class="text-white font-medium">Set as Default</span>
                    </label>
                  </div>
                </div>

                <div class="flex gap-2 mt-4">
                  <button 
                    class="btn btn-primary btn-sm"
                    on:click={editingPaymentType ? updatePaymentType : addPaymentType}
                  >
                    {editingPaymentType ? 'Update' : 'Add'} Payment Type
                  </button>
                  <button 
                    class="btn btn-ghost btn-sm"
                    on:click={resetPaymentTypeForm}
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>
          {/if}

          <!-- Payment Types Table -->
          <div class="overflow-x-auto">
            <table class="table table-zebra bg-white/5">
              <thead>
                <tr class="text-white">
                  <th>Name</th>
                  <th>Name (Arabic)</th>
                  <th>Code</th>
                  <th>Description</th>
                  <th>Default</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {#each paymentTypes as paymentType}
                  <tr class="text-white/80">
                    <td class="font-medium">{paymentType.name}</td>
                    <td class="font-medium" dir="rtl">{paymentType.name_arabic}</td>
                    <td><span class="badge badge-outline">{paymentType.code}</span></td>
                    <td>{paymentType.description}</td>
                    <td>
                      {#if paymentType.is_default}
                        <span class="badge badge-primary">Default</span>
                      {:else}
                        <button 
                          class="btn btn-xs btn-ghost text-white/60"
                          on:click={() => setDefaultPaymentType(paymentType.id)}
                        >
                          Set Default
                        </button>
                      {/if}
                    </td>
                    <td>
                      <span class="badge {paymentType.is_active ? 'badge-success' : 'badge-error'}">
                        {paymentType.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td>
                      <div class="flex gap-1">
                        <button 
                          class="btn btn-xs btn-primary"
                          on:click={() => editPaymentType(paymentType)}
                        >
                          Edit
                        </button>
                        <button 
                          class="btn btn-xs btn-error"
                          on:click={() => deletePaymentType(paymentType.id)}
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>

        {:else if activeTab === 'sales-categories'}
          <!-- Sales Categories Tab -->
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-xl font-bold text-white">Sales Categories</h2>
            <button 
              class="btn btn-primary btn-sm"
              on:click={() => showSalesCategoryForm = true}
            >
              Add Sales Category
            </button>
          </div>

          {#if showSalesCategoryForm}
            <div class="card bg-white/10 border border-white/20 mb-6">
              <div class="card-body">
                <h3 class="text-lg font-semibold text-white mb-4">
                  {editingSalesCategory ? 'Edit Sales Category' : 'Add New Sales Category'}
                </h3>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Category Name</label>
                    <input 
                      type="text" 
                      bind:value={newSalesCategory.name}
                      placeholder="Enter category name" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                    />
                  </div>

                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Category Name (Arabic)</label>
                    <input 
                      type="text" 
                      bind:value={newSalesCategory.name_arabic}
                      placeholder="أدخل اسم الفئة" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                      dir="rtl"
                    />
                  </div>

                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Code</label>
                    <input 
                      type="text" 
                      bind:value={newSalesCategory.code}
                      placeholder="Enter category code" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                    />
                  </div>

                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Description</label>
                    <input 
                      type="text" 
                      bind:value={newSalesCategory.description}
                      placeholder="Enter description" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                    />
                  </div>

                  <div class="form-control md:col-span-2">
                    <label class="text-white font-medium mb-2 block">Description (Arabic)</label>
                    <input 
                      type="text" 
                      bind:value={newSalesCategory.description_arabic}
                      placeholder="أدخل الوصف" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                      dir="rtl"
                    />
                  </div>

                  <div class="form-control">
                    <label class="cursor-pointer flex items-center gap-3 mt-8">
                      <input 
                        type="checkbox" 
                        bind:checked={newSalesCategory.is_default}
                        class="checkbox checkbox-primary" 
                      />
                      <span class="text-white font-medium">Set as Default</span>
                    </label>
                  </div>
                </div>

                <div class="flex gap-2 mt-4">
                  <button 
                    class="btn btn-primary btn-sm"
                    on:click={editingSalesCategory ? updateSalesCategory : addSalesCategory}
                  >
                    {editingSalesCategory ? 'Update' : 'Add'} Sales Category
                  </button>
                  <button 
                    class="btn btn-ghost btn-sm"
                    on:click={resetSalesCategoryForm}
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>
          {/if}

          <!-- Sales Categories Table -->
          <div class="overflow-x-auto">
            <table class="table table-zebra bg-white/5">
              <thead>
                <tr class="text-white">
                  <th>Name</th>
                  <th>Name (Arabic)</th>
                  <th>Code</th>
                  <th>Description</th>
                  <th>Default</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {#each salesCategories as salesCategory}
                  <tr class="text-white/80">
                    <td class="font-medium">{salesCategory.name}</td>
                    <td class="font-medium" dir="rtl">{salesCategory.name_arabic}</td>
                    <td><span class="badge badge-outline">{salesCategory.code}</span></td>
                    <td>{salesCategory.description}</td>
                    <td>
                      {#if salesCategory.is_default}
                        <span class="badge badge-primary">Default</span>
                      {:else}
                        <button 
                          class="btn btn-xs btn-ghost text-white/60"
                          on:click={() => setDefaultSalesCategory(salesCategory.id)}
                        >
                          Set Default
                        </button>
                      {/if}
                    </td>
                    <td>
                      <span class="badge {salesCategory.is_active ? 'badge-success' : 'badge-error'}">
                        {salesCategory.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td>
                      <div class="flex gap-1">
                        <button 
                          class="btn btn-xs btn-primary"
                          on:click={() => editSalesCategory(salesCategory)}
                        >
                          Edit
                        </button>
                        <button 
                          class="btn btn-xs btn-error"
                          on:click={() => deleteSalesCategory(salesCategory.id)}
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>

        {:else if activeTab === 'system'}
          <!-- System Settings Tab -->
          <h2 class="text-xl font-bold text-white mb-6">System Settings</h2>
          
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div class="space-y-4">
              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Default Currency</label>
                <select 
                  bind:value={systemSettings.currency}
                  class="select select-bordered bg-white/10 border-white/20 text-white"
                >
                  <option value="SAR">Saudi Riyal (SAR)</option>
                  <option value="USD">US Dollar (USD)</option>
                  <option value="EUR">Euro (EUR)</option>
                  <option value="AED">UAE Dirham (AED)</option>
                </select>
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Language</label>
                <select 
                  bind:value={systemSettings.language}
                  class="select select-bordered bg-white/10 border-white/20 text-white"
                >
                  <option value="en">English</option>
                  <option value="ar">العربية</option>
                </select>
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Timezone</label>
                <select 
                  bind:value={systemSettings.timezone}
                  class="select select-bordered bg-white/10 border-white/20 text-white"
                >
                  <option value="Asia/Riyadh">Asia/Riyadh</option>
                  <option value="Asia/Dubai">Asia/Dubai</option>
                  <option value="UTC">UTC</option>
                </select>
              </div>
            </div>

            <div class="space-y-4">
              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Date Format</label>
                <select 
                  bind:value={systemSettings.date_format}
                  class="select select-bordered bg-white/10 border-white/20 text-white"
                >
                  <option value="DD/MM/YYYY">DD/MM/YYYY</option>
                  <option value="MM/DD/YYYY">MM/DD/YYYY</option>
                  <option value="YYYY-MM-DD">YYYY-MM-DD</option>
                </select>
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Backup Frequency</label>
                <select 
                  bind:value={systemSettings.backup_frequency}
                  class="select select-bordered bg-white/10 border-white/20 text-white"
                >
                  <option value="daily">Daily</option>
                  <option value="weekly">Weekly</option>
                  <option value="monthly">Monthly</option>
                </select>
              </div>

              <div class="space-y-3">
                <div class="form-control">
                  <label class="cursor-pointer flex items-center gap-3">
                    <input 
                      type="checkbox" 
                      bind:checked={systemSettings.zatca_enabled}
                      class="checkbox checkbox-primary" 
                    />
                    <span class="text-white font-medium">Enable ZATCA Integration</span>
                  </label>
                </div>

                <div class="form-control">
                  <label class="cursor-pointer flex items-center gap-3">
                    <input 
                      type="checkbox" 
                      bind:checked={systemSettings.auto_backup}
                      class="checkbox checkbox-primary" 
                    />
                    <span class="text-white font-medium">Enable Auto Backup</span>
                  </label>
                </div>
              </div>
            </div>
          </div>

        {:else if activeTab === 'taxrates'}
          <!-- Tax Rates Tab -->
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-xl font-bold text-white">Tax Rates</h2>
            <button 
              class="btn btn-primary btn-sm"
              on:click={() => showTaxRateForm = true}
            >
              Add Tax Rate
            </button>
          </div>

          {#if showTaxRateForm}
            <div class="card bg-white/10 border border-white/20 mb-6">
              <div class="card-body">
                <h3 class="text-lg font-semibold text-white mb-4">
                  {editingTaxRate ? 'Edit Tax Rate' : 'Add New Tax Rate'}
                </h3>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Tax Rate Name</label>
                    <input 
                      type="text" 
                      bind:value={newTaxRate.name}
                      placeholder="Enter tax rate name" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                    />
                  </div>

                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Tax Rate Name (Arabic)</label>
                    <input 
                      type="text" 
                      bind:value={newTaxRate.name_arabic}
                      placeholder="أدخل اسم معدل الضريبة" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                      dir="rtl"
                    />
                  </div>

                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Rate (%)</label>
                    <input 
                      type="number" 
                      bind:value={newTaxRate.rate}
                      placeholder="0.00" 
                      step="0.01"
                      min="0"
                      max="100"
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                    />
                  </div>

                  <div class="form-control">
                    <label class="cursor-pointer flex items-center gap-3 mt-8">
                      <input 
                        type="checkbox" 
                        bind:checked={newTaxRate.is_default}
                        class="checkbox checkbox-primary" 
                      />
                      <span class="text-white font-medium">Set as Default</span>
                    </label>
                  </div>
                </div>

                <div class="flex gap-2 mt-4">
                  <button 
                    class="btn btn-primary btn-sm"
                    on:click={editingTaxRate ? updateTaxRate : addTaxRate}
                  >
                    {editingTaxRate ? 'Update' : 'Add'} Tax Rate
                  </button>
                  <button 
                    class="btn btn-ghost btn-sm"
                    on:click={resetTaxRateForm}
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>
          {/if}

          <div class="overflow-x-auto">
            <table class="table table-zebra bg-white/5">
              <thead>
                <tr class="text-white">
                  <th>Name</th>
                  <th>Name (Arabic)</th>
                  <th>Rate (%)</th>
                  <th>Default</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {#each taxRates as taxRate}
                  <tr class="text-white/80">
                    <td>{taxRate.name}</td>
                    <td dir="rtl">{taxRate.name_arabic}</td>
                    <td>{taxRate.rate}%</td>
                    <td>
                      {#if taxRate.is_default}
                        <span class="badge badge-primary">Default</span>
                      {:else}
                        <button 
                          class="btn btn-xs btn-outline"
                          on:click={() => setDefaultTaxRate(taxRate.id)}
                        >
                          Set Default
                        </button>
                      {/if}
                    </td>
                    <td>
                      <span class="badge {taxRate.is_active ? 'badge-success' : 'badge-error'}">
                        {taxRate.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td>
                      <div class="flex gap-1">
                        <button 
                          class="btn btn-xs btn-primary"
                          on:click={() => editTaxRate(taxRate)}
                        >
                          Edit
                        </button>
                        <button 
                          class="btn btn-xs btn-error"
                          on:click={() => deleteTaxRate(taxRate.id)}
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>

        {:else if activeTab === 'units'}
          <!-- Units Tab -->
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-xl font-bold text-white">Units of Measurement</h2>
            <button 
              class="btn btn-primary btn-sm"
              on:click={() => showUnitForm = true}
            >
              Add Unit
            </button>
          </div>

          {#if showUnitForm}
            <div class="card bg-white/10 border border-white/20 mb-6">
              <div class="card-body">
                <h3 class="text-lg font-semibold text-white mb-4">
                  {editingUnit ? 'Edit Unit' : 'Add New Unit'}
                </h3>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Unit Code</label>
                    <input 
                      type="text" 
                      bind:value={newUnit.value}
                      placeholder="e.g., kg, pcs, m" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                    />
                  </div>

                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Unit Label</label>
                    <input 
                      type="text" 
                      bind:value={newUnit.label}
                      placeholder="e.g., Kilograms, Pieces" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                    />
                  </div>

                  <div class="form-control">
                    <label class="text-white font-medium mb-2 block">Arabic Label</label>
                    <input 
                      type="text" 
                      bind:value={newUnit.arabic}
                      placeholder="مثال: كيلوغرام، قطعة" 
                      class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                      dir="rtl"
                    />
                  </div>

                  <div class="form-control">
                    <label class="cursor-pointer flex items-center gap-3 mt-8">
                      <input 
                        type="checkbox" 
                        bind:checked={newUnit.is_default}
                        class="checkbox checkbox-primary" 
                      />
                      <span class="text-white font-medium">Set as Default</span>
                    </label>
                  </div>
                </div>

                <div class="flex gap-2 mt-4">
                  <button 
                    class="btn btn-primary btn-sm"
                    on:click={editingUnit ? updateUnit : addUnit}
                  >
                    {editingUnit ? 'Update' : 'Add'} Unit
                  </button>
                  <button 
                    class="btn btn-ghost btn-sm"
                    on:click={resetUnitForm}
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>
          {/if}

          <div class="overflow-x-auto">
            <table class="table table-zebra bg-white/5">
              <thead>
                <tr class="text-white">
                  <th>Code</th>
                  <th>Label</th>
                  <th>Arabic Label</th>
                  <th>Default</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {#each units as unit}
                  <tr class="text-white/80">
                    <td><code class="bg-white/10 px-2 py-1 rounded">{unit.value}</code></td>
                    <td>{unit.label}</td>
                    <td dir="rtl">{unit.arabic}</td>
                    <td>
                      {#if unit.is_default}
                        <span class="badge badge-primary">Default</span>
                      {:else}
                        <button 
                          class="btn btn-xs btn-outline"
                          on:click={() => setDefaultUnit(unit.id)}
                        >
                          Set Default
                        </button>
                      {/if}
                    </td>
                    <td>
                      <span class="badge {unit.is_active ? 'badge-success' : 'badge-error'}">
                        {unit.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td>
                      <div class="flex gap-1">
                        <button 
                          class="btn btn-xs btn-primary"
                          on:click={() => editUnit(unit)}
                        >
                          Edit
                        </button>
                        <button 
                          class="btn btn-xs btn-error"
                          on:click={() => deleteUnit(unit.id)}
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>

        {:else if activeTab === 'quantity'}
          <!-- Default Quantity Tab -->
          <h2 class="text-xl font-bold text-white mb-6">Default Quantity Settings</h2>
          
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div class="space-y-4">
              <div class="form-control">
                <label class="cursor-pointer flex items-center gap-3">
                  <input 
                    type="checkbox" 
                    bind:checked={defaultQuantitySettings.enable_default_quantity}
                    class="checkbox checkbox-primary" 
                  />
                  <span class="text-white font-medium">Enable Default Quantity</span>
                </label>
                <p class="text-white/60 text-sm mt-1">When enabled, new products will have a default quantity value</p>
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Default Quantity Value</label>
                <input 
                  type="number" 
                  bind:value={defaultQuantitySettings.default_quantity_value}
                  placeholder="1" 
                  step="0.01"
                  min="0"
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                  disabled={!defaultQuantitySettings.enable_default_quantity}
                />
              </div>

              <div class="form-control">
                <label class="cursor-pointer flex items-center gap-3">
                  <input 
                    type="checkbox" 
                    bind:checked={defaultQuantitySettings.allow_quantity_override}
                    class="checkbox checkbox-primary" 
                  />
                  <span class="text-white font-medium">Allow Quantity Override</span>
                </label>
                <p class="text-white/60 text-sm mt-1">Allow users to change quantity when creating invoices</p>
              </div>
            </div>

            <div class="space-y-4">
              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Minimum Quantity</label>
                <input 
                  type="number" 
                  bind:value={defaultQuantitySettings.min_quantity}
                  placeholder="0.01" 
                  step="0.01"
                  min="0"
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                />
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Maximum Quantity</label>
                <input 
                  type="number" 
                  bind:value={defaultQuantitySettings.max_quantity}
                  placeholder="9999" 
                  step="1"
                  min="1"
                  class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/50"
                />
              </div>

              <div class="form-control">
                <label class="text-white font-medium mb-2 block">Decimal Places</label>
                <select 
                  bind:value={defaultQuantitySettings.quantity_decimal_places}
                  class="select select-bordered bg-white/10 border-white/20 text-white"
                >
                  <option value={0}>0 (Whole numbers)</option>
                  <option value={1}>1 decimal place</option>
                  <option value={2}>2 decimal places</option>
                  <option value={3}>3 decimal places</option>
                </select>
              </div>
            </div>
          </div>

          <div class="mt-6 p-4 bg-blue-500/20 border border-blue-500/30 rounded-lg">
            <h3 class="text-white font-semibold mb-2">Preview</h3>
            <p class="text-white/80 text-sm">
              With current settings, new products will have a default quantity of 
              <strong>{defaultQuantitySettings.default_quantity_value}</strong>
              {#if defaultQuantitySettings.allow_quantity_override}
                (can be changed by users)
              {:else}
                (fixed value)
              {/if}
              with {defaultQuantitySettings.quantity_decimal_places} decimal places.
            </p>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>