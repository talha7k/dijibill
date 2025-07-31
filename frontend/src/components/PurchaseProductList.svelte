<script>
  import { createEventDispatcher } from 'svelte'
  import DataTable from './DataTable.svelte'
  import ConfirmationModal from './ConfirmationModal.svelte'
  
  /** @type {Array<any>} */
  export let purchaseProducts = []
  /** @type {boolean} */
  export let loading = false
  /** @type {string} */
  export let searchTerm = ''

  let showDeleteConfirm = false
  let purchaseProductToDelete = null
  let confirmLoading = false

  const dispatch = createEventDispatcher()

  function editPurchaseProduct(purchaseProduct) {
    dispatch('edit', purchaseProduct)
  }

  function showDeleteConfirmation(purchaseProduct) {
    purchaseProductToDelete = purchaseProduct
    showDeleteConfirm = true
  }

  function cancelDelete() {
    showDeleteConfirm = false
    purchaseProductToDelete = null
  }

  function confirmDelete() {
    if (purchaseProductToDelete) {
      confirmLoading = true
      dispatch('delete', purchaseProductToDelete)
      showDeleteConfirm = false
      purchaseProductToDelete = null
      confirmLoading = false
    }
  }

  function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'SAR'
    }).format(amount)
  }

  $: filteredPurchaseProducts = Array.isArray(purchaseProducts) ? purchaseProducts.filter(product => {
    if (!product) return false
    const searchLower = (searchTerm || '').toLowerCase()
    if (!searchLower) return true
    
    return (
      (product.name || '').toLowerCase().includes(searchLower) ||
      (product.sku || '').toLowerCase().includes(searchLower) ||
      (product.barcode || '').toLowerCase().includes(searchLower)
    )
  }) : []

  // Table configuration
  /** @type {Array<{label: string, key?: string, class?: string, render?: Function, actions?: Array<{key: string, text: string, icon?: string, class?: string, title?: string}>}>} */
  const columns = [
    { 
      key: 'name', 
      label: 'Purchase Product',
      render: (item) => {
        let html = `<div class="flex items-center gap-3">`
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
    { key: 'unit_price', label: 'Unit Price', render: (item) => formatCurrency(item.unit_price) },
    { 
      key: 'unit', 
      label: 'Unit',
      render: (item) => {
        let html = `<div class="flex flex-col">`
        html += `<span class="text-white/80">${item.unit}</span>`
        if (item.unit_arabic) {
          html += `<span class="text-xs text-white/60" dir="rtl">${item.unit_arabic}</span>`
        }
        html += `</div>`
        return html
      }
    },
    { 
      key: 'vat_rate', 
      label: 'VAT Rate',
      render: (item) => `${item.vat_rate}%`
    },
    { key: 'is_active', label: 'Status' },
    {
      label: 'Actions',
      actions: [
        { key: 'edit', text: 'Edit', icon: 'fa-edit', class: 'btn-warning', title: 'Edit Purchase Product' },
        { key: 'delete', text: 'Delete', icon: 'fa-trash', class: 'btn-error', title: 'Delete Purchase Product' }
      ]
    }
  ]

  // Primary and secondary actions
  /** @type {{text: string, icon: string}} */
  const primaryAction = {
    text: 'Add Purchase Product',
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
      editPurchaseProduct(item)
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
  data={filteredPurchaseProducts}
  {columns}
  {loading}
  {searchTerm}
  searchPlaceholder="Search purchase products by name, SKU, or barcode..."
  emptyStateTitle="No purchase products found"
  emptyStateMessage="Get started by adding your first purchase product"
  emptyStateIcon="fa-shopping-cart"
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
    title="Delete Purchase Product"
    message={`Are you sure you want to delete "${purchaseProductToDelete?.name}"? This action cannot be undone.`}
    confirmText="Delete"
    cancelText="Cancel"
    confirmClass="btn-error"
    loading={confirmLoading}
    on:confirm={confirmDelete}
    on:cancel={cancelDelete}
  />
{/if}