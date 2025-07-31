<script>
  import { createEventDispatcher } from 'svelte'
  import Modal from './Modal.svelte'
  import FormField from './FormField.svelte'

  const dispatch = createEventDispatcher()

  export let show = false
  export let user = null // null for new user, user object for editing
  export let currentUser = {
    id: 0,
    username: '',
    email: '',
    first_name: '',
    last_name: '',
    role: '',
    is_active: true,
    company_id: 0
  } // current logged-in user for permission checks

  let formData = {
    username: '',
    email: '',
    first_name: '',
    last_name: '',
    password: '',
    confirmPassword: '',
    role: 'user',
    is_active: true
  }

  let errors = {}
  let isLoading = false

  // Reset form when modal opens/closes or user changes
  $: if (show) {
    resetForm()
  }

  function resetForm() {
    if (user) {
      // Editing existing user
      formData = {
        username: user.username || '',
        email: user.email || '',
        first_name: user.first_name || '',
        last_name: user.last_name || '',
        password: '',
        confirmPassword: '',
        role: user.role || 'user',
        is_active: user.is_active !== undefined ? user.is_active : true
      }
    } else {
      // Adding new user
      formData = {
        username: '',
        email: '',
        first_name: '',
        last_name: '',
        password: '',
        confirmPassword: '',
        role: 'user',
        is_active: true
      }
    }
    errors = {}
  }

  function validateForm() {
    errors = {}

    if (!formData.username.trim()) {
      errors.username = 'Username is required'
    }

    if (!formData.email.trim()) {
      errors.email = 'Email is required'
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      errors.email = 'Please enter a valid email address'
    }

    if (!formData.first_name.trim()) {
      errors.first_name = 'First name is required'
    }

    if (!formData.last_name.trim()) {
      errors.last_name = 'Last name is required'
    }

    // Password validation for new users or when changing password
    if (!user || formData.password) {
      if (!formData.password) {
        errors.password = 'Password is required'
      } else if (formData.password.length < 6) {
        errors.password = 'Password must be at least 6 characters'
      }

      if (formData.password !== formData.confirmPassword) {
        errors.confirmPassword = 'Passwords do not match'
      }
    }

    return Object.keys(errors).length === 0
  }

  async function handleSubmit() {
    if (!validateForm()) return

    isLoading = true
    try {
      const userData = {
        username: formData.username,
        email: formData.email,
        first_name: formData.first_name,
        last_name: formData.last_name,
        role: formData.role,
        is_active: formData.is_active
      }

      // Only include password if it's provided
      if (formData.password) {
        userData.password = formData.password
      }

      if (user) {
        // Editing existing user
        userData.id = user.id
        dispatch('update-user', userData)
      } else {
        // Creating new user
        dispatch('create-user', userData)
      }

      closeModal()
    } catch (error) {
      console.error('Error saving user:', error)
    } finally {
      isLoading = false
    }
  }

  function closeModal() {
    show = false
    dispatch('close')
  }

  // Check if current user can manage roles
  $: canManageRoles = currentUser && currentUser.role === 'admin'
</script>

<Modal bind:show title="{user ? 'Edit User' : 'Add New User'}" size="lg" on:close={closeModal}>
  <form on:submit|preventDefault={handleSubmit} class="space-y-4">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Username -->
      <FormField
        label="Username"
        type="text"
        bind:value={formData.username}
        error={errors.username}
        required
        placeholder="Enter username"
      />

      <!-- Email -->
      <FormField
        label="Email"
        type="email"
        bind:value={formData.email}
        error={errors.email}
        required
        placeholder="Enter email address"
      />

      <!-- First Name -->
      <FormField
        label="First Name"
        type="text"
        bind:value={formData.first_name}
        error={errors.first_name}
        required
        placeholder="Enter first name"
      />

      <!-- Last Name -->
      <FormField
        label="Last Name"
        type="text"
        bind:value={formData.last_name}
        error={errors.last_name}
        required
        placeholder="Enter last name"
      />
    </div>

    <!-- Password Section -->
    <div class="divider">
      {user ? 'Change Password (Optional)' : 'Password'}
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Password -->
      <FormField
        label="Password"
        type="password"
        bind:value={formData.password}
        error={errors.password}
        required={!user}
        placeholder={user ? 'Leave blank to keep current password' : 'Enter password'}
      />

      <!-- Confirm Password -->
      <FormField
        label="Confirm Password"
        type="password"
        bind:value={formData.confirmPassword}
        error={errors.confirmPassword}
        required={Boolean(!user && formData.password)}
        placeholder="Confirm password"
      />
    </div>

    <!-- Role and Status Section -->
    <div class="divider">Permissions & Status</div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Role Selection -->
      <div class="form-control">
        <label class="label">
          <span class="label-text">Role</span>
        </label>
        <select 
          class="select select-bordered w-full" 
          bind:value={formData.role}
          disabled={!canManageRoles}
        >
          <option value="user">Standard User</option>
          <option value="admin">Administrator</option>
        </select>
        {#if !canManageRoles}
          <label class="label">
            <span class="label-text-alt text-warning">Only administrators can change roles</span>
          </label>
        {/if}
      </div>

      <!-- Active Status -->
      <div class="form-control">
        <label class="label">
          <span class="label-text">Status</span>
        </label>
        <div class="flex items-center gap-2 mt-2">
          <input 
            type="checkbox" 
            class="toggle toggle-success" 
            bind:checked={formData.is_active}
          />
          <span class="text-sm">
            {formData.is_active ? 'Active' : 'Inactive'}
          </span>
        </div>
      </div>
    </div>

    <!-- Role Description -->
    <div class="alert {formData.role === 'admin' ? 'alert-warning' : 'alert-info'}">
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
      </svg>
      <div>
        <h3 class="font-bold">{formData.role === 'admin' ? 'Administrator' : 'Standard User'}</h3>
        <div class="text-xs">
          {#if formData.role === 'admin'}
            Full access to all features including user management, system settings, and administration tools.
          {:else}
            Access to core features like POS, invoices, customers, and products. Cannot access user management or system administration.
          {/if}
        </div>
      </div>
    </div>

    <!-- Form Actions -->
    <div class="modal-action">
      <button type="button" class="btn btn-ghost" on:click={closeModal}>
        Cancel
      </button>
      <button 
        type="submit" 
        class="btn btn-primary"
        class:loading={isLoading}
        disabled={isLoading}
      >
        {#if isLoading}
          Saving...
        {:else}
          {user ? 'Update User' : 'Create User'}
        {/if}
      </button>
    </div>
  </form>
</Modal>