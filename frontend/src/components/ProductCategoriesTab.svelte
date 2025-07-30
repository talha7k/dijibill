<script>
  import { createEventDispatcher } from 'svelte'
  import { CreateProductCategory } from '../../wailsjs/go/main/App.js'
  import { database } from '../../wailsjs/go/models'
  import DataTable from './DataTable.svelte'
  import ProductCategoryModal from './ProductCategoryModal.svelte'

  // Props
  export let categories = []
  export let isLoading = false

  // Event dispatcher
  const dispatch = createEventDispatcher()

  // State
  let showProductCategoryModal = false
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
    showProductCategoryModal = true
    console.log('showProductCategoryModal set to:', showProductCategoryModal)
  }

  function handleCloseProductCategoryModal() {
    console.log('handleCloseProductCategoryModal called')
    showProductCategoryModal = false
    resetCategoryForm()
  }

  async function handleSaveCategory() {
    try {
      isLoading = true
      console.log('Saving category:', categoryForm)
      const category = database.ProductCategory.createFrom(categoryForm)
      await CreateProductCategory(category)
      
      dispatch('categoriesUpdated')
      showProductCategoryModal = false
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
</script>

<DataTable
  data={filteredCategories}
  {columns}
  loading={isLoading}
  {primaryAction}
  searchPlaceholder="Search categories..."
  emptyStateTitle="No categories found"
  emptyStateMessage="Get started by creating your first category."
  on:primary-action={handlePrimaryAction}
  on:search={handleSearch}
/>

<!-- Category Modal -->
<ProductCategoryModal
  show={showProductCategoryModal}
  bind:categoryForm
  {isLoading}
  on:close={handleCloseProductCategoryModal}
  on:save={handleSaveCategory}
/>