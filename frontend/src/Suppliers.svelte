<script>
  import PageLayout from './components/PageLayout.svelte'
  import DataTable from './components/DataTable.svelte'
  import SupplierModal from './components/SupplierModal.svelte'
  import StatusBadge from './components/StatusBadge.svelte'

  let suppliers = [
    {
      id: 1,
      company_name: 'Tech Solutions Ltd',
      company_name_arabic: 'شركة الحلول التقنية المحدودة',
      contact_person: 'Ahmed Al-Rashid',
      contact_person_arabic: 'أحمد الراشد',
      email: 'ahmed@techsolutions.com',
      phone: '+966 11 234 5678',
      vat_number: '300123456789003',
      address: '123 King Fahd Road, Riyadh',
      address_arabic: '123 طريق الملك فهد، الرياض',
      city: 'Riyadh',
      city_arabic: 'الرياض',
      country: 'SA',
      country_arabic: 'المملكة العربية السعودية',
      payment_terms: 'net_30',
      active: true,
      created_at: '2024-01-15'
    },
    {
      id: 2,
      company_name: 'Global Supplies Co',
      company_name_arabic: 'شركة الإمدادات العالمية',
      contact_person: 'Sarah Johnson',
      contact_person_arabic: 'سارة جونسون',
      email: 'sarah@globalsupplies.com',
      phone: '+966 12 345 6789',
      vat_number: '300987654321003',
      address: '456 Business District, Jeddah',
      address_arabic: '456 الحي التجاري، جدة',
      city: 'Jeddah',
      city_arabic: 'جدة',
      country: 'SA',
      country_arabic: 'المملكة العربية السعودية',
      payment_terms: 'net_15',
      active: true,
      created_at: '2024-02-01'
    }
  ]

  let searchTerm = ''
  let showSupplierModal = false
  let editingSupplier = null
  let loading = false

  // Filter suppliers based on search term
  $: filteredSuppliers = suppliers.filter(supplier => 
    supplier.company_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    supplier.company_name_arabic.includes(searchTerm) ||
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

  function handleDeleteSupplier(supplier) {
    if (confirm(`Are you sure you want to delete ${supplier.company_name}?`)) {
      suppliers = suppliers.filter(s => s.id !== supplier.id)
    }
  }

  function handleSupplierSave(event) {
    loading = true
    const supplierData = event.detail

    setTimeout(() => {
      if (editingSupplier) {
        // Update existing supplier
        const index = suppliers.findIndex(s => s.id === editingSupplier.id)
        if (index !== -1) {
          suppliers[index] = { ...suppliers[index], ...supplierData }
          suppliers = [...suppliers]
        }
      } else {
        // Add new supplier
        const newSupplier = {
          id: Math.max(...suppliers.map(s => s.id)) + 1,
          ...supplierData,
          created_at: new Date().toISOString().split('T')[0]
        }
        suppliers = [...suppliers, newSupplier]
      }
      
      showSupplierModal = false
      editingSupplier = null
      loading = false
    }, 1000)
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
    searchTerm = event.detail
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
      actions: [
        { key: 'edit', text: 'Edit', icon: 'fa-edit', class: 'btn-primary' },
        { key: 'delete', text: 'Delete', icon: 'fa-trash', class: 'btn-danger' }
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