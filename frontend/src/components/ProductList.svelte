<script>
  import { createEventDispatcher } from 'svelte'
  import DataTable from './DataTable.svelte'
  import ConfirmationModal from './ConfirmationModal.svelte'
  
  /** @type {Array<any>} */
  export let products = []
  /** @type {boolean} */
  export let loading = false
  /** @type {string} */
  export let searchTerm = ''

  let showDeleteConfirm = false
  let productToDelete = null
  let confirmLoading = false

  const dispatch = createEventDispatcher()

  function editProduct(product) {
    dispatch('edit', product)
  }

  function showDeleteConfirmation(product) {
    productToDelete = product
    showDeleteConfirm = true
  }

  function cancelDelete() {
    showDeleteConfirm = false
    productToDelete = null
  }

  function confirmDelete() {
    if (productToDelete) {
      confirmLoading = true
      dispatch('delete', productToDelete)
      showDeleteConfirm = false
      productToDelete = null
      confirmLoading = false
    }
  }

  function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'SAR'
    }).format(amount)
  }

  $: filteredProducts = (products || []).filter(product => 
    product.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.sku?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.barcode?.toLowerCase().includes(searchTerm.toLowerCase())
  )

  // Table configuration
  /** @type {Array<{label: string, key?: string, class?: string, render?: Function, actions?: Array<{key: string, text: string, icon?: string, class?: string, title?: string}>}>} */
  const columns = [
    { 
      key: 'image', 
      label: 'Image',
      render: (item) => {
        if (item.image_url) {
          return `<div class="avatar"><div class="w-12 h-12 rounded-lg"><img src="${item.image_url}" alt="${item.name}" /></div></div>`
        }
        return `<div class="w-12 h-12 bg-white/10 rounded-lg flex items-center justify-center"><i class="fas fa-image text-white/30 text-xl"></i></div>`
      }
    },
    { 
      key: 'name', 
      label: 'Product',
      render: (item) => {
        let html = `<div class="flex items-center gap-3">`
        if (item.color) {
          html += `<div class="w-4 h-4 rounded-full border border-white/20" style="background-color: ${item.color}"></div>`
        }
        html += `<div><div class="font-bold text-white">${item.name}</div>`
        if (item.name_arabic) {
          html += `<div class="text-sm text-white/60" dir="rtl">${item.name_arabic}</div>`
        }
        html += `</div></div>`
        return html
      }
    },
    { key: 'sku', label: 'SKU', render: (item) => item.sku || '-' },
    { key: 'category_name', label: 'Category', render: (item) => item.category_name || 'Uncategorized' },
    { key: 'unit_price', label: 'Price', render: (item) => formatCurrency(item.unit_price) },
    { 
      key: 'stock', 
      label: 'Stock',
      render: (item) => {
        let html = `<div class="flex items-center gap-2"><span class="text-white/80">${item.stock}</span>`
        if (item.stock <= item.min_stock) {
          html += `<div class="badge badge-error badge-sm">Low</div>`
        }
        if (item.service_not_using_stock) {
          html += `<div class="badge badge-info badge-sm">Service</div>`
        }
        html += `</div>`
        return html
      }
    },
    { key: 'is_active', label: 'Status' },
    {
      label: 'Actions',
      actions: [
        { key: 'edit', text: 'Edit', icon: 'fa-edit', class: 'btn-warning', title: 'Edit Product' },
        { key: 'delete', text: 'Delete', icon: 'fa-trash', class: 'btn-error', title: 'Delete Product' }
      ]
    }
  ]

  // Primary and secondary actions
  /** @type {{text: string, icon: string}} */
  const primaryAction = {
    text: 'Add Product',
    icon: 'fa-plus'
  }

  /** @type {Array<{text: string, icon: string}>} */
  const secondaryActions = [
    {
      text: 'Reports',
      icon: 'fa-chart-bar'
    }
  ]

  function handlePrimaryAction() {
    dispatch('add')
  }

  function handleSecondaryAction(event) {
    const { action } = event.detail
    if (action.text === 'Reports') {
      dispatch('reports')
    }
  }

  function handleRowAction(event) {
    const { action, item } = event.detail
    if (action === 'edit') {
      editProduct(item)
    } else if (action === 'delete') {
      showDeleteConfirmation(item)
    }
  }

  function handleSearch(event) {
    searchTerm = event.detail.searchTerm
    dispatch('search', { searchTerm })
  }
</script>

<DataTable
  data={filteredProducts}
  {columns}
  {loading}
  {searchTerm}
  searchPlaceholder="Search products by name, SKU, or barcode..."
  emptyStateTitle="No products found"
  emptyStateMessage="Get started by adding your first product"
  emptyStateIcon="fa-box"
  showSearch={true}
  showActions={true}
  {primaryAction}
  {secondaryActions}
  on:primary-action={handlePrimaryAction}
  on:secondary-action={handleSecondaryAction}
  on:row-action={handleRowAction}
  on:search={handleSearch}
>
  <svelte:fragment slot="cell" let:item let:column>
    {#if column.key === 'is_active'}
      <div class="badge {item.is_active ? 'badge-success' : 'badge-error'} badge-sm">
        {item.is_active ? 'Active' : 'Inactive'}
      </div>
    {:else if column.render}
      {@html column.render(item)}
    {:else}
      {item[column.key] || '-'}
    {/if}
  </svelte:fragment>
</DataTable>

{#if showDeleteConfirm}
  <ConfirmationModal
    show={showDeleteConfirm}
    title="Delete Product"
    message={`Are you sure you want to delete "${productToDelete?.name}"? This action cannot be undone.`}
    confirmText="Delete"
    cancelText="Cancel"
    confirmClass="btn-error"
    loading={confirmLoading}
    on:confirm={confirmDelete}
    on:cancel={cancelDelete}
  />
{/if}