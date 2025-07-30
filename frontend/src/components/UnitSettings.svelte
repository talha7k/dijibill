<script>
  import { GetUnitsOfMeasurement, CreateUnitOfMeasurement, UpdateUnitOfMeasurement, DeleteUnitOfMeasurement } from '../../wailsjs/go/main/App.js'
   import { database } from '../../wailsjs/go/models'
  import { createEventDispatcher } from 'svelte'
  import DataTable from './DataTable.svelte'
  import FormField from './FormField.svelte'

  const dispatch = createEventDispatcher()

  export let units = []
  export let isLoading = false

  let newUnit = { value: '', label: '', arabic: '', is_default: false, is_active: true }
  let showUnitForm = false
  let editingUnit = null
  let searchTerm = ''

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
    }
  ]

  const actions = [
    {
      key: 'setDefault',
      label: 'Set Default',
      labelArabic: 'تعيين كافتراضي',
      icon: 'fa-star',
      variant: 'secondary',
      condition: (unit) => !unit.is_default
    },
    {
      key: 'edit',
      label: 'Edit',
      labelArabic: 'تعديل',
      icon: 'fa-edit',
      variant: 'primary'
    },
    {
      key: 'delete',
      label: 'Delete',
      labelArabic: 'حذف',
      icon: 'fa-trash',
      variant: 'danger'
    }
  ]

  const primaryAction = {
    label: 'Add Unit',
    labelArabic: 'إضافة وحدة',
    icon: 'fa-plus'
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
    showUnitForm = true
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
        deleteUnit(item.id)
        break
      case 'setDefault':
        setDefaultUnit(item.id)
        break
    }
  }

  async function addUnit() {
    if (newUnit.value && newUnit.label) {
      try {
        const unitData = new database.UnitOfMeasurement({
          id: 0,
          value: newUnit.value,
          label: newUnit.label,
          arabic: newUnit.arabic || '',
          is_default: newUnit.is_default || false,
          is_active: true,
          created_at: null,
          updated_at: null
        });
        await CreateUnitOfMeasurement(unitData);
        dispatch('reload')
        resetUnitForm();
      } catch (error) {
        console.error('Error adding unit:', error);
        dispatch('error', { message: 'Error adding unit: ' + error.message })
      }
    }
  }

  function editUnit(unit) {
    editingUnit = unit.id
    newUnit = { ...unit }
    showUnitForm = true
  }

  async function updateUnit() {
    try {
      const unitData = new database.UnitOfMeasurement({
        id: editingUnit,
        value: newUnit.value,
        label: newUnit.label,
        arabic: newUnit.arabic || '',
        is_default: newUnit.is_default || false,
        is_active: newUnit.is_active !== undefined ? newUnit.is_active : true,
        created_at: newUnit.created_at || null,
        updated_at: newUnit.updated_at || null
      });
      await UpdateUnitOfMeasurement(unitData);
      dispatch('reload')
      resetUnitForm();
    } catch (error) {
      console.error('Error updating unit:', error);
      dispatch('error', { message: 'Error updating unit: ' + error.message })
    }
  }

  async function deleteUnit(id) {
    if (confirm('Are you sure you want to delete this unit?')) {
      try {
        await DeleteUnitOfMeasurement(id);
        dispatch('reload')
      } catch (error) {
        console.error('Error deleting unit:', error);
        dispatch('error', { message: 'Error deleting unit: ' + error.message })
      }
    }
  }

  function resetUnitForm() {
    newUnit = { value: '', label: '', arabic: '', is_default: false, is_active: true }
    showUnitForm = false
    editingUnit = null
  }

  async function setDefaultUnit(id) {
    try {
      // First, update all units to not be default
      for (const unit of units) {
        if (unit.is_default) {
          const updatedUnit = new database.UnitOfMeasurement({ ...unit, is_default: false });
          await UpdateUnitOfMeasurement(updatedUnit);
        }
      }
      // Then set the selected one as default
      const selectedUnit = units.find(u => u.id === id);
      if (selectedUnit) {
        const updatedUnit = new database.UnitOfMeasurement({ ...selectedUnit, is_default: true });
        await UpdateUnitOfMeasurement(updatedUnit);
        dispatch('reload')
      }
    } catch (error) {
      console.error('Error setting default unit:', error);
      dispatch('error', { message: 'Error setting default unit: ' + error.message })
    }
  }
</script>

<div class="space-y-6">
  <!-- Unit Form -->
  {#if showUnitForm}
    <div class="bg-gray-50 p-6 rounded-lg border">
      <h4 class="text-lg font-medium text-gray-100 mb-6">
        {editingUnit ? 'Edit Unit' : 'Add New Unit'}
      </h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <FormField
          label="Value"
          labelArabic="القيمة"
          type="text"
          bind:value={newUnit.value}
          placeholder="Enter unit value (e.g., kg, pcs)"
          required
        />

        <FormField
          label="Label"
          labelArabic="التسمية"
          type="text"
          bind:value={newUnit.label}
          placeholder="Enter unit label (e.g., Kilograms, Pieces)"
          required
        />

        <div class="col-span-2">
          <FormField
            label="Arabic Label"
            labelArabic="التسمية العربية"
            type="text"
            bind:value={newUnit.arabic}
            placeholder="التسمية العربية"
            dir="rtl"
          />
        </div>

        <div class="flex items-center space-x-6 col-span-2">
          <FormField
            type="checkbox"
            bind:checked={newUnit.is_default}
            placeholder="Set as default"
          />
          <FormField
            type="checkbox"
            bind:checked={newUnit.is_active}
            placeholder="Active"
          />
        </div>
      </div>
      
      <div class="flex justify-end space-x-3 mt-6">
        <button
          on:click={resetUnitForm}
          class="px-4 py-2 border border-gray-300 text-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          Cancel
        </button>
        <button
          on:click={editingUnit ? updateUnit : addUnit}
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          {editingUnit ? 'Update' : 'Add'} Unit
        </button>
      </div>
    </div>
  {/if}

  <!-- Units DataTable -->
   <DataTable
     data={filteredUnits}
     {columns}
     secondaryActions={actions}
     {primaryAction}
     loading={isLoading}
     searchTerm={searchTerm}
     searchPlaceholder="Search units..."
     emptyStateTitle="No units found"
     emptyStateMessage="Start by adding your first unit of measurement"
     emptyStateIcon="fa-ruler"
     on:primary-action={handlePrimaryAction}
     on:search={handleSearch}
     on:row-action={handleRowAction}
   />
</div>