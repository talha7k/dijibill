<script>
  export let title = ''
  export let icon = ''
  export let showIndicator = false
</script>

<div class="page-container">
  <div class="page-header">
    <h1 class="page-title">
      {#if icon}
        <i class="fas {icon}"></i>
      {/if}
      {title}
      {#if showIndicator}
        <div class="status-indicator">
          <div class="indicator-dot"></div>
        </div>
      {/if}
    </h1>
    <div class="header-actions">
      <slot name="actions" />
    </div>
  </div>

  <div class="glass-card">
    <slot />
  </div>
</div>

<style>
  .page-container {
    padding: 1.5rem;
    min-height: 100vh;
  }

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  .page-title {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-size: 2rem;
    font-weight: 600;
    color: var(--color-text-primary);
    margin: 0;
  }

  .page-title i {
    font-size: 1.5rem;
  }

  .status-indicator {
    margin-left: auto;
  }

  .indicator-dot {
    width: 1rem;
    height: 1rem;
    background: var(--color-accent);
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .indicator-dot::after {
    content: '';
    width: 0.5rem;
    height: 0.5rem;
    background: white;
    border-radius: var(--radius-full);
  }

  .header-actions {
    display: flex;
    gap: 0.75rem;
    align-items: center;
  }

  .glass-card {
    background: var(--color-glass-white);
    backdrop-filter: blur(16px);
    border: 1px solid var(--color-border-light);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-glass);
    overflow: hidden;
  }
</style>