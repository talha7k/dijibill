<script>
  import { createEventDispatcher } from 'svelte'

  export let show = false
  export let title = 'Inventory Updated'
  export let message = ''
  export let updatedProducts = []
  export let loading = false

  const dispatch = createEventDispatcher()

  function handleClose() {
    dispatch('close')
  }

  function handleConfirm() {
    dispatch('confirm')
  }
</script>

{#if show}
  <div class="modal-overlay" on:click={handleClose}>
    <div class="modal-content" on:click|stopPropagation>
      <div class="modal-header">
        <div class="flex items-center">
          <div class="icon-container success">
            <i class="fas fa-check"></i>
          </div>
          <h3 class="modal-title">{title}</h3>
        </div>
        <button class="modal-close" on:click={handleClose}>
          <i class="fas fa-times"></i>
        </button>
      </div>

      <div class="modal-body">
        <p class="message">{message}</p>
        
        {#if updatedProducts.length > 0}
          <div class="updated-products">
            <h4 class="products-title">Updated Products:</h4>
            <div class="products-list">
              {#each updatedProducts as product}
                <div class="product-item">
                  <div class="product-info">
                    <span class="product-name">{product.name}</span>
                    <span class="quantity-change">+{product.quantity} units</span>
                  </div>
                  <div class="stock-info">
                    <span class="new-stock">New Stock: {product.newStock}</span>
                  </div>
                </div>
              {/each}
            </div>
          </div>
        {/if}
      </div>

      <div class="modal-footer">
        <button 
          class="btn btn-primary" 
          on:click={handleConfirm}
          disabled={loading}
        >
          {#if loading}
            <i class="fas fa-spinner fa-spin"></i>
          {/if}
          OK
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(4px);
  }

  .modal-content {
    background: var(--color-surface);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-xl);
    max-width: 500px;
    width: 90%;
    max-height: 80vh;
    overflow-y: auto;
    border: 1px solid var(--color-border);
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1.5rem;
    border-bottom: 1px solid var(--color-border);
  }

  .modal-header .flex {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .icon-container {
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.25rem;
  }

  .icon-container.success {
    background: var(--color-success);
    color: white;
  }

  .modal-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--color-text);
    margin: 0;
  }

  .modal-close {
    background: none;
    border: none;
    color: var(--color-text-secondary);
    font-size: 1.25rem;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: var(--radius);
    transition: all 0.2s ease;
  }

  .modal-close:hover {
    background: var(--color-surface-hover);
    color: var(--color-text);
  }

  .modal-body {
    padding: 1.5rem;
  }

  .message {
    color: var(--color-text-secondary);
    margin: 0 0 1.5rem 0;
    line-height: 1.5;
  }

  .updated-products {
    margin-top: 1rem;
  }

  .products-title {
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text);
    margin: 0 0 1rem 0;
  }

  .products-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .product-item {
    background: var(--color-surface-secondary);
    border: 1px solid var(--color-border);
    border-radius: var(--radius);
    padding: 1rem;
  }

  .product-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  .product-name {
    font-weight: 500;
    color: var(--color-text);
  }

  .quantity-change {
    background: var(--color-success);
    color: white;
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius-sm);
    font-size: 0.875rem;
    font-weight: 500;
  }

  .stock-info {
    display: flex;
    justify-content: flex-end;
  }

  .new-stock {
    color: var(--color-text-secondary);
    font-size: 0.875rem;
  }

  .modal-footer {
    padding: 1.5rem;
    border-top: 1px solid var(--color-border);
    display: flex;
    justify-content: flex-end;
  }

  .btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: var(--radius);
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
  }

  .btn-primary {
    background: var(--color-primary);
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background: var(--color-primary-hover);
  }

  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>