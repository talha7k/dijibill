<script>
  // Props are kept for backward compatibility but not used
  export let title = ''
  export let icon = ''
  export let showIndicator = false
</script>

<div class="page-container">
  <div class="glass-card">
    <slot />
  </div>
</div>

<style>
  .page-container {
    padding: 1.5rem;
    min-height: 100vh;
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