<script>
  import { GetUnitsOfMeasurement, CreateUnitOfMeasurement, UpdateUnitOfMeasurement, DeleteUnitOfMeasurement } from '../../wailsjs/go/main/App.js'
  import * as main from '../../wailsjs/go/models'
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  export let units = []
  export let isLoading = false

  let newUnit = { value: '', label: '', arabic: '', is_default: false, is_active: true }
  let showUnitForm = false
  let editingUnit = null

  async function addUnit() {
    if (newUnit.value && newUnit.label) {
      try {
        const unitData = new main.main.UnitOfMeasurement({
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
      const unitData = new main.main.UnitOfMeasurement({
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
          const updatedUnit = new main.main.UnitOfMeasurement({ ...unit, is_default: false });
          await UpdateUnitOfMeasurement(updatedUnit);
        }
      }
      // Then set the selected one as default
      const selectedUnit = units.find(u => u.id === id);
      if (selectedUnit) {
        const updatedUnit = new main.main.UnitOfMeasurement({ ...selectedUnit, is_default: true });
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
  <!-- Add Unit Button -->
  <div class="flex justify-between items-center">
    <h3 class="text-lg font-medium text-gray-900">Units of Measurement</h3>
    <button
      on:click={() => showUnitForm = true}
      class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
    >
      Add Unit
    </button>
  </div>

  <!-- Unit Form -->
  {#if showUnitForm}
    <div class="bg-gray-50 p-4 rounded-lg">
      <h4 class="text-md font-medium text-gray-900 mb-4">
        {editingUnit ? 'Edit Unit' : 'Add New Unit'}
      </h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Value</label>
          <input
            type="text"
            bind:value={newUnit.value}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter unit value (e.g., kg, pcs)"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Label</label>
          <input
            type="text"
            bind:value={newUnit.label}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter unit label (e.g., Kilograms, Pieces)"
          />
        </div>
        <div class="col-span-2">
          <label class="block text-sm font-medium text-gray-700 mb-2">Arabic Label</label>
          <input
            type="text"
            bind:value={newUnit.arabic}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="التسمية العربية"
            dir="rtl"
          />
        </div>
        <div class="flex items-center space-x-4 col-span-2">
          <label class="flex items-center">
            <input
              type="checkbox"
              bind:checked={newUnit.is_default}
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span class="ml-2 text-sm text-gray-700">Default</span>
          </label>
          <label class="flex items-center">
            <input
              type="checkbox"
              bind:checked={newUnit.is_active}
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span class="ml-2 text-sm text-gray-700">Active</span>
          </label>
        </div>
      </div>
      <div class="flex justify-end space-x-2 mt-4">
        <button
          on:click={resetUnitForm}
          class="px-4 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
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

  <!-- Units Table -->
  <div class="bg-white shadow overflow-hidden sm:rounded-md">
    <ul class="divide-y divide-gray-200">
      {#each units as unit (unit.id)}
        <li class="px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex-1">
              <div class="flex items-center">
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900">
                    {unit.label} ({unit.value})
                    {#if unit.arabic}
                      <span class="text-gray-500">- {unit.arabic}</span>
                    {/if}
                  </p>
                </div>
                <div class="flex items-center space-x-2">
                  {#if unit.is_default}
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      Default
                    </span>
                  {/if}
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {unit.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
                    {unit.is_active ? 'Active' : 'Inactive'}
                  </span>
                </div>
              </div>
            </div>
            <div class="flex items-center space-x-2">
              {#if !unit.is_default}
                <button
                  on:click={() => setDefaultUnit(unit.id)}
                  class="text-sm text-blue-600 hover:text-blue-900"
                >
                  Set Default
                </button>
              {/if}
              <button
                on:click={() => editUnit(unit)}
                class="text-sm text-indigo-600 hover:text-indigo-900"
              >
                Edit
              </button>
              <button
                on:click={() => deleteUnit(unit.id)}
                class="text-sm text-red-600 hover:text-red-900"
              >
                Delete
              </button>
            </div>
          </div>
        </li>
      {/each}
    </ul>
  </div>
</div>