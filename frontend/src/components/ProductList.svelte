<script>
  import { createEventDispatcher } from 'svelte'
  import DataTable from './DataTable.svelte'
  
  /** @type {Array<any>} */
  export let products = []
  /** @type {boolean} */
  export let loading = false
  /** @type {string} */
  export let searchTerm = ''

  const dispatch = createEventDispatcher()

  function editProduct(product) {
    dispatch('edit', product)
  }

  function deleteProduct(product) {
    if (confirm(`Are you sure you want to delete "${product.name}"?`)) {
      dispatch('delete', product)
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
        return `<div class="w-12 h-12 bg-white/10 rounded-lg flex items-center justify-center"><svg class="w-6 h-6 text-white/30" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg></div>`
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
      deleteProduct(item)
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