<script>
  import { createEventDispatcher } from 'svelte';
  
  export let isAuthenticated = false;
  /** @type {{ username: string, role: string } | null} */
  export let currentUser = null;
  export let currentPageIcon = '';
  export let currentPageTitle = '';
  export let wailsReady = false;
  export let wailsError = null;
  export let backupStatus = 'unknown';
  
  const dispatch = createEventDispatcher();
  
  function toggleSidebar() {
    dispatch('toggle-sidebar');
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
        <div class="dropdown dropdown-end">
          <div tabindex="0" role="button" class="btn btn-ghost btn-circle text-white">
            <i class="fas fa-user w-5 h-5"></i>
          </div>
          <ul tabindex="0" class="menu menu-sm dropdown-content mt-3 z-[9999] p-2 shadow bg-base-100 rounded-box w-52">
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
        </div>
      {/if}
    </div>
  </div>
</header>