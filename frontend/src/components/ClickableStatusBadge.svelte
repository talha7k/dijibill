<script>
  import { createEventDispatcher } from 'svelte'
  
  export let status = ''
  export let variant = 'default'
  export let size = 'md'
  export let clickable = false
  export let disabled = false
  export let loading = false

  const dispatch = createEventDispatcher()

  function getStatusClass(status, variant) {
    if (variant !== 'default') {
      return `badge-${variant}`
    }

    switch (status?.toLowerCase()) {
      case 'completed':
      case 'active':
      case 'paid':
      case 'success':
      case 'received':
        return 'badge-success'
      case 'pending':
      case 'processing':
      case 'warning':
        return 'badge-warning'
      case 'failed':
      case 'error':
      case 'cancelled':
      case 'inactive':
        return 'badge-error'
      case 'refunded':
        return 'badge-secondary'
      case 'draft':
      case 'info':
        return 'badge-info'
      default:
        return 'badge-neutral'
    }
  }

  function handleClick() {
    if (clickable && !disabled && !loading) {
      dispatch('click', { status })
    }
  }

  $: badgeClass = getStatusClass(status, variant)
  $: sizeClass = size === 'sm' ? 'badge-sm' : size === 'lg' ? 'badge-lg' : ''
  $: cursorClass = clickable && !disabled ? 'cursor-pointer' : ''
  $: disabledClass = disabled ? 'opacity-50' : ''
</script>

<span 
  class="badge {badgeClass} {sizeClass} {cursorClass} {disabledClass} {clickable ? 'hover:opacity-80 transition-opacity' : ''}"
  on:click={handleClick}
  on:keydown={(e) => e.key === 'Enter' && handleClick()}
  role={clickable ? 'button' : undefined}
  tabindex={clickable && !disabled ? 0 : undefined}
  title={clickable && status === 'pending' ? 'Click to mark as received' : undefined}
>
  {#if loading}
    <i class="fas fa-spinner fa-spin mr-1"></i>
  {/if}
  <slot>{status}</slot>
</span>

<style>
  .badge {
    display: inline-flex;
    align-items: center;
    padding: 0.25rem 0.75rem;
    font-size: 0.75rem;
    font-weight: 500;
    border-radius: var(--radius-full);
    text-transform: capitalize;
    white-space: nowrap;
  }

  .badge-sm {
    padding: 0.125rem 0.5rem;
    font-size: 0.625rem;
  }

  .badge-lg {
    padding: 0.375rem 1rem;
    font-size: 0.875rem;
  }

  .badge-success {
    background: var(--color-success);
    color: white;
  }

  .badge-warning {
    background: var(--color-warning);
    color: white;
  }

  .badge-error {
    background: var(--color-error);
    color: white;
  }

  .badge-info {
    background: var(--color-info);
    color: white;
  }

  .badge-neutral {
    background: var(--color-neutral);
    color: white;
  }

  .badge-primary {
    background: var(--color-primary);
    color: white;
  }

  .badge-secondary {
    background: var(--color-secondary);
    color: white;
  }

  .cursor-pointer {
    cursor: pointer;
  }

  .opacity-50 {
    opacity: 0.5;
  }

  .hover\:opacity-80:hover {
    opacity: 0.8;
  }

  .transition-opacity {
    transition: opacity 0.2s ease-in-out;
  }
</style>