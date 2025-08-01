<script>
  import { onMount } from 'svelte';
  import { wailsReady, wailsError, initializeWailsRuntime } from './stores/wailsStore.js';
  
  // Layout Components
  import Sidebar from './components/layout/Sidebar.svelte';
  import Header from './components/layout/Header.svelte';
  import LoadingScreen from './components/layout/LoadingScreen.svelte';
  import AuthenticationScreen from './components/layout/AuthenticationScreen.svelte';
  import MainContent from './components/layout/MainContent.svelte';
  
  // Modal Components
  import IntroSlider from './components/IntroSlider.svelte';
  import UserProfile from './components/UserProfile.svelte';
  import AccessControl from './components/AccessControl.svelte';
  
  // Wails functions
  import { 
    Greet, 
    GetSystemSettings, 
    GetAuthContext, 
    GetCurrentUser, 
    MarkIntroAsViewed, 
    Logout
  } from '../wailsjs/go/main/App.js'

  let currentView = 'dashboard' // Default to dashboard
  let resultText = "Welcome to DijiBill - ZATCA Compliant Invoicing"
  let sidebarVisible = false // Sidebar visibility state - hidden by default
  let sidebarAsOverlay = true // Whether sidebar should appear as overlay
  
  // Authentication state
  let isAuthenticated = false
  let authContext = null
  let currentUser = null
  let currentCompany = null
  let authMode = 'login' // 'login' or 'signup'
  
  // Intro slider state
  let showIntroSlider = false
  let userHasViewedIntro = false
  
  // Modal states
  let showUserProfile = false
  let showAccessControl = false
  
  // Get runtime status for display
  $: runtimeStatus = $wailsReady ? 'Runtime Ready' : $wailsError ? 'Runtime Error' : 'Initializing...'
  let name = ''
  let backupStatus = 'unknown' // 'recent', 'overdue', 'unknown'
  let lastBackupTime = null
  let autoBackupEnabled = false

  // Check backup status
  async function checkBackupStatus() {
    try {
      const settings = await GetSystemSettings()
      if (settings) {
        autoBackupEnabled = settings.auto_backup
        lastBackupTime = settings.last_backup_time
        
        if (!autoBackupEnabled) {
          backupStatus = 'disabled'
          return
        }
        
        if (!lastBackupTime) {
          backupStatus = 'never'
          return
        }
        
        // Convert Go time to JavaScript Date
        const lastBackup = new Date(lastBackupTime.toString())
        const now = new Date()
        const hoursSinceBackup = (now.getTime() - lastBackup.getTime()) / (1000 * 60 * 60)
        
        // Check backup frequency from settings
        const backupFrequency = Number(settings.backup_frequency) || 24 // default 24 hours
        
        if (hoursSinceBackup <= backupFrequency) {
          backupStatus = 'recent'
        } else {
          backupStatus = 'overdue'
        }
      }
    } catch (error) {
      console.error('Error checking backup status:', error)
      backupStatus = 'unknown'
    }
  }

  // Check authentication status
  async function checkAuthStatus() {
    try {
      const context = await GetAuthContext()
      if (context && context.user_id && context.company_id) {
        isAuthenticated = true
        authContext = context
        currentUser = {
          id: context.user_id,
          username: context.username,
          role: context.role
        }
        currentCompany = {
          id: context.company_id
        }
        
        // Check if user has viewed intro
        await checkIntroStatus()
      } else {
        isAuthenticated = false
        authContext = null
        currentUser = null
        currentCompany = null
        // Show intro slider for non-authenticated users
        showIntroSlider = true
        userHasViewedIntro = false
      }
    } catch (error) {
      console.error('Error checking auth status:', error)
      isAuthenticated = false
      authContext = null
      currentUser = null
      currentCompany = null
      // Show intro slider for non-authenticated users
      showIntroSlider = true
      userHasViewedIntro = false
    }
  }

  // Check intro status
  async function checkIntroStatus() {
    try {
      const user = await GetCurrentUser()
      if (user) {
        userHasViewedIntro = user.intro_viewed || false
        // Show intro slider if user hasn't viewed it yet
        if (!userHasViewedIntro) {
          showIntroSlider = true
        }
      }
    } catch (error) {
      console.error('Error checking intro status:', error)
      userHasViewedIntro = false
    }
  }

  // Handle login success
  function handleLoginSuccess(event) {
    const context = event.detail
    isAuthenticated = true
    authContext = context
    currentUser = {
      id: context.user_id,
      username: context.username,
      role: context.role
    }
    currentCompany = {
      id: context.company_id
    }
    
    // Check intro status after login
    checkIntroStatus()
  }

  // Handle signup success
  function handleSignupSuccess(event) {
    const context = event.detail
    isAuthenticated = true
    authContext = context
    currentUser = {
      id: context.user_id,
      username: context.username,
      role: context.role
    }
    currentCompany = {
      id: context.company_id
    }
    
    // New users should see the intro
    showIntroSlider = true
    userHasViewedIntro = false
  }

  // Switch between login and signup modes
  function switchToSignup() {
    authMode = 'signup'
  }

  function switchToLogin() {
    authMode = 'login'
  }

  // Handle intro completion
  async function handleIntroComplete() {
    try {
      if (currentUser && isAuthenticated) {
        // For authenticated users, mark intro as viewed in the database
        await MarkIntroAsViewed(currentUser.id)
        userHasViewedIntro = true
      }
      // For both authenticated and non-authenticated users, hide the intro slider
      showIntroSlider = false
    } catch (error) {
      console.error('Error marking intro as viewed:', error)
      // Still hide the intro slider even if there's an error
      showIntroSlider = false
    }
  }

  // Show intro slider manually (for administration page)
  function showIntro() {
    showIntroSlider = true
  }

  // Handle logout
  async function handleLogout() {
    try {
      await Logout()
      isAuthenticated = false
      authContext = null
      currentUser = null
      currentCompany = null
    } catch (error) {
      console.error('Error during logout:', error)
    }
  }

  // Handle profile update
  async function handleProfileUpdate(event) {
    try {
      const { UpdateUser } = await import('../wailsjs/go/main/App')
      const { database } = await import('../wailsjs/go/models')
      const updatedData = event.detail
      
      // Create updated user object using the User constructor
      const updatedUser = new database.User({
        id: currentUser.id,
        username: updatedData.username,
        email: updatedData.email,
        first_name: updatedData.first_name,
        last_name: updatedData.last_name,
        role: currentUser.role,
        company_id: currentCompany.id,
        password: currentUser.password || '',
        is_active: currentUser.is_active !== undefined ? currentUser.is_active : true,
        intro_viewed: currentUser.intro_viewed || false,
        created_at: currentUser.created_at || new Date().toISOString(),
        updated_at: new Date().toISOString()
      })
      
      await UpdateUser(updatedUser)
      
      // Update current user state
      currentUser = {
        ...currentUser,
        username: updatedData.username,
        email: updatedData.email,
        first_name: updatedData.first_name,
        last_name: updatedData.last_name
      }
      
      showUserProfile = false
      console.log('Profile updated successfully')
    } catch (error) {
      console.error('Error updating profile:', error)
    }
  }

  // Handle password change
  async function handlePasswordChange(event) {
    try {
      const { UpdateUser } = await import('../wailsjs/go/main/App')
      const { database } = await import('../wailsjs/go/models')
      const { newPassword } = event.detail
      
      // Create updated user object using the User constructor
      const updatedUser = new database.User({
        id: currentUser.id,
        username: currentUser.username,
        email: currentUser.email,
        first_name: currentUser.first_name,
        last_name: currentUser.last_name,
        role: currentUser.role,
        company_id: currentCompany.id,
        password: newPassword,
        is_active: currentUser.is_active !== undefined ? currentUser.is_active : true,
        intro_viewed: currentUser.intro_viewed || false,
        created_at: currentUser.created_at || new Date().toISOString(),
        updated_at: new Date().toISOString()
      })
      
      await UpdateUser(updatedUser)
      showUserProfile = false
      console.log('Password changed successfully')
    } catch (error) {
      console.error('Error changing password:', error)
    }
  }



  // Initialize Wails runtime on mount
  onMount(async () => {
    console.log('🎯 App component mounted, initializing Wails runtime...');
    await initializeWailsRuntime();
    
    // Check authentication status when app loads
    await checkAuthStatus();
    
    // Check backup status when app loads
    checkBackupStatus();
    
    // Check backup status every 5 minutes
    const backupStatusInterval = setInterval(checkBackupStatus, 5 * 60 * 1000);
    
    return () => {
      clearInterval(backupStatusInterval);
    };
  });

  // Retry function for manual retry
  async function retryInitialization() {
    console.log('🔄 Retrying Wails runtime initialization...');
    await initializeWailsRuntime();
  }

  function greet() {
    Greet(name).then(result => resultText = result)
  }

  // Menu configuration
  const menuItems = [
    {
      id: 'dashboard',
      label: 'Dashboard',
      icon: 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2 2v-2m0 0V5a2 2 0 012-2h6l2 2h6a2 2 0 012 2v2M7 13h10M7 17h4'
    },
    {
      category: 'Sales',
      items: [
        {
          id: 'pos',
          label: 'Point of Sale',
          icon: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
        },
        {
          id: 'sales-invoices',
          label: 'Sales Invoices',
          icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z'
        },
        {
          id: 'customer',
          label: 'Customer',
          icon: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z'
        },
        {
          id: 'payments',
          label: 'Payments',
          icon: 'M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z'
        }
      ]
    },
    {
      category: 'Purchases',
      items: [
        {
          id: 'purchase-invoices',
          label: 'Purchase Invoices',
          icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z'
        },
        {
          id: 'purchase-products',
          label: 'Purchase Products',
          icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4'
        },
        {
          id: 'suppliers',
          label: 'Suppliers',
          icon: 'M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4'
        }
      ]
    },
    {
      category: 'Other',
      items: [
        {
          id: 'products',
          label: 'Products',
          icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4'
        },
        {
          id: 'users',
          label: 'Users',
          icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z'
        },
        {
          id: 'general-settings',
          label: 'General settings',
          icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z'
        },
        {
          id: 'administration',
          label: 'Administration',
          icon: 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z'
        },
        {
          id: 'qr-validation',
          label: 'QR Validation',
          icon: 'M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z'
        }
      ]
    }
  ]

  // Event handlers
  function switchView(view) {
    currentView = view;
  }

  function toggleSidebar() {
    sidebarVisible = !sidebarVisible;
  }

  function handleViewChange(event) {
    switchView(event.detail);
  }

  function handleSidebarClose() {
    sidebarVisible = false;
  }

  function handleShowUserProfile() {
    showUserProfile = true;
  }

  function handleShowAccessControl() {
    showAccessControl = true;
  }

  function handleRetryInitialization() {
    retryInitialization();
  }

  function handleIntroShow() {
    showIntroSlider = true;
  }

  // Reactive functions to get the current page icon and title
  $: currentPageIcon = (() => {
    // Dashboard icon
    if (currentView === 'dashboard') {
      return menuItems[0].icon
    }
    
    // Search through categories for the current view
    for (const category of menuItems.slice(1)) {
      const item = category.items.find(item => item.id === currentView)
      if (item) {
        return item.icon
      }
    }
    
    // Default icon if not found
    return 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z'
  })()

  $: currentPageTitle = (() => {
    if (currentView === 'dashboard') return 'Dashboard'
    if (currentView === 'pos') return 'Point of Sale'
    if (currentView === 'sales-invoices') return 'Sales Invoices'
    if (currentView === 'customer') return 'Customer'
    if (currentView === 'payments') return 'Payments'
    if (currentView === 'purchase-invoices') return 'Purchase Invoices'
    if (currentView === 'purchase-products') return 'Purchase Products'
    if (currentView === 'suppliers') return 'Suppliers'
    if (currentView === 'products') return 'Products'
    if (currentView === 'users') return 'Users'
    if (currentView === 'general-settings') return 'General Settings'
    if (currentView === 'administration') return 'Administration'
    if (currentView === 'qr-validation') return 'QR Validation'
    return 'Page'
  })()
