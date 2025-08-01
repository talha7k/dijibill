<script>
  import { onMount } from 'svelte'
  import { CreateSupplier, GetSuppliers, UpdateSupplier, DeleteSupplier } from '../wailsjs/go/main/App.js'
  import PageLayout from './components/PageLayout.svelte'
  import DataTable from './components/DataTable.svelte'
  import SupplierModal from './components/SupplierModal.svelte'
  import StatusBadge from './components/StatusBadge.svelte'
  import ConfirmationModal from './components/ConfirmationModal.svelte'

  /** @type {Array<{id: number, company_name: string, company_name_arabic: string, contact_person: string, contact_person_arabic: string, email: string, phone: string, vat_number: string, address: string, address_arabic: string, city: string, city_arabic: string, country: string, country_arabic: string, payment_terms: string, active: boolean, created_at: string}>} */
  let suppliers = []

  /** @type {string} */
  let searchTerm = ''
  /** @type {boolean} */
  let showSupplierModal = false
  /** @type {any} */
  let editingSupplier = null
  /** @type {boolean} */
  let loading = false
  
  // Confirmation modal state
  let showDeleteConfirm = false
  let supplierToDelete = null
  let confirmLoading = false

  // Load suppliers on component mount
  onMount(async () => {
    await loadSuppliers()
  })

  async function loadSuppliers() {
    try {
      loading = true
      console.log('Loading suppliers...')
      const result = await GetSuppliers()
      console.log('Suppliers loaded:', result)
      suppliers = result || []
    } catch (error) {
      console.error('Error loading suppliers:', error)
      suppliers = []
      // Show user-friendly error message
      alert('Failed to load suppliers. Please check your connection and try again.')
    } finally {
      loading = false
      console.log('Loading finished. Suppliers count:', suppliers.length)
    }
  }

  // Filter suppliers based on search term
  $: filteredSuppliers = suppliers.filter(supplier => {
    if (!searchTerm) return true
    
    const term = searchTerm.toLowerCase()
    return (
      (supplier.company_name || '').toLowerCase().includes(term) ||
      (supplier.company_name_arabic || '').includes(searchTerm) ||
      (supplier.contact_person || '').toLowerCase().includes(term) ||
      (supplier.email || '').toLowerCase().includes(term) ||
      (supplier.phone || '').includes(searchTerm)
    )
  })

  function handleAddSupplier() {
    editingSupplier = null
    showSupplierModal = true
  }

  function handleEditSupplier(supplier) {
    editingSupplier = supplier
    showSupplierModal = true
  }

  function showDeleteConfirmation(supplier) {
    supplierToDelete = supplier;
    showDeleteConfirm = true;
  }

  function cancelDelete() {
    showDeleteConfirm = false;
    supplierToDelete = null;
  }

  async function confirmDelete() {
    if (!supplierToDelete) return;
    
    try {
      confirmLoading = true;
      await DeleteSupplier(supplierToDelete.id);
      await loadSuppliers(); // Reload the list
      showDeleteConfirm = false;
      supplierToDelete = null;
    } catch (error) {
      console.error('Error deleting supplier:', error);
      alert('Error deleting supplier. Please try again.');
    } finally {
      confirmLoading = false;
    }
  }

  async function handleDeleteSupplier(supplier) {
    showDeleteConfirmation(supplier);
  }

  async function handleSupplierSave(event) {
    loading = true
    const supplierData = event.detail

    try {
      if (editingSupplier) {
        // Update existing supplier
        await UpdateSupplier({ ...editingSupplier, ...supplierData })
      } else {
        // Add new supplier
        await CreateSupplier(supplierData)
      }
      
      await loadSuppliers() // Reload the list
      showSupplierModal = false
      editingSupplier = null
    } catch (error) {
      console.error('Error saving supplier:', error)
      alert('Error saving supplier. Please try again.')
    } finally {
      loading = false
    }
  }

  function handleSupplierModalClose() {
    showSupplierModal = false
    editingSupplier = null
    loading = false
  }

  function handleImport() {
    // TODO: Implement import functionality
    alert('Import functionality will be implemented soon')
  }

  function handleSearch(event) {
    searchTerm = event.detail.searchTerm
  }

  function formatPaymentTerms(terms) {
    const termMap = {
      'net_15': 'Net 15 days',
      'net_30': 'Net 30 days',
      'net_45': 'Net 45 days',
      'net_60': 'Net 60 days',
      'cash_on_delivery': 'Cash on Delivery',
      'advance_payment': 'Advance Payment'
    }
    return termMap[terms] || terms
  }

  /** @type {Array<{label: string, key?: string, labelArabic?: string, sortable?: boolean, class?: string, render?: Function, actions?: Array<{key: string, text: string, icon?: string, class?: string, title?: string}>}>} */
  const columns = [
    {
      key: 'company_name',
      label: 'Company Name',
      labelArabic: 'اسم الشركة',
      sortable: true,
      render: (supplier) => `
        <div>
          <div class="font-medium text-white">${supplier.company_name}</div>
          <div class="text-sm text-white/60">${supplier.company_name_arabic || ''}</div>
        </div>
      `
    },
    {
      key: 'contact_person',
      label: 'Contact Person',
      labelArabic: 'الشخص المسؤول',
      sortable: true,
      render: (supplier) => `
        <div>
          <div class="font-medium text-white">${supplier.contact_person}</div>
          <div class="text-sm text-white/60">${supplier.email}</div>
        </div>
      `
    },
    {
      key: 'phone',
      label: 'Phone',
      labelArabic: 'الهاتف',
      sortable: true
    },
    {
      key: 'city',
      label: 'City',
      labelArabic: 'المدينة',
      sortable: true,
      render: (supplier) => `
        <div>
          <div class="text-white">${supplier.city}</div>
          <div class="text-sm text-white/60">${supplier.city_arabic || ''}</div>
        </div>
      `
    },
    {
      key: 'payment_terms',
      label: 'Payment Terms',
      labelArabic: 'شروط الدفع',
      sortable: true,
      render: (supplier) => formatPaymentTerms(supplier.payment_terms)
    },
    {
      key: 'active',
      label: 'Status',
      labelArabic: 'الحالة',
      sortable: true,
      render: (supplier) => supplier.active 
        ? '<span class="status-badge status-badge-success">Active</span>'
        : '<span class="status-badge status-badge-error">Inactive</span>'
    },
    {
      label: 'Actions',
      actions: [
        { key: 'edit', text: 'Edit', icon: 'edit', class: 'btn-primary' },
        { key: 'delete', text: 'Delete', icon: 'delete', class: 'btn-error' }
      ]
    }
  ]
