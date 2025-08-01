<script>
  import { GetUnitsOfMeasurement, CreateUnitOfMeasurement, UpdateUnitOfMeasurement, DeleteUnitOfMeasurement } from '../../wailsjs/go/main/App.js'
  import { database } from '../../wailsjs/go/models'
  import { createEventDispatcher } from 'svelte'
  import DataTable from './DataTable.svelte'
  import UnitModal from './UnitModal.svelte'
  import ConfirmationModal from './ConfirmationModal.svelte'
  import { showDbSuccess, showDbError } from '../stores/notificationStore.js'

  const dispatch = createEventDispatcher()

  export let units = []
  export let isLoading = false

  let showUnitModal = false
  let editingUnit = null
  let unitToEdit = null
  let searchTerm = ''
  
  // Confirmation modal state
  let showDeleteConfirm = false
  let unitToDelete = null
  let loading = false

  // DataTable configuration
  const columns = [
    {
      key: 'label',
      label: 'Label',
      labelArabic: 'التسمية',
      sortable: true,
      render: (unit) => `
        <div>
          <div class="font-medium text-gray-100">${unit.label}</div>
          ${unit.arabic ? `<div class="text-sm text-gray-300" dir="rtl">${unit.arabic}</div>` : ''}
        </div>
      `
    },
    {
      key: 'value',
      label: 'Value',
      labelArabic: 'القيمة',
      sortable: true,
      render: (unit) => `<span class="font-mono text-sm bg-gray-100 px-2 py-1 rounded">${unit.value}</span>`
    },
    {
      key: 'status',
      label: 'Status',
      labelArabic: 'الحالة',
      render: (unit) => `
        <div class="flex items-center space-x-2">
          ${unit.is_default ? '<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">Default</span>' : ''}
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${unit.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
            ${unit.is_active ? 'Active' : 'Inactive'}
          </span>
        </div>
      `
    },
    {
      label: 'Actions',
      labelArabic: 'الإجراءات',
      actions: [
        {
          key: 'edit',
          text: 'Edit',
          icon: 'edit',
          class: 'btn-secondary',
          title: 'Edit unit'
        },
        {
          key: 'delete',
          text: 'Delete',
          icon: 'delete',
          class: 'btn-error',
          title: 'Delete unit'
        }
      ]
    }
  ]

  const primaryAction = {
    label: 'Add Unit',
    icon: 'add'
  }

  // Reactive declarations
  $: filteredUnits = units.filter(unit => {
    if (!searchTerm) return true
    const search = searchTerm.toLowerCase()
    return (
      unit.label?.toLowerCase().includes(search) ||
      unit.arabic?.toLowerCase().includes(search) ||
      unit.value?.toLowerCase().includes(search)
    )
  })

  // Event handlers
  function handlePrimaryAction() {
    showUnitModal = true
  }

  function handleSearch(event) {
    searchTerm = event.detail.searchTerm
  }

  function handleRowAction(event) {
    const { action, item } = event.detail
    
    switch (action) {
      case 'edit':
        editUnit(item)
        break
      case 'delete':
        showDeleteConfirmation(item)
        break
    }
  }

  function showDeleteConfirmation(unit) {
    unitToDelete = unit;
    showDeleteConfirm = true;
  }

  function cancelDelete() {
    showDeleteConfirm = false;
    unitToDelete = null;
  }

  async function confirmDelete() {
    if (!unitToDelete) return;
    
    try {
      loading = true;
      await DeleteUnitOfMeasurement(unitToDelete.id);
      showDbSuccess('delete', 'Unit of Measurement');
      dispatch('reload');
      showDeleteConfirm = false;
      unitToDelete = null;
    } catch (error) {
      console.error('Error deleting unit:', error);
      showDbError('delete', 'Unit of Measurement', error);
    } finally {
      loading = false;
    }
  }

  async function handleUnitSave(event) {
    const unitData = event.detail;
    
    try {
      loading = true;
      
      if (editingUnit) {
        // Update existing unit
        const updatedUnit = new database.UnitOfMeasurement({
          id: editingUnit,
          value: unitData.value,
          label: unitData.label,
          arabic: unitData.arabic || '',
          is_default: unitData.is_default || false,
          is_active: unitData.is_active !== undefined ? unitData.is_active : true,
          created_at: unitToEdit.created_at || null,
          updated_at: unitToEdit.updated_at || null
        });
        await UpdateUnitOfMeasurement(updatedUnit);
        showDbSuccess('update', 'Unit of Measurement');
      } else {
        // Create new unit
        const newUnit = new database.UnitOfMeasurement({
          id: 0,
          value: unitData.value,
          label: unitData.label,
          arabic: unitData.arabic || '',
          is_default: unitData.is_default || false,
          is_active: true,
          created_at: null,
          updated_at: null
        });
        await CreateUnitOfMeasurement(newUnit);
        showDbSuccess('create', 'Unit of Measurement');
      }
      
      dispatch('reload');
      closeUnitModal();
    } catch (error) {
      console.error('Error saving unit:', error);
      showDbError('save', 'Unit of Measurement', error);
    } finally {
      loading = false;
    }
  }

  function editUnit(unit) {
    editingUnit = unit.id
    unitToEdit = { ...unit }
    showUnitModal = true
  }

  function closeUnitModal() {
    showUnitModal = false
    editingUnit = null
    unitToEdit = null
  }
</script>

<div class="space-y-6">
  <!-- Units DataTable -->
   <DataTable
     data={filteredUnits}
     {columns}
     {primaryAction}
     loading={isLoading}
     searchTerm={searchTerm}
     searchPlaceholder="Search units..."
     emptyStateTitle="No units found"
     emptyStateMessage="Start by adding your first unit of measurement"
     emptyStateIcon="grid"
     on:primary-action={handlePrimaryAction}
     on:search={handleSearch}
     on:row-action={handleRowAction}
   />

   <!-- Unit Modal -->
   <UnitModal
     show={showUnitModal}
     unit={unitToEdit}
     {loading}
     on:save={handleUnitSave}
     on:close={closeUnitModal}
   />

   <!-- Delete Confirmation Modal -->
    <ConfirmationModal
      show={showDeleteConfirm}
      title="Delete Unit"
      message={unitToDelete ? `Are you sure you want to delete "${unitToDelete.label}"? This action cannot be undone.` : ''}
      confirmText="Delete"
      cancelText="Cancel"
      loading={loading}
      on:confirm={confirmDelete}
      on:cancel={cancelDelete}
    />
</div>