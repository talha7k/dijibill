<script>
  import { createEventDispatcher } from 'svelte'
  import { CreatePurchaseProductCategory } from '../../wailsjs/go/main/App.js'
  import { database } from '../../wailsjs/go/models'
  import DataTable from './DataTable.svelte'
  import PurchaseProductCategoryModal from './PurchaseProductCategoryModal.svelte'
  import { showDbSuccess, showDbError } from '../stores/notificationStore.js'

  // Props
  export let categories = []
  export let isLoading = false

  // Event dispatcher
  const dispatch = createEventDispatcher()

  // State
  let showPurchaseProductCategoryModal = false
  let searchTerm = ''
  let categoryForm = {
    name: '',
    name_arabic: '',
    description: '',
    description_arabic: ''
  }

  // DataTable configuration
  const columns = [
    {
      key: 'name',
      label: 'Category Name',
      sortable: true,
      searchable: true
    },
    {
      key: 'name_arabic',
      label: 'Arabic Name',
      sortable: true,
      searchable: true
    },
    {
      key: 'description',
      label: 'Description',
      sortable: false,
      searchable: true
    },
    {
      key: 'description_arabic',
      label: 'Arabic Description',
      sortable: false,
      searchable: true
    }
  ]

  const primaryAction = {
    label: 'Add Category',
    icon: 'M12 6v6m0 0v6m0-6h6m-6 0H6'
  }

  // Computed
  $: filteredCategories = categories.filter(category => {
    if (!searchTerm) return true
    const term = searchTerm.toLowerCase()
    return (
      category.name?.toLowerCase().includes(term) ||
      category.name_arabic?.toLowerCase().includes(term) ||
      category.description?.toLowerCase().includes(term) ||
      category.description_arabic?.toLowerCase().includes(term)
    )
  })

  // Event handlers
  function handlePrimaryAction() {
    console.log('Primary action clicked - opening modal')
    handleAddCategory()
  }

  function handleSearch(event) {
    searchTerm = event.detail.searchTerm
  }

  function handleAddCategory() {
    console.log('handleAddCategory called')
    resetCategoryForm()
    showPurchaseProductCategoryModal = true
    console.log('showPurchaseProductCategoryModal set to:', showPurchaseProductCategoryModal)
  }

  function handleClosePurchaseProductCategoryModal() {
    console.log('handleClosePurchaseProductCategoryModal called')
    showPurchaseProductCategoryModal = false
    resetCategoryForm()
  }

  async function handleSaveCategory() {
    try {
      isLoading = true
      console.log('Saving purchase product category:', categoryForm)
      const category = database.PurchaseProductCategory.createFrom(categoryForm)
      await CreatePurchaseProductCategory(category)
      
      dispatch('categoriesUpdated')
      showPurchaseProductCategoryModal = false
      resetCategoryForm()
    } catch (error) {
      console.error('Error saving purchase product category:', error)
      showDbError('Error saving purchase product category: ' + error.message)
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
</script>

<DataTable
  data={filteredCategories}
  {columns}
  loading={isLoading}
  {primaryAction}
  searchPlaceholder="Search purchase product categories..."
  emptyStateTitle="No purchase product categories found"
  emptyStateMessage="Get started by creating your first purchase product category."
  on:primary-action={handlePrimaryAction}
  on:search={handleSearch}
/>

<!-- Purchase Product Category Modal -->
<PurchaseProductCategoryModal
  show={showPurchaseProductCategoryModal}
  bind:categoryForm
  {isLoading}
  on:close={handleClosePurchaseProductCategoryModal}
  on:save={handleSaveCategory}
/>