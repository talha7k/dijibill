<script>
  import { onMount } from 'svelte'
  import { GetProducts, CreateProduct, GetProductCategories, CreateProductCategory, GetDefaultProductSettings } from '../wailsjs/go/main/App.js'
  import * as main from '../wailsjs/go/models'
  
  // Import components
  import ProductList from './components/ProductList.svelte'
  import ProductCategoriesTab from './components/ProductCategoriesTab.svelte'
  import ProductModal from './components/ProductModal.svelte'
  import HorizontalTabs from './components/HorizontalTabs.svelte'

  // State variables
  let products = []
  let categories = []
  let defaultProductSettings = null
  let isLoading = false
  let searchTerm = ''
  let showProductModal = false
  let editingProduct = null
  let activeTab = 'products'

  // Tab configuration
  const tabs = [
    { id: 'products', label: 'Products', icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4' },
    { id: 'categories', label: 'Categories', icon: 'M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10' }
  ]

  // Handle tab change
  function handleTabChange(event) {
    activeTab = event.detail.tabId
  }

  // Form data with default values
  let productForm = {
    name: '',
    name_arabic: '',
    description: '',
    description_arabic: '',
    category_id: 0,
    unit_price: 0,
    cost: 0,
    markup: 0,
    vat_rate: 15.0,
    unit: 'pcs',
    unit_arabic: 'قطعة',
    sku: '',
    barcode: '',
    stock: 0,
    min_stock: 0,
    is_active: true,
    service_not_using_stock: false,
    price_includes_tax: true,
    price_change_allowed: false,
    image_url: '',
    color: '#000000'
  }

  // Tax rates (these would normally come from backend)
  const taxRates = [
    { id: 1, name: 'Standard VAT', rate: 15.0, is_default: true },
    { id: 2, name: 'Zero VAT', rate: 0.0, is_default: false },
    { id: 3, name: 'Exempt', rate: 0.0, is_default: false }
  ]

  // Units of measurement
  const units = [
    { value: 'pcs', label: 'Pieces', arabic: 'قطعة' },
    { value: 'kg', label: 'Kilograms', arabic: 'كيلوغرام' },
    { value: 'g', label: 'Grams', arabic: 'غرام' },
    { value: 'l', label: 'Liters', arabic: 'لتر' },
    { value: 'ml', label: 'Milliliters', arabic: 'مليلتر' },
    { value: 'm', label: 'Meters', arabic: 'متر' },
    { value: 'cm', label: 'Centimeters', arabic: 'سنتيمتر' },
    { value: 'box', label: 'Box', arabic: 'صندوق' },
    { value: 'pack', label: 'Pack', arabic: 'عبوة' }
  ]

  // Load data on mount
  onMount(async () => {
    await loadData()
  })

  async function loadData() {
    try {
      isLoading = true
      await Promise.all([
        loadProducts(),
        loadCategories(),
        loadDefaultProductSettings()
      ])
    } catch (error) {
      console.error('Error loading data:', error)
    } finally {
      isLoading = false
    }
  }

  async function loadProducts() {
    try {
      products = await GetProducts()
    } catch (error) {
      console.error('Error loading products:', error)
      products = [] // Ensure products is always an array
    }
  }

  async function loadCategories() {
    try {
      categories = await GetProductCategories()
    } catch (error) {
      console.error('Error loading categories:', error)
      categories = [] // Ensure categories is always an array
    }
  }

  async function loadDefaultProductSettings() {
    try {
      defaultProductSettings = await GetDefaultProductSettings()
    } catch (error) {
      console.error('Error loading default product settings:', error)
      defaultProductSettings = null
    }
  }

  // Product handlers
  function handleAddProduct() {
    editingProduct = null
    resetProductForm()
    applyDefaultSettings()
    showProductModal = true
  }

  function handleEditProduct(event) {
    editingProduct = event.detail
    productForm = { ...event.detail }
    showProductModal = true
  }

  function handleDeleteProduct(event) {
    // TODO: Implement delete functionality
    console.log('Delete product:', event.detail)
  }

  function handleCloseProductModal() {
    showProductModal = false
    resetProductForm()
  }

  async function handleSaveProduct() {
    try {
      isLoading = true
      const product = main.main.Product.createFrom(productForm)
      await CreateProduct(product)
      await loadProducts()
      showProductModal = false
      resetProductForm()
    } catch (error) {
      console.error('Error saving product:', error)
      alert('Error saving product: ' + error.message)
    } finally {
      isLoading = false
    }
  }

  function resetProductForm() {
    productForm = {
      name: '',
      name_arabic: '',
      description: '',
      description_arabic: '',
      category_id: 0,
      unit_price: 0,
      cost: 0,
      markup: 0,
      vat_rate: 15.0,
      unit: 'pcs',
      unit_arabic: 'قطعة',
      sku: '',
      barcode: '',
      stock: 0,
      min_stock: 0,
      is_active: true,
      service_not_using_stock: false,
      price_includes_tax: true,
      price_change_allowed: false,
      image_url: '',
      color: '#000000'
    }
  }

  function applyDefaultSettings() {
    if (defaultProductSettings) {
      // Set current stock from default stock setting
      productForm.stock = defaultProductSettings.default_stock || 0
      productForm.markup = defaultProductSettings.default_markup || 0
      productForm.is_active = defaultProductSettings.default_product_status !== false
      productForm.price_includes_tax = defaultProductSettings.default_price_includes_tax !== false
      productForm.price_change_allowed = defaultProductSettings.default_price_change_allowed !== false
      
      // Apply default tax rate if available
      if (defaultProductSettings.default_tax_rate_id) {
        const defaultTaxRate = taxRates.find(rate => rate.id === defaultProductSettings.default_tax_rate_id)
        if (defaultTaxRate) {
          productForm.vat_rate = defaultTaxRate.rate
        }
      }
      
      // Apply default unit if available
      if (defaultProductSettings.default_unit_id) {
        const defaultUnit = units.find(unit => unit.id === defaultProductSettings.default_unit_id)
        if (defaultUnit) {
          productForm.unit = defaultUnit.value
        }
      }
      
      // Apply default category if available
      if (defaultProductSettings.default_category_id) {
        productForm.category_id = defaultProductSettings.default_category_id
      }
    }
  }

  // Utility functions for product form
  function generateSKU() {
    const prefix = productForm.name.substring(0, 3).toUpperCase()
    const timestamp = Date.now().toString().slice(-6)
    productForm.sku = `${prefix}${timestamp}`
  }

  function generateBarcode() {
    const timestamp = Date.now().toString()
    productForm.barcode = timestamp.slice(-10)
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-gray-900 via-blue-900 to-purple-900 p-6">
  <div class="max-w-7xl mx-auto">
    <!-- Header -->
    <div class="glass-card mb-8">
      <div class="p-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="heading-1 mb-2">Product Management</h1>
            <p class="text-secondary">Manage your product inventory and categories</p>
          </div>
          {#if activeTab === 'products'}
            <button 
              class="btn-primary"
              on:click={handleAddProduct}
            >
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
              </svg>
              Add Product
            </button>
          {/if}
        </div>

        <!-- Horizontal Tabs -->
        <div class="mt-6">
          <HorizontalTabs 
            {tabs}
            {activeTab}
            variant="glass"
            showScrollButtons={false}
            on:tabChange={handleTabChange}
          />
        </div>
      </div>
    </div>

    <!-- Tab Content -->
    <div class="space-y-6">
      {#if activeTab === 'products'}
        <ProductList 
          {products}
          {isLoading}
          bind:searchTerm
          on:edit={handleEditProduct}
          on:delete={handleDeleteProduct}
        />
      {:else if activeTab === 'categories'}
        <ProductCategoriesTab 
          {categories}
          {isLoading}
          on:categoriesUpdated={loadCategories}
        />
      {/if}
    </div>
  </div>
</div>

<!-- Product Modal -->
<ProductModal 
  bind:show={showProductModal}
  {editingProduct}
  bind:productForm
  {categories}
  {taxRates}
  {units}
  {isLoading}
  {generateSKU}
  {generateBarcode}
  on:close={handleCloseProductModal}
  on:save={handleSaveProduct}
/>