</script>

<PageLayout
  title="Suppliers"
  icon="fa-truck"
>
  <svelte:fragment slot="actions">
    <button class="btn btn-secondary-outline" on:click={handleImport}>
      <i class="fas fa-upload"></i>
      Import
    </button>
    <button class="btn btn-primary" on:click={handleAddSupplier}>
      <i class="fas fa-plus"></i>
      Add Supplier
    </button>
  </svelte:fragment>

  <DataTable
    data={filteredSuppliers}
    {columns}
    {loading}
    searchPlaceholder="Search suppliers..."
    emptyStateTitle="No suppliers found"
    emptyStateMessage="Start by adding your first supplier"
    emptyStateIcon="truck"
    primaryAction={{ text: 'Add Supplier', icon: 'add' }}
    on:primary-action={handleAddSupplier}
    on:search={handleSearch}
    on:row-action={(e) => {
      const { action, item } = e.detail
      if (action === 'edit') {
        handleEditSupplier(item)
      } else if (action === 'delete') {
        handleDeleteSupplier(item)
      }
    }}
  />
</PageLayout>

<SupplierModal
  show={showSupplierModal}
  {editingSupplier}
  {loading}
  on:save={handleSupplierSave}
  on:close={handleSupplierModalClose}
/>

<!-- Confirmation Modal -->
<ConfirmationModal
  show={showDeleteConfirm}
  title="Delete Supplier"
  message={supplierToDelete ? `Are you sure you want to delete "${supplierToDelete.company_name}"? This action cannot be undone.` : ''}
  confirmText="Delete"
  cancelText="Cancel"
  loading={confirmLoading}
  on:confirm={confirmDelete}
  on:cancel={cancelDelete}
/>