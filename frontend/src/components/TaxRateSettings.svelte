<script>
  import { GetTaxRates, CreateTaxRate, UpdateTaxRate, DeleteTaxRate } from '../../wailsjs/go/main/App.js'
  import { main } from '../../wailsjs/go/models'
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  export let taxRates = []
  export let isLoading = false

  let newTaxRate = { name: '', name_arabic: '', rate: 0, is_default: false, is_active: true }
  let showTaxRateForm = false
  let editingTaxRate = null

  async function addTaxRate() {
    if (newTaxRate.name && newTaxRate.rate >= 0) {
      try {
        const taxRateData = new main.TaxRate({
          id: 0,
          name: newTaxRate.name,
          name_arabic: newTaxRate.name_arabic || '',
          rate: newTaxRate.rate,
          is_default: newTaxRate.is_default || false,
          is_active: true,
          description: '',
          created_at: null,
          updated_at: null
        });
        await CreateTaxRate(taxRateData)
        dispatch('reload')
        resetTaxRateForm()
      } catch (error) {
        console.error('Error adding tax rate:', error)
        dispatch('error', { message: 'Error adding tax rate: ' + error.message })
      }
    }
  }

  function editTaxRate(taxRate) {
    editingTaxRate = taxRate.id
    newTaxRate = { ...taxRate }
    showTaxRateForm = true
  }

  async function updateTaxRate() {
    try {
      const taxRateData = new main.TaxRate({
        id: editingTaxRate,
        name: newTaxRate.name,
        name_arabic: newTaxRate.name_arabic || '',
        rate: newTaxRate.rate,
        is_default: newTaxRate.is_default || false,
        is_active: newTaxRate.is_active !== undefined ? newTaxRate.is_active : true,
        description: newTaxRate.description || '',
        created_at: newTaxRate.created_at || null,
        updated_at: newTaxRate.updated_at || null
      });
      await UpdateTaxRate(taxRateData)
      dispatch('reload')
      resetTaxRateForm()
    } catch (error) {
      console.error('Error updating tax rate:', error)
      dispatch('error', { message: 'Error updating tax rate: ' + error.message })
    }
  }

  async function deleteTaxRate(id) {
    if (confirm('Are you sure you want to delete this tax rate?')) {
      try {
        await DeleteTaxRate(id)
        dispatch('reload')
      } catch (error) {
        console.error('Error deleting tax rate:', error)
        dispatch('error', { message: 'Error deleting tax rate: ' + error.message })
      }
    }
  }

  function resetTaxRateForm() {
    newTaxRate = { name: '', name_arabic: '', rate: 0, is_default: false, is_active: true }
    showTaxRateForm = false
    editingTaxRate = null
  }

  async function setDefaultTaxRate(id) {
    try {
      // First, update all tax rates to not be default
      for (const taxRate of taxRates) {
        if (taxRate.is_default) {
          const updatedTaxRate = new main.TaxRate({ ...taxRate, is_default: false });
          await UpdateTaxRate(updatedTaxRate);
        }
      }
      // Then set the selected one as default
      const selectedTaxRate = taxRates.find(tr => tr.id === id);
      if (selectedTaxRate) {
        const updatedTaxRate = new main.TaxRate({ ...selectedTaxRate, is_default: true });
        await UpdateTaxRate(updatedTaxRate);
        dispatch('reload')
      }
    } catch (error) {
      console.error('Error setting default tax rate:', error)
      dispatch('error', { message: 'Error setting default tax rate: ' + error.message })
    }
  }
</script>

<div class="space-y-6">
  <!-- Add Tax Rate Button -->
  <div class="flex justify-between items-center">
    <h3 class="text-lg font-medium text-gray-900">Tax Rates</h3>
    <button
      on:click={() => showTaxRateForm = true}
      class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
    >
      Add Tax Rate
    </button>
  </div>

  <!-- Tax Rate Form -->
  {#if showTaxRateForm}
    <div class="bg-gray-50 p-4 rounded-lg">
      <h4 class="text-md font-medium text-gray-900 mb-4">
        {editingTaxRate ? 'Edit Tax Rate' : 'Add New Tax Rate'}
      </h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Name</label>
          <input
            type="text"
            bind:value={newTaxRate.name}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter tax rate name"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Name (Arabic)</label>
          <input
            type="text"
            bind:value={newTaxRate.name_arabic}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="اسم معدل الضريبة"
            dir="rtl"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Rate (%)</label>
          <input
            type="number"
            step="0.01"
            min="0"
            max="100"
            bind:value={newTaxRate.rate}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter tax rate percentage"
          />
        </div>
        <div class="flex items-center space-x-4">
          <label class="flex items-center">
            <input
              type="checkbox"
              bind:checked={newTaxRate.is_default}
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span class="ml-2 text-sm text-gray-700">Default</span>
          </label>
          <label class="flex items-center">
            <input
              type="checkbox"
              bind:checked={newTaxRate.is_active}
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span class="ml-2 text-sm text-gray-700">Active</span>
          </label>
        </div>
      </div>
      <div class="flex justify-end space-x-2 mt-4">
        <button
          on:click={resetTaxRateForm}
          class="px-4 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          Cancel
        </button>
        <button
          on:click={editingTaxRate ? updateTaxRate : addTaxRate}
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          {editingTaxRate ? 'Update' : 'Add'} Tax Rate
        </button>
      </div>
    </div>
  {/if}

  <!-- Tax Rates Table -->
  <div class="bg-white shadow overflow-hidden sm:rounded-md">
    <ul class="divide-y divide-gray-200">
      {#each taxRates as taxRate (taxRate.id)}
        <li class="px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex-1">
              <div class="flex items-center">
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900">
                    {taxRate.name}
                    {#if taxRate.name_arabic}
                      <span class="text-gray-500">({taxRate.name_arabic})</span>
                    {/if}
                  </p>
                  <p class="text-sm text-gray-500">Rate: {taxRate.rate}%</p>
                </div>
                <div class="flex items-center space-x-2">
                  {#if taxRate.is_default}
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      Default
                    </span>
                  {/if}
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {taxRate.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
                    {taxRate.is_active ? 'Active' : 'Inactive'}
                  </span>
                </div>
              </div>
            </div>
            <div class="flex items-center space-x-2">
              {#if !taxRate.is_default}
                <button
                  on:click={() => setDefaultTaxRate(taxRate.id)}
                  class="text-sm text-blue-600 hover:text-blue-900"
                >
                  Set Default
                </button>
              {/if}
              <button
                on:click={() => editTaxRate(taxRate)}
                class="text-sm text-indigo-600 hover:text-indigo-900"
              >
                Edit
              </button>
              <button
                on:click={() => deleteTaxRate(taxRate.id)}
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