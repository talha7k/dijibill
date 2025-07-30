<script>
  import { onMount } from 'svelte'
  import { GetProducts, CreateProduct, GetProductCategories, CreateProductCategory } from '../wailsjs/go/main/App.js'
  import { main } from '../wailsjs/go/models'
  
  // Import components
  import ProductList from './components/ProductList.svelte'
  import ProductModal from './components/ProductModal.svelte'
  import CategoryModal from './components/CategoryModal.svelte'

  // State variables
  let products = []
  let categories = []
  let isLoading = false
  let searchTerm = ''
  let showProductModal = false
  let showCategoryModal = false
  let editingProduct = null
  let editingCategory = null

  // Form data
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
    default_quantity: false,
    service_not_using_stock: false,
    price_includes_tax: true,
    price_change_allowed: false,
    image_url: '',
    color: '#000000'
  }

  let categoryForm = {
    name: '',
    name_arabic: '',
    description: '',
    description_arabic: ''
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
    await loadProducts()
    await loadCategories()
  })

  async function loadProducts() {
    try {
      isLoading = true
      products = await GetProducts()
    } catch (error) {
      console.error('Error loading products:', error)
    } finally {
      isLoading = false
    }
  }

  async function loadCategories() {
    try {
      categories = await GetProductCategories()
    } catch (error) {
      console.error('Error loading categories:', error)
    }
  }

  // Product handlers
  function handleAddProduct() {
    editingProduct = null
    resetProductForm()
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
      const product = main.Product.createFrom(productForm)
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
      default_quantity: false,
      service_not_using_stock: false,
      price_includes_tax: true,
      price_change_allowed: false,
      image_url: '',
      color: '#000000'
    }
  }

  // Category handlers
  function handleAddCategory() {
    editingCategory = null
    resetCategoryForm()
    showCategoryModal = true
  }

  function handleCloseCategoryModal() {
    showCategoryModal = false
    resetCategoryForm()
  }

  async function handleSaveCategory() {
    try {
      isLoading = true
      const category = main.ProductCategory.createFrom(categoryForm)
      await CreateProductCategory(category)
      await loadCategories()
      showCategoryModal = false
      resetCategoryForm()
    } catch (error) {
      console.error('Error saving category:', error)
      alert('Error saving category: ' + error.message)
    } finally {
      isLoading = false
    }
  }

  function resetCategoryForm() {
    categoryForm = {
      name: '',
      name_arabic: '',
      description: '',
      description_arabic: ''
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
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-bold text-white mb-2">Products</h1>
        <p class="text-white/60">Manage your product inventory and categories</p>
      </div>
      <div class="flex gap-3">
        <button 
          class="btn btn-outline text-white border-white/20 hover:bg-white/10"
          on:click={handleAddCategory}
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
          </svg>
          Add Category
        </button>
      </div>
    </div>

    <!-- Product List -->
    <ProductList 
      {products}
      {isLoading}
      bind:searchTerm
      on:add={handleAddProduct}
      on:edit={handleEditProduct}
      on:delete={handleDeleteProduct}
    />
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

<!-- Category Modal -->
<CategoryModal 
  bind:show={showCategoryModal}
  {editingCategory}
  bind:categoryForm
  {isLoading}
  on:close={handleCloseCategoryModal}
  on:save={handleSaveCategory}
/>