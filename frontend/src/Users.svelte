<script>
  import { onMount } from 'svelte'
  import { GetUsersByCompany, CreateUser, UpdateUser, DeleteUser, GetCurrentUser } from '../wailsjs/go/main/App'
  import { database } from '../wailsjs/go/models'
  import FormField from './components/FormField.svelte'
  import UserModal from './components/UserModal.svelte'
  import AccessControl from './components/AccessControl.svelte'
  import ConfirmationModal from './components/ConfirmationModal.svelte'
  import ActionButton from './components/ActionButton.svelte'
  import IconButton from './components/IconButton.svelte'
  import { showDbSuccess, showDbError } from './stores/notificationStore.js'
  
  let users = []
  let currentUser = null
  let isLoading = false
  let searchTerm = ''
  let showUserModal = false
  let showAccessControl = false
  let selectedUser = null
  let selectedUserForAccess = null
  let showDeleteConfirm = false
  let userToDelete = null
  let confirmLoading = false

  // Reactive filtered users
  $: filteredUsers = users.filter(user => {
    if (!searchTerm) return true
    const search = searchTerm.toLowerCase()
    return (
      user.username.toLowerCase().includes(search) ||
      user.email.toLowerCase().includes(search) ||
      user.first_name.toLowerCase().includes(search) ||
      user.last_name.toLowerCase().includes(search) ||
      user.role.toLowerCase().includes(search)
    )
  })

  onMount(async () => {
    await loadCurrentUser()
    await loadUsers()
  })

  async function loadCurrentUser() {
    try {
      currentUser = await GetCurrentUser()
    } catch (error) {
      console.error('Error loading current user:', error)
    }
  }

  async function loadUsers() {
    if (!currentUser) return
    
    isLoading = true
    try {
      users = await GetUsersByCompany(currentUser.company_id)
    } catch (error) {
      console.error('Error loading users:', error)
      showDbError('load', 'users', error)
    } finally {
      isLoading = false
    }
  }

  function openAddUserModal() {
    selectedUser = null
    showUserModal = true
  }

  function openEditUserModal(user) {
    selectedUser = user
    showUserModal = true
  }

  function openAccessControl(user) {
    selectedUserForAccess = user
    showAccessControl = true
  }

  async function handleCreateUser(event) {
    try {
      const userData = event.detail
      
      // Create user object with all required fields
      const newUser = new database.User({
        username: userData.username,
        email: userData.email,
        password: userData.password,
        first_name: userData.first_name,
        last_name: userData.last_name,
        role: userData.role,
        is_active: userData.is_active,
        company_id: currentUser.company_id,
        intro_viewed: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      })

      await CreateUser(newUser)
      await loadUsers()
      showDbSuccess('create', 'User')
    } catch (error) {
      console.error('Error creating user:', error)
      showDbError('create', 'User', error)
    }
  }

  async function handleUpdateUser(event) {
    try {
      const userData = event.detail
      
      // Find the existing user to preserve fields
      const existingUser = users.find(u => u.id === userData.id)
      if (!existingUser) {
        throw new Error('User not found')
      }

      // Create updated user object
      const updatedUser = new database.User({
        id: userData.id,
        username: userData.username,
        email: userData.email,
        password: userData.password || existingUser.password,
        first_name: userData.first_name,
        last_name: userData.last_name,
        role: userData.role,
        is_active: userData.is_active,
        company_id: existingUser.company_id,
        intro_viewed: existingUser.intro_viewed,
        last_login: existingUser.last_login,
        created_at: existingUser.created_at,
        updated_at: new Date().toISOString()
      })

      await UpdateUser(updatedUser)
      await loadUsers()
      showDbSuccess('update', 'User')
    } catch (error) {
      console.error('Error updating user:', error)
      showDbError('update', 'User', error)
    }
  }

  function handleDeleteUser(user) {
    if (user.id === currentUser.id) {
      showDbError('delete', 'User', new Error('You cannot delete your own account'))
      return
    }

    userToDelete = user
    showDeleteConfirm = true
  }

  function cancelDelete() {
    showDeleteConfirm = false
    userToDelete = null
  }

  async function confirmDelete() {
    if (!userToDelete) return
    
    try {
      confirmLoading = true
      await DeleteUser(userToDelete.id)
      await loadUsers()
      showDbSuccess('delete', 'User')
      showDeleteConfirm = false
      userToDelete = null
    } catch (error) {
      console.error('Error deleting user:', error)
      showDbError('delete', 'User', error)
    } finally {
      confirmLoading = false
    }
  }

  async function toggleUserStatus(user) {
    if (user.id === currentUser.id) {
      showDbError('update', 'User', new Error('You cannot deactivate your own account'))
      return
    }

    try {
      const updatedUser = new database.User({
        ...user,
        is_active: !user.is_active,
        updated_at: new Date().toISOString()
      })

      await UpdateUser(updatedUser)
      await loadUsers()
      showDbSuccess('update', 'User status')
    } catch (error) {
      console.error('Error updating user status:', error)
      showDbError('update', 'User status', error)
    }
  }

  function formatDate(dateString) {
    if (!dateString) return 'Never'
    try {
      return new Date(dateString).toLocaleDateString()
    } catch {
      return 'Invalid date'
    }
  }

  // Check if current user is admin
  $: isAdmin = currentUser && currentUser.role === 'admin'
