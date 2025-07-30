<script>
  import { onMount } from 'svelte'
  import { CreateSupplier, GetSuppliers, UpdateSupplier, DeleteSupplier } from '../wailsjs/go/main/App.js'
  import PageLayout from './components/PageLayout.svelte'
  import DataTable from './components/DataTable.svelte'
  import SupplierModal from './components/SupplierModal.svelte'
  import StatusBadge from './components/StatusBadge.svelte'

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

  // Load suppliers on component mount
  onMount(async () => {
    await loadSuppliers()
  })

  async function loadSuppliers() {
    try {
      loading = true
      suppliers = await GetSuppliers()
    } catch (error) {
      console.error('Error loading suppliers:', error)
      suppliers = []
    } finally {
      loading = false
    }
  }

  // Filter suppliers based on search term
  $: filteredSuppliers = suppliers.filter(supplier => 
    supplier.company_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (supplier.company_name_arabic && supplier.company_name_arabic.includes(searchTerm)) ||
    supplier.contact_person.toLowerCase().includes(searchTerm.toLowerCase()) ||
    supplier.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
    supplier.phone.includes(searchTerm)
  )

  function handleAddSupplier() {
    editingSupplier = null
    showSupplierModal = true
  }

  function handleEditSupplier(supplier) {
    editingSupplier = supplier
    showSupplierModal = true
  }

  async function handleDeleteSupplier(supplier) {
    if (confirm(`Are you sure you want to delete ${supplier.company_name}?`)) {
      try {
        loading = true
        await DeleteSupplier(supplier.id)
        await loadSuppliers() // Reload the list
      } catch (error) {
        console.error('Error deleting supplier:', error)
        alert('Error deleting supplier. Please try again.')
      } finally {
        loading = false
      }
    }
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
        { key: 'edit', text: 'Edit', icon: 'fa-edit', class: 'btn-primary' },
        { key: 'delete', text: 'Delete', icon: 'fa-trash', class: 'btn-error' }
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
    emptyStateIcon="fa-truck"
    primaryAction={{ text: 'Add Supplier', icon: 'fa-plus' }}
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