<script>
  import { GetPaymentTypes, CreatePaymentType, UpdatePaymentType, DeletePaymentType } from '../../wailsjs/go/main/App.js'
  import { main } from '../../wailsjs/go/models'
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  export let paymentTypes = []
  export let isLoading = false

  let newPaymentType = { name: '', name_arabic: '', code: '', description: '', is_default: false, is_active: true }
  let showPaymentTypeForm = false
  let editingPaymentType = null

  async function addPaymentType() {
    if (newPaymentType.name && newPaymentType.code) {
      try {
        const paymentTypeData = new main.PaymentType({
          id: 0,
          name: newPaymentType.name,
          name_arabic: newPaymentType.name_arabic || '',
          code: newPaymentType.code,
          description: newPaymentType.description || '',
          is_default: newPaymentType.is_default || false,
          is_active: true,
          created_at: null,
          updated_at: null
        });
        await CreatePaymentType(paymentTypeData)
        dispatch('reload')
        resetPaymentTypeForm()
      } catch (error) {
        console.error('Error adding payment type:', error)
        dispatch('error', { message: 'Error adding payment type: ' + error.message })
      }
    }
  }

  function editPaymentType(paymentType) {
    editingPaymentType = paymentType.id
    newPaymentType = { ...paymentType }
    showPaymentTypeForm = true
  }

  async function updatePaymentType() {
    try {
      const paymentTypeData = new main.PaymentType({
        id: editingPaymentType,
        name: newPaymentType.name,
        name_arabic: newPaymentType.name_arabic || '',
        code: newPaymentType.code,
        description: newPaymentType.description || '',
        is_default: newPaymentType.is_default || false,
        is_active: newPaymentType.is_active !== undefined ? newPaymentType.is_active : true,
        created_at: newPaymentType.created_at || null,
        updated_at: newPaymentType.updated_at || null
      });
      await UpdatePaymentType(paymentTypeData)
      dispatch('reload')
      resetPaymentTypeForm()
    } catch (error) {
      console.error('Error updating payment type:', error)
      dispatch('error', { message: 'Error updating payment type: ' + error.message })
    }
  }

  async function deletePaymentType(id) {
    if (confirm('Are you sure you want to delete this payment type?')) {
      try {
        await DeletePaymentType(id)
        dispatch('reload')
      } catch (error) {
        console.error('Error deleting payment type:', error)
        dispatch('error', { message: 'Error deleting payment type: ' + error.message })
      }
    }
  }

  function resetPaymentTypeForm() {
    newPaymentType = { name: '', name_arabic: '', code: '', description: '', is_default: false, is_active: true }
    showPaymentTypeForm = false
    editingPaymentType = null
  }

  async function setDefaultPaymentType(id) {
    try {
      // First, update all payment types to not be default
      for (const paymentType of paymentTypes) {
        if (paymentType.is_default) {
          const updatedPaymentType = new main.PaymentType({ ...paymentType, is_default: false });
          await UpdatePaymentType(updatedPaymentType);
        }
      }
      // Then set the selected one as default
      const selectedPaymentType = paymentTypes.find(pt => pt.id === id);
      if (selectedPaymentType) {
        const updatedPaymentType = new main.PaymentType({ ...selectedPaymentType, is_default: true });
        await UpdatePaymentType(updatedPaymentType);
        dispatch('reload')
      }
    } catch (error) {
      console.error('Error setting default payment type:', error)
      dispatch('error', { message: 'Error setting default payment type: ' + error.message })
    }
  }
</script>

<div class="space-y-6">
  <!-- Add Payment Type Button -->
  <div class="flex justify-between items-center">
    <h3 class="text-lg font-medium text-gray-900">Payment Types</h3>
    <button
      on:click={() => showPaymentTypeForm = true}
      class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
    >
      Add Payment Type
    </button>
  </div>

  <!-- Payment Type Form -->
  {#if showPaymentTypeForm}
    <div class="bg-gray-50 p-4 rounded-lg">
      <h4 class="text-md font-medium text-gray-900 mb-4">
        {editingPaymentType ? 'Edit Payment Type' : 'Add New Payment Type'}
      </h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Name</label>
          <input
            type="text"
            bind:value={newPaymentType.name}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter payment type name"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Name (Arabic)</label>
          <input
            type="text"
            bind:value={newPaymentType.name_arabic}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="اسم نوع الدفع"
            dir="rtl"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Code</label>
          <input
            type="text"
            bind:value={newPaymentType.code}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter payment type code"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
          <input
            type="text"
            bind:value={newPaymentType.description}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter description"
          />
        </div>
        <div class="flex items-center space-x-4 col-span-2">
          <label class="flex items-center">
            <input
              type="checkbox"
              bind:checked={newPaymentType.is_default}
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span class="ml-2 text-sm text-gray-700">Default</span>
          </label>
          <label class="flex items-center">
            <input
              type="checkbox"
              bind:checked={newPaymentType.is_active}
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span class="ml-2 text-sm text-gray-700">Active</span>
          </label>
        </div>
      </div>
      <div class="flex justify-end space-x-2 mt-4">
        <button
          on:click={resetPaymentTypeForm}
          class="px-4 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          Cancel
        </button>
        <button
          on:click={editingPaymentType ? updatePaymentType : addPaymentType}
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          {editingPaymentType ? 'Update' : 'Add'} Payment Type
        </button>
      </div>
    </div>
  {/if}

  <!-- Payment Types Table -->
  <div class="bg-white shadow overflow-hidden sm:rounded-md">
    <ul class="divide-y divide-gray-200">
      {#each paymentTypes as paymentType (paymentType.id)}
        <li class="px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex-1">
              <div class="flex items-center">
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900">
                    {paymentType.name}
                    {#if paymentType.name_arabic}
                      <span class="text-gray-500">({paymentType.name_arabic})</span>
                    {/if}
                  </p>
                  <p class="text-sm text-gray-500">Code: {paymentType.code}</p>
                  {#if paymentType.description}
                    <p class="text-sm text-gray-500">{paymentType.description}</p>
                  {/if}
                </div>
                <div class="flex items-center space-x-2">
                  {#if paymentType.is_default}
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      Default
                    </span>
                  {/if}
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {paymentType.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
                    {paymentType.is_active ? 'Active' : 'Inactive'}
                  </span>
                </div>
              </div>
            </div>
            <div class="flex items-center space-x-2">
              {#if !paymentType.is_default}
                <button
                  on:click={() => setDefaultPaymentType(paymentType.id)}
                  class="text-sm text-blue-600 hover:text-blue-900"
                >
                  Set Default
                </button>
              {/if}
              <button
                on:click={() => editPaymentType(paymentType)}
                class="text-sm text-indigo-600 hover:text-indigo-900"
              >
                Edit
              </button>
              <button
                on:click={() => deletePaymentType(paymentType.id)}
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