</script>

<div class="min-h-screen bg-gradient-to-br from-primary to-secondary flex relative" data-theme="dijibill">
  <!-- Sidebar Component -->
  <Sidebar 
    isVisible={sidebarVisible}
    asOverlay={sidebarAsOverlay}
    {currentView}
    {menuItems}
    on:view-change={handleViewChange}
    on:close={handleSidebarClose}
  />

  <!-- Main Content Area -->
  <div class="flex-1 min-w-0 flex flex-col">
    <!-- Header Component -->
    <Header 
      {isAuthenticated}
      {currentUser}
      {currentPageIcon}
      {currentPageTitle}
      wailsReady={$wailsReady}
      wailsError={$wailsError}
      {backupStatus}
      on:toggle-sidebar={toggleSidebar}
      on:show-user-profile={handleShowUserProfile}
      on:show-access-control={handleShowAccessControl}
      on:view-change={handleViewChange}
      on:logout={handleLogout}
      on:retry-initialization={handleRetryInitialization}
    />

    <!-- Page Content -->
    {#if !$wailsReady}
      <!-- Loading Screen Component -->
      <LoadingScreen 
        wailsReady={$wailsReady}
        wailsError={$wailsError}
        on:retry-initialization={handleRetryInitialization}
      />
    {:else if !isAuthenticated}
      <!-- Authentication Screen Component -->
      <AuthenticationScreen 
        {authMode}
        on:login-success={handleLoginSuccess} 
        on:signup-success={handleSignupSuccess}
        on:switch-to-signup={switchToSignup}
        on:switch-to-login={switchToLogin}
      />
    {:else}
      <!-- Main Content Component -->
      <MainContent 
        {currentView}
        showIntro={handleIntroShow}
      />
    {/if}
  </div>
  
  <!-- Intro Slider Overlay -->
  {#if showIntroSlider}
    <IntroSlider show={showIntroSlider} on:complete={handleIntroComplete} />
  {/if}
  
  <!-- User Profile Modal -->
  {#if showUserProfile && currentUser}
    <UserProfile 
      {currentUser}
      show={showUserProfile}
      on:close={() => showUserProfile = false}
      on:update-profile={handleProfileUpdate}
      on:change-password={handlePasswordChange}
    />
  {/if}
  
  <!-- Access Control Modal -->
  {#if showAccessControl && currentUser}
    <AccessControl 
      {currentUser}
      show={showAccessControl}
      on:close={() => showAccessControl = false}
    />
  {/if}
</div>
