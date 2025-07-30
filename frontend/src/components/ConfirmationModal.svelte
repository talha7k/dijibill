<script>
  import { createEventDispatcher } from 'svelte'
  
  export let show = false
  export let title = 'Confirm Action'
  export let message = 'Are you sure you want to proceed?'
  export let confirmText = 'Confirm'
  export let cancelText = 'Cancel'
  export let confirmClass = 'btn-error'
  export let loading = false

  const dispatch = createEventDispatcher()

  function handleConfirm() {
    dispatch('confirm')
  }

  function handleCancel() {
    dispatch('cancel')
  }

  function handleBackdropClick(event) {
    if (event.target === event.currentTarget) {
      handleCancel()
    }
  }
</script>

{#if show}
  <div class="modal-backdrop" on:click={handleBackdropClick}>
    <div class="modal-content">
      <div class="modal-header">
        <h3>{title}</h3>
      </div>
      
      <div class="modal-body">
        <p>{message}</p>
      </div>
      
      <div class="modal-footer">
        <button 
          class="btn btn-secondary" 
          on:click={handleCancel}
          disabled={loading}
        >
          {cancelText}
        </button>
        <button 
          class="btn {confirmClass}" 
          on:click={handleConfirm}
          disabled={loading}
        >
          {#if loading}
            <i class="fas fa-spinner fa-spin"></i>
          {/if}
          {confirmText}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal-content {
    border-radius: 8px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
    max-width: 400px;
    width: 90%;
    max-height: 90vh;
    overflow: hidden;
  }

  .modal-header {
    padding: 1.5rem 1.5rem 0 1.5rem;
  }

  .modal-header h3 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--color-text-primary);
  }

  .modal-body {
    padding: 1rem 1.5rem;
  }

  .modal-body p {
    margin: 0;
    color: var(--color-text-secondary);
    line-height: 1.5;
  }

  .modal-footer {
    padding: 1rem 1.5rem 1.5rem 1.5rem;
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
  }

  .btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-secondary {
    background-color: var(--color-bg-secondary);
    color: var(--color-text-primary);
    border: 1px solid var(--color-border);
  }

  .btn-secondary:hover:not(:disabled) {
    background-color: var(--color-bg-tertiary);
  }

  .btn-error {
    background-color: #dc2626;
    color: white;
  }

  .btn-error:hover:not(:disabled) {
    background-color: #b91c1c;
  }

  .fa-spinner {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
</style>