<script>
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import ActionButton from './ActionButton.svelte'

  const dispatch = createEventDispatcher()

  export let show = false
  export let currentUser = {
    id: 0,
    username: '',
    email: '',
    first_name: '',
    last_name: '',
    role: '',
    is_active: true,
    company_id: 0
  }

  // Define page permissions based on role
  const pagePermissions = {
    admin: [
      { id: 'dashboard', name: 'Dashboard', icon: 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z', access: true },
      { id: 'pos', name: 'Point of Sale', icon: 'M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z', access: true },
      { id: 'sales-invoices', name: 'Sales Invoices', icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z', access: true },
      { id: 'customer', name: 'Customer Management', icon: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z', access: true },
      { id: 'payments', name: 'Payments', icon: 'M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z', access: true },
      { id: 'purchase-invoices', name: 'Purchase Invoices', icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z', access: true },
      { id: 'purchase-products', name: 'Purchase Products', icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4', access: true },
      { id: 'suppliers', name: 'Suppliers', icon: 'M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4', access: true },
      { id: 'products', name: 'Products', icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4', access: true },
      { id: 'users', name: 'User Management', icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z', access: true },
      { id: 'general-settings', name: 'General Settings', icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z', access: true },
      { id: 'administration', name: 'Administration', icon: 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z', access: true },
      { id: 'qr-validation', name: 'QR Validation', icon: 'M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z', access: true }
    ],
    user: [
      { id: 'dashboard', name: 'Dashboard', icon: 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z', access: true },
      { id: 'pos', name: 'Point of Sale', icon: 'M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z', access: true },
      { id: 'sales-invoices', name: 'Sales Invoices', icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z', access: true },
      { id: 'customer', name: 'Customer Management', icon: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z', access: true },
      { id: 'payments', name: 'Payments', icon: 'M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z', access: false },
      { id: 'purchase-invoices', name: 'Purchase Invoices', icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z', access: false },
      { id: 'purchase-products', name: 'Purchase Products', icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4', access: false },
      { id: 'suppliers', name: 'Suppliers', icon: 'M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4', access: false },
      { id: 'products', name: 'Products', icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4', access: true },
      { id: 'users', name: 'User Management', icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z', access: false },
      { id: 'general-settings', name: 'General Settings', icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z', access: false },
      { id: 'administration', name: 'Administration', icon: 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z', access: false },
      { id: 'qr-validation', name: 'QR Validation', icon: 'M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z', access: true }
    ]
  }

  $: userPermissions = pagePermissions[currentUser.role] || pagePermissions.user
  $: allowedPages = userPermissions.filter(page => page.access)
  $: restrictedPages = userPermissions.filter(page => !page.access)

  function closeModal() {
    show = false
    dispatch('close')
  }
</script>

<Modal bind:show title="Access Control & Permissions" size="lg" on:close={closeModal}>
  <div class="space-y-6">
    <!-- User Role Info -->
    <div class="bg-base-200 p-4 rounded-lg">
      <div class="flex items-center gap-4">
        <div class="avatar placeholder">
          <div class="bg-primary text-primary-content rounded-full w-12">
            <span class="text-lg">{currentUser.first_name?.[0] || 'U'}{currentUser.last_name?.[0] || ''}</span>
          </div>
        </div>
        <div>
          <h3 class="text-lg font-semibold">{currentUser.first_name} {currentUser.last_name}</h3>
          <p class="text-sm opacity-70">@{currentUser.username}</p>
          <div class="flex items-center gap-2 mt-1">
            <div class="badge badge-primary badge-sm capitalize">{currentUser.role}</div>
            {#if currentUser.role === 'admin'}
              <div class="badge badge-success badge-sm">Full Access</div>
            {:else}
              <div class="badge badge-warning badge-sm">Limited Access</div>
            {/if}
          </div>
        </div>
      </div>
    </div>

    <!-- Role Description -->
    <div class="alert {currentUser.role === 'admin' ? 'alert-success' : 'alert-info'}">
      {#if currentUser.role === 'admin'}
        <i class="fas fa-shield-alt w-6 h-6"></i>
      {:else}
        <i class="fas fa-info-circle w-6 h-6"></i>
      {/if}
      <div>
        <h3 class="font-bold">
          {currentUser.role === 'admin' ? 'Administrator' : 'Standard User'}
        </h3>
        <div class="text-xs">
          {#if currentUser.role === 'admin'}
            You have full administrative access to all features and can manage other users.
          {:else}
            Your access is limited to specific features as assigned by the administrator.
          {/if}
        </div>
      </div>
    </div>

    <!-- Accessible Pages -->
    <div>
      <h4 class="text-lg font-semibold mb-3 flex items-center gap-2">
        <i class="fas fa-check-circle w-5 h-5 text-success"></i>
        Accessible Pages ({allowedPages.length})
      </h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        {#each allowedPages as page}
          <div class="flex items-center gap-3 p-3 bg-success/10 border border-success/20 rounded-lg">
            <div class="w-8 h-8 bg-success/20 rounded-lg flex items-center justify-center">
              <i class="fas fa-check w-4 h-4 text-success"></i>
            </div>
            <div>
              <p class="font-medium text-sm">{page.name}</p>
              <p class="text-xs opacity-70">Full access granted</p>
            </div>
          </div>
        {/each}
      </div>
    </div>

    <!-- Restricted Pages -->
    {#if restrictedPages.length > 0}
      <div>
        <h4 class="text-lg font-semibold mb-3 flex items-center gap-2">
          <i class="fas fa-times-circle w-5 h-5 text-error"></i>
          Restricted Pages ({restrictedPages.length})
        </h4>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          {#each restrictedPages as page}
            <div class="flex items-center gap-3 p-3 bg-error/10 border border-error/20 rounded-lg">
              <div class="w-8 h-8 bg-error/20 rounded-lg flex items-center justify-center">
                <i class="fas fa-times w-4 h-4 text-error"></i>
              </div>
              <div>
                <p class="font-medium text-sm">{page.name}</p>
                <p class="text-xs opacity-70">Access denied</p>
              </div>
            </div>
          {/each}
        </div>
        <div class="mt-3 p-3 bg-warning/10 border border-warning/20 rounded-lg">
          <p class="text-sm text-warning">
            <i class="fas fa-exclamation-triangle w-4 h-4 inline mr-1"></i>
            Contact your administrator to request access to restricted pages.
          </p>
        </div>
      </div>
    {/if}

    <!-- Close Button -->
    <div class="flex justify-end pt-4">
      <ActionButton
        variant="primary"
        text="Close"
        on:click={closeModal}
      />
    </div>
  </div>
</Modal>