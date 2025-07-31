<script>
  export let status = ''
  export let variant = 'default'
  export let size = 'md'

  function getStatusClass(status, variant) {
    if (variant !== 'default') {
      return `badge-${variant}`
    }

    switch (status?.toLowerCase()) {
      case 'completed':
      case 'active':
      case 'paid':
      case 'success':
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

  $: badgeClass = getStatusClass(status, variant)
  $: sizeClass = size === 'sm' ? 'badge-sm' : size === 'lg' ? 'badge-lg' : ''
</script>

<span class="badge {badgeClass} {sizeClass}">
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
</style>