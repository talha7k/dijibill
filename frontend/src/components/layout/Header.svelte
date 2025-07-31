<script>
  import { createEventDispatcher } from 'svelte';
  
  export let isAuthenticated = false;
  export let currentUser = null; // { username: string, role: string } | null
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
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
          </svg>
        </button>
      {/if}
      
      <div class="flex items-center gap-3">
        <!-- Page Icon -->
        <div class="w-10 h-10 bg-primary/20 rounded-lg flex items-center justify-center">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="{currentPageIcon}"></path>
          </svg>
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
            <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
            </svg>
            Connected
          </div>
        {:else if wailsError}
          <div class="badge badge-error badge-sm">
            <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
            </svg>
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
            <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z"></path>
            </svg>
            Backup OK
          </div>
        {:else if backupStatus === 'overdue'}
          <div class="badge badge-warning badge-sm">
            <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
            </svg>
            Backup Overdue
          </div>
        {:else if backupStatus === 'never'}
          <div class="badge badge-error badge-sm">
            <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
            </svg>
            No Backup
          </div>
        {:else if backupStatus === 'disabled'}
          <div class="badge badge-neutral badge-sm">
            <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M13.477 14.89A6 6 0 015.11 6.524l8.367 8.368zm1.414-1.414L6.524 5.11a6 6 0 018.367 8.367zM18 10a8 8 0 11-16 0 8 8 0 0116 0z" clip-rule="evenodd"></path>
            </svg>
            Backup Disabled
          </div>
        {:else}
          <div class="badge badge-ghost badge-sm">
            <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path d="M9 12a1 1 0 002 0V8a1 1 0 10-2 0v4zm1-7a1 1 0 100 2 1 1 0 000-2z"></path>
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm0-2a6 6 0 100-12 6 6 0 000 12z" clip-rule="evenodd"></path>
            </svg>
            Backup Status Unknown
          </div>
        {/if}
      </div>
      
      <!-- User Menu -->
      {#if isAuthenticated && currentUser}
        <div class="dropdown dropdown-end">
          <div tabindex="0" role="button" class="btn btn-ghost btn-circle text-white">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
            </svg>
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
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                </svg>
                Profile
              </button>
            </li>
            <li>
              <button class="btn btn-ghost justify-start" on:click={showAccessControl}>
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
                </svg>
                Access Control
              </button>
            </li>
            <li>
              <button class="btn btn-ghost justify-start" on:click={switchToSettings}>
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                </svg>
                Settings
              </button>
            </li>
            <div class="divider my-1"></div>
            <li>
              <button class="btn btn-ghost justify-start text-error" on:click={handleLogout}>
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013 3v1"></path>
                </svg>
                Logout
              </button>
            </li>
          </ul>
        </div>
      {/if}
    </div>
  </div>
</header>