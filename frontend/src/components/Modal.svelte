<script>
  import { createEventDispatcher } from 'svelte'
  import ActionButton from './ActionButton.svelte'
  import IconButton from './IconButton.svelte'
  
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
        <h3 class="font-bold text-lg text-white">{title}</h3>
        <IconButton
          icon="times"
          variant="ghost"
          size="sm"
          customClass="text-white modal-close-btn"
          on:click={closeModal}
        />
      </div>

      <!-- Body -->
      <div class="modal-body">
        <slot></slot>
      </div>

      <!-- Footer -->
      {#if showPrimaryButton || showSecondaryButton}
        <div class="modal-footer">
          {#if showSecondaryButton}
            <ActionButton
              variant="ghost"
              text={secondaryButtonText}
              disabled={loading}
              on:click={handleSecondaryAction}
            />
          {/if}
          {#if showPrimaryButton}
            <ActionButton
              variant="primary"
              text={primaryButtonText}
              disabled={primaryButtonDisabled || loading}
              {loading}
              on:click={handlePrimaryAction}
            />
          {/if}
        </div>
      {/if}
    </div>
  </div>
{/if}