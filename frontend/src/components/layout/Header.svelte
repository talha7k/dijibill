<script>
  import { createEventDispatcher } from 'svelte';
  import UserMenuDropdown from './UserMenuDropdown.svelte';
  
  export let isAuthenticated = false;
  /** @type {{ username: string, role: string } | null} */
  export let currentUser = null;
  export let currentPageIcon = '';
  export let currentPageTitle = '';
  export let wailsReady = false;
  export let wailsError = null;
  export let backupStatus = 'unknown';
  
  const dispatch = createEventDispatcher();
  
  let userMenuOpen = false;
  let userMenuButton;
  let userMenuButtonRect = null;
  
  function toggleSidebar() {
    dispatch('toggle-sidebar');
  }
  
  function toggleUserMenu(event) {
    event.stopPropagation();
    if (!userMenuOpen && userMenuButton) {
      userMenuButtonRect = userMenuButton.getBoundingClientRect();
    }
    userMenuOpen = !userMenuOpen;
  }
  
  function closeUserMenu() {
    userMenuOpen = false;
  }
  
  function showUserProfile() {
    dispatch('show-user-profile');
  }
  
  function showAccessControl() {
    dispatch('show-access-control');
  }
  
  function switchToSettings() {
    dispatch('view-change', 'general-settings');
  }
  
  function handleLogout() {
    dispatch('logout');
  }
  
  function retryInitialization() {
    dispatch('retry-initialization');
  }
</script>

<header class="bg-base-100/10 backdrop-blur-lg border-b border-white/20 p-4">
  <div class="flex items-center justify-between">
    <div class="flex items-center gap-4">
      <!-- Sidebar Toggle Button - Only visible when authenticated -->
      {#if isAuthenticated}
        <button
          class="btn btn-ghost btn-circle text-white hover:bg-white/10"
          on:click={toggleSidebar}
          title="Toggle Menu"
        >
          <i class="fas fa-bars w-5 h-5"></i>
        </button>
      {/if}
      
      <div class="flex items-center gap-3">
        <!-- Page Icon -->
        <div class="w-10 h-10 bg-primary/20 rounded-lg flex items-center justify-center">
          <i class="fas fa-{currentPageIcon} w-6 h-6 text-white"></i>
        </div>
        
        <!-- Page Title -->
        <h2 class="text-2xl font-bold text-white">
          {currentPageTitle}
        </h2>
      </div>
    </div>
    
    <div class="flex items-center gap-4">
      <!-- Runtime Status -->
      <div class="flex items-center gap-2">
        {#if wailsReady}
          <div class="badge badge-success badge-sm">
            <i class="fas fa-check-circle w-3 h-3 mr-1"></i>
            Connected
          </div>
        {:else if wailsError}
          <div class="badge badge-error badge-sm">
            <i class="fas fa-exclamation-circle w-3 h-3 mr-1"></i>
            Error
          </div>
        {:else}
          <div class="badge badge-warning badge-sm">
            <div class="loading loading-spinner loading-xs mr-1"></div>
            Connecting...
          </div>
        {/if}
      </div>
      
      <!-- Backup Status -->
      <div class="flex items-center gap-2">
        {#if backupStatus === 'recent'}
          <div class="badge badge-success badge-sm">
            <i class="fas fa-database w-3 h-3 mr-1"></i>
            Backup OK
          </div>
        {:else if backupStatus === 'overdue'}
          <div class="badge badge-warning badge-sm">
            <i class="fas fa-exclamation-triangle w-3 h-3 mr-1"></i>
            Backup Overdue
          </div>
        {:else if backupStatus === 'never'}
          <div class="badge badge-error badge-sm">
            <i class="fas fa-exclamation-circle w-3 h-3 mr-1"></i>
            No Backup
          </div>
        {:else if backupStatus === 'disabled'}
          <div class="badge badge-neutral badge-sm">
            <i class="fas fa-ban w-3 h-3 mr-1"></i>
            Backup Disabled
          </div>
        {:else}
          <div class="badge badge-ghost badge-sm">
            <i class="fas fa-question-circle w-3 h-3 mr-1"></i>
            Backup Status Unknown
          </div>
        {/if}
      </div>
      
      <!-- User Menu -->
      {#if isAuthenticated && currentUser}
        <button 
          bind:this={userMenuButton}
          class="btn btn-ghost btn-circle text-white"
          on:click={toggleUserMenu}
        >
          <i class="fas fa-user w-5 h-5"></i>
        </button>
      {/if}
    </div>
  </div>
</header>

<!-- User Menu Dropdown - Rendered outside of header stacking context -->
<UserMenuDropdown 
  {currentUser}
  isOpen={userMenuOpen}
  buttonRect={userMenuButtonRect}
  on:close={closeUserMenu}
  on:show-user-profile={showUserProfile}
  on:show-access-control={showAccessControl}
  on:view-change={switchToSettings}
  on:logout={handleLogout}
/>