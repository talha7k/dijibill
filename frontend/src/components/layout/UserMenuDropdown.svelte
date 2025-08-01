<script>
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';
  import ActionButton from '../ActionButton.svelte';
  
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
    dropdownElement.style.right = '16px'; // Fixed right margin
    dropdownElement.style.left = 'unset';
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

{#if isOpen && buttonRect && currentUser}
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
      <ActionButton
        variant="ghost"
        customClass="justify-start"
        icon="user"
        text="Profile"
        on:click={showUserProfile}
      />
    </li>
    <li>
      <ActionButton
        variant="ghost"
        customClass="justify-start"
        icon="shield"
        text="Access Control"
        on:click={showAccessControl}
      />
    </li>
    <li>
      <ActionButton
        variant="ghost"
        customClass="justify-start"
        icon="cog"
        text="Settings"
        on:click={switchToSettings}
      />
    </li>
    <div class="divider my-1"></div>
    <li>
      <ActionButton
        variant="ghost"
        customClass="justify-start text-error"
        icon="logout"
        text="Logout"
        on:click={handleLogout}
      />
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