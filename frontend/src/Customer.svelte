<script>
  import { onMount } from 'svelte'
  import CustomerModal from './components/CustomerModal.svelte'

  let customers = []
  let loading = false
  let searchTerm = ''
  let showModal = false
  let editingCustomer = null
  let modalLoading = false

  onMount(async () => {
    await loadCustomers()
  })

  async function loadCustomers() {
    loading = true
    try {
      const response = await fetch('/api/customers')
      if (response.ok) {
        customers = await response.json()
      } else {
        console.error('Failed to load customers:', response.statusText)
      }
    } catch (error) {
      console.error('Error loading customers:', error)
    } finally {
      loading = false
    }
  }

  function openModal(customer = null) {
    editingCustomer = customer
    showModal = true
  }

  function closeModal() {
    showModal = false
    editingCustomer = null
    modalLoading = false
  }

  async function handleSaveCustomer(event) {
    const { customer, isEditing } = event.detail
    modalLoading = true

    try {
      const url = isEditing ? `/api/customers/${editingCustomer.id}` : '/api/customers'
      const method = isEditing ? 'PUT' : 'POST'
      
      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(customer)
      })

      if (response.ok) {
        await loadCustomers()
        closeModal()
      } else {
        console.error('Failed to save customer:', response.statusText)
      }
    } catch (error) {
      console.error('Error saving customer:', error)
    } finally {
      modalLoading = false
    }
  }

  async function deleteCustomer(customerId) {
    if (!confirm('Are you sure you want to delete this customer?')) {
      return
    }

    try {
      const response = await fetch(`/api/customers/${customerId}`, {
        method: 'DELETE'
      })

      if (response.ok) {
        await loadCustomers()
      } else {
        console.error('Failed to delete customer:', response.statusText)
      }
    } catch (error) {
      console.error('Error deleting customer:', error)
    }
  }

  $: filteredCustomers = customers.filter(customer => {
    if (!searchTerm) return true
    const search = searchTerm.toLowerCase()
    return (
      customer.name?.toLowerCase().includes(search) ||
      customer.name_arabic?.toLowerCase().includes(search) ||
      customer.email?.toLowerCase().includes(search) ||
      customer.phone?.toLowerCase().includes(search) ||
      customer.vat_number?.toLowerCase().includes(search)
    )
  })
</script>

<div class="page-container">
  <div class="page-header">
    <h1 class="page-title">Customers</h1>
    <div class="header-actions">
      <button class="btn btn-secondary" type="button">
        <i class="fas fa-upload"></i>
        Import
      </button>
      <button class="btn btn-primary" type="button" on:click={() => openModal()}>
        <i class="fas fa-plus"></i>
        Add Customer
      </button>
    </div>
  </div>

  <div class="glass-card">
    <div class="card-header">
      <div class="search-container">
        <i class="fas fa-search search-icon"></i>
        <input
          type="text"
          placeholder="Search customers..."
          class="search-input"
          bind:value={searchTerm}
        />
      </div>
    </div>

    <div class="table-container">
      {#if loading}
        <div class="loading-state">
          <i class="fas fa-spinner fa-spin"></i>
          Loading customers...
        </div>
      {:else if filteredCustomers.length === 0}
        <div class="empty-state">
          <i class="fas fa-users"></i>
          <h3>No customers found</h3>
          <p>
            {searchTerm ? 'No customers match your search criteria.' : 'Start by adding your first customer.'}
          </p>
          {#if !searchTerm}
            <button class="btn btn-primary" on:click={() => openModal()}>
              <i class="fas fa-plus"></i>
              Add Customer
            </button>
          {/if}
        </div>
      {:else}
        <table class="data-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Phone</th>
              <th>VAT Number</th>
              <th>City</th>
              <th>Country</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each filteredCustomers as customer (customer.id)}
              <tr>
                <td>
                  <div class="customer-name">
                    <div class="name-primary">{customer.name}</div>
                    {#if customer.name_arabic}
                      <div class="name-secondary">{customer.name_arabic}</div>
                    {/if}
                  </div>
                </td>
                <td>{customer.email || '-'}</td>
                <td>{customer.phone || '-'}</td>
                <td>{customer.vat_number || '-'}</td>
                <td>
                  {#if customer.city}
                    <div class="city-info">
                      <div>{customer.city}</div>
                      {#if customer.city_arabic}
                        <div class="city-arabic">{customer.city_arabic}</div>
                      {/if}
                    </div>
                  {:else}
                    -
                  {/if}
                </td>
                <td>
                  {#if customer.country}
                    <div class="country-info">
                      <div>{customer.country}</div>
                      {#if customer.country_arabic}
                        <div class="country-arabic">{customer.country_arabic}</div>
                      {/if}
                    </div>
                  {:else}
                    -
                  {/if}
                </td>
                <td>
                  <div class="action-buttons">
                    <button
                      class="btn btn-sm btn-secondary"
                      on:click={() => openModal(customer)}
                      title="Edit customer"
                    >
                      <i class="fas fa-edit"></i>
                    </button>
                    <button
                      class="btn btn-sm btn-danger"
                      on:click={() => deleteCustomer(customer.id)}
                      title="Delete customer"
                    >
                      <i class="fas fa-trash"></i>
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </div>
  </div>
</div>

<CustomerModal
  show={showModal}
  {editingCustomer}
  loading={modalLoading}
  on:save={handleSaveCustomer}
  on:close={closeModal}
/>

<style>
  .customer-name .name-primary {
    font-weight: 500;
    color: var(--text-primary);
  }

  .customer-name .name-secondary {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-top: 2px;
  }

  .city-info, .country-info {
    line-height: 1.4;
  }

  .city-arabic, .country-arabic {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-top: 2px;
  }

  .action-buttons {
    display: flex;
    gap: 0.5rem;
  }

  .loading-state {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    color: var(--text-secondary);
    gap: 0.75rem;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    text-align: center;
  }

  .empty-state i {
    font-size: 3rem;
    color: var(--text-muted);
    margin-bottom: 1rem;
  }

  .empty-state h3 {
    margin: 0 0 0.5rem 0;
    color: var(--text-primary);
  }

  .empty-state p {
    margin: 0 0 1.5rem 0;
    color: var(--text-secondary);
  }
</style>