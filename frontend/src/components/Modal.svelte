<script>
  import { createEventDispatcher } from 'svelte'
  
  export let show = false
  export let title = ''
  export let size = 'md' // sm, md, lg, xl
  export let loading = false
  export let primaryButtonText = 'Save'
  export let secondaryButtonText = 'Cancel'
  export let showPrimaryButton = true
  export let showSecondaryButton = true
  export let primaryButtonDisabled = false

  const dispatch = createEventDispatcher()

  function closeModal() {
    dispatch('close')
  }

  function handlePrimaryAction() {
    dispatch('primary')
  }

  function handleSecondaryAction() {
    dispatch('secondary')
  }

  function handleBackdropClick(event) {
    if (event.target === event.currentTarget) {
      closeModal()
    }
  }

  function handleKeydown(event) {
    if (event.key === 'Escape') {
      closeModal()
    }
  }

  $: sizeClass = {
    sm: 'max-w-md',
    md: 'max-w-2xl',
    lg: 'max-w-4xl',
    xl: 'max-w-6xl'
  }[size] || 'max-w-2xl'
</script>

{#if show}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div 
    class="modal-glass" 
    on:click={handleBackdropClick}
    on:keydown={handleKeydown}
  >
    <div class="modal-backdrop"></div>
    
    <div class="modal-content {sizeClass}">
      <!-- Header -->
      <div class="modal-header">
        <h3 class="heading-3">{title}</h3>
        <button 
          class="btn-icon text-white/70 hover:text-white hover:bg-white/10"
          on:click={closeModal}
          aria-label="Close modal"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>

      <!-- Body -->
      <div class="modal-body">
        <slot></slot>
      </div>

      <!-- Footer -->
      {#if showPrimaryButton || showSecondaryButton}
        <div class="modal-footer">
          {#if showSecondaryButton}
            <button 
              class="btn-glass-outline" 
              on:click={handleSecondaryAction}
              disabled={loading}
            >
              {secondaryButtonText}
            </button>
          {/if}
          
          {#if showPrimaryButton}
            <button 
              class="btn-primary" 
              on:click={handlePrimaryAction}
              disabled={primaryButtonDisabled || loading}
            >
              {#if loading}
                <span class="loading-spinner-sm mr-2"></span>
              {/if}
              {primaryButtonText}
            </button>
          {/if}
        </div>
      {/if}
    </div>
  </div>
{/if}