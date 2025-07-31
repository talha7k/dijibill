<script>
  import { createEventDispatcher } from 'svelte';
  
  export let wailsReady = false;
  export let wailsError = null;
  
  const dispatch = createEventDispatcher();
  
  function retryInitialization() {
    dispatch('retry-initialization');
  }
</script>

<div class="flex items-center justify-center h-full">
  <div class="text-center">
    {#if wailsError}
      <!-- Error State -->
      <div class="alert alert-error max-w-md mx-auto">
        <i class="fas fa-exclamation-triangle w-6 h-6"></i>
        <div>
          <h3 class="font-bold">Initialization Error</h3>
          <div class="text-xs">{wailsError}</div>
        </div>
        <button class="btn btn-sm" on:click={retryInitialization}>
          Retry
        </button>
      </div>
    {:else}
      <!-- Loading State -->
      <div class="flex flex-col items-center gap-4">
        <div class="loading loading-spinner loading-lg text-primary"></div>
        <div class="text-white">
          <h3 class="text-lg font-semibold">Initializing Application</h3>
          <p class="text-sm opacity-70">Please wait while we connect to the backend...</p>
        </div>
      </div>
    {/if}
  </div>
</div>