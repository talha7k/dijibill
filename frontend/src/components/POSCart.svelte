<script>
  import { 
    currentSale,
    showRefundModal,
    showTransferModal,
    showPaymentModal,
    loading,
    removeItem,
    updateQuantity
  } from '../stores/posStore.js'
  import POSHeader from './POSHeader.svelte'
  import { openPaymentModal } from '../services/posService.js'

  export let onSaveSale = () => {}
  export let onRefund = () => {}
  export let onTransfer = () => {}

  async function handlePayment() {
    await openPaymentModal()
  }
</script>

<div class="w-1/3 bg-gray-800 flex flex-col">
  <!-- Header with action buttons -->
  <POSHeader />

  <!-- Items List -->
  <div class="flex-1 overflow-y-auto p-4">
    {#if $currentSale.items.length === 0}
      <div class="text-center text-gray-400 mt-8">
        <p>No Items</p>
        <p class="text-sm">Select products to add to sale</p>
      </div>
    {:else}
      <div class="space-y-2">
        {#each $currentSale.items as item, index}
          <div class="bg-gray-700 p-3 rounded-lg">
            <div class="flex justify-between items-start mb-2">
              <div class="flex-1">
                <h4 class="font-medium">{item.product_name}</h4>
                <p class="text-sm text-gray-300">{item.unit_price.toFixed(2)} SAR</p>
              </div>
              <button
                class="text-red-400 hover:text-red-300"
                on:click={() => removeItem(index)}
              >
                <i class="fas fa-times w-4 h-4"></i>
              </button>
            </div>
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <button
                  class="w-8 h-8 bg-gray-600 hover:bg-gray-500 rounded flex items-center justify-center"
                  on:click={() => updateQuantity(index, item.quantity - 1)}
                >
                  -
                </button>
                <span class="w-8 text-center">{item.quantity}</span>
                <button
                  class="w-8 h-8 bg-gray-600 hover:bg-gray-500 rounded flex items-center justify-center"
                  on:click={() => updateQuantity(index, item.quantity + 1)}
                >
                  +
                </button>
              </div>
              <div class="text-right">
                <div class="font-medium">{item.total.toFixed(2)} SAR</div>
              </div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <!-- Totals -->
  <div class="p-4 border-t border-gray-700">
    <div class="space-y-2 mb-4">
      <div class="flex justify-between">
        <span>Subtotal:</span>
        <span>{$currentSale.subtotal.toFixed(2)} SAR</span>
      </div>
      <div class="flex justify-between">
        <span>VAT:</span>
        <span>{$currentSale.vat_total.toFixed(2)} SAR</span>
      </div>
      <div class="flex justify-between text-xl font-bold border-t border-gray-600 pt-2">
        <span>TOTAL:</span>
        <span>{$currentSale.total.toFixed(2)} SAR</span>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="space-y-2">
      <!-- Primary Action - Payment -->
      <button
        class="bg-green-600 hover:bg-green-700 text-white py-3 px-4 rounded-lg font-medium transition-colors w-full"
        disabled={$currentSale.items.length === 0 || $loading}
        on:click={handlePayment}
      >
        {$loading ? 'Processing...' : 'Make Payment'}
      </button>
      
      <!-- Secondary Actions -->
      <div class="grid grid-cols-3 gap-2">
        <button
          class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-3 rounded-lg font-medium transition-colors text-sm"
          disabled={$currentSale.items.length === 0}
          on:click={() => { $showTransferModal = true; onTransfer(); }}
          title="Transfer to another invoice"
        >
          Transfer
        </button>
        <button
          class="bg-red-600 hover:bg-red-700 text-white py-2 px-3 rounded-lg font-medium transition-colors text-sm"
          disabled={$currentSale.items.length === 0}
          on:click={() => { $showRefundModal = true; onRefund(); }}
          title="Process refund"
        >
          Refund
        </button>
        <button
          class="bg-gray-600 hover:bg-gray-700 text-white py-2 px-3 rounded-lg font-medium transition-colors text-sm"
          disabled={$currentSale.items.length === 0 || $loading}
          on:click={onSaveSale}
          title="Save as draft"
        >
          {$loading ? 'Saving...' : 'Save Draft'}
        </button>
      </div>
    </div>
  </div>
</div>