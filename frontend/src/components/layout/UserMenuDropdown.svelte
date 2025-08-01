<script>
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';
  
  export let isOpen = false;
  /** @type {{ username: string, role: string } | null} */
  export let currentUser = null;
  export let buttonRect = null;
  
  const dispatch = createEventDispatcher();
  
  let dropdownElement;
  let mounted = false;
  
  onMount(() => {
    mounted = true;
    if (isOpen && buttonRect) {
      positionDropdown();
    }
  });
  
  onDestroy(() => {
    mounted = false;
  });
  
  $: if (mounted && isOpen && buttonRect) {
    positionDropdown();
  }
  
  function positionDropdown() {
    if (!dropdownElement || !buttonRect) return;
    
    const rect = buttonRect;
    dropdownElement.style.position = 'fixed';
    dropdownElement.style.top = `${rect.bottom + 8}px`;
    dropdownElement.style.left = `${rect.right - 208}px`; // 208px = w-52 (13rem * 16px)
    dropdownElement.style.zIndex = '999999';
  }
  
  function showUserProfile() {
    dispatch('show-user-profile');
    dispatch('close');
  }
  
  function showAccessControl() {
    dispatch('show-access-control');
    dispatch('close');
  }
  
  function switchToSettings() {
    dispatch('view-change', 'general-settings');
    dispatch('close');
  }
  
  function handleLogout() {
    dispatch('logout');
    dispatch('close');
  }
  
  let clickOutsideTimeout;
  
  function handleClickOutside(event) {
    if (dropdownElement && !dropdownElement.contains(event.target)) {
      dispatch('close');
    }
  }
  
  $: if (isOpen) {
    // Add a small delay to prevent the opening click from immediately closing the menu
    clickOutsideTimeout = setTimeout(() => {
      document.addEventListener('click', handleClickOutside);
    }, 50);
  } else {
    if (clickOutsideTimeout) {
      clearTimeout(clickOutsideTimeout);
    }
    document.removeEventListener('click', handleClickOutside);
  }
</script>

{#if isOpen && currentUser}
  <ul 
    bind:this={dropdownElement}
    class="menu menu-sm dropdown-content mt-3 p-2 shadow rounded-box w-52 glass-card"
    style="position: fixed; z-index: 999999; pointer-events: auto;"
  >
    <li class="menu-title">
      <span class="text-xs text-gray-500">Signed in as</span>
      <span class="font-semibold">{currentUser.username}</span>
      <span class="text-xs text-gray-500 capitalize">{currentUser.role}</span>
    </li>
    <div class="divider my-1"></div>
    <li>
      <button class="btn btn-ghost justify-start" on:click={showUserProfile}>
        <i class="fas fa-user w-4 h-4 mr-2"></i>
        Profile
      </button>
    </li>
    <li>
      <button class="btn btn-ghost justify-start" on:click={showAccessControl}>
        <i class="fas fa-shield-alt w-4 h-4 mr-2"></i>
        Access Control
      </button>
    </li>
    <li>
      <button class="btn btn-ghost justify-start" on:click={switchToSettings}>
        <i class="fas fa-cog w-4 h-4 mr-2"></i>
        Settings
      </button>
    </li>
    <div class="divider my-1"></div>
    <li>
      <button class="btn btn-ghost justify-start text-error" on:click={handleLogout}>
        <i class="fas fa-sign-out-alt w-4 h-4 mr-2"></i>
        Logout
      </button>
    </li>
  </ul>
{/if}

<style>
  :global(.glass-card) {
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(16px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  }
</style>