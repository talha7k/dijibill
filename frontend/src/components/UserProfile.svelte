<script>
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'

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

  let isLoading = false
  let error = ''
  let success = ''
  let showPasswordChange = false

  // Form data
  let profileForm = {
    username: '',
    email: '',
    first_name: '',
    last_name: ''
  }

  let passwordForm = {
    current_password: '',
    new_password: '',
    confirm_password: ''
  }

  // Initialize form when user data is available
  $: if (currentUser && show) {
    profileForm = {
      username: currentUser.username || '',
      email: currentUser.email || '',
      first_name: currentUser.first_name || '',
      last_name: currentUser.last_name || ''
    }
  }

  function closeModal() {
    show = false
    error = ''
    success = ''
    showPasswordChange = false
    passwordForm = {
      current_password: '',
      new_password: '',
      confirm_password: ''
    }
    dispatch('close')
  }

  function validateProfileForm() {
    const errors = []
    
    if (!profileForm.username.trim()) {
      errors.push('Username is required')
    }
    
    if (!profileForm.email.trim()) {
      errors.push('Email is required')
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(profileForm.email)) {
      errors.push('Please enter a valid email address')
    }
    
    if (!profileForm.first_name.trim()) {
      errors.push('First name is required')
    }
    
    if (!profileForm.last_name.trim()) {
      errors.push('Last name is required')
    }
    
    return errors
  }

  function validatePasswordForm() {
    const errors = []
    
    if (!passwordForm.current_password) {
      errors.push('Current password is required')
    }
    
    if (!passwordForm.new_password) {
      errors.push('New password is required')
    } else if (passwordForm.new_password.length < 6) {
      errors.push('New password must be at least 6 characters long')
    }
    
    if (passwordForm.new_password !== passwordForm.confirm_password) {
      errors.push('New passwords do not match')
    }
    
    return errors
  }

  async function handleUpdateProfile() {
    error = ''
    success = ''
    
    const validationErrors = validateProfileForm()
    if (validationErrors.length > 0) {
      error = validationErrors.join(', ')
      return
    }

    isLoading = true
    
    try {
      // Dispatch the update request to parent component
      dispatch('update-profile', {
        username: profileForm.username,
        email: profileForm.email,
        first_name: profileForm.first_name,
        last_name: profileForm.last_name
      })
      
      success = 'Profile updated successfully!'
      
      setTimeout(() => {
        success = ''
      }, 3000)
      
    } catch (err) {
      console.error('Error updating profile:', err)
      error = err.message || 'Failed to update profile'
    } finally {
      isLoading = false
    }
  }

  async function handleChangePassword() {
    error = ''
    success = ''
    
    const validationErrors = validatePasswordForm()
    if (validationErrors.length > 0) {
      error = validationErrors.join(', ')
      return
    }

    isLoading = true
    
    try {
      // Dispatch the password change request to parent component
      dispatch('change-password', {
        current_password: passwordForm.current_password,
        new_password: passwordForm.new_password
      })
      
      success = 'Password changed successfully!'
      showPasswordChange = false
      passwordForm = {
        current_password: '',
        new_password: '',
        confirm_password: ''
      }
      
      setTimeout(() => {
        success = ''
      }, 3000)
      
    } catch (err) {
      console.error('Error changing password:', err)
      error = err.message || 'Failed to change password'
    } finally {
      isLoading = false
    }
  }
</script>

<Modal bind:show title="User Profile" size="lg" on:close={closeModal}>
  <div class="space-y-6">
    <!-- Success Message -->
    {#if success}
      <div class="alert alert-success">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <span>{success}</span>
      </div>
    {/if}

    <!-- Error Message -->
    {#if error}
      <div class="alert alert-error">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
        </svg>
        <span>{error}</span>
      </div>
    {/if}

    <!-- User Info Display -->
    {#if currentUser}
      <div class="bg-base-200 p-4 rounded-lg">
        <div class="flex items-center gap-4">
          <div class="avatar placeholder">
            <div class="bg-primary text-primary-content rounded-full w-16">
              <span class="text-xl">{currentUser.first_name?.[0] || 'U'}{currentUser.last_name?.[0] || ''}</span>
            </div>
          </div>
          <div>
            <h3 class="text-lg font-semibold">{currentUser.first_name} {currentUser.last_name}</h3>
            <p class="text-sm opacity-70">@{currentUser.username}</p>
            <div class="badge badge-primary badge-sm mt-1 capitalize">{currentUser.role}</div>
          </div>
        </div>
      </div>
    {/if}

    <!-- Tabs -->
    <div class="tabs tabs-boxed">
      <button 
        class="tab {!showPasswordChange ? 'tab-active' : ''}" 
        on:click={() => showPasswordChange = false}
      >
        Profile Information
      </button>
      <button 
        class="tab {showPasswordChange ? 'tab-active' : ''}" 
        on:click={() => showPasswordChange = true}
      >
        Change Password
      </button>
    </div>

    {#if !showPasswordChange}
      <!-- Profile Form -->
      <form on:submit|preventDefault={handleUpdateProfile} class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField
            type="text"
            label="Username"
            placeholder="Enter username"
            bind:value={profileForm.username}
            required
          />
          
          <FormField
            type="email"
            label="Email"
            placeholder="Enter email address"
            bind:value={profileForm.email}
            required
          />
          
          <FormField
            type="text"
            label="First Name"
            placeholder="Enter first name"
            bind:value={profileForm.first_name}
            required
          />
          
          <FormField
            type="text"
            label="Last Name"
            placeholder="Enter last name"
            bind:value={profileForm.last_name}
            required
          />
        </div>

        <div class="flex justify-end gap-2 pt-4">
          <button type="button" class="btn btn-ghost" on:click={closeModal}>
            Cancel
          </button>
          <button type="submit" class="btn btn-primary" disabled={isLoading}>
            {#if isLoading}
              <span class="loading loading-spinner loading-sm"></span>
              Updating...
            {:else}
              Update Profile
            {/if}
          </button>
        </div>
      </form>
    {:else}
      <!-- Password Change Form -->
      <form on:submit|preventDefault={handleChangePassword} class="space-y-4">
        <FormField
          type="password"
          label="Current Password"
          placeholder="Enter current password"
          bind:value={passwordForm.current_password}
          required
        />
        
        <FormField
          type="password"
          label="New Password"
          placeholder="Enter new password"
          bind:value={passwordForm.new_password}
          required
        />
        
        <FormField
          type="password"
          label="Confirm New Password"
          placeholder="Confirm new password"
          bind:value={passwordForm.confirm_password}
          required
        />

        <div class="flex justify-end gap-2 pt-4">
          <button type="button" class="btn btn-ghost" on:click={closeModal}>
            Cancel
          </button>
          <button type="submit" class="btn btn-primary" disabled={isLoading}>
            {#if isLoading}
              <span class="loading loading-spinner loading-sm"></span>
              Changing...
            {:else}
              Change Password
            {/if}
          </button>
        </div>
      </form>
    {/if}
  </div>
</Modal>