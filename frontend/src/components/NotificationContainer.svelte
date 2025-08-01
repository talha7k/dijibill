<script>
  import { notifications, removeNotification, NOTIFICATION_TYPES } from '../stores/notificationStore.js';
  import { fly } from 'svelte/transition';
  import { flip } from 'svelte/animate';

  // Icons for different notification types
  const icons = {
    [NOTIFICATION_TYPES.SUCCESS]: '✓',
    [NOTIFICATION_TYPES.ERROR]: '✕',
    [NOTIFICATION_TYPES.WARNING]: '⚠',
    [NOTIFICATION_TYPES.INFO]: 'ℹ'
  };

  // CSS classes for different notification types
  const typeClasses = {
    [NOTIFICATION_TYPES.SUCCESS]: 'bg-green-50 border-green-200 text-green-800',
    [NOTIFICATION_TYPES.ERROR]: 'bg-red-50 border-red-200 text-red-800',
    [NOTIFICATION_TYPES.WARNING]: 'bg-yellow-50 border-yellow-200 text-yellow-800',
    [NOTIFICATION_TYPES.INFO]: 'bg-blue-50 border-blue-200 text-blue-800'
  };

  const iconClasses = {
    [NOTIFICATION_TYPES.SUCCESS]: 'bg-green-100 text-green-600',
    [NOTIFICATION_TYPES.ERROR]: 'bg-red-100 text-red-600',
    [NOTIFICATION_TYPES.WARNING]: 'bg-yellow-100 text-yellow-600',
    [NOTIFICATION_TYPES.INFO]: 'bg-blue-100 text-blue-600'
  };

  function handleClose(id) {
    removeNotification(id);
  }
</script>

<!-- Notification Container -->
<div class="fixed top-4 right-4 z-50 space-y-2 max-w-sm w-full">
  {#each $notifications as notification (notification.id)}
    <div
      class="notification-item border rounded-lg shadow-lg p-4 {typeClasses[notification.type]}"
      in:fly="{{ x: 300, duration: 300 }}"
      out:fly="{{ x: 300, duration: 200 }}"
      animate:flip="{{ duration: 200 }}"
    >
      <div class="flex items-start">
        <!-- Icon -->
        <div class="flex-shrink-0 mr-3">
          <div class="w-6 h-6 rounded-full flex items-center justify-center text-sm font-bold {iconClasses[notification.type]}">
            {icons[notification.type]}
          </div>
        </div>
        
        <!-- Content -->
        <div class="flex-1 min-w-0">
          {#if notification.title}
            <h4 class="text-sm font-semibold mb-1">{notification.title}</h4>
          {/if}
          {#if notification.message}
            <p class="text-sm opacity-90">{notification.message}</p>
          {/if}
        </div>
        
        <!-- Close Button -->
        <button
          class="flex-shrink-0 ml-2 text-current opacity-60 hover:opacity-100 transition-opacity"
          on:click={() => handleClose(notification.id)}
          aria-label="Close notification"
        >
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path>
          </svg>
        </button>
      </div>
      
      <!-- Progress bar for auto-dismiss -->
      {#if !notification.persistent}
        <div class="mt-2 w-full bg-black bg-opacity-10 rounded-full h-1">
          <div 
            class="bg-current h-1 rounded-full transition-all ease-linear notification-progress"
            style="animation-duration: {notification.duration}ms"
          ></div>
        </div>
      {/if}
    </div>
  {/each}
</div>

<style>
  .notification-item {
    backdrop-filter: blur(8px);
  }
  
  .notification-progress {
    animation: progress-bar linear forwards;
    width: 100%;
  }
  
  @keyframes progress-bar {
    from {
      width: 100%;
    }
    to {
      width: 0%;
    }
  }
  
  /* Ensure notifications are above modals */
  .notification-item {
    position: relative;
    z-index: 9999;
  }
</style>