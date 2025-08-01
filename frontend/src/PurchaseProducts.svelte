<script>
  import { onMount } from 'svelte'
  import { GetPurchaseProducts, CreatePurchaseProduct, UpdatePurchaseProduct, DeletePurchaseProduct, GetPurchaseProductCategories, CreatePurchaseProductCategory } from '../wailsjs/go/main/App.js'
  import {database} from '../wailsjs/go/models'
  
  // Import components
  import PurchaseProductList from './components/PurchaseProductList.svelte'
  import PurchaseProductCategoriesTab from './components/PurchaseProductCategoriesTab.svelte'
  import PurchaseProductModal from './components/PurchaseProductModal.svelte'
  import HorizontalTabs from './components/HorizontalTabs.svelte'
  import PageLayout from './components/PageLayout.svelte'
  import ConfirmationModal from './components/ConfirmationModal.svelte'

  // State variables
  let purchaseProducts = []
  let categories = []
  let isLoading = true
  let searchTerm = ''
  let showPurchaseProductModal = false
  let editingPurchaseProduct = null
  let activeTab = 'purchase-products'
  
  // Confirmation modal state
  let showDeleteConfirm = false
  let purchaseProductToDelete = null
  let confirmLoading = false

  // Tab configuration
  const tabs = [
    { id: 'purchase-products', label: 'Purchase Products', icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4' },
    { id: 'categories', label: 'Categories', icon: 'M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10' }
  ]

  // Handle tab change
  function handleTabChange(event) {
    activeTab = event.detail.tabId
  }

  // Form data with default values
  let purchaseProductForm = {
    name: '',
    name_arabic: '',
    description: '',
    description_arabic: '',
    category_id: 0,
    unit_price: 0,
    vat_rate: 15.0,
    unit: 'pcs',
    unit_arabic: 'قطعة',
    sku: '',
    barcode: '',
    is_active: true,
    notes: '',
    notes_arabic: ''
  }

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
        loadPurchaseProducts(),
        loadCategories()
      ])
    } catch (error) {
      console.error('Error loading data:', error)
    } finally {
      isLoading = false
    }
  }

  async function loadPurchaseProducts() {
    try {
      const result = await GetPurchaseProducts()
      purchaseProducts = Array.isArray(result) ? result : []
    } catch (error) {
      console.error('Error loading purchase products:', error)
      purchaseProducts = [] // Ensure purchaseProducts is always an array
    }
  }

  async function loadCategories() {
    try {
      const result = await GetPurchaseProductCategories()
      categories = Array.isArray(result) ? result : []
    } catch (error) {
      console.error('Error loading categories:', error)
      categories = [] // Ensure categories is always an array
    }
  }

  // Purchase product handlers
  function handleAddPurchaseProduct() {
    editingPurchaseProduct = null
    resetPurchaseProductForm()
    showPurchaseProductModal = true
  }

  function handlePurchaseProductReports() {
    // TODO: Implement purchase product reports functionality
    console.log('Purchase product reports functionality - to be implemented')
    alert('Purchase product reports functionality will be implemented soon!')
  }

  function handleEditPurchaseProduct(event) {
    editingPurchaseProduct = event.detail
    purchaseProductForm = { ...event.detail }
    showPurchaseProductModal = true
  }

  function handleDeletePurchaseProduct(event) {
    const purchaseProduct = event.detail;
    showDeleteConfirmation(purchaseProduct);
  }

  function showDeleteConfirmation(purchaseProduct) {
    purchaseProductToDelete = purchaseProduct;
    showDeleteConfirm = true;
  }

  function cancelDelete() {
    showDeleteConfirm = false;
    purchaseProductToDelete = null;
  }

  async function confirmDelete() {
    if (!purchaseProductToDelete) return;
    
    try {
      confirmLoading = true;
      await DeletePurchaseProduct(purchaseProductToDelete.id);
      await loadPurchaseProducts();
      showDeleteConfirm = false;
      purchaseProductToDelete = null;
    } catch (error) {
      console.error('Error deleting purchase product:', error);
      alert('Error deleting purchase product: ' + error.message);
    } finally {
      confirmLoading = false;
    }
  }

  function handleClosePurchaseProductModal() {
    showPurchaseProductModal = false
    resetPurchaseProductForm()
  }

  async function handleSavePurchaseProduct() {
    try {
      isLoading = true
      const purchaseProduct = database.PurchaseProduct.createFrom(purchaseProductForm)
      
      if (editingPurchaseProduct) {
        // We're editing an existing purchase product, so update it
        await UpdatePurchaseProduct(purchaseProduct)
      } else {
        // We're creating a new purchase product
        await CreatePurchaseProduct(purchaseProduct)
      }
      
      await loadPurchaseProducts()
      showPurchaseProductModal = false
      resetPurchaseProductForm()
    } catch (error) {
      console.error('Error saving purchase product:', error)
      alert('Error saving purchase product: ' + error.message)
    } finally {
      isLoading = false
    }
  }

  function resetPurchaseProductForm() {
    purchaseProductForm = {
      name: '',
      name_arabic: '',
      description: '',
      description_arabic: '',
      category_id: 0,
      unit_price: 0,
      vat_rate: 15.0,
      unit: 'pcs',
      unit_arabic: 'قطعة',
      sku: '',
      barcode: '',
      is_active: true,
      notes: '',
      notes_arabic: ''
    }
  }

  // Utility functions for purchase product form
  function generateSKU() {
    const prefix = purchaseProductForm.name.substring(0, 3).toUpperCase()
    const timestamp = Date.now().toString().slice(-6)
    purchaseProductForm.sku = `${prefix}${timestamp}`
  }

  function generateBarcode() {
    const timestamp = Date.now().toString()
    purchaseProductForm.barcode = timestamp.slice(-10)
  }
</script>

<PageLayout 
  title="Purchase Product Management" 
  icon="cart" 
  showIndicator={true}
>
  <!-- Horizontal Tabs -->
  <div class="mb-6">
    <HorizontalTabs 
      {tabs}
      {activeTab}
      variant="glass"
      showScrollButtons={false}
      on:tabChange={handleTabChange}
    />
  </div>

  <!-- Tab Content -->
  <div class="space-y-6">
    {#if activeTab === 'purchase-products'}
      <PurchaseProductList 
        {purchaseProducts}
        loading={isLoading}
        bind:searchTerm
        on:edit={handleEditPurchaseProduct}
        on:delete={handleDeletePurchaseProduct}
        on:add={handleAddPurchaseProduct}
        on:reports={handlePurchaseProductReports}
      />
    {:else if activeTab === 'categories'}
      <PurchaseProductCategoriesTab 
        {categories}
        {isLoading}
        on:categoriesUpdated={loadCategories}
      />
    {/if}
  </div>
</PageLayout>

<!-- Purchase Product Modal -->
<PurchaseProductModal 
  bind:show={showPurchaseProductModal}
  {editingPurchaseProduct}
  bind:purchaseProductForm
  {categories}
  {units}
  {isLoading}
  {generateSKU}
  {generateBarcode}
  on:close={handleClosePurchaseProductModal}
  on:save={handleSavePurchaseProduct}
/>

<!-- Confirmation Modal -->
<ConfirmationModal
  show={showDeleteConfirm}
  title="Delete Purchase Product"
  message={purchaseProductToDelete ? `Are you sure you want to delete "${purchaseProductToDelete.name}"? This action cannot be undone.` : ''}
  confirmText="Delete"
  cancelText="Cancel"
  loading={confirmLoading}
  on:confirm={confirmDelete}
  on:cancel={cancelDelete}
/>