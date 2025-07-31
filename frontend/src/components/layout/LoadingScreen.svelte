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
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
        </svg>
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