</script>

<div class="p-6">
  <div class="card bg-base-100/20 backdrop-blur-lg border border-white/20 shadow-xl">
    <div class="card-body">
      <h2 class="card-title text-white text-2xl mb-4">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
        </svg>
        User Management
      </h2>
      
      <div class="flex justify-between items-center mb-6">
        <div class="flex gap-2">
          {#if isAdmin}
            <ActionButton
              variant="primary"
              icon="plus"
              text="Add User"
              on:click={openAddUserModal}
            />
          {/if}
          <ActionButton
            variant="outline"
            icon="refresh"
            text="Refresh"
            on:click={() => loadUsers()}
          />
        </div>
        <div class="form-control">
          <FormField
            type="text"
            placeholder="Search users..."
            bind:value={searchTerm}
            label=""
          />
        </div>
      </div>

      {#if isLoading}
        <div class="flex justify-center items-center py-8">
          <div class="loading loading-spinner loading-lg text-primary"></div>
        </div>
      {:else}
        <div class="overflow-x-auto">
          <table class="table table-zebra">
            <thead>
              <tr class="text-white/80">
                <th>User</th>
                <th>Email</th>
                <th>Role</th>
                <th>Last Login</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody class="text-white/70">
              {#each filteredUsers as user (user.id)}
                <tr class="hover:bg-white/5">
                  <td>
                    <div class="flex items-center gap-3">
                      <div class="avatar placeholder">
                        <div class="bg-primary text-primary-content rounded-full w-8 h-8">
                          <span class="text-xs">{user.first_name?.[0]}{user.last_name?.[0]}</span>
                        </div>
                      </div>
                      <div>
                        <div class="font-bold">{user.first_name} {user.last_name}</div>
                        <div class="text-sm opacity-50">@{user.username}</div>
                      </div>
                    </div>
                  </td>
                  <td>{user.email}</td>
                  <td>
                    <div class="badge {user.role === 'admin' ? 'badge-warning' : 'badge-info'} badge-sm">
                      {user.role}
                    </div>
                  </td>
                  <td>{formatDate(user.last_login)}</td>
                  <td>
                    <div class="badge {user.is_active ? 'badge-success' : 'badge-error'} badge-sm">
                      {user.is_active ? 'Active' : 'Inactive'}
                    </div>
                  </td>
                  <td>
                    <div class="flex gap-1">
                      <!-- View Permissions -->
                      <IconButton
                        icon="shield-check"
                        variant="ghost"
                        size="xs"
                        title="View Permissions"
                        on:click={() => openAccessControl(user)}
                      />
                      
                      {#if isAdmin}
                        <!-- Edit User -->
                        <IconButton
                          icon="edit"
                          variant="ghost"
                          size="xs"
                          title="Edit User"
                          on:click={() => openEditUserModal(user)}
                        />
                        
                        <!-- Toggle Status -->
                        {#if user.id !== currentUser.id}
                          <IconButton
                            icon={user.is_active ? "times" : "check"}
                            variant="ghost"
                            size="xs"
                            title={user.is_active ? 'Deactivate User' : 'Activate User'}
                            on:click={() => toggleUserStatus(user)}
                          />
                          
                          <!-- Delete User -->
                          <IconButton
                            icon="trash"
                            variant="danger"
                            size="xs"
                            title="Delete User"
                            on:click={() => handleDeleteUser(user)}
                          />
                        {/if}
                      {/if}
                    </div>
                  </td>
                </tr>
              {:else}
                <tr>
                  <td colspan="6" class="text-center py-8">
                    <div class="flex flex-col items-center gap-2">
                      <svg class="w-12 h-12 text-white/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
                      </svg>
                      <p class="text-white/60">
                        {searchTerm ? 'No users found matching your search' : 'No users found'}
                      </p>
                      {#if isAdmin && !searchTerm}
                        <ActionButton
                          variant="primary"
                          size="sm"
                          text="Add your first user"
                          on:click={openAddUserModal}
                        />
                      {/if}
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}

      {#if !isAdmin}
        <div class="alert alert-info mt-4">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <div>
            <h3 class="font-bold">Limited Access</h3>
            <div class="text-xs">You can view user information but cannot add, edit, or delete users. Contact an administrator for user management tasks.</div>
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>

<!-- User Modal -->
{#if showUserModal && currentUser}
  <UserModal 
    bind:show={showUserModal}
    user={selectedUser}
    currentUser={currentUser}
    on:close={() => showUserModal = false}
    on:create-user={handleCreateUser}
    on:update-user={handleUpdateUser}
  />
{/if}

<!-- Access Control Modal -->
{#if showAccessControl && selectedUserForAccess}
  <AccessControl 
    bind:show={showAccessControl}
    currentUser={selectedUserForAccess}
    on:close={() => showAccessControl = false}
  />
{/if}

<!-- Confirmation Modal -->
<ConfirmationModal
  show={showDeleteConfirm}
  title="Delete User"
  message={userToDelete ? `Are you sure you want to delete user "${userToDelete.username}"? This action cannot be undone.` : ''}
  confirmText="Delete"
  cancelText="Cancel"
  loading={confirmLoading}
  on:confirm={confirmDelete}
  on:cancel={cancelDelete}